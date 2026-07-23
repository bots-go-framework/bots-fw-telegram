package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-go-core/botkb"
)

// GetTelegramKeyboard converts a platform-neutral botkb.Keyboard into its
// Telegram representation. Exported so code that builds Telegram messages
// OUTSIDE a webhook responder — e.g. a delayed/background task that edits a
// message via tgbotapi directly — can reuse the same conversion the responder
// uses. A nil keyboard yields a nil Telegram keyboard.
func GetTelegramKeyboard(keyboard botkb.Keyboard) tgbotapi.Keyboard {
	return getTelegramKeyboard(keyboard)
}

// GetInlineKeyboard converts a botkb.MessageKeyboard into a Telegram inline
// keyboard markup, for use as an editMessageText/ sendMessage ReplyMarkup from
// off-webhook code. See GetTelegramKeyboard for the general case.
func GetInlineKeyboard(kb *botkb.MessageKeyboard) *tgbotapi.InlineKeyboardMarkup {
	return getInlineKeyboard(kb)
}
