package commonerror

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ICommonError interface {
	error
	Code() int32
	Msg() string
}

type CommonError struct {
	code int32
	msg  string
}

func (ce *CommonError) Error() string {
	return fmt.Sprintf(ErrorFormat, ce.code, ce.msg)
}

func (ce *CommonError) Code() int32 {
	return ce.code
}

func (ce *CommonError) Msg() string {
	return ce.msg
}

func New(code int32, msg string) ICommonError {
	if code == CodeOk {
		return nil
	}

	return &CommonError{
		code: code,
		msg:  msg,
	}
}

func Convert(err error) ICommonError {
	if err == nil {
		return &CommonError{}
	}

	if ce, ok := err.(*CommonError); ok {
		return ce
	}

	status := status.Convert(err)
	code, msg := status.Code(), status.Message()
	errCode := grpcToCommonErrCode(code)

	return &CommonError{
		code: errCode,
		msg:  msg,
	}
}

func grpcToCommonErrCode(code codes.Code) int32 {
	commonErrCode := int32(code)

	switch code {
	case codes.Unknown:
		commonErrCode = ErrCodeUnknown
	case codes.Internal:
		commonErrCode = ErrCodeServer
	case codes.DeadlineExceeded:
		commonErrCode = ErrCodeTimeout
	default:
	}

	return commonErrCode
}
