package telegram

import (
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
