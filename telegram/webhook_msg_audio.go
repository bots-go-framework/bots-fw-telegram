package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
)

type tgWebhookAudioMessage struct {
	tgWebhookMessage
	TgMessageType TgMessageType
}

var _ botinput.AudioMessage = (*tgWebhookAudioMessage)(nil)

func (tgWebhookAudioMessage) InputType() botinput.Type {
	return botinput.TypeAudio
}

func newTgWebhookAudioMessage(input tgWebhookInput, tgMessageType TgMessageType, tgMessage *tgbotapi.Message) tgWebhookAudioMessage {
	return tgWebhookAudioMessage{
		tgWebhookMessage: newTelegramWebhookMessage(input, tgMessage),
		TgMessageType:    tgMessageType,
	}
}
