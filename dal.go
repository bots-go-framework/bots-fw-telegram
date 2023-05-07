package telegram

import (
	"context"
	botsfwtgmodels "github.com/bots-go-framework/bots-fw-telegram-models"
	"github.com/dal-go/dalgo/dal"
)

// TgChatInstanceDal is Data Access Layer for telegram chat instance Data
type TgChatInstanceDal interface {
	GetTelegramChatInstanceByID(c context.Context, tx dal.ReadTransaction, id string) (tgChatInstance botsfwtgmodels.ChatInstance, err error)
	NewTelegramChatInstance(chatInstanceID string, chatID int64, preferredLanguage string) (tgChatInstance botsfwtgmodels.ChatInstance)
	SaveTelegramChatInstance(c context.Context, tgChatInstance botsfwtgmodels.ChatInstance) (err error)
}

var getDatabase func(context.Context) (dal.Database, error)
var tgChatInstanceDal TgChatInstanceDal

func Init(getDb func(context.Context) (dal.Database, error)) {
	if getDb == nil {
		panic("getDb == nil")
	}
	getDatabase = getDb
}
