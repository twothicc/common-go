package pool

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_pool "github.com/processout/grpc-go-pool"
	"github.com/twothicc/common-go/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// allowOverwrite indicates whether existing connection pool and its connections
// should be closed and overwritten by the new connection pool to be created.
type PoolCreatorFunc func(
	ctx context.Context,
	selector *PoolSelector,
	allowOverwrite bool,
) error

// PoolCreator - creates a PoolCreatorFunc that handles creating and setting
// a connection pool for a specific server to a PoolSelector.
//
// configs specified will take precedence over the PoolSelector's default connection configs.
//
// extraUnaryClientInterceptors specified will be appended to the default unaryClientIntercetor chain.
//
// extraStreamClientInterceptors specified will be appended to the default streamClientIntercetor chain.
func PoolCreator(
	configs *ConnPoolConfigs,
	extraUnaryClientInterceptors []grpc.UnaryClientInterceptor,
	extraStreamClientInterceptors []grpc.StreamClientInterceptor,
) PoolCreatorFunc {
	return func(
		ctx context.Context,
		selector *PoolSelector,
		allowOverwrite bool,
	) error {
		connFactory := func(ctx context.Context) (*grpc.ClientConn, error) {
			ctx, cancel := context.WithTimeout(ctx, configs.CreateTimeout)
			defer cancel()

			logger.WithContext(ctx).Debug("creating connection", zap.String("server", configs.Server))

			unaryClientInterceptors := selector.defaultUnaryClientInterceptors
			unaryClientInterceptors = append(unaryClientInterceptors, extraUnaryClientInterceptors...)

			streamClientInterceptors := selector.defaultStreamClientInterceptors
			streamClientInterceptors = append(streamClientInterceptors, extraStreamClientInterceptors...)

			dialOptions := []grpc.DialOption{
				grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(unaryClientInterceptors...)),
				grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(streamClientInterceptors...)),
			}

			if !configs.EnableTLS {
				dialOptions = append(dialOptions,
					grpc.WithTransportCredentials(insecure.NewCredentials()),
					grpc.WithBlock(),
				)
			}

			conn, err := grpc.DialContext(ctx, configs.Server, dialOptions...)
			if err != nil {
				return nil, err
			}

			return conn, nil
		}

		pool, err := grpc_pool.NewWithContext(
			ctx,
			connFactory,
			configs.InitConn,
			configs.MaxConn,
			configs.IdleTimeout,
			configs.MaxLifeDuration,
		)

		if err != nil {
			return err
		}

		selector.mu.Lock()

		existingPool, ok := selector.pools[configs.Server]
		if !ok {
			selector.pools[configs.Server] = pool
		} else if ok {
			if allowOverwrite {
				logger.WithContext(ctx).Debug("overwriting connection pool", zap.String("server", configs.Server))
				existingPool.Close()
				selector.pools[configs.Server] = pool
			}
		}
		selector.mu.Unlock()

		return nil
	}
}
