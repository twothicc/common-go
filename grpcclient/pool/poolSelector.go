package pool

import (
	"context"
	"sync"

	grpc_pool "github.com/processout/grpc-go-pool"
	"github.com/twothicc/common-go/commonerror"
	"github.com/twothicc/common-go/logger"
	"go.uber.org/zap"

	"google.golang.org/grpc"
)

// PoolSelector - selects a connection pool by server
type PoolSelector struct {
	pools                           map[string]*grpc_pool.Pool
	defaultConnConfigs              *ConnConfigs
	defaultUnaryClientInterceptors  []grpc.UnaryClientInterceptor
	defaultStreamClientInterceptors []grpc.StreamClientInterceptor
	mu                              sync.RWMutex
}

// NewPoolSelector - creates a connection pool selector.
//
// poolCreators provide a way for users to provide configurations and add grpc call options
// specifically for each connection pool. It is not necessary to provide, but
// not providing means the default connection configs will be used.
//
// if multiple poolCreators intend to add connection pools for a server, only
// the first poolCreator will succeed in adding its connection pool.
//
// defaultUnaryClientInterceptors and defaultStreamClientInterceptors provided will
// be set up for every connection created.
func NewPoolSelector(
	ctx context.Context,
	defaultUnaryClientInterceptors []grpc.UnaryClientInterceptor,
	defaultStreamClientInterceptors []grpc.StreamClientInterceptor,
	poolCreators []PoolCreatorFunc,
) *PoolSelector {
	selector := &PoolSelector{
		pools:                           make(map[string]*grpc_pool.Pool),
		defaultConnConfigs:              GetDefaultConnConfigs(),
		defaultUnaryClientInterceptors:  defaultUnaryClientInterceptors,
		defaultStreamClientInterceptors: defaultStreamClientInterceptors,
	}

	for _, creatorFunc := range poolCreators {
		if err := creatorFunc(ctx, selector, false); err != nil {
			logger.WithContext(ctx).Error("fail to create and set connection pool", zap.Error(err))
		}
	}

	return selector
}

// Close - closes all existing connections.
//
// While Close() is in effect, no connections may be used or created.
func (ps *PoolSelector) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, pool := range ps.pools {
		pool.Close()
	}

	ps.pools = nil
}

// Get - if connection pool for server exists, returns an existing connection object
// if possible, otherwise creates a new connection object.
//
// Take special note that the returned connection object embeds grpc.ClientConn. It's
// Close() method is overridden to return the connection to the pool, so it is necessary
// to call Close() after it is invoked.
//
// server: <domain>:<port>
//
// createIfNotExist: Indicates whether to create new connection pool for
// the server (if not existing) and return a connection from the pool. New connection
// pool will use this PoolSelector's default connection configs.
func (ps *PoolSelector) Get(
	ctx context.Context,
	server string,
	createIfNotExist bool,
) (*grpc_pool.ClientConn, error) {
	ps.mu.RLock()
	pool := ps.pools[server]
	ps.mu.RUnlock()

	if pool != nil {
		clientConn, err := pool.Get(ctx)
		if err != nil {
			return nil, err
		}

		return clientConn, nil
	}

	if !createIfNotExist {
		logger.WithContext(ctx).Debug("missing connection pool", zap.String("server", server))
		return nil, commonerror.New(commonerror.ErrCodeServer, "pool not initialized")
	}

	err := ps.SetPool(ctx, ps.getDefaultConnPoolConfigs(server), nil, nil, false)
	if err != nil {
		logger.WithContext(ctx).Debug("fail to set connection pool", zap.String("server", server))
		return nil, commonerror.New(commonerror.ErrCodeServer, "fail to initialize pool")
	}

	return ps.Get(ctx, server, false)
}

// SetPool - set adds a new connection pool based on configs and
// any extra interceptors.
//
// allowOverride: should be set to false. Setting to true could cause
// existing connection pools to the same server to close and be overwritten.
func (ps *PoolSelector) SetPool(
	ctx context.Context,
	configs *ConnPoolConfigs,
	extraUnaryClientInterceptors []grpc.UnaryClientInterceptor,
	extraStreamClientInterceptors []grpc.StreamClientInterceptor,
	allowOverride bool,
) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	return PoolCreator(
		configs,
		extraUnaryClientInterceptors,
		extraStreamClientInterceptors,
	)(ctx, ps, allowOverride)
}

// SetDefaultConnConfigs - changes the default connection configs of this PoolSelector.
func (ps *PoolSelector) SetDefaultConnConfigs(
	connConfigs *ConnConfigs,
) {
	ps.defaultConnConfigs = connConfigs
}

// getDefaultConnPoolConfigs - gets a connection pool configs with this PoolSelector's
// default connection configs.
func (ps *PoolSelector) getDefaultConnPoolConfigs(server string) *ConnPoolConfigs {
	return GetConnPoolConfigs(
		server,
		ps.defaultConnConfigs.IdleTimeout,
		ps.defaultConnConfigs.CreateTimeout,
		ps.defaultConnConfigs.MaxLifeDuration,
		ps.defaultConnConfigs.InitConn,
		ps.defaultConnConfigs.MaxConn,
	)
}
