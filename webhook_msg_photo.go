package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
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
