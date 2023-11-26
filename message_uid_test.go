package telegram

import (
	"testing"
)

func TestCallbackCurrent(t *testing.T) {
	if CallbackCurrent == nil {
		t.Error("CallbackCurrent is nil")
	}
}
