package pool

import "time"

type ConnConfigs struct {
	IdleTimeout     time.Duration
	CreateTimeout   time.Duration // timeout for establishing connection
	MaxLifeDuration time.Duration
	InitConn        int
	MaxConn         int
}

type ConnPoolConfigs struct {
	*ConnConfigs
	Server string
}

func GetConnConfigs(
	idleTimeout, createTimeout, maxLifeDuration time.Duration,
	init, capacity int,
) *ConnConfigs {
	return &ConnConfigs{
		IdleTimeout:     idleTimeout,
		CreateTimeout:   createTimeout,
		MaxLifeDuration: maxLifeDuration,
		InitConn:        init,
		MaxConn:         capacity,
	}
}

func GetDefaultConnConfigs() *ConnConfigs {
	return &ConnConfigs{
		IdleTimeout:     DEFAULT_IDLE_TIMEOUT,
		CreateTimeout:   DEFAULT_CREATE_TIMEOUT,
		MaxLifeDuration: DEFAULT_MAX_LIFE_DURATION,
		InitConn:        DEFAULT_INIT_CONN,
		MaxConn:         DEFAULT_MAX_CONN,
	}
}

func GetConnPoolConfigs(
	server string,
	idleTimeout, createTimeout, maxLifeDuration time.Duration,
	init, capacity int,
) *ConnPoolConfigs {
	return &ConnPoolConfigs{
		Server: server,
		ConnConfigs: GetConnConfigs(
			idleTimeout, createTimeout, maxLifeDuration,
			init, capacity,
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
