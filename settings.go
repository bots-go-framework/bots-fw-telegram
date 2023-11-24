package telegram

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/strongo/i18n"
	"github.com/strongo/strongoapp"
)

// NewTelegramBot creates definition of new telegram bot
func NewTelegramBot(
	mode strongoapp.Environment,
	profile botsfw.BotProfile,
	code, token, paymentTestToken, paymentToken, gaToken string,
	locale i18n.Locale,
	getDatabase botsfw.DbGetter,
	getAppUser botsfw.AppUserGetter,
) botsfw.BotSettings {
	settings := botsfw.NewBotSettings(botsfw.PlatformTelegram, mode, profile, code, "", token, gaToken, locale, getDatabase, getAppUser)
	settings.PaymentTestToken = paymentTestToken
	settings.PaymentToken = paymentToken
	return settings
}
