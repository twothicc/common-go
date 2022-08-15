package pool

import "time"

const (
	DEFAULT_IDLE_TIMEOUT      = 5 * time.Minute
	DEFAULT_CREATE_TIMEOUT    = 3 * time.Second
	DEFAULT_MAX_LIFE_DURATION = 1 * time.Hour
	DEFAULT_INIT_CONN         = 0
	DEFAULT_MAX_CONN          = 5
	DEFAULT_ENABLE_TLS        = false
)
