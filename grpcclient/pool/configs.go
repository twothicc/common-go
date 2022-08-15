package pool

import "time"

type ConnConfigs struct {
	IdleTimeout     time.Duration
	CreateTimeout   time.Duration // timeout for establishing connection
	MaxLifeDuration time.Duration
	InitConn        int
	MaxConn         int
	EnableTLS       bool
}

type ConnPoolConfigs struct {
	*ConnConfigs
	Server string
}

func GetConnConfigs(
	idleTimeout, createTimeout, maxLifeDuration time.Duration,
	init, capacity int,
	enableTLS bool,
) *ConnConfigs {
	return &ConnConfigs{
		IdleTimeout:     idleTimeout,
		CreateTimeout:   createTimeout,
		MaxLifeDuration: maxLifeDuration,
		InitConn:        init,
		MaxConn:         capacity,
		EnableTLS:       enableTLS,
	}
}

func GetDefaultConnConfigs() *ConnConfigs {
	return &ConnConfigs{
		IdleTimeout:     DEFAULT_IDLE_TIMEOUT,
		CreateTimeout:   DEFAULT_CREATE_TIMEOUT,
		MaxLifeDuration: DEFAULT_MAX_LIFE_DURATION,
		InitConn:        DEFAULT_INIT_CONN,
		MaxConn:         DEFAULT_MAX_CONN,
		EnableTLS:       DEFAULT_ENABLE_TLS,
	}
}

func GetConnPoolConfigs(
	server string,
	idleTimeout, createTimeout, maxLifeDuration time.Duration,
	init, capacity int,
	enableTLS bool,
) *ConnPoolConfigs {
	return &ConnPoolConfigs{
		Server: server,
		ConnConfigs: GetConnConfigs(
			idleTimeout, createTimeout, maxLifeDuration,
			init, capacity,
			enableTLS,
		),
	}
}

func GetDefaultConnPoolConfigs(
	server string,
) *ConnPoolConfigs {
	return &ConnPoolConfigs{
		Server:      server,
		ConnConfigs: GetDefaultConnConfigs(),
	}
}
