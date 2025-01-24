package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"strconv"
)

var _ botinput.WebhookSharedUserMessage = (*tgWebhookUsersSharedMessage)(nil)

type tgWebhookUsersSharedMessage struct {
	tgWebhookMessage
	TgMessageType TgMessageType
}

func (m tgWebhookUsersSharedMessage) GetSharedUsers() []botinput.SharedUserMessageItem {
	//TODO implement me
	panic("implement me")
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
	return v.Username
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
