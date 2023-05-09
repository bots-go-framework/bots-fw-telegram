package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

var _ botsfw.WebhookSender = (*tgWebhookSender)(nil)

type tgWebhookSender struct {
	tgUser *tgbotapi.User
}

func (tgWebhookSender) IsBotUser() bool { // TODO: Can we get rid of it here?
	return false
}

func (s tgWebhookSender) GetID() interface{} {
	return s.tgUser.ID
}

func (s tgWebhookSender) GetFirstName() string {
	return s.tgUser.FirstName
}

func (s tgWebhookSender) GetLastName() string {
	return s.tgUser.LastName
}

func (s tgWebhookSender) GetUserName() string {
	return s.tgUser.UserName
}

func (tgWebhookSender) Platform() string {
	return PlatformID
}

func (tgWebhookSender) GetAvatar() string {
	return ""
}

func (s tgWebhookSender) GetLanguage() string {
	return s.tgUser.LanguageCode
}
