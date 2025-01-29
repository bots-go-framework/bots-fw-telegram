package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
)

var (
	_ botinput.WebhookInput       = (*tgWebhookTextMessage)(nil)
	_ botinput.WebhookMessage     = (*tgWebhookTextMessage)(nil)
	_ botinput.WebhookTextMessage = (*tgWebhookTextMessage)(nil)
)

type tgWebhookTextMessage struct {
	tgWebhookMessage
	TgMessageType TgMessageType
}

func (tgWebhookTextMessage) InputType() botinput.WebhookInputType {
	return botinput.WebhookInputText
}

func newTgWebhookTextMessage(input tgWebhookInput, tgMessageType TgMessageType, tgMessage *tgbotapi.Message) tgWebhookTextMessage {
	return tgWebhookTextMessage{
		tgWebhookMessage: newTelegramWebhookMessage(input, tgMessage),
		TgMessageType:    tgMessageType,
	}
}

func (whm tgWebhookTextMessage) Text() string {
	return whm.message.Text
}

func (whm tgWebhookTextMessage) IsEdited() bool {
	return whm.TgMessageType == TgMessageTypeEdited || whm.TgMessageType == TgMessageTypeEditedChannelPost
}
