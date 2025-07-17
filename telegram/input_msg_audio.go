package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
)

type tgWebhookAudioMessage struct {
	tgInputMessage
	TgMessageType MessageType
}

var _ botinput.AudioMessage = (*tgWebhookAudioMessage)(nil)

func (tgWebhookAudioMessage) InputType() botinput.Type {
	return botinput.TypeAudio
}

func newTgWebhookAudioMessage(input tgInput, tgMessageType MessageType, tgMessage *tgbotapi.Message) tgWebhookAudioMessage {
	return tgWebhookAudioMessage{
		tgInputMessage: newTelegramWebhookMessage(input, tgMessage),
		TgMessageType:  tgMessageType,
	}
}
