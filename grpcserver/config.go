package grpcserver

import (
	"time"

	"google.golang.org/grpc"
)

type ServerConfigs struct {
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
	domain, port string,
	timeout time.Duration,
	maxIdleConn time.Duration,
	keepAliveInterval time.Duration,
	isTest bool,
	disableProm bool,
	registerServerHandlers ...RegisterServerHandler,
) *ServerConfigs {
	return &ServerConfigs{
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
	domain, port string,
	isTest bool,
	registerServerHandlers ...RegisterServerHandler,
) *ServerConfigs {
	return &ServerConfigs{
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
