package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
)

type tgWebhookPhotoMessage struct {
	tgWebhookMessage
	TgMessageType TgMessageType
}

var _ botinput.WebhookPhotoMessage = (*tgWebhookPhotoMessage)(nil)

func (tgWebhookPhotoMessage) InputType() botinput.WebhookInputType {
	return botinput.WebhookInputPhoto
}

func newTgWebhookPhotoMessage(input tgWebhookInput, tgMessageType TgMessageType, tgMessage *tgbotapi.Message) tgWebhookPhotoMessage {
	return tgWebhookPhotoMessage{
		tgWebhookMessage: newTelegramWebhookMessage(input, tgMessage),
		TgMessageType:    tgMessageType,
	}
}
