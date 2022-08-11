package commonerror

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCode0Nil(t *testing.T) {
	commonError := New(0, ErrMsgServer)

	assert.Nil(t, commonError)
}

func TestTypecastCommonErrorToInbuiltError(t *testing.T) {
	commonError := New(ErrCodeServer, ErrMsgServer)
	_, ok := commonError.(error)

	assert.True(t, ok)
}

func TestConvertSameCodeAndMsg(t *testing.T) {
	commonError := New(ErrCodeServer, ErrMsgServer)
	err, _ := commonError.(error)
	commonErrorConvert := Convert(err)

	assert.Equal(t, commonError.Code(), commonErrorConvert.Code())
	assert.Equal(t, commonError.Msg(), commonErrorConvert.Msg())
	assert.Equal(t, commonError.Error(), commonErrorConvert.Error())
}
