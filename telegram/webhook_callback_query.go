package telegram

import (
	"github.com/bots-go-framework/bots-fw/botinput"
)

type WebhookCallbackQuery interface {
	botinput.CallbackQuery
	GetInlineMessageID() string // Telegram only?
	GetChatInstanceID() string  // Telegram only?
}
