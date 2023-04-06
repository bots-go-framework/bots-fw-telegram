package telegram

import "github.com/bots-go-framework/bots-fw/botsfw"

type tgWebhookChosenInlineResult struct {
	tgWebhookInput
}

var _ botsfw.WebhookChosenInlineResult = (*tgWebhookChosenInlineResult)(nil)

func (tgWebhookChosenInlineResult) InputType() botsfw.WebhookInputType {
	return botsfw.WebhookInputChosenInlineResult
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

func (q tgWebhookChosenInlineResult) GetFrom() botsfw.WebhookSender {
	return tgSender{tgUser: q.update.ChosenInlineResult.From}
}

func (q tgWebhookChosenInlineResult) BotChatID() (string, error) {
	return "", nil
}
