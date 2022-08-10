package errortype

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var dummyErrorType11 = ErrorType{code: 1000, pkg: "dummypackage1"}
var dummyErrorType22 = ErrorType{code: 2000, pkg: "dummypackage2"}
var dummyErrorType13 = ErrorType{code: 1000, pkg: "dummypackage3"}
var dummyErrorType21 = ErrorType{code: 2000, pkg: "dummypackage1"}

func TestSameCodeDiffPkgIsFalse(t *testing.T) {
	dummyError11 := dummyErrorType11.New("one")
	dummyError13 := dummyErrorType13.New("one")

	assert.False(t, dummyErrorType13.Is(dummyError11))
	assert.False(t, dummyErrorType11.Is(dummyError13))
}

func TestDiffCodeSamePkgIsFalse(t *testing.T) {
	dummyError11 := dummyErrorType11.New("one")
	dummyError21 := dummyErrorType21.New("one")

	assert.False(t, dummyErrorType21.Is(dummyError11))
	assert.False(t, dummyErrorType11.Is(dummyError21))
}

func TestSameErrorTypeIsTrue(t *testing.T) {
	dummyError11 := dummyErrorType11.New("one")

	assert.True(t, dummyErrorType11.Is(dummyError11))
}

func TestDiffErrorTypeIsFalse(t *testing.T) {
	dummyError11 := dummyErrorType11.New("one")

	assert.False(t, dummyErrorType22.Is(dummyError11))
}

func TestInbuiltErrorIsFalse(t *testing.T) {
	dummyError := errors.New("one")

	assert.False(t, dummyErrorType11.Is(dummyError))
}

func TestSameErrorTypeWrapping(t *testing.T) {
	dummyError := dummyErrorType11.New("one")
	dummyErrorWrapped := dummyErrorType11.Wrap(dummyError)

	assert.Equal(t, dummyErrorWrapped.Error(), dummyError.Error())
}

func TestDiffErrorTypeWrapping(t *testing.T) {
	dummyError := dummyErrorType11.New("one")
	dummyErrorWrapped := dummyErrorType22.Wrap(dummyError)

	assert.Equal(t, dummyErrorWrapped.Error(),
		"error: code=2000, pkg=dummypackage2, msg=one | error: code=1000, pkg=dummypackage1, msg=one")
}

func TestSameWrappedErrorTypeIsFalse(t *testing.T) {
	dummyError := dummyErrorType11.New("one")
	dummyErrorWrapped := dummyErrorType22.Wrap(dummyError)

	assert.False(t, dummyErrorType11.Is(dummyErrorWrapped))
}
