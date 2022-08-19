package errortype

import (
	"errors"
	"fmt"
)

type IError interface {
	error
	Msg() string
}

// Error - contains ErrorType and error msg providing stack trace error.
type Error struct {
	msg    string
	detail ErrorType
}

// ErrorType - contains details to differentiate between Errors.
type ErrorType struct {
	Pkg  string
	Code int32
}

// Error - returns formatted string containing error details and error msg.
func (e *Error) Error() string {
	return fmt.Sprintf(ErrFormat, e.detail.Code, e.detail.Pkg, e.msg)
}

// Msg - returns error msg.
func (e *Error) Msg() string {
	return e.msg
}

// New - constructor for custom Error
func (e ErrorType) New(msg string) IError {
	return &Error{
		detail: e,
		msg:    msg,
	}
}

// Is - checks if err is of same code and package as ErrorType
func (e ErrorType) Is(err error) bool {
	otherErr := &Error{}

	if !errors.As(err, &otherErr) {
		return false
	}

	return otherErr.detail.Code == e.Code && otherErr.detail.Pkg == e.Pkg
}

// Wrap - if err is of same ErrorType, then no wrapping is done.
func (e ErrorType) Wrap(err error) IError {
	otherErr := &Error{}

	if !errors.As(err, &otherErr) {
		return e.New(err.Error())
	}

	if e.Code == otherErr.detail.Code &&
		e.Pkg == otherErr.detail.Pkg {
		return otherErr
	}

	return &Error{
		detail: e,
		msg:    fmt.Sprintf("%s | %s", otherErr.msg, otherErr.Error()),
	}
}

// WrapWithMsg - if err is of same ErrorType, err's msg is changed
// to the provided msg.
func (e ErrorType) WrapWithMsg(err error, msg string) IError {
	otherErr := &Error{}

	if !errors.As(err, &otherErr) {
		return e.New(err.Error())
	}

	if e.Code == otherErr.detail.Code &&
		e.Pkg == otherErr.detail.Pkg {
		otherErr.msg = msg
		return otherErr
	}

	return &Error{
		detail: e,
		msg:    fmt.Sprintf("%s | %s", msg, otherErr.Error()),
	}
}
