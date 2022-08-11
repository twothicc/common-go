package grpcserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/soheilhy/cmux"
	"github.com/twothicc/common-go/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	configs    *ServerConfigs
	grpcServer *grpc.Server
	httpServer *http.Server
	connMux    cmux.CMux
}

// InitAndRunGrpcServer - Configures, initializes and runs the grpc server,
// along with prometheus monitoring if specified.
func InitAndRunGrpcServer(ctx context.Context, config *ServerConfigs) {
	server := InitGrpcServer(ctx, config)

	go server.ListenSignals(ctx)
	server.Run(ctx)
}

// InitGrpcServer - configures and initializes a grpc server
func InitGrpcServer(ctx context.Context, config *ServerConfigs) *Server {
	serverOptions := parseServerOptions(ctx, config)
	grpcServer := grpc.NewServer(serverOptions...)

	for _, registerServerHandler := range config.registerServerHandlers {
		registerServerHandler(grpcServer)
	}

	var httpServer *http.Server

	if !config.disableProm {
		grpc_prometheus.Register(grpcServer)

		httpHandler := promhttp.Handler()

		http.Handle("/metrics", httpHandler)

		httpServer = &http.Server{
			Handler:           httpHandler,
			ReadHeaderTimeout: HTTP_READ_HEADER_TIMEOUT,
		}
	}

	return &Server{
		grpcServer: grpcServer,
		httpServer: httpServer,
		configs:    config,
	}
}

// Run - starts the grpc server
func (g *Server) Run(ctx context.Context) {
	logger.WithContext(ctx).Info("start grpc server")

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", g.configs.port))
	if err != nil {
		logger.WithContext(ctx).Fatal("failed to listen", zap.Error(err))
	}

	// cmux multiplexes connections based on payload, allowing various protocols to run on the same TCP listener
	m := cmux.New(lis)
	g.connMux = m
	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())

	go func() {
		if err := g.grpcServer.Serve(grpcL); err != nil {
			logger.WithContext(ctx).Error("fail to serve grpc server", zap.Error(err))
		}
	}()

	go func() {
		if err := g.httpServer.Serve(httpL); err != nil {
			logger.WithContext(ctx).Error("fail to serve http server", zap.Error(err))
		}
	}()

	if err := m.Serve(); err != nil {
		logger.WithContext(ctx).Error("failed to serve grpc and http server", zap.Error(err))
	}
}

// ListenSignals - listens for os signals to gracefully stop server
func (g *Server) ListenSignals(ctx context.Context) {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-signalChan

	logger.WithContext(ctx).Info("receive signal, stop server", zap.String("signal", sig.String()))
	time.Sleep(1 * time.Second)

	if g.connMux != nil {
		g.connMux.Close()
		logger.WithContext(ctx).Info("stop cmux server")
	}

	if g.grpcServer != nil {
		g.grpcServer.GracefulStop()
		logger.WithContext(ctx).Info("stop grpc server")
	}

	if g.httpServer != nil {
		if err := g.httpServer.Shutdown(ctx); err != nil {
			logger.WithContext(ctx).Error("fail to gracefully shutdown http server", zap.Error(err))
		} else {
			logger.WithContext(ctx).Info("stop http server")
		}
	}

	logger.Sync()
	logger.WithContext(ctx).Info("stop server gracefully")
}

func parseServerOptions(ctx context.Context, config *ServerConfigs) []grpc.ServerOption {
	keepAliveParams := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: config.maxIdleConn,
		Timeout:           config.timeout,
		Time:              config.keepAliveInterval,
	})

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_zap.UnaryServerInterceptor(logger.WithContext(ctx)),
		grpc_recovery.UnaryServerInterceptor(),
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(),
		grpc_zap.StreamServerInterceptor(logger.WithContext(ctx)),
		grpc_recovery.StreamServerInterceptor(),
	}

	if !config.disableProm {
		unaryInterceptors = insertIntoUnaryServerInterceptors(
			unaryInterceptors,
			grpc_prometheus.UnaryServerInterceptor,
			PROMETHEUS_INTERCEPTOR_IDX,
		)

		streamInterceptors = insertIntoStreamServerInterceptors(
			streamInterceptors,
			grpc_prometheus.StreamServerInterceptor,
			PROMETHEUS_INTERCEPTOR_IDX,
		)
	}

	return []grpc.ServerOption{
		keepAliveParams,
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	}
}
