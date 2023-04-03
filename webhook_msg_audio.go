package telegram

import (
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo/bots-framework/botsfw"
)

type tgWebhookAudioMessage struct {
	tgWebhookMessage
	TgMessageType TgMessageType
}

var _ botsfw.WebhookAudioMessage = (*tgWebhookAudioMessage)(nil)

func (tgWebhookAudioMessage) InputType() botsfw.WebhookInputType {
	return botsfw.WebhookInputAudio
}

func newTgWebhookAudioMessage(input tgWebhookInput, tgMessageType TgMessageType, tgMessage *tgbotapi.Message) tgWebhookAudioMessage {
	return tgWebhookAudioMessage{
		tgWebhookMessage: newTelegramWebhookMessage(input, tgMessage),
		TgMessageType:    tgMessageType,
	}
}
