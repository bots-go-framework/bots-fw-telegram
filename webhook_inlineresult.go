package telegram

import (
	"github.com/bots-go-framework/bots-fw/botinput"
)

type tgWebhookChosenInlineResult struct {
	tgWebhookInput
}

var _ botinput.WebhookChosenInlineResult = (*tgWebhookChosenInlineResult)(nil)

func (tgWebhookChosenInlineResult) InputType() botinput.WebhookInputType {
	return botinput.WebhookInputChosenInlineResult
}

func newTelegramWebhookChosenInlineResult(input tgWebhookInput) tgWebhookChosenInlineResult {
	return tgWebhookChosenInlineResult{tgWebhookInput: input}
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

func (q tgWebhookChosenInlineResult) GetFrom() botinput.WebhookSender {
	return tgWebhookUser{tgUser: q.update.ChosenInlineResult.From}
}

func (q tgWebhookChosenInlineResult) BotChatID() (string, error) {
	return "", nil
}
