package telegram

import (
	"github.com/bots-go-framework/bots-fw/botinput"
)

type tgWebhookChosenInlineResult struct {
	tgInput
}

var _ botinput.InputMessage = (*tgWebhookChosenInlineResult)(nil)
var _ botinput.ChosenInlineResult = (*tgWebhookChosenInlineResult)(nil)

func (tgWebhookChosenInlineResult) InputType() botinput.Type {
	return botinput.TypeChosenInlineResult
}

func newTelegramWebhookChosenInlineResult(input tgInput) tgWebhookChosenInlineResult {
	return tgWebhookChosenInlineResult{tgInput: input}
}

func (q tgWebhookChosenInlineResult) GetResultID() string {
	return q.update.ChosenInlineResult.ResultID
}

func (q tgWebhookChosenInlineResult) GetQuery() string {
	return q.update.ChosenInlineResult.Query
}

func (q tgWebhookChosenInlineResult) GetInlineMessageID() string {
	if q.update.ChosenInlineResult != nil {
		return q.update.ChosenInlineResult.InlineMessageID
	}
	return ""
}

func (q tgWebhookChosenInlineResult) GetFrom() botinput.Sender {
	return tgWebhookUser{tgUser: q.update.ChosenInlineResult.From}
}

func (q tgWebhookChosenInlineResult) BotChatID() (string, error) {
	return "", nil
}
