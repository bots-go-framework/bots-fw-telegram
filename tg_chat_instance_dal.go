package telegram

import (
	"context"
	"github.com/bots-go-framework/bots-fw-telegram/store"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
)

type tgChatInstanceDalgo struct {
	db dal.Database
}

var _ TgChatInstanceDal = (*tgChatInstanceDalgo)(nil)

func (tgChatInstanceDal tgChatInstanceDalgo) GetTelegramChatInstanceByID(c context.Context, tx dal.ReadTransaction, id string) (tgChatInstance store.ChatInstance, err error) {
	tgChatInstance = tgChatInstanceDal.NewTelegramChatInstance(id, 0, "")

	var session dal.ReadSession
	if tx == nil {
		session = tgChatInstanceDal.db
	} else {
		session = tx
	}
	if err = session.Get(c, tgChatInstance.Record); dal.IsNotFound(err) {
		tgChatInstance.SetEntity(nil)
		return
	}
	return
}

func (tgChatInstanceDal tgChatInstanceDalgo) SaveTelegramChatInstance(c context.Context, tgChatInstance store.ChatInstance) (err error) {
	err = tgChatInstanceDal.db.RunReadwriteTransaction(c, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		return tx.Set(ctx, tgChatInstance.Record)
	})
	return
}

func (tgChatInstanceDalgo) NewTelegramChatInstance(chatInstanceID string, chatID int64, preferredLanguage string) (tgChatInstance store.ChatInstance) {
	key := dal.NewKeyWithID(store.ChatInstanceKind, chatInstanceID)
	var chatInstance store.ChatInstanceEntity = &store.ChatInstanceEntityBase{
		TgChatID:          chatID,
		PreferredLanguage: preferredLanguage,
	}
	return store.ChatInstance{
		WithID: record.NewWithID(chatInstanceID, key, chatInstance),
		Data:   chatInstance,
	}
}

func init() {
	DAL.TgChatInstance = tgChatInstanceDalgo{}
}
