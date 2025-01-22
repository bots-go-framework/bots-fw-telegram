package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
)

var _ botinput.WebhookUser = (*tgWebhookUser)(nil)

type tgWebhookUser struct {
	tgUser *tgbotapi.User
}

func (s tgWebhookUser) GetCountry() string {
	return ""
}

func (tgWebhookUser) IsBotUser() bool { // TODO: Can we get rid of it here?
	return false
}

func (s tgWebhookUser) GetID() interface{} {
	return s.tgUser.ID
}

func (s tgWebhookUser) GetFirstName() string {
	return s.tgUser.FirstName
}

func (s tgWebhookUser) GetLastName() string {
	return s.tgUser.LastName
}

func (s tgWebhookUser) GetUserName() string {
	return s.tgUser.UserName
}

func (tgWebhookUser) Platform() string {
	return PlatformID
}

func (tgWebhookUser) GetAvatar() string {
	return ""
}

func (s tgWebhookUser) GetLanguage() string {
	return s.tgUser.LanguageCode
}
