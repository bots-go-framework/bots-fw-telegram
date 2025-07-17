package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"strconv"
)

type tgInputMessage struct {
	tgInput
	message *tgbotapi.Message // Can be either whi.update.Message or whi.update.CallbackQuery.Message
}

//func (v *tgInputMessage) GetMessage() botinput.Message {
//	return v.message
//}

func (whm *tgInputMessage) IntID() int {
	return whm.message.MessageID
}

func newTelegramWebhookMessage(input tgInput, message *tgbotapi.Message) tgInputMessage {
	if message == nil {
		panic("message == nil")
	}
	return tgInputMessage{tgInput: input, message: message}
}

func (whm tgInputMessage) BotChatID() (string, error) {
	return strconv.FormatInt(whm.message.Chat.ID, 10), nil
}
