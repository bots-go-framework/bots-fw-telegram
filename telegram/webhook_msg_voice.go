package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
)

type tgWebhookVoiceMessage struct {
	tgWebhookMessage
	TgMessageType TgMessageType
}

var _ botinput.WebhookVoiceMessage = (*tgWebhookVoiceMessage)(nil)

func (tgWebhookVoiceMessage) InputType() botinput.WebhookInputType {
	return botinput.WebhookInputVoice
}

func newTgWebhookVoiceMessage(input tgWebhookInput, tgMessageType TgMessageType, tgMessage *tgbotapi.Message) tgWebhookVoiceMessage {
	return tgWebhookVoiceMessage{
		tgWebhookMessage: newTelegramWebhookMessage(input, tgMessage),
		TgMessageType:    tgMessageType,
	}
}
