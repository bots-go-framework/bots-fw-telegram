package telegram

import (
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo/bots-framework/botsfw"
)

type tgWebhookVoiceMessage struct {
	tgWebhookMessage
	TgMessageType TgMessageType
}

var _ botsfw.WebhookVoiceMessage = (*tgWebhookVoiceMessage)(nil)

func (tgWebhookVoiceMessage) InputType() botsfw.WebhookInputType {
	return botsfw.WebhookInputVoice
}

func newTgWebhookVoiceMessage(input tgWebhookInput, tgMessageType TgMessageType, tgMessage *tgbotapi.Message) tgWebhookVoiceMessage {
	return tgWebhookVoiceMessage{
		tgWebhookMessage: newTelegramWebhookMessage(input, tgMessage),
		TgMessageType:    tgMessageType,
	}
}
