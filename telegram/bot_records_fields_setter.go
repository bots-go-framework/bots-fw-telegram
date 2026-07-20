package telegram

import (
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

func NewBotRecordsFieldsSetter() botsfw.BotRecordsFieldsSetter { return tgBotRecordsFieldsSetter{} }

type tgBotRecordsFieldsSetter struct{}

func (b tgBotRecordsFieldsSetter) Platform() string {
	return string(PlatformID)
}

func (b tgBotRecordsFieldsSetter) SetBotUserFields(botUser botsfwmodels.PlatformUserData, sender botinput.Sender, botID, botUserID, appUserID string) error {
	//tgSender := sender.(tgWebhookUser)
	tgBotUser := botUser.(botsfwtgmodels.TgPlatformUser)
	tgBotUserBaseData := tgBotUser.TgPlatformUserBaseDbo()
	botUserBaseData := tgBotUserBaseData.BaseData()
	//botUserBaseData.AppUserIntID = tgSender.tgUser.ID
	botUserBaseData.FirstName = sender.GetFirstName()
	botUserBaseData.LastName = sender.GetLastName()
	return nil
}

func (b tgBotRecordsFieldsSetter) SetBotChatFields(botChat botsfwmodels.BotChatData, chat botinput.Chat, botID, botUserID, appUserID string, isAccessGranted bool) error {
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
