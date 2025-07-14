package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"strconv"
)

var (
	_ TgWebhookInput         = (*TgWebhookCallbackQuery)(nil)
	_ botinput.InputMessage  = (*TgWebhookCallbackQuery)(nil)
	_ botinput.CallbackQuery = (*TgWebhookCallbackQuery)(nil)
	_ WebhookCallbackQuery   = (*TgWebhookCallbackQuery)(nil)
)

// TgWebhookCallbackQuery is wrapper on callback query
type TgWebhookCallbackQuery struct { // TODO: make non-exportable
	tgWebhookInput
	//callbackQuery *tgbotapi.CallbackQuery
	//message       botsfw.WebhookMessage
}

func (whi TgWebhookCallbackQuery) GetChatInstanceID() string {
	return whi.update.CallbackQuery.ChatInstance
}

func newTelegramWebhookCallbackQuery(input tgWebhookInput) TgWebhookCallbackQuery {
	if input.update.CallbackQuery == nil {
		panic("update.CallbackQuery == nil")
	}
	return TgWebhookCallbackQuery{
		tgWebhookInput: input,
	}
}

// GetID returns update ID
func (whi TgWebhookCallbackQuery) GetID() string {
	return whi.update.CallbackQuery.ID
}

// Sequence returns update ID
func (whi TgWebhookCallbackQuery) Sequence() int {
	return whi.update.UpdateID
}

// GetMessage returns message
func (whi TgWebhookCallbackQuery) GetMessage() botinput.Message {
	return newTelegramWebhookMessage(whi.tgWebhookInput, whi.update.CallbackQuery.Message)
}

// TelegramCallbackMessage returns message
func (whi TgWebhookCallbackQuery) TelegramCallbackMessage() *tgbotapi.Message {
	return whi.update.CallbackQuery.Message
}

// GetFrom returns sender
func (whi TgWebhookCallbackQuery) GetFrom() botinput.Sender {
	return tgWebhookUser{tgUser: whi.update.CallbackQuery.From}
}

// GetData returns callback query data
func (whi TgWebhookCallbackQuery) GetData() string {
	return whi.update.CallbackQuery.Data
}

// GetInlineMessageID returns callback query inline message ID
func (whi TgWebhookCallbackQuery) GetInlineMessageID() string {
	return whi.update.CallbackQuery.InlineMessageID
}

// BotChatID returns bot chat ID
func (whi TgWebhookCallbackQuery) BotChatID() (string, error) {
	if cbq := whi.update.CallbackQuery; cbq.Message != nil && cbq.Message.Chat != nil {
		return strconv.FormatInt(cbq.Message.Chat.ID, 10), nil
	}
	return "", nil
}

// EditMessageOnCallbackQuery creates edit message
func EditMessageOnCallbackQuery(whcbq botinput.CallbackQuery, parseMode, text string) *tgbotapi.EditMessageTextConfig {
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
