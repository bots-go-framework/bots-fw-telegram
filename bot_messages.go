package telegram

import (
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo/bots-framework/botsfw"
)

// InlineBotMessage is wrapper for Telegram bot message
type InlineBotMessage tgbotapi.InlineConfig

// BotMessageType returns BotMessageTypeInlineResults
func (InlineBotMessage) BotMessageType() botsfw.BotMessageType {
	return botsfw.BotMessageTypeInlineResults
}

// CallbackAnswer is callback answer message
type CallbackAnswer tgbotapi.AnswerCallbackQueryConfig

// BotMessageType returns BotMessageTypeCallbackAnswer
func (CallbackAnswer) BotMessageType() botsfw.BotMessageType {
	return botsfw.BotMessageTypeCallbackAnswer
}

// LeaveChat is leave chat message from bot
type LeaveChat tgbotapi.LeaveChatConfig

// BotMessageType return BotMessageTypeLeaveChat
func (LeaveChat) BotMessageType() botsfw.BotMessageType {
	return botsfw.BotMessageTypeLeaveChat
}

// ExportChatInviteLink is TG message
type ExportChatInviteLink tgbotapi.ExportChatInviteLink

// BotMessageType returns BotMessageTypeExportChatInviteLink
func (ExportChatInviteLink) BotMessageType() botsfw.BotMessageType {
	return botsfw.BotMessageTypeExportChatInviteLink
}
