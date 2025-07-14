package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
)

var _ botinput.PhotoMessage = (*tgWebhookPhotoMessage)(nil)

type tgWebhookPhotoMessage struct {
	tgWebhookMessage
	TgMessageType TgMessageType
}

func (tgWebhookPhotoMessage) InputType() botinput.Type {
	return botinput.TypePhoto
}

func newTgWebhookPhotoMessage(input tgWebhookInput, tgMessageType TgMessageType, tgMessage *tgbotapi.Message) tgWebhookPhotoMessage {
	return tgWebhookPhotoMessage{
		tgWebhookMessage: newTelegramWebhookMessage(input, tgMessage),
		TgMessageType:    tgMessageType,
	}
}
