package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
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
