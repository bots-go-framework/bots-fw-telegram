package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
)

type tgWebhookTextMessage struct {
	tgWebhookMessage
	TgMessageType TgMessageType
}

var _ botinput.WebhookTextMessage = (*tgWebhookTextMessage)(nil)

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
