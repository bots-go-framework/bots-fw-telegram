package telegram

import (
	"github.com/strongo/app"
	"github.com/strongo/bots-framework/botsfw"
)

// NewTelegramBot creates definition of new telegram bot
func NewTelegramBot(mode strongo.Environment, profile, code, token, paymentTestToken, paymentToken, gaToken string, locale strongo.Locale) botsfw.BotSettings {
	settings := botsfw.NewBotSettings(mode, profile, code, "", token, gaToken, locale)
	settings.PaymentTestToken = paymentTestToken
	settings.PaymentToken = paymentToken
	return settings
}
