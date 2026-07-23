package telegram

import (
	"testing"

	"github.com/bots-go-framework/bots-fw/botmsg"
)

func TestChatActionSendable(t *testing.T) {
	const defaultChat int64 = 999

	t.Run("empty chat id falls back to the current chat", func(t *testing.T) {
		got := chatActionSendable(botmsg.ChatAction{Action: botmsg.ChatActionTyping}, defaultChat)
		if got.Action != botmsg.ChatActionTyping {
			t.Errorf("Action = %q, want %q", got.Action, botmsg.ChatActionTyping)
		}
		if got.ChatID != defaultChat {
			t.Errorf("ChatID = %d, want fallback %d", got.ChatID, defaultChat)
		}
		if got.ChannelUsername != "" {
			t.Errorf("ChannelUsername = %q, want empty", got.ChannelUsername)
		}
	})

	t.Run("numeric chat id is used as the chat id", func(t *testing.T) {
		got := chatActionSendable(botmsg.ChatAction{Action: botmsg.ChatActionTyping, ChatID: "42"}, defaultChat)
		if got.ChatID != 42 {
			t.Errorf("ChatID = %d, want 42", got.ChatID)
		}
		if got.ChannelUsername != "" {
			t.Errorf("ChannelUsername = %q, want empty", got.ChannelUsername)
		}
	})

	t.Run("non-numeric chat id is used as a channel username", func(t *testing.T) {
		got := chatActionSendable(botmsg.ChatAction{Action: botmsg.ChatActionTyping, ChatID: "@mychannel"}, defaultChat)
		if got.ChannelUsername != "@mychannel" {
			t.Errorf("ChannelUsername = %q, want @mychannel", got.ChannelUsername)
		}
		if got.ChatID != 0 {
			t.Errorf("ChatID = %d, want 0", got.ChatID)
		}
	})

	t.Run("Values carries the action for sendChatAction", func(t *testing.T) {
		got := chatActionSendable(botmsg.ChatAction{Action: botmsg.ChatActionTyping}, defaultChat)
		if method := got.TelegramMethod(); method != "sendChatAction" {
			t.Errorf("TelegramMethod() = %q, want sendChatAction", method)
		}
		values, err := got.Values()
		if err != nil {
			t.Fatalf("Values() error: %v", err)
		}
		if values.Get("action") != botmsg.ChatActionTyping {
			t.Errorf("values[action] = %q, want %q", values.Get("action"), botmsg.ChatActionTyping)
		}
	})
}
