package store

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
)

type Chat struct {
	record.WithID[string]
	//TgChatRecord
	*ChatEntity
}

//var _ dal.EntityHolder = (*Chat)(nil)

func NewChat(id string) Chat {
	key := dal.NewKeyWithID(TgChatCollection, id)
	dto := new(ChatEntity)
	return Chat{
		WithID: record.WithID[string]{
			ID:     id,
			Record: dal.NewRecordWithData(key, dto),
		},
		ChatEntity: dto,
	}
}

func (Chat) Kind() string {
	return TgChatCollection
}

//func (tgChat Chat) Entity() interface{} {
//	return tgChat.ChatEntity
//}

//func (Chat) NewEntity() interface{} {
//	return new(ChatEntity)
//}

//func (tgChat *Chat) SetEntity(Data interface{}) {
//	if Data == nil {
//		tgChat.ChatEntity = nil
//	} else {
//		tgChat.ChatEntity = Data.(*ChatEntity)
//	}
//}

type ChatEntity struct {
	UserGroupID string `datastore:",index,omitempty"` // Do index
	TgChatData
}

//func (entity *ChatEntity) Load(ps []datastore.Property) error {
//	return datastore.LoadStruct(entity, ps)
//}
//
//func (entity *ChatEntity) Save() (properties []datastore.Property, err error) {
//	if properties, err = datastore.SaveStruct(entity); err != nil {
//		return properties, err
//	}
//	if properties, err = entity.TgChatData.CleanProperties(properties); err != nil {
//		return
//	}
//	//if properties, err = gaedb.CleanProperties(properties, map[string]gaedb.IsOkToRemove{
//	//	"TgChatInstanceID": gaedb.IsEmptyString,
//	//}); err != nil {
//	//	return
//	//}
//	return
//}
