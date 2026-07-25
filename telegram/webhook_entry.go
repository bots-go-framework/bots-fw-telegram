package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"strconv"
	"time"
)

type tgWebhookEntry struct {
	update *tgbotapi.Update
}

var _ botinput.Entry = (*tgWebhookEntry)(nil)
var _ botinput.DurableWebhookEntry = (*tgWebhookEntry)(nil)

func (e tgWebhookEntry) GetID() interface{} {
	return e.update.UpdateID
}

// WebhookUpdateID declares Telegram's UpdateID as the stable delivery ID used
// by the framework's durable webhook inbox. UpdateID zero is Telegram's absent
// default and must not be used as a deduplication key.
func (e tgWebhookEntry) WebhookUpdateID() (string, bool) {
	if e.update == nil || e.update.UpdateID == 0 {
		return "", false
	}
	return strconv.Itoa(e.update.UpdateID), true
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
