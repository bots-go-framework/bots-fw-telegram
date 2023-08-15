package telegram

import "testing"

func TestNewTelegramWebhookHandler(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("NewTelegramWebhookHandler() did not panic")
		}
	}()
	NewTelegramWebhookHandler(nil, nil, nil, nil)
}
