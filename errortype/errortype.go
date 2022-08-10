package errortype

import (
	"errors"
	"fmt"
)

type IError interface {
	error
	Msg() string
}

type Error struct {
	detail ErrorType
	msg    string
}

type ErrorType struct {
	code int32
	pkg  string
}

func (e *Error) Error() string {
	return fmt.Sprintf(ErrFormat, e.detail.code, e.detail.pkg, e.msg)
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) Code() int32 {
	return e.detail.code
}

func (e *Error) Pkg() string {
	return e.detail.pkg
}

func (e ErrorType) New(msg string) IError {
	return &Error{
		detail: e,
		msg:    msg,
	}
}

func (e ErrorType) Is(err error) bool {
	otherErr := &Error{}

	if !errors.As(err, &otherErr) {
		return false
	}

	return otherErr.detail.code == e.code && otherErr.detail.pkg == e.pkg
}

func (e ErrorType) Wrap(err error) IError {
	otherErr := &Error{}

	if !errors.As(err, &otherErr) {
		return e.New(err.Error())
	}

	if e.code == otherErr.detail.code &&
		e.pkg == otherErr.detail.pkg {
		return otherErr
	}

	return &Error{
		detail: e,
		msg:    fmt.Sprintf("%s | %s", otherErr.msg, otherErr.Error()),
	}
}

func (e ErrorType) WrapWithMsg(err error, msg string) IError {
	otherErr := &Error{}

	if !errors.As(err, &otherErr) {
		return e.New(err.Error())
	}

	if e.code == otherErr.detail.code &&
		e.pkg == otherErr.detail.pkg {
		otherErr.msg = msg
		return otherErr
	}

	return &Error{
		detail: e,
		msg:    fmt.Sprintf("%s | %s", msg, otherErr.Error()),
	}
}
