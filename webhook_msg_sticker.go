package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
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
