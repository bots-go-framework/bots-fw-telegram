package telegram

import (
	"github.com/bots-go-framework/bots-fw-telegram/models"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

var _ botsfw.BotRecordsMaker = botRecordsMaker{}

type botRecordsMaker struct{}

func (b botRecordsMaker) MakeBotUserDto() botsfw.BotUser {
	return new(models.TgBotUserData)
}

func (b botRecordsMaker) MakeBotChatDto() botsfw.BotChat {
	return new(models.TgChatBase)
}
