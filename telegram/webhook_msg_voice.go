package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
)

type tgWebhookVoiceMessage struct {
	tgWebhookMessage
	TgMessageType TgMessageType
}

var _ botinput.VoiceMessage = (*tgWebhookVoiceMessage)(nil)

func (tgWebhookVoiceMessage) InputType() botinput.Type {
	return botinput.TypeVoice
}

func newTgWebhookVoiceMessage(input tgWebhookInput, tgMessageType TgMessageType, tgMessage *tgbotapi.Message) tgWebhookVoiceMessage {
	return tgWebhookVoiceMessage{
		tgWebhookMessage: newTelegramWebhookMessage(input, tgMessage),
		TgMessageType:    tgMessageType,
	}
}
