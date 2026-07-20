package telegram

import (
	"context"

	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
)

// ChatInstanceStore persists the Telegram callback-query to chat mapping. A
// callback that has no Message uses this mapping to recover the chat and locale.
// It is intentionally independent from the framework state store because it is
// Telegram-specific metadata.
type ChatInstanceStore interface {
	Get(ctx context.Context, botCode, chatInstanceID string) (data botsfwtgmodels.TgChatInstanceData, found bool, err error)
	Save(ctx context.Context, botCode, chatInstanceID string, data botsfwtgmodels.TgChatInstanceData) error
}

func NewTelegramChatInstance(chatInstanceID string, chatID int64, preferredLanguage string) botsfwtgmodels.TgChatInstanceData {
	_ = chatInstanceID
	return &botsfwtgmodels.TgChatInstanceBaseData{
		TgChatID:          chatID,
		PreferredLanguage: preferredLanguage,
	}
}
