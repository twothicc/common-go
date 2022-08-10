package logger

// config constants
const (
	CONSOLE_SEPARATOR = "|"
)

// log file constants
const (
	LOG_FILENAME = "server.log"
)

type LoggerKey int

// logger context key
const (
	LOGGER_KEY LoggerKey = iota
)
