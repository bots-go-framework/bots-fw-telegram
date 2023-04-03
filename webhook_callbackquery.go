package telegram

import (
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo/bots-framework/botsfw"
	"strconv"
)

// TgWebhookCallbackQuery is wrapper on callback query
type TgWebhookCallbackQuery struct { // TODO: make non-exportable
	tgWebhookInput
	//callbackQuery *tgbotapi.CallbackQuery
	//message       botsfw.WebhookMessage
}

var (
	_ botsfw.WebhookCallbackQuery = (*TgWebhookCallbackQuery)(nil)
	_ TgWebhookInput              = (*TgWebhookCallbackQuery)(nil)
	_ botsfw.WebhookInput         = (*TgWebhookCallbackQuery)(nil)
)

// InputType return WebhookInputCallbackQuery
func (twhcbq TgWebhookCallbackQuery) InputType() botsfw.WebhookInputType {
	return botsfw.WebhookInputCallbackQuery
}

func newTelegramWebhookCallbackQuery(input tgWebhookInput) TgWebhookCallbackQuery {
	callbackQuery := input.update.CallbackQuery
	if callbackQuery == nil {
		panic("update.CallbackQuery == nil")
	}
	q := TgWebhookCallbackQuery{
		tgWebhookInput: input,
	}
	return q
}

// GetID returns update ID
func (twhcbq TgWebhookCallbackQuery) GetID() string {
	return twhcbq.update.CallbackQuery.ID
}

// Sequence returns update ID
func (twhcbq TgWebhookCallbackQuery) Sequence() int {
	return twhcbq.update.UpdateID
}

// GetMessage returns message
func (twhcbq TgWebhookCallbackQuery) GetMessage() botsfw.WebhookMessage {
	return newTelegramWebhookMessage(twhcbq.tgWebhookInput, twhcbq.update.CallbackQuery.Message)
}

// TelegramCallbackMessage returns message
func (twhcbq TgWebhookCallbackQuery) TelegramCallbackMessage() *tgbotapi.Message {
	return twhcbq.update.CallbackQuery.Message
}

// GetFrom returns sender
func (twhcbq TgWebhookCallbackQuery) GetFrom() botsfw.WebhookSender {
	return tgSender{tgUser: twhcbq.update.CallbackQuery.From}
}

// GetData returns callback query data
func (twhcbq TgWebhookCallbackQuery) GetData() string {
	return twhcbq.update.CallbackQuery.Data
}

// GetInlineMessageID returns callback query inline message ID
func (twhcbq TgWebhookCallbackQuery) GetInlineMessageID() string {
	return twhcbq.update.CallbackQuery.InlineMessageID
}

// BotChatID returns bot chat ID
func (twhcbq TgWebhookCallbackQuery) BotChatID() (string, error) {
	if cbq := twhcbq.update.CallbackQuery; cbq.Message != nil && cbq.Message.Chat != nil {
		return strconv.FormatInt(cbq.Message.Chat.ID, 10), nil
	}
	return "", nil
}

// EditMessageOnCallbackQuery creates edit message
func EditMessageOnCallbackQuery(whcbq botsfw.WebhookCallbackQuery, parseMode, text string) *tgbotapi.EditMessageTextConfig {
	twhcbq := whcbq.(TgWebhookCallbackQuery)
	callbackQuery := twhcbq.update.CallbackQuery

	emc := tgbotapi.EditMessageTextConfig{
		Text:      text,
		ParseMode: parseMode,
		BaseEdit: tgbotapi.BaseEdit{
			InlineMessageID: callbackQuery.InlineMessageID,
		},
	}
	if emc.InlineMessageID == "" {
		emc.ChatID = callbackQuery.Message.Chat.ID
		emc.MessageID = callbackQuery.Message.MessageID
	}
	return &emc
}
