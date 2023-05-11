package telegram

import (
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

func NewBotRecordsFieldsSetter(
	setAppUserFields func(appUser botsfwmodels.AppUserData, sender botsfw.WebhookSender) error,
) botsfw.BotRecordsFieldsSetter {
	if setAppUserFields == nil {
		panic("setAppUserFields is nil")
	}
	return botRecordsFieldsSetter{
		setAppUserFields: setAppUserFields,
	}
}

type botRecordsFieldsSetter struct {
	setAppUserFields func(appUser botsfwmodels.AppUserData, sender botsfw.WebhookSender) error
}

func (b botRecordsFieldsSetter) Platform() string {
	return PlatformID
}

func (b botRecordsFieldsSetter) SetAppUserFields(appUser botsfwmodels.AppUserData, sender botsfw.WebhookSender) error {
	return b.setAppUserFields(appUser, sender)
}

func (b botRecordsFieldsSetter) SetBotUserFields(botUser botsfwmodels.BotUser, botID, botUserID, appUserID string, sender botsfw.WebhookSender) error {
	//tgSender := sender.(tgWebhookSender)
	tgBotUser := botUser.(botsfwtgmodels.TgBotUser)
	tgBotUserBaseData := tgBotUser.TgBotUserBaseData()
	botUserBaseData := tgBotUserBaseData.BaseData()
	//botUserBaseData.AppUserIntID = tgSender.tgUser.ID
	botUserBaseData.FirstName = sender.GetFirstName()
	botUserBaseData.LastName = sender.GetLastName()
	return nil
}

func (b botRecordsFieldsSetter) SetBotChatFields(botChat botsfwmodels.ChatData, botID, botUserID, appUserID string, chat botsfw.WebhookChat, isAccessGranted bool) error {
	return nil
}
