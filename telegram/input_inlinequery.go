package telegram

import (
	"github.com/bots-go-framework/bots-fw/botinput"
)

var (
	_ botinput.InputMessage = (*TgWebhookInlineQuery)(nil)
	_ botinput.InlineQuery  = (*TgWebhookInlineQuery)(nil)
)

// TgWebhookInlineQuery is wrapper
type TgWebhookInlineQuery struct {
	tgInput
}

func newTelegramWebhookInlineQuery(input tgInput) TgWebhookInlineQuery {
	return TgWebhookInlineQuery{tgInput: input}
}

// GetInlineQueryID return inline query ID
func (iq TgWebhookInlineQuery) GetInlineQueryID() string {
	return iq.update.InlineQuery.ID
}

// GetQuery returns query string
func (iq TgWebhookInlineQuery) GetQuery() string {
	return iq.update.InlineQuery.Query
}

// GetFrom returns recipient
func (iq TgWebhookInlineQuery) GetFrom() botinput.Sender {
	return tgWebhookUser{tgUser: iq.update.InlineQuery.From}
}

// GetOffset returns offset
func (iq TgWebhookInlineQuery) GetOffset() string {
	return iq.update.InlineQuery.Offset
}

// BotChatID returns bot chat ID
func (iq TgWebhookInlineQuery) BotChatID() (string, error) {
	return "", nil
}
