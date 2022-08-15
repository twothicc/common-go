package grpcclient

import (
	"time"

	"github.com/twothicc/common-go/grpcclient/pool"
)

type clientConfigs struct {
	defaultConnConfigs *pool.ConnConfigs
	serviceName        string
	poolCreators       []pool.PoolCreatorFunc
	isTest             bool
}

func GetClientConfigs(
	serviceName string,
	isTest bool,
	idleTimeout, createTimeout, maxLifeDuration time.Duration,
	init, capacity int,
	enableTLS bool,
	poolCreators ...pool.PoolCreatorFunc,
) *clientConfigs {
	return &clientConfigs{
		serviceName: serviceName,
		isTest:      isTest,
		defaultConnConfigs: pool.GetConnConfigs(
			idleTimeout, createTimeout, maxLifeDuration,
			init, capacity,
			enableTLS,
		),
		poolCreators: poolCreators,
	}
}

func GetDefaultClientConfigs(
	serviceName string,
	isTest bool,
	poolCreators ...pool.PoolCreatorFunc,
) *clientConfigs {
	return &clientConfigs{
		serviceName:        serviceName,
		isTest:             isTest,
		defaultConnConfigs: pool.GetDefaultConnConfigs(),
		poolCreators:       poolCreators,
	}
}
