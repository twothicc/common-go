package commonerror

// Common Error Codes
const (
	ErrCodeGRPC    = 1
	ErrCodeServer  = 2
	ErrCodeUnknown = 3
	ErrCodeTimeout = 4
)

const (
	ErrMsgServer  = "server error"
	ErrMsgUnknown = "unknown error"
	ErrMsgTimeout = "request timed out"
)
