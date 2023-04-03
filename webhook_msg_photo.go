package telegram

import (
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo/bots-framework/botsfw"
)

type tgWebhookPhotoMessage struct {
	tgWebhookMessage
	TgMessageType TgMessageType
}

var _ botsfw.WebhookPhotoMessage = (*tgWebhookPhotoMessage)(nil)

func (tgWebhookPhotoMessage) InputType() botsfw.WebhookInputType {
	return botsfw.WebhookInputPhoto
}

func newTgWebhookPhotoMessage(input tgWebhookInput, tgMessageType TgMessageType, tgMessage *tgbotapi.Message) tgWebhookPhotoMessage {
	return tgWebhookPhotoMessage{
		tgWebhookMessage: newTelegramWebhookMessage(input, tgMessage),
		TgMessageType:    tgMessageType,
	}
}
