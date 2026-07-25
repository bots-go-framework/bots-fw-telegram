package telegram

import (
	"testing"

	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
)

func TestTgWebhookEntryWebhookUpdateID(t *testing.T) {
	tests := []struct {
		name  string
		entry tgWebhookEntry
		id    string
		ok    bool
	}{
		{name: "nil update"},
		{name: "zero update ID", entry: tgWebhookEntry{update: &tgbotapi.Update{}}},
		{name: "telegram update ID", entry: tgWebhookEntry{update: &tgbotapi.Update{UpdateID: 42}}, id: "42", ok: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, ok := tt.entry.WebhookUpdateID()
			if id != tt.id || ok != tt.ok {
				t.Fatalf("WebhookUpdateID() = %q, %t; want %q, %t", id, ok, tt.id, tt.ok)
			}
		})
	}
}
