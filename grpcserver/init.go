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
	"github.com/twothicc/common-go/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type grpcServer struct {
	configs *ServerConfigs
	server  *grpc.Server
}

// InitAndRunGrpcServer - Configures, initializes and runs the grpc server,
// along with prometheus monitoring if specified.
func InitAndRunGrpcServer(ctx context.Context, config *ServerConfigs) {
	server := InitGrpcServer(ctx, config)

	if !server.configs.disableProm {
		go server.StartPrometheusMonitoring(ctx)
	}

	go server.ListenSignals(ctx)
	server.Run(ctx)
}

// InitGrpcServer - configures and initializes a grpc server
func InitGrpcServer(ctx context.Context, config *ServerConfigs) *grpcServer {
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
		insertIntoUnaryServerInterceptors(
			unaryInterceptors,
			grpc_prometheus.UnaryServerInterceptor,
			PROMETHEUS_INTERCEPTOR_IDX,
		)
		insertIntoStreamServerInterceptors(
			streamInterceptors,
			grpc_prometheus.StreamServerInterceptor,
			PROMETHEUS_INTERCEPTOR_IDX,
		)
	}

	serverOptions := []grpc.ServerOption{
		keepAliveParams,
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	}

	server := grpc.NewServer(serverOptions...)

	for _, registerServerHandler := range config.registerServerHandlers {
		registerServerHandler(server)
	}

	return &grpcServer{
		server:  server,
		configs: config,
	}
}

// Run - starts the grpc server
func (g *grpcServer) Run(ctx context.Context) {
	logger.WithContext(ctx).Info("start grpc server")

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", g.configs.port))
	if err != nil {
		logger.WithContext(ctx).Fatal("failed to listen", zap.Error(err))
	}

	if err := g.server.Serve(lis); err != nil {
		logger.WithContext(ctx).Fatal("failed to init grpc server", zap.Error(err))
	}
}

// ListenSignals - listens for os signals to gracefully stop grpc server
func (g *grpcServer) ListenSignals(ctx context.Context) {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-signalChan

	logger.WithContext(ctx).Info("receive signal, stop server", zap.String("signal", sig.String()))
	time.Sleep(1 * time.Second)

	if g.server != nil {
		g.server.GracefulStop()
	}

	logger.Sync()
}

// StartPrometheusMonitoring
func (g *grpcServer) StartPrometheusMonitoring(ctx context.Context) {
	grpc_prometheus.Register(g.server)
	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(
		fmt.Sprintf("localhost:%s", PROMETHEUS_PORT),
		nil,
	); err != nil {
		logger.WithContext(ctx).Error("fail to start prometheus monitor", zap.Error(err))
	}
}
