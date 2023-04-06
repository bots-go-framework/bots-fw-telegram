package telegram

import (
	"context"
	"github.com/bots-go-framework/bots-fw-telegram/store"
	"github.com/dal-go/dalgo/dal"
)

// TgChatInstanceDal is DAL for telegram chat instance Data
type TgChatInstanceDal interface {
	GetTelegramChatInstanceByID(c context.Context, tx dal.ReadTransaction, id string) (tgChatInstance store.ChatInstance, err error)
	NewTelegramChatInstance(chatInstanceID string, chatID int64, preferredLanguage string) (tgChatInstance store.ChatInstance)
	SaveTelegramChatInstance(c context.Context, tgChatInstance store.ChatInstance) (err error)
}

type dal1 struct {
	DB             dal.Database
	TgChatInstance TgChatInstanceDal
}

// DAL is data access layer
var DAL dal1
