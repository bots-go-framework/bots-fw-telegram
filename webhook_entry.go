package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"time"
)

type tgWebhookEntry struct {
	update *tgbotapi.Update
}

var _ botinput.WebhookEntry = (*tgWebhookEntry)(nil)

func (e tgWebhookEntry) GetID() interface{} {
	return e.update.UpdateID
}

func (e tgWebhookEntry) GetTime() time.Time {
	if e.update.Message != nil {
		return e.update.Message.Time()
	}
	if e.update.EditedMessage != nil {
		return e.update.EditedMessage.Time()
	}
	panic("Both `update.Message` & `update.EditedMessage` are nil.")
}
