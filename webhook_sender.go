package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

type tgSender struct {
	tgUser *tgbotapi.User
}

func (tgSender) IsBotUser() bool { // TODO: Can we get rid of it here?
	return false
}

var _ botsfw.WebhookSender = (*tgSender)(nil)

func (s tgSender) GetID() interface{} {
	return s.tgUser.ID
}

func (s tgSender) GetFirstName() string {
	return s.tgUser.FirstName
}

func (s tgSender) GetLastName() string {
	return s.tgUser.LastName
}

func (s tgSender) GetUserName() string {
	return s.tgUser.UserName
}

func (tgSender) Platform() string {
	return PlatformID
}

func (tgSender) GetAvatar() string {
	return ""
}

func (s tgSender) GetLanguage() string {
	return s.tgUser.LanguageCode
}
