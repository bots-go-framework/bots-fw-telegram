package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"strconv"
)

type WebhookCallbackQuery interface {
	botinput.CallbackQuery
	GetInlineMessageID() string // Telegram only?
	GetChatInstanceID() string  // Telegram only?
}

var (
	_ TgWebhookInput         = (*callbackQueryInput)(nil)
	_ botinput.InputMessage  = (*callbackQueryInput)(nil)
	_ botinput.CallbackQuery = (*callbackQueryInput)(nil)
	_ WebhookCallbackQuery   = (*callbackQueryInput)(nil)
)

// callbackQueryInput is wrapper on callback query
type callbackQueryInput struct { // TODO: make non-exportable
	tgInput
	//callbackQuery *tgbotapi.CallbackQuery
	//message       botsfw.WebhookMessage
}

func (whi callbackQueryInput) GetChatInstanceID() string {
	return whi.update.CallbackQuery.ChatInstance
}

func newTelegramWebhookCallbackQuery(input tgInput) callbackQueryInput {
	if input.update.CallbackQuery == nil {
		panic("update.CallbackQuery == nil")
	}
	return callbackQueryInput{
		tgInput: input,
	}
}

// GetID returns update ID
func (whi callbackQueryInput) GetID() string {
	return whi.update.CallbackQuery.ID
}

// Sequence returns update ID
func (whi callbackQueryInput) Sequence() int {
	return whi.update.UpdateID
}

// GetMessage returns message
func (whi callbackQueryInput) GetMessage() botinput.Message {
	return newTelegramWebhookMessage(whi.tgInput, whi.update.CallbackQuery.Message)
}

// TelegramCallbackMessage returns message
func (whi callbackQueryInput) TelegramCallbackMessage() *tgbotapi.Message {
	return whi.update.CallbackQuery.Message
}

// GetFrom returns sender
func (whi callbackQueryInput) GetFrom() botinput.Sender {
	return tgWebhookUser{tgUser: whi.update.CallbackQuery.From}
}

// GetData returns callback query data
func (whi callbackQueryInput) GetData() string {
	return whi.update.CallbackQuery.Data
}

// GetInlineMessageID returns callback query inline message ID
func (whi callbackQueryInput) GetInlineMessageID() string {
	return whi.update.CallbackQuery.InlineMessageID
}

// BotChatID returns bot chat ID
func (whi callbackQueryInput) BotChatID() (string, error) {
	if cbq := whi.update.CallbackQuery; cbq.Message != nil && cbq.Message.Chat != nil {
		return strconv.FormatInt(cbq.Message.Chat.ID, 10), nil
	}
	return "", nil
}

// EditMessageOnCallbackQuery creates edit message
func EditMessageOnCallbackQuery(whcbq botinput.CallbackQuery, parseMode, text string) *tgbotapi.EditMessageTextConfig {
	twhcbq := whcbq.(callbackQueryInput)
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
