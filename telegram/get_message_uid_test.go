package telegram

import (
	"testing"
)

func TestGetEditMessageUID(t *testing.T) {
	if _, err := GetMessageUID(nil); err == nil {
		t.Errorf("GetMessageUID() expected to return error if WebhookContext is nil")
		return
	}
}
