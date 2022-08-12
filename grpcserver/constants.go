package grpcserver

import "time"

const (
	PROMETHEUS_INTERCEPTOR_IDX = 2
	PROMETHEUS_PORT            = "9090"
)

const (
	DEFAULT_KEEPALIVE_TIMEOUT  = 10 * time.Second
	DEFAULT_MAX_IDLE_CONN      = 5 * time.Minute
	DEFAULT_KEEPALIVE_INTERVAL = 1 * time.Hour
)

const (
	HTTP_READ_HEADER_TIMEOUT = 20 * time.Second
)
