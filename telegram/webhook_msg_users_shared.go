package telegram

import (
	"strconv"

	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
)

var (
	_ botinput.WebhookInput              = (*tgWebhookUsersSharedMessage)(nil)
	_ botinput.WebhookMessage            = (*tgWebhookUsersSharedMessage)(nil)
	_ botinput.WebhookSharedUsersMessage = (*tgWebhookUsersSharedMessage)(nil)
)

type tgWebhookUsersSharedMessage struct {
	tgWebhookMessage
	TgMessageType TgMessageType
}

func (m tgWebhookUsersSharedMessage) GetRequestID() int {
	return m.message.UsersShared.RequestID
}

func (m tgWebhookUsersSharedMessage) GetSharedUsers() (sharedUsers []botinput.SharedUserMessageItem) {
	if m.message == nil {
		panic("m.message is nil")
	}
	if m.message.UsersShared == nil {
		panic("m.message.UsersShared is nil")
	}
	sharedUsers = make([]botinput.SharedUserMessageItem, 0, len(m.message.UsersShared.Users))
	for _, sharedUser := range m.message.UsersShared.Users {
		sharedUsers = append(sharedUsers, tgSharedUser{SharedUser: sharedUser})
	}
	return
}

func (tgWebhookUsersSharedMessage) InputType() botinput.WebhookInputType {
	return botinput.WebhookInputSharedUsers
}

func newTgWebhookUsersSharedMessage(input tgWebhookInput, tgMessageType TgMessageType, tgMessage *tgbotapi.Message) tgWebhookUsersSharedMessage {
	return tgWebhookUsersSharedMessage{
		tgWebhookMessage: newTelegramWebhookMessage(input, tgMessage),
		TgMessageType:    tgMessageType,
	}
}

var _ botinput.SharedUserMessageItem = (*tgSharedUser)(nil)

type tgSharedUser struct {
	tgbotapi.SharedUser
}

func (v tgSharedUser) GetBotUserID() string {
	return strconv.Itoa(v.UserID)
}

func (v tgSharedUser) GetUsername() string {
	return v.Username
}

func (v tgSharedUser) GetFirstName() string {
	return v.FirstName
}

func (v tgSharedUser) GetLastName() string {
	return v.LastName
}

func (v tgSharedUser) GetPhotos() (photos []botinput.PhotoMessageItem) {
	photos = make([]botinput.PhotoMessageItem, len(v.Photo))
	for i, photo := range v.Photo {
		photos[i] = photoSize{PhotoSize: photo}
	}
	return
}
