package telegram

import (
	"fmt"
	"github.com/strongo/bots-framework/botsfw"
)

type callbackCurrent struct {
}

// CallbackCurrent is what?
var CallbackCurrent botsfw.MessageUID = &callbackCurrent{}

func (callbackCurrent) UID() string {
	return "callbackCurrent"
}

// InlineMessageUID is inline message UID
type InlineMessageUID struct {
	InlineMessageID string
}

var _ botsfw.MessageUID = (*InlineMessageUID)(nil)

// NewInlineMessageUID creates new inline message UID
func NewInlineMessageUID(inlineMessageID string) *InlineMessageUID {
	return &InlineMessageUID{InlineMessageID: inlineMessageID}
}

// UID is unique ID of the message
func (m InlineMessageUID) UID() string {
	return m.InlineMessageID
}

// NewChatMessageUID create new ChatMessageUID
func NewChatMessageUID(chatID int64, messageID int) *ChatMessageUID {
	return &ChatMessageUID{ChatID: chatID, MessageID: messageID}
}

// ChatMessageUID is what?
type ChatMessageUID struct {
	ChatID    int64
	MessageID int
}

var _ botsfw.MessageUID = (*ChatMessageUID)(nil)

// UID return unique ID of the message
func (m ChatMessageUID) UID() string {
	return fmt.Sprintf("%d:%d", m.ChatID, m.MessageID)
}
