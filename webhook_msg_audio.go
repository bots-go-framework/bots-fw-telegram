package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
)

type tgWebhookAudioMessage struct {
	tgWebhookMessage
	TgMessageType TgMessageType
}

var _ botinput.WebhookAudioMessage = (*tgWebhookAudioMessage)(nil)

func (tgWebhookAudioMessage) InputType() botinput.WebhookInputType {
	return botinput.WebhookInputAudio
}

func newTgWebhookAudioMessage(input tgWebhookInput, tgMessageType TgMessageType, tgMessage *tgbotapi.Message) tgWebhookAudioMessage {
	return tgWebhookAudioMessage{
		tgWebhookMessage: newTelegramWebhookMessage(input, tgMessage),
		TgMessageType:    tgMessageType,
	}
}
