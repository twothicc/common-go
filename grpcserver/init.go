package grpcserver

import (
	"context"
	"errors"
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
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/twothicc/common-go/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// Server - contains fields necessary for initializing and running servers
type Server struct {
	configs    *ServerConfigs
	grpcServer *grpc.Server
	httpServer *http.Server
}

// InitAndRunGrpcServer - initializes and runs the grpc server,
// Also initializes and runs http server for prometheus monitoring if specified.
func InitAndRunGrpcServer(ctx context.Context, config *ServerConfigs) {
	server := InitGrpcServer(ctx, config)

	go server.ListenSignals(ctx)
	server.Run(ctx)
}

// InitGrpcServer - initializes a grpc server.
// Also initializes a http server for prometheus monitoring if specified.
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
			Addr:              fmt.Sprintf("%s:%s", config.domain, PROMETHEUS_PORT),
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

// Run - starts the grpc server.
// Also starts http server for prometheus monitoring if specified.
func (g *Server) Run(ctx context.Context) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", g.configs.domain, g.configs.port))
	if err != nil {
		logger.WithContext(ctx).Fatal("fail to listen", zap.Error(err))
	}

	if g.httpServer != nil {
		go func() {
			logger.WithContext(ctx).Info("start http server")

			if err := g.httpServer.ListenAndServe(); err != nil {
				if !errors.Is(err, http.ErrServerClosed) {
					logger.WithContext(ctx).Error("fail to serve http server", zap.Error(err))
				}
			}
		}()
	}

	if g.grpcServer != nil {
		logger.WithContext(ctx).Info("start grpc server")

		if err := g.grpcServer.Serve(lis); err != nil {
			logger.WithContext(ctx).Error("fail to serve grpc server", zap.Error(err))
		}
	} else {
		logger.WithContext(ctx).Error("missing grpc server")
	}
}

// ListenSignals - listens for os signals to gracefully stop server.
// http server for prometheus monitoring is first stopped, followed by grpc server.
func (g *Server) ListenSignals(ctx context.Context) {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-signalChan

	logger.WithContext(ctx).Info("receive signal, stop server", zap.String("signal", sig.String()))
	time.Sleep(1 * time.Second)

	if g.httpServer != nil {
		if err := g.httpServer.Shutdown(ctx); err != nil {
			logger.WithContext(ctx).Error("fail to gracefully stop http server", zap.Error(err))
		} else {
			logger.WithContext(ctx).Info("http server gracefully stopped")
		}
	}

	if g.grpcServer != nil {
		g.grpcServer.GracefulStop()
		logger.WithContext(ctx).Info("grpc server gracefully stopped")
	}

	logger.Sync()
}

func parseServerOptions(ctx context.Context, config *ServerConfigs) []grpc.ServerOption {
	keepAliveParams := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: config.maxIdleConn,
		Timeout:           config.timeout,
		Time:              config.keepAliveInterval,
	})

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_opentracing.UnaryServerInterceptor(),
		grpc_zap.UnaryServerInterceptor(logger.WithContext(ctx)),
		grpc_recovery.UnaryServerInterceptor(),
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(),
		grpc_opentracing.StreamServerInterceptor(),
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
