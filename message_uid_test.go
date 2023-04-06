package telegram

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCallbackCurrent(t *testing.T) {
	assert.NotNil(t, CallbackCurrent)
}
