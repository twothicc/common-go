package grpcclient

import (
	"time"

	"github.com/twothicc/common-go/grpcclient/pool"
)

type clientConfigs struct {
	defaultConnConfigs *pool.ConnConfigs
	serviceName        string
	domain             string
	port               string
	poolCreators       []pool.PoolCreatorFunc
	isTest             bool
	disableProm        bool
}

func GetClientConfigs(
	serviceName, domain, port string,
	isTest, disableProm bool,
	idleTimeout, createTimeout, maxLifeDuration time.Duration,
	init, capacity int,
	poolCreators []pool.PoolCreatorFunc,
) *clientConfigs {
	return &clientConfigs{
		domain:      domain,
		port:        port,
		isTest:      isTest,
		disableProm: disableProm,
		defaultConnConfigs: pool.GetConnConfigs(
			idleTimeout, createTimeout, maxLifeDuration,
			init, capacity,
		),
		poolCreators: poolCreators,
	}
}

func GetDefaultClientConfigs(
	serviceName, domain, port string,
	isTest bool,
	poolCreators []pool.PoolCreatorFunc,
) *clientConfigs {
	return &clientConfigs{
		domain:             domain,
		port:               port,
		isTest:             isTest,
		disableProm:        false,
		defaultConnConfigs: pool.GetDefaultConnConfigs(),
		poolCreators:       poolCreators,
	}
}
