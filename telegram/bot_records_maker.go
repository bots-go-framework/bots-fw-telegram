package telegram

import (
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
)

func BaseTgUserDtoMaker(botID string) (botsfwmodels.PlatformUserData, error) {
	tgPlatformUserData := botsfwtgmodels.TgPlatformUserBaseDbo{}
	return &tgPlatformUserData, nil
}

func BaseTgChatDtoMaker(botID string) (botChat botsfwmodels.BotChatData, err error) {
	tgChat := botsfwtgmodels.TgChatBaseData{}
	return &tgChat, nil
}
