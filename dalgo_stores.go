package store

import (
	"github.com/bots-go-framework/bots-fw-dalgo/dalgo4botsfw"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

func NewDalgoStores() (botsfw.BotChatStore, botsfw.BotUserStore) {
	return NewBotChatStore(), nil
}

func NewBotChatStore() botsfw.BotChatStore {
	return dalgo4botsfw.NewBotChatStore(TgChatCollection, nil, func() botsfw.BotChat {
		return new(TgChatData)
	})
}
