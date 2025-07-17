package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
)

type tgWebhookVoiceMessage struct {
	tgInputMessage
	TgMessageType MessageType
}

var _ botinput.VoiceMessage = (*tgWebhookVoiceMessage)(nil)

func (tgWebhookVoiceMessage) InputType() botinput.Type {
	return botinput.TypeVoice
}

func newTgWebhookVoiceMessage(input tgInput, tgMessageType MessageType, tgMessage *tgbotapi.Message) tgWebhookVoiceMessage {
	return tgWebhookVoiceMessage{
		tgInputMessage: newTelegramWebhookMessage(input, tgMessage),
		TgMessageType:  tgMessageType,
	}
}
