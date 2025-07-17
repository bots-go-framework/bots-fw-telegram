package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
)

var _ botinput.PhotoMessage = (*tgWebhookPhotoMessage)(nil)

type tgWebhookPhotoMessage struct {
	tgInputMessage
	TgMessageType MessageType
}

func (tgWebhookPhotoMessage) InputType() botinput.Type {
	return botinput.TypePhoto
}

func newTgWebhookPhotoMessage(input tgInput, tgMessageType MessageType, tgMessage *tgbotapi.Message) tgWebhookPhotoMessage {
	return tgWebhookPhotoMessage{
		tgInputMessage: newTelegramWebhookMessage(input, tgMessage),
		TgMessageType:  tgMessageType,
	}
}
