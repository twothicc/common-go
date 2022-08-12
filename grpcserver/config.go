package grpcserver

import (
	"time"

	"google.golang.org/grpc"
)

type ServerConfigs struct {
	serviceName            string
	domain                 string
	port                   string
	registerServerHandlers []RegisterServerHandler
	timeout                time.Duration
	maxIdleConn            time.Duration
	keepAliveInterval      time.Duration
	isTest                 bool
	disableProm            bool
}

type RegisterServerHandler func(s *grpc.Server)

func GetServerConfigs(
	serviceName, domain, port string,
	timeout time.Duration,
	maxIdleConn time.Duration,
	keepAliveInterval time.Duration,
	isTest bool,
	disableProm bool,
	registerServerHandlers ...RegisterServerHandler,
) *ServerConfigs {
	return &ServerConfigs{
		serviceName:            serviceName,
		domain:                 domain,
		port:                   port,
		timeout:                timeout,
		maxIdleConn:            maxIdleConn,
		keepAliveInterval:      keepAliveInterval,
		isTest:                 isTest,
		disableProm:            disableProm,
		registerServerHandlers: registerServerHandlers,
	}
}

func GetDefaultServerConfigs(
	serviceName, domain, port string,
	isTest bool,
	registerServerHandlers ...RegisterServerHandler,
) *ServerConfigs {
	return &ServerConfigs{
		serviceName:            serviceName,
		domain:                 domain,
		port:                   port,
		timeout:                DEFAULT_KEEPALIVE_TIMEOUT,
		maxIdleConn:            DEFAULT_MAX_IDLE_CONN,
		keepAliveInterval:      DEFAULT_KEEPALIVE_INTERVAL,
		isTest:                 isTest,
		disableProm:            false,
		registerServerHandlers: registerServerHandlers,
	}
}
