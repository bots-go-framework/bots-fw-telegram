package gaehost

import (
	"context"
	"github.com/strongo/app/user"
	"github.com/strongo/bots-framework/botsfw"
	"github.com/strongo/bots-fw-dalgo/dalgo4botsfw"
	telegram "github.com/strongo/bots-fw-telegram"
	gae "github.com/strongo/bots-host-gae"
	"time"
)

// gaeTelegramUserStore is DAL to telegram user entity
type gaeTelegramUserStore struct {
	gae.GaeBotUserStore
}

var _ botsfw.BotUserStore = (*gaeTelegramUserStore)(nil) // Check for interface implementation at compile time

// newGaeTelegramUserStore create DAL to Telegram user entity
func newGaeTelegramUserStore(gaeAppUserStore gae.GaeAppUserStore) botsfw.BotUserStore {
	newBotUserEntity := func(_ context.Context, botID string, apiUser botsfw.WebhookActor) (botsfw.BotUser, error) {
		if apiUser == nil {
			return &telegram.TgUserEntity{}, nil
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
		return &telegram.TgUserEntity{BotUserEntity: botUserEntity}, nil
	}
	//botUserKey := func(c context.Context, botUserId interface{}) *datastore.Key {
	//	if intID, ok := botUserId.(int); ok {
	//		if intID == 0 {
	//			panic("botUserKey(): intID == 0")
	//		}
	//		return datastore.NewKey(c, telegram.TgUserKind, "", (int64)(intID), nil)
	//	}
	//	panic(fmt.Sprintf("Expected botUserId as int, got: %T", botUserId))
	//}
	//validateBotUserEntityType := func(entity botsfw.BotUser) {
	//	if _, ok := entity.(*telegram.TgUserEntity); !ok {
	//		panic(fmt.Sprintf("Expected *telegram.TgUser but received %T", entity))
	//	}
	//}

	newBotUser := func() botsfw.BotUser {
		return new(telegram.TgUserEntity)
	}
	return dalgo4botsfw.NewBotUserStore(telegram.TgUserKind, nil, newBotUser, newBotUserEntity)
	//return gaeTelegramUserStore{
	//	GaeBotUserStore: gae.GaeBotUserStore{
	//		GaeBaseStore:              gae.NewGaeBaseStore(telegram.TgUserKind),
	//		gaeAppUserStore:           gaeAppUserStore,
	//		validateBotUserEntityType: validateBotUserEntityType,
	//		botUserKey:                botUserKey,
	//		newN
	//	},
	//}
}
