package grpcclient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/twothicc/common-go/commonerror"
	"github.com/twothicc/common-go/grpcclient/pool"
	"github.com/twothicc/common-go/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	opentracing "github.com/opentracing/opentracing-go"
	grpc_pool "github.com/processout/grpc-go-pool"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

type Client struct {
	Pools        *pool.PoolSelector
	configs      *clientConfigs
	tracerCloser io.Closer
}

func NewClient(
	ctx context.Context,
	configs *clientConfigs,
) *Client {
	unaryClientInterceptors, streamClientInterceptors, tracerCloser := parseInterceptors(ctx, configs)

	return &Client{
		Pools: pool.NewPoolSelector(
			ctx,
			unaryClientInterceptors,
			streamClientInterceptors,
			configs.poolCreators,
		),
		configs:      configs,
		tracerCloser: tracerCloser,
	}
}

func (gc *Client) Call(
	ctx context.Context,
	server, fullMethod string,
	req interface{},
	resp interface{},
) error {
	if gc == nil {
		return commonerror.New(commonerror.ErrCodeServer, "grpc client not initialized")
	}

	conn, err := gc.Pools.Get(ctx, server, true)
	if err != nil {
		logger.WithContext(ctx).Debug("fail to get connection pool", zap.String("server", server))

		code := commonerror.ErrCodeGRPC
		msg := fmt.Sprintf("DialContext error, server = %s, err = %v", server, err)

		if errors.Is(err, context.DeadlineExceeded) {
			code = commonerror.ErrCodeTimeout
		}

		return commonerror.New(int32(code), msg)
	}

	defer returnOrCloseConnection(ctx, server, conn)

	err = conn.Invoke(ctx, fullMethod, req, resp)
	if err != nil {
		return commonerror.Convert(err)
	}

	return nil
}

// ListenSignals - listens for os signals and closes client if necessary.
func (gc *Client) ListenSignals(ctx context.Context) {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-signalChan

	logger.WithContext(ctx).Info("receive signal, stop server", zap.String("signal", sig.String()))
	time.Sleep(1 * time.Second)

	if gc.tracerCloser != nil {
		if err := gc.tracerCloser.Close(); err != nil {
			logger.WithContext(ctx).Error("fail to close jaeger tracer")
		} else {
			logger.WithContext(ctx).Info("jaeger tracer closed")
		}
	}

	if gc.Pools != nil {
		gc.Pools.Close()
		logger.WithContext(ctx).Info("all connections closed and grpc client stopped")
	}

	logger.Sync()
}

// returnOrCloseConnection - returns connection obj to pool, or close underlying connection if pool is full
func returnOrCloseConnection(ctx context.Context, server string, conn *grpc_pool.ClientConn) {
	if err := conn.Close(); err != nil {
		if errors.Is(err, grpc_pool.ErrFullPool) {
			logger.WithContext(ctx).Debug("pool capacity reached, closing connection", zap.String("server", server))
			conn.ClientConn.Close()
		} else {
			logger.WithContext(ctx).Error("Fail to return connection to connection pool",
				zap.String("server", server),
			)
		}
	}
}

func parseInterceptors(
	ctx context.Context,
	configs *clientConfigs,
) (
	unaryClientInterceptors []grpc.UnaryClientInterceptor,
	streamClientInterceptors []grpc.StreamClientInterceptor,
	tracerCloser io.Closer) {
	if !opentracing.IsGlobalTracerRegistered() {
		tracerCfg := jaegercfg.Configuration{}

		var err error

		if tracerCloser, err = tracerCfg.InitGlobalTracer(
			configs.serviceName,
		); err != nil {
			logger.WithContext(ctx).Error("fail to initialize jaeger tracer", zap.Error(err))
		}
	}

	unaryClientInterceptors = []grpc.UnaryClientInterceptor{
		grpc_opentracing.UnaryClientInterceptor(),
		grpc_zap.UnaryClientInterceptor(logger.WithContext(ctx)),
	}

	streamClientInterceptors = []grpc.StreamClientInterceptor{
		grpc_opentracing.StreamClientInterceptor(),
		grpc_zap.StreamClientInterceptor(logger.WithContext(ctx)),
	}

	if !configs.disableProm {
		unaryClientInterceptors = insertIntoUnaryClientInterceptors(
			unaryClientInterceptors,
			grpc_prometheus.UnaryClientInterceptor,
			PROMETHEUS_INTERCEPTOR_IDX,
		)

		streamClientInterceptors = insertIntoStreamClientInterceptors(
			streamClientInterceptors,
			grpc_prometheus.StreamClientInterceptor,
			PROMETHEUS_INTERCEPTOR_IDX,
		)
	}

	return unaryClientInterceptors, streamClientInterceptors, tracerCloser
}
