package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"strconv"
)

// TgWebhookChat is wrapper for Telegram chat
type TgWebhookChat struct {
	chat *tgbotapi.Chat
}

var _ botsfw.WebhookChat = (*TgWebhookChat)(nil)

// GetID returns telegram chat ID
func (wh TgWebhookChat) GetID() string {
	return strconv.FormatInt(wh.chat.ID, 10)
}

// GetType returns telegram chat type
func (wh TgWebhookChat) GetType() string {
	return wh.chat.Type
}

// IsGroupChat indicates type of chat (group or private)
func (wh TgWebhookChat) IsGroupChat() bool {
	return !wh.chat.IsPrivate()
}
