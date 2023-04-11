package store

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
)

// Chat is a Telegram chat record with a strongly typed ID
// It can be used on its own or as a base for a record with fields specific to a bot
type Chat struct {
	record.WithID[string]
	Data TgChatData
}

//var _ dal.EntityHolder = (*Chat)(nil)

func NewChat(id string, data TgChatData) Chat {
	if data == nil {
		panic("data is nil")
	}
	key := dal.NewKeyWithID(TgChatCollection, id)
	//if data == nil {
	//	data = new(TgChatBase)
	//}
	return Chat{
		WithID: record.NewWithID(id, key, data),
		Data:   data,
	}
}

// TgChatData must be implemented by a data struct that stores chat data related to specific app/bot.
type TgChatData interface {
	BaseChatData() *TgChatBase
}

//func (entity *Data) Load(ps []datastore.Property) error {
//	return datastore.LoadStruct(entity, ps)
//}
//
//func (entity *Data) Save() (properties []datastore.Property, err error) {
//	if properties, err = datastore.SaveStruct(entity); err != nil {
//		return properties, err
//	}
//	if properties, err = entity.TgChatBase.CleanProperties(properties); err != nil {
//		return
//	}
//	//if properties, err = gaedb.CleanProperties(properties, map[string]gaedb.IsOkToRemove{
//	//	"TgChatInstanceID": gaedb.IsEmptyString,
//	//}); err != nil {
//	//	return
//	//}
//	return
//}
