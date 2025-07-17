package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
)

type tgWebhookStickerMessage struct {
	tgInputMessage
	TgMessageType MessageType
}

var _ botinput.StickerMessage = (*tgWebhookStickerMessage)(nil)

func (tgWebhookStickerMessage) InputType() botinput.Type {
	return botinput.TypeSticker
}

func newTgWebhookStickerMessage(input tgInput, tgMessageType MessageType, tgMessage *tgbotapi.Message) tgWebhookStickerMessage {
	return tgWebhookStickerMessage{
		tgInputMessage: newTelegramWebhookMessage(input, tgMessage),
		TgMessageType:  tgMessageType,
	}
}

//func (whm tgWebhookStickerMessage) IsEdited() bool {
//	return whm.MessageType == MessageTypeEdited || whm.MessageType == MessageTypeEditedChannelPost
//}
