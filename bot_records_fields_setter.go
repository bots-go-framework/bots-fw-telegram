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
	return tgBotRecordsFieldsSetter{
		setAppUserFields: setAppUserFields,
	}
}

type tgBotRecordsFieldsSetter struct {
	setAppUserFields func(appUser botsfwmodels.AppUserData, sender botsfw.WebhookSender) error
}

func (b tgBotRecordsFieldsSetter) Platform() string {
	return PlatformID
}

func (b tgBotRecordsFieldsSetter) SetAppUserFields(appUser botsfwmodels.AppUserData, sender botsfw.WebhookSender) error {
	return b.setAppUserFields(appUser, sender)
}

func (b tgBotRecordsFieldsSetter) SetBotUserFields(botUser botsfwmodels.BotUserData, sender botsfw.WebhookSender, botID, botUserID, appUserID string) error {
	//tgSender := sender.(tgWebhookSender)
	tgBotUser := botUser.(botsfwtgmodels.TgBotUser)
	tgBotUserBaseData := tgBotUser.TgBotUserBaseData()
	botUserBaseData := tgBotUserBaseData.BaseData()
	//botUserBaseData.AppUserIntID = tgSender.tgUser.ID
	botUserBaseData.FirstName = sender.GetFirstName()
	botUserBaseData.LastName = sender.GetLastName()
	return nil
}

func (b tgBotRecordsFieldsSetter) SetBotChatFields(botChat botsfwmodels.BotChatData, chat botsfw.WebhookChat, botID, botUserID, appUserID string, isAccessGranted bool) error {
	_ = botID
	_ = chat
	tgBotChatData := botChat.(botsfwtgmodels.TgChatData)
	baseTgChatData := tgBotChatData.BaseTgChatData()
	//baseTgChatData.BotID = botID
	baseTgChatData.SetBotUserID(botUserID)
	baseTgChatData.AppUserID = appUserID
	baseTgChatData.SetAccessGranted(isAccessGranted) // TODO(help-wanted): can be set outside, no need to pass as parameter
	return nil
}
