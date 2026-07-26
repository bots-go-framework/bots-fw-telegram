package telegram

import (
	"errors"
	"net/url"

	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botmsg"
)

var _ botmsg.BotMessage = (*InlineBotMessage)(nil)

// InlineBotMessage is a wrapper for Telegram bot message
type InlineBotMessage tgbotapi.InlineConfig

// BotMessageType returns BotMessageTypeInlineResults
func (InlineBotMessage) BotMessageType() botmsg.Type {
	return botmsg.BotMessageTypeInlineResults
}

func (m InlineBotMessage) BotEndpoint() string {
	return (tgbotapi.InlineConfig)(m).TelegramMethod()
}

var _ botmsg.BotMessage = (*CallbackAnswer)(nil)

// CallbackAnswer is a callback answer message
type CallbackAnswer tgbotapi.AnswerCallbackQueryConfig

// BotMessageType returns BotMessageTypeCallbackAnswer
func (CallbackAnswer) BotMessageType() botmsg.Type {
	return botmsg.TypeCallbackAnswer
}

func (a CallbackAnswer) BotEndpoint() string {
	return (tgbotapi.AnswerCallbackQueryConfig)(a).TelegramMethod()
}

var _ botmsg.BotMessage = (*LeaveChat)(nil)

// LeaveChat is a leave chat message from bot
type LeaveChat tgbotapi.LeaveChatConfig

// BotMessageType return BotMessageTypeLeaveChat
func (LeaveChat) BotMessageType() botmsg.Type {
	return botmsg.TypeLeaveChat
}

func (m LeaveChat) BotEndpoint() string {
	return (tgbotapi.LeaveChatConfig)(m).TelegramMethod()
}

var _ botmsg.BotMessage = (*ExportChatInviteLink)(nil)

// ExportChatInviteLink is a TG message
type ExportChatInviteLink tgbotapi.ExportChatInviteLink

// BotMessageType returns BotMessageTypeExportChatInviteLink
func (ExportChatInviteLink) BotMessageType() botmsg.Type {
	return botmsg.TypeExportChatInviteLink
}

func (m ExportChatInviteLink) BotEndpoint() string {
	return (tgbotapi.ExportChatInviteLink)(m).TelegramMethod()
}

var _ botmsg.BotMessage = (*Invoice)(nil)

type Invoice tgbotapi.InvoiceConfig

func (Invoice) BotMessageType() botmsg.Type {
	return botmsg.TypeSendInvoice
}

func (m Invoice) BotEndpoint() string {
	i := (tgbotapi.InvoiceConfig)(m)
	return i.TelegramMethod()
}

type PreCheckoutQueryAnswer tgbotapi.AnswerPreCheckoutQueryConfig

func (PreCheckoutQueryAnswer) BotMessageType() botmsg.Type {
	return botmsg.TypeAnswerPreCheckoutQuery
}

var _ botmsg.BotMessage = SetBotDescription{}

type SetBotDescription tgbotapi.SetMyDescription

func (SetBotDescription) BotMessageType() botmsg.Type {
	return botmsg.TypeSetDescription
}

func (m SetBotDescription) BotEndpoint() string {
	return (tgbotapi.SetMyDescription)(m).TelegramMethod()
}

type SetBotShortDescription tgbotapi.SetMyShortDescription

func (SetBotShortDescription) BotMessageType() botmsg.Type {
	return botmsg.TypeSetShortDescription
}

func (m SetBotShortDescription) BotEndpoint() string {
	return (tgbotapi.SetMyShortDescription)(m).TelegramMethod()
}

type SetBotCommands tgbotapi.SetMyCommandsConfig

func (SetBotCommands) BotMessageType() botmsg.Type {
	return botmsg.TypeSetCommands
}

func (m SetBotCommands) BotEndpoint() string {
	return (tgbotapi.SetMyCommandsConfig)(m).TelegramMethod()
}

type SendPhoto tgbotapi.PhotoConfig

func (SendPhoto) BotMessageType() botmsg.Type {
	return botmsg.TypeSendPhoto
}

func (m SendPhoto) BotEndpoint() string {
	return (tgbotapi.PhotoConfig)(m).TelegramMethod()
}

// SendRichMessage is a handler-friendly Telegram BotMessage for a persistent
// native rich message. MessageFromBot.Keyboard is applied as its reply markup
// by the Telegram responder.
type SendRichMessage tgbotapi.RichMessageConfig

var _ botmsg.BotMessage = SendRichMessage{}

func (SendRichMessage) BotMessageType() botmsg.Type {
	return botmsg.TypeSendRichMessage
}

func (m SendRichMessage) BotEndpoint() string {
	return (tgbotapi.RichMessageConfig)(m).TelegramMethod()
}

// NewSendRichMessage constructs a persistent rich message. chatID may be zero:
// the responder then uses MessageFromBot.ToChat or the current update's chat.
func NewSendRichMessage(chatID int64, richMessage tgbotapi.InputRichMessage) SendRichMessage {
	return SendRichMessage(tgbotapi.RichMessageConfig{
		BaseChat:    tgbotapi.BaseChat{ChatID: chatID},
		RichMessage: richMessage,
	})
}

// SendRichMessageDraft is a handler-friendly Telegram BotMessage for streaming
// a temporary rich-message draft. Telegram returns bool for this method.
type SendRichMessageDraft tgbotapi.RichMessageDraftConfig

var _ botmsg.BotMessage = SendRichMessageDraft{}

func (SendRichMessageDraft) BotMessageType() botmsg.Type {
	return botmsg.TypeSendRichMessageDraft
}

func (m SendRichMessageDraft) BotEndpoint() string {
	return (tgbotapi.RichMessageDraftConfig)(m).TelegramMethod()
}

// NewSendRichMessageDraft constructs a streaming draft. chatID may be zero and
// is resolved by the responder in the same way as NewSendRichMessage.
func NewSendRichMessageDraft(chatID, draftID int64, richMessage tgbotapi.InputRichMessage) SendRichMessageDraft {
	return SendRichMessageDraft(tgbotapi.RichMessageDraftConfig{
		ChatID:      chatID,
		DraftID:     draftID,
		RichMessage: richMessage,
	})
}

// EditRichMessage is a handler-friendly Telegram BotMessage for replacing the
// content of an existing message with a native rich message.
type EditRichMessage tgbotapi.EditMessageTextConfig

var _ botmsg.BotMessage = EditRichMessage{}

func (EditRichMessage) BotMessageType() botmsg.Type {
	return botmsg.TypeEditRichMessage
}

func (m EditRichMessage) BotEndpoint() string {
	return (tgbotapi.EditMessageTextConfig)(m).TelegramMethod()
}

// NewEditRichMessage constructs an edit for an arbitrary chat message, which
// is useful for independently maintained per-player game cards.
func NewEditRichMessage(chatID int64, messageID int, richMessage tgbotapi.InputRichMessage) EditRichMessage {
	return EditRichMessage(tgbotapi.EditMessageTextConfig{
		BaseEdit:    tgbotapi.NewChatMessageEdit(chatID, messageID),
		RichMessage: &richMessage,
	})
}

// NewEditInlineRichMessage constructs a rich edit for an inline message.
func NewEditInlineRichMessage(inlineMessageID string, richMessage tgbotapi.InputRichMessage) EditRichMessage {
	return EditRichMessage(tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			InlineMessageID: inlineMessageID,
		},
		RichMessage: &richMessage,
	})
}

// TelegramRequest exposes any low-level Telegram Sendable through the standard
// MessageFromBot handler path. Prefer a dedicated wrapper such as
// SendRichMessage when one exists; use this for Telegram-only methods such as
// ephemeral-message edits and deletes.
type TelegramRequest struct {
	Request        tgbotapi.Sendable
	ReturnsMessage bool
}

var (
	_ botmsg.BotMessage = TelegramRequest{}
	_ tgbotapi.Sendable = TelegramRequest{}
)

// BotMessageType is intentionally undefined: tgWebhookResponder recognizes
// TelegramRequest by concrete type before the platform-neutral type switch.
func (TelegramRequest) BotMessageType() botmsg.Type {
	return botmsg.TypeUndefined
}

func (m TelegramRequest) TelegramMethod() string {
	if m.Request == nil {
		return ""
	}
	return m.Request.TelegramMethod()
}

func (m TelegramRequest) Values() (url.Values, error) {
	if m.Request == nil {
		return nil, errors.New("TelegramRequest.Request is required")
	}
	return m.Request.Values()
}

// NewTelegramMessageRequest wraps a Telegram method that returns Message.
func NewTelegramMessageRequest(request tgbotapi.Sendable) TelegramRequest {
	return TelegramRequest{Request: request, ReturnsMessage: true}
}

// NewTelegramBooleanRequest wraps a Telegram method that returns True.
func NewTelegramBooleanRequest(request tgbotapi.Sendable) TelegramRequest {
	return TelegramRequest{Request: request}
}
