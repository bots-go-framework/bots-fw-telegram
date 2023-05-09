package telegram

import (
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
)

func BaseTgUserDtoMaker(botID string) (botsfwmodels.BotUser, error) {
	tgBotUserData := botsfwtgmodels.TgBotUserBaseData{}
	return &tgBotUserData, nil
}

func BaseTgChatDtoMaker(botID string) (botChat botsfwmodels.BotChat, err error) {
	tgChat := botsfwtgmodels.TgChatBase{}
	return &tgChat, nil
}
