package telegram

import (
	"context"
	"github.com/bots-go-framework/bots-fw-telegram/models"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/bots-go-framework/dalgo4botsfw"
	"github.com/strongo/app/user"
	"time"
)

func NewDalgoStores(db dalgo4botsfw.DbProvider) (botsfw.BotChatStore, botsfw.BotUserStore) {
	return newDalgoBotChatStore(db), newDalgoBotUserStore(db)
}

func newDalgoBotChatStore(db dalgo4botsfw.DbProvider) botsfw.BotChatStore {
	newChatData := func() botsfw.BotChat {
		return new(models.TgChatBase)
	}
	return dalgo4botsfw.NewBotChatStore(models.TgChatCollection, db, newChatData)
}

func newDalgoBotUserStore(db dalgo4botsfw.DbProvider) botsfw.BotUserStore {

	newUserData := func() botsfw.BotUser {
		return new(models.TgBotUserData)
	}

	createBotUser := func(c context.Context, botID string, apiUser botsfw.WebhookActor) (botsfw.BotUser, error) {
		if apiUser == nil {
			return &models.TgBotUserData{}, nil
		}
		return &models.TgBotUserData{
			BotUserData: botsfw.BotUserData{
				BotEntity: botsfw.BotEntity{
					OwnedByUserWithID: user.NewOwnedByUserWithIntID(0, time.Now()),
				},
				FirstName: apiUser.GetFirstName(),
				LastName:  apiUser.GetLastName(),
				UserName:  apiUser.GetUserName(),
			},
		}, nil
	}

	return dalgo4botsfw.NewBotUserStore(models.BotUserCollection, db, newUserData, createBotUser)
}
