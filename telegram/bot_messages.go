package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

var _ botsfw.BotMessage = (*InlineBotMessage)(nil)

// InlineBotMessage is a wrapper for Telegram bot message
type InlineBotMessage tgbotapi.InlineConfig

// BotMessageType returns BotMessageTypeInlineResults
func (InlineBotMessage) BotMessageType() botsfw.BotMessageType {
	return botsfw.BotMessageTypeInlineResults
}

func (m InlineBotMessage) BotEndpoint() string {
	return (tgbotapi.InlineConfig)(m).TelegramMethod()
}

var _ botsfw.BotMessage = (*CallbackAnswer)(nil)

// CallbackAnswer is a callback answer message
type CallbackAnswer tgbotapi.AnswerCallbackQueryConfig

// BotMessageType returns BotMessageTypeCallbackAnswer
func (CallbackAnswer) BotMessageType() botsfw.BotMessageType {
	return botsfw.BotMessageTypeCallbackAnswer
}

func (a CallbackAnswer) BotEndpoint() string {
	return (tgbotapi.AnswerCallbackQueryConfig)(a).TelegramMethod()
}

var _ botsfw.BotMessage = (*LeaveChat)(nil)

// LeaveChat is a leave chat message from bot
type LeaveChat tgbotapi.LeaveChatConfig

// BotMessageType return BotMessageTypeLeaveChat
func (LeaveChat) BotMessageType() botsfw.BotMessageType {
	return botsfw.BotMessageTypeLeaveChat
}

func (m LeaveChat) BotEndpoint() string {
	return (tgbotapi.LeaveChatConfig)(m).TelegramMethod()
}

var _ botsfw.BotMessage = (*ExportChatInviteLink)(nil)

// ExportChatInviteLink is a TG message
type ExportChatInviteLink tgbotapi.ExportChatInviteLink

// BotMessageType returns BotMessageTypeExportChatInviteLink
func (ExportChatInviteLink) BotMessageType() botsfw.BotMessageType {
	return botsfw.BotMessageTypeExportChatInviteLink
}

func (m ExportChatInviteLink) BotEndpoint() string {
	return (tgbotapi.ExportChatInviteLink)(m).TelegramMethod()
}

var _ botsfw.BotMessage = (*Invoice)(nil)

type Invoice tgbotapi.InvoiceConfig

func (Invoice) BotMessageType() botsfw.BotMessageType {
	return botsfw.BotMessageTypeSendInvoice
}

func (m Invoice) BotEndpoint() string {
	i := (tgbotapi.InvoiceConfig)(m)
	return i.TelegramMethod()
}

type PreCheckoutQueryAnswer tgbotapi.AnswerPreCheckoutQueryConfig

func (PreCheckoutQueryAnswer) BotMessageType() botsfw.BotMessageType {
	return botsfw.BotMessageTypeAnswerPreCheckoutQuery
}

var _ botsfw.BotMessage = SetBotDescription{}

type SetBotDescription tgbotapi.SetMyDescription

func (SetBotDescription) BotMessageType() botsfw.BotMessageType {
	return botsfw.BotMessageTypeSetDescription
}

func (m SetBotDescription) BotEndpoint() string {
	return (tgbotapi.SetMyDescription)(m).TelegramMethod()
}

type SetBotShortDescription tgbotapi.SetMyShortDescription

func (SetBotShortDescription) BotMessageType() botsfw.BotMessageType {
	return botsfw.BotMessageTypeSetShortDescription
}

func (m SetBotShortDescription) BotEndpoint() string {
	return (tgbotapi.SetMyShortDescription)(m).TelegramMethod()
}

type SetBotCommands tgbotapi.SetMyCommandsConfig

func (SetBotCommands) BotMessageType() botsfw.BotMessageType {
	return botsfw.BotMessageTypeSetCommands
}

func (m SetBotCommands) BotEndpoint() string {
	return (tgbotapi.SetMyCommandsConfig)(m).TelegramMethod()
}
