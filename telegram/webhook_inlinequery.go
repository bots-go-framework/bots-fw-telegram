package telegram

import (
	"github.com/bots-go-framework/bots-fw/botinput"
)

var (
	_ botinput.WebhookInput       = (*TgWebhookInlineQuery)(nil)
	_ botinput.WebhookInlineQuery = (*TgWebhookInlineQuery)(nil)
)

// TgWebhookInlineQuery is wrapper
type TgWebhookInlineQuery struct {
	tgWebhookInput
}

var _ botinput.WebhookInlineQuery = (*TgWebhookInlineQuery)(nil)

func newTelegramWebhookInlineQuery(input tgWebhookInput) TgWebhookInlineQuery {
	return TgWebhookInlineQuery{tgWebhookInput: input}
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
func (iq TgWebhookInlineQuery) GetFrom() botinput.WebhookSender {
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
