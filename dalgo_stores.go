package telegram

import (
	"context"
	"github.com/bots-go-framework/bots-fw-dalgo/dalgo4botsfw"
	"github.com/bots-go-framework/bots-fw-telegram/store"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/strongo/app/user"
	"time"
)

func NewDalgoStores(db dalgo4botsfw.DbProvider) (botsfw.BotChatStore, botsfw.BotUserStore) {
	return newDalgoBotChatStore(db), newDalgoBotUserStore(db)
}

func newDalgoBotChatStore(db dalgo4botsfw.DbProvider) botsfw.BotChatStore {
	newChatData := func() botsfw.BotChat {
		return new(store.TgChatData)
	}
	return dalgo4botsfw.NewBotChatStore(store.TgChatCollection, db, newChatData)
}

func newDalgoBotUserStore(db dalgo4botsfw.DbProvider) botsfw.BotUserStore {

	newUserData := func() botsfw.BotUser {
		return new(store.TgBotUserData)
	}

	createBotUser := func(c context.Context, botID string, apiUser botsfw.WebhookActor) (botsfw.BotUser, error) {
		if apiUser == nil {
			return &store.TgBotUserData{}, nil
		}
		botEntity := botsfw.BotEntity{
			OwnedByUserWithIntID: user.NewOwnedByUserWithIntID(0, time.Now()),
		}
		botUserEntity := botsfw.BotUserEntity{
			BotEntity: botEntity,
			FirstName: apiUser.GetFirstName(),
			LastName:  apiUser.GetLastName(),
			UserName:  apiUser.GetUserName(),
		}
		return &store.TgBotUserData{BotUserEntity: botUserEntity}, nil

		return newUserData(), nil
	}

	return dalgo4botsfw.NewBotUserStore(store.BotUserCollection, db, newUserData, createBotUser)
}
