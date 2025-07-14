package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
)

type tgWebhookStickerMessage struct {
	tgWebhookMessage
	TgMessageType TgMessageType
}

var _ botinput.StickerMessage = (*tgWebhookStickerMessage)(nil)

func (tgWebhookStickerMessage) InputType() botinput.Type {
	return botinput.TypeSticker
}

func newTgWebhookStickerMessage(input tgWebhookInput, tgMessageType TgMessageType, tgMessage *tgbotapi.Message) tgWebhookStickerMessage {
	return tgWebhookStickerMessage{
		tgWebhookMessage: newTelegramWebhookMessage(input, tgMessage),
		TgMessageType:    tgMessageType,
	}
}

//func (whm tgWebhookStickerMessage) IsEdited() bool {
//	return whm.TgMessageType == TgMessageTypeEdited || whm.TgMessageType == TgMessageTypeEditedChannelPost
//}
