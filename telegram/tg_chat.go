package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"strconv"
)

// TgChat is wrapper for Telegram chat
type TgChat struct {
	chat *tgbotapi.Chat
}

var _ botinput.Chat = (*TgChat)(nil)

// GetID returns telegram chat ID
func (wh TgChat) GetID() string {
	return strconv.FormatInt(wh.chat.ID, 10)
}

// GetType returns telegram chat type
func (wh TgChat) GetType() string {
	return wh.chat.Type
}

// IsGroupChat indicates type of chat (group or private)
func (wh TgChat) IsGroupChat() bool {
	return !wh.chat.IsPrivate()
}
