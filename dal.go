package telegram

import (
	"context"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/dal-go/dalgo/dal"
)

type DataStore interface {
	TgChatInstanceDal
}

// TgChatInstanceDal is Data Access Layer for telegram chat instance Data
type TgChatInstanceDal interface {
	GetTelegramChatInstanceByID(c context.Context, tx dal.ReadTransaction, id string) (tgChatInstance botsfwtgmodels.TgChatInstanceData, err error)
	NewTelegramChatInstance(chatInstanceID string, chatID int64, preferredLanguage string) (tgChatInstance botsfwtgmodels.TgChatInstanceData)
	SaveTelegramChatInstance(c context.Context, tgChatInstance botsfwtgmodels.TgChatInstanceData) (err error)
}

var getDatabase func(context.Context) (dal.Database, error)
var tgChatInstanceDal TgChatInstanceDal

func Init(getDb func(context.Context) (dal.Database, error)) {
	if getDb == nil {
		panic("getDb == nil")
	}
	getDatabase = getDb
}
