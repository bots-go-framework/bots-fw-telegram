package store

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/strongo/app/user"
	"strconv"
	"time"
)

// TgChatEntity is Telegram chat TgChatData interface
type TgChatEntity interface {
	SetTgChatInstanceID(v string)
	GetTgChatInstanceID() string
	GetPreferredLanguage() string
}

// TgChatRecord holds base properties of Telegram chat TgChatData
type TgChatRecord struct { // TODO: Do we need this struct at all?
	record.WithID[string]
	//Data *TgChatBase
}

// SetID sets ID
func (v *TgChatRecord) SetID(tgBotID string, tgChatID int64) {
	v.ID = tgBotID + ":" + strconv.FormatInt(tgChatID, 10) // TODO: Should we migrate to format "id@bot"?
	v.Key = dal.NewKeyWithID(TgChatCollection, v.ID)
}

var _ TgChatData = (*TgChatBase)(nil)

// TgChatBase holds base properties of Telegram chat TgChatData
type TgChatBase struct {
	botsfw.BotChatData
	UserGroupID           string  `datastore:",index,omitempty" firestore:",omitempty" dalgo:",index,omitempty"` // Do index
	TelegramUserID        int64   `datastore:",noindex,omitempty" firestore:",noindex,omitempty"`
	TelegramUserIDs       []int64 `datastore:",noindex" firestore:",noindex"` // For groups
	LastProcessedUpdateID int     `datastore:",noindex,omitempty" firestore:",noindex,omitempty"`
	TgChatInstanceID      string  // Do index // TODO: document here what is chat instance
}

func (data *TgChatBase) BaseChatData() *TgChatBase {
	return data
}

// SetTgChatInstanceID is what it is
func (data *TgChatBase) SetTgChatInstanceID(v string) {
	data.TgChatInstanceID = v
}

// GetTgChatInstanceID is what it is
func (data *TgChatBase) GetTgChatInstanceID() string {
	return data.TgChatInstanceID
}

// GetPreferredLanguage returns preferred language for the chat
func (data *TgChatBase) GetPreferredLanguage() string {
	return data.PreferredLanguage
}

var _ botsfw.BotChat = (*TgChatBase)(nil)

// NewTelegramChatEntity create new telegram chat TgChatData
func NewTelegramChatEntity() *TgChatBase {
	return &TgChatBase{
		BotChatData: botsfw.BotChatData{
			BotEntity: botsfw.BotEntity{OwnedByUserWithIntID: user.NewOwnedByUserWithIntID(0, time.Now())},
		},
	}
}

// SetAppUserIntID sets app user int ID
func (data *TgChatBase) SetAppUserIntID(id int64) {
	if data.IsGroup && id != 0 {
		panic("TgChatBase.IsGroup && id != 0")
	}
	data.AppUserIntID = id
}

// SetBotUserID sets bot user int ID
func (data *TgChatBase) SetBotUserID(id interface{}) {
	switch id := id.(type) {
	case string:
		var err error
		data.TelegramUserID, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			panic(err.Error())
		}
	case int:
		data.TelegramUserID = int64(id)
	case int64:
		data.TelegramUserID = id
	default:
		panic(fmt.Sprintf("Expected string, got: %T=%v", id, id))
	}
}

// Load loads Data from datastore
//func (data *TgChatBase) Load(ps []datastore.Property) error {
//	return datastore.LoadStruct(data, ps)
//}
//
//// Save saves Data to datastore
//func (data *TgChatBase) Save() (properties []datastore.Property, err error) {
//	if properties, err = datastore.SaveStruct(data); err != nil {
//		return
//	}
//	if properties, err = data.CleanProperties(properties); err != nil {
//		return
//	}
//	return
//}
//
//// CleanProperties cleans properties
//func (data *TgChatBase) CleanProperties(properties []datastore.Property) ([]datastore.Property, error) {
//	if data.IsGroup && data.AppUserIntID != 0 {
//		for _, userID := range data.AppUserIntIDs {
//			if userID == data.AppUserIntID {
//				goto found
//			}
//		}
//		data.AppUserIntIDs = append(data.AppUserIntIDs, data.AppUserIntID)
//		data.AppUserIntID = 0
//	found:
//	}
//
//	for i, userID := range data.AppUserIntIDs {
//		if userID == 0 {
//			panic(fmt.Sprintf("*TgChatBase.AppUserIntIDs[%d] == 0", i))
//		}
//	}
//
//	var err error
//	//if properties, err = gaedb.CleanProperties(properties, map[string]gaedb.IsOkToRemove{
//	//	"AppUserIntID":          gaedb.IsZeroInt,
//	//	"AccessGranted":         gaedb.IsFalse,
//	//	"AwaitingReplyTo":       gaedb.IsEmptyString,
//	//	"DtForbidden":           gaedb.IsZeroTime,
//	//	"DtForbiddenLast":       gaedb.IsZeroTime,
//	//	"GaClientID":            gaedb.IsEmptyByteArray,
//	//	"TelegramUserID":        gaedb.IsZeroInt,
//	//	"LastProcessedUpdateID": gaedb.IsZeroInt,
//	//	"PreferredLanguage":     gaedb.IsEmptyString,
//	//	"Title":                 gaedb.IsEmptyString, // TODO: Is it obsolete?
//	//	"Type":                  gaedb.IsEmptyString, // TODO: Is it obsolete?
//	//}); err != nil {
//	//	return properties, err
//	//}
//	return properties, err
//}
