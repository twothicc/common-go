package logger

// config constants
const (
	CONSOLE_SEPARATOR = "|"
)

// log file constants
const (
	LOG_FILENAME   = "server.log"
	LOG_PERMISSION = 0o644
)

type LoggerKey int

// logger context key
const (
	LOGGER_KEY LoggerKey = iota
)
