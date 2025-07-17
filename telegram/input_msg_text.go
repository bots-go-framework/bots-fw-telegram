package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
)

var (
	_ botinput.Message      = (*tgWebhookTextMessage)(nil)
	_ botinput.InputMessage = (*tgWebhookTextMessage)(nil)
	_ botinput.TextMessage  = (*tgWebhookTextMessage)(nil)
)

type tgWebhookTextMessage struct {
	tgInputMessage
	TgMessageType MessageType
}

func (tgWebhookTextMessage) InputType() botinput.Type {
	return botinput.TypeText
}

func newTgWebhookTextMessage(input tgInput, tgMessageType MessageType, tgMessage *tgbotapi.Message) tgWebhookTextMessage {
	return tgWebhookTextMessage{
		tgInputMessage: newTelegramWebhookMessage(input, tgMessage),
		TgMessageType:  tgMessageType,
	}
}

func (whm tgWebhookTextMessage) Text() string {
	return whm.message.Text
}

func (whm tgWebhookTextMessage) IsEdited() bool {
	return whm.TgMessageType == MessageTypeEdited || whm.TgMessageType == MessageTypeEditedChannelPost
}
