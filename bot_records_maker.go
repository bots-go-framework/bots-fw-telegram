package telegram

import (
	"github.com/bots-go-framework/bots-fw-models/botsfwmodels"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

var _ botsfw.BotRecordsMaker = botRecordsMaker{}

type botRecordsMaker struct{}

func (b botRecordsMaker) MakeBotUserDto() botsfwmodels.BotUser {
	return new(botsfwtgmodels.TgBotUserData)
}

func (b botRecordsMaker) MakeBotChatDto() botsfwmodels.BotChat {
	return new(botsfwtgmodels.TgChatBase)
}
