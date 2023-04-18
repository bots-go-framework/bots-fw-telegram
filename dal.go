package telegram

import (
	"context"
	"github.com/bots-go-framework/bots-fw-telegram/store"
	"github.com/dal-go/dalgo/dal"
)

// TgChatInstanceDal is Data Access Layer for telegram chat instance Data
type TgChatInstanceDal interface {
	GetTelegramChatInstanceByID(c context.Context, tx dal.ReadTransaction, id string) (tgChatInstance store.ChatInstance, err error)
	NewTelegramChatInstance(chatInstanceID string, chatID int64, preferredLanguage string) (tgChatInstance store.ChatInstance)
	SaveTelegramChatInstance(c context.Context, tgChatInstance store.ChatInstance) (err error)
}

var getDatabase func(context.Context) (dal.Database, error)
var tgChatInstanceDal TgChatInstanceDal

func Init(getDb func(context.Context) (dal.Database, error)) {
	if getDb == nil {
		panic("getDb == nil")
	}
	getDatabase = getDb
}
