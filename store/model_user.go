package store

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/strongo/app/user"
	"strconv"
)

// TgBotUserData is Telegram user DB Data (without ID)
type TgBotUserData struct {
	botsfw.BotUserEntity
	//TgChatID int64
}

func (entity *TgBotUserData) GetAppUserStrID() string {
	return strconv.FormatInt(entity.BotUserEntity.GetAppUserIntID(), 10)
}

var _ botsfw.BotUser = (*TgBotUserData)(nil)
var _ user.AccountEntity = (*TgBotUserData)(nil)

// TgUser is Telegram user DB record (with ID)
type TgUser struct {
	ID int64
	TgBotUserData
}

// GetEmail returns empty string
func (TgUser) GetEmail() string {
	return ""
}

// Name returns full display name cmbined from (first+last, nick) name
func (entity *TgBotUserData) Name() string {
	if entity.FirstName == "" && entity.LastName == "" {
		return "@" + entity.UserName
	}
	name := entity.FirstName
	if name != "" {
		name += " " + entity.LastName
	} else {
		name = entity.LastName
	}
	if entity.UserName == "" {
		return name
	}
	return "@" + entity.UserName + " - " + name
}

// GetNames return user names
func (entity *TgBotUserData) GetNames() user.Names {
	return user.Names{
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		NickName:  entity.UserName,
	}
}

// IsEmailConfirmed returns false
func (entity *TgBotUserData) IsEmailConfirmed() bool {
	return false
}

//// Load is for datastore
//func (entity *TgBotUserData) Load(ps []datastore.Property) error {
//	return datastore.LoadStruct(entity, ps)
//}
//
//// Save is for datastore
//func (entity *TgBotUserData) Save() (properties []datastore.Property, err error) {
//	if properties, err = datastore.SaveStruct(entity); err != nil {
//		return properties, err
//	}
//
//	//if properties, err = gaedb.CleanProperties(properties, map[string]gaedb.IsOkToRemove{
//	//	"AccessGranted": gaedb.IsFalse,
//	//	"FirstName":     gaedb.IsEmptyString,
//	//	"LastName":      gaedb.IsEmptyString,
//	//	"UserName":      gaedb.IsEmptyString,
//	//}); err != nil {
//	//	return
//	//}
//
//	return
//}