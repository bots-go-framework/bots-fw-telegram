package store

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
)

// TgChat is a Telegram chat record with a strongly typed ID
// It can be used on its own or as a base for a record with fields specific to a bot
type TgChat struct {
	record.WithID[string]
	Data TgChatData
}

func NewTgChat(id string, data TgChatData) TgChat {
	if data == nil {
		panic("data is nil")
	}
	key := dal.NewKeyWithID(TgChatCollection, id)
	//if data == nil {
	//	data = new(TgChatBase)
	//}
	return TgChat{
		WithID: record.NewWithID(id, key, data),
		Data:   data,
	}
}

// TgChatData must be implemented by a data struct that stores chat data related to specific app/bot.
type TgChatData interface {
	BaseChatData() *TgChatBase
}
