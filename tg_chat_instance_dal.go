package telegram

import (
	"context"
	botsfwtgmodels "github.com/bots-go-framework/bots-fw-telegram-models"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
)

type tgChatInstanceDalgo struct {
	db dal.Database
}

var _ TgChatInstanceDal = (*tgChatInstanceDalgo)(nil)

func (tgChatInstanceDal tgChatInstanceDalgo) GetTelegramChatInstanceByID(c context.Context, tx dal.ReadTransaction, id string) (tgChatInstance botsfwtgmodels.ChatInstance, err error) {
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

func (tgChatInstanceDal tgChatInstanceDalgo) SaveTelegramChatInstance(c context.Context, tgChatInstance botsfwtgmodels.ChatInstance) (err error) {
	err = tgChatInstanceDal.db.RunReadwriteTransaction(c, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		return tx.Set(ctx, tgChatInstance.Record)
	})
	return
}

func (tgChatInstanceDalgo) NewTelegramChatInstance(chatInstanceID string, chatID int64, preferredLanguage string) (tgChatInstance botsfwtgmodels.ChatInstance) {
	key := dal.NewKeyWithID(botsfwtgmodels.ChatInstanceKind, chatInstanceID)
	var chatInstance botsfwtgmodels.ChatInstanceEntity = &botsfwtgmodels.ChatInstanceEntityBase{
		TgChatID:          chatID,
		PreferredLanguage: preferredLanguage,
	}
	return botsfwtgmodels.ChatInstance{
		WithID: record.NewWithID(chatInstanceID, key, chatInstance),
		Data:   chatInstance,
	}
}

func init() {
	tgChatInstanceDal = tgChatInstanceDalgo{}
}
