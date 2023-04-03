package telegram

import (
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo/bots-framework/botsfw"
)

type tgWebhookStickerMessage struct {
	tgWebhookMessage
	TgMessageType TgMessageType
}

var _ botsfw.WebhookStickerMessage = (*tgWebhookStickerMessage)(nil)

func (tgWebhookStickerMessage) InputType() botsfw.WebhookInputType {
	return botsfw.WebhookInputSticker
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
