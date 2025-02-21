package telegram

import (
	"context"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/bots-go-framework/bots-fw/botsdal"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
)

type TgChatInstance struct {
	record.WithID[string]
	Data botsfwtgmodels.TgChatInstanceData
}

func getTelegramChatInstanceByID(c context.Context, tx dal.ReadSession, botCode, chatInstanceID string) (tgChatInstanceData botsfwtgmodels.TgChatInstanceData, err error) {
	tgChatInstanceData = NewChatInstanceData(botCode)

	key := NewTgChatInstanceKey(botCode, chatInstanceID)
	tgChatInstance := TgChatInstance{
		WithID: record.NewWithID(chatInstanceID, key, tgChatInstanceData),
		Data:   tgChatInstanceData,
	}

	if err = tx.Get(c, tgChatInstance.Record); dal.IsNotFound(err) {
		return
	}
	return
}

func saveTelegramChatInstance(ctx context.Context, tx dal.ReadwriteTransaction, botCode, chatInstanceID string, tgChatInstanceData botsfwtgmodels.TgChatInstanceData) (err error) {
	key := NewTgChatInstanceKey(botCode, chatInstanceID)
	chatInstance := record.NewWithID(chatInstanceID, key, tgChatInstanceData)
	return tx.Set(ctx, chatInstance.Record)
}

func NewTelegramChatInstance(chatInstanceID string, chatID int64, preferredLanguage string) (tgChatInstanceData botsfwtgmodels.TgChatInstanceData) {
	_ = chatInstanceID
	tgChatInstanceData = &botsfwtgmodels.TgChatInstanceBaseData{
		TgChatID:          chatID,
		PreferredLanguage: preferredLanguage,
	}
	return tgChatInstanceData
}

const ChatInstancesCollection = "chatInstances"

func NewTgChatInstanceKey(botCode, chatInstanceID string) *dal.Key {
	platformKey := botsdal.NewPlatformKey(PlatformID)
	return dal.NewKeyWithParentAndID(platformKey, ChatInstancesCollection, chatInstanceID)
}
