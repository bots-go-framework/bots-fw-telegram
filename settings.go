package telegram

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/strongo/app"
)

// NewTelegramBot creates definition of new telegram bot
func NewTelegramBot(mode strongo.Environment, profile, code, token, paymentTestToken, paymentToken, gaToken string, locale strongo.Locale) botsfw.BotSettings {
	settings := botsfw.NewBotSettings(botsfw.PlatformTelegram, mode, profile, code, "", token, gaToken, locale)
	settings.PaymentTestToken = paymentTestToken
	settings.PaymentToken = paymentToken
	return settings
}
