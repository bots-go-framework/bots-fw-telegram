package telegram

import (
	"strconv"

	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botmsg"
)

// chatActionSendable maps a platform-neutral botmsg.ChatAction to a Telegram
// sendChatAction config. A numeric ChatID is used as the chat id, a non-numeric
// one as a channel username, and an empty ChatID falls back to defaultChatID
// (the current chat).
func chatActionSendable(chatAction botmsg.ChatAction, defaultChatID int64) tgbotapi.ChatActionConfig {
	config := tgbotapi.ChatActionConfig{Action: chatAction.Action}
	if chatAction.ChatID != "" {
		if id, err := strconv.ParseInt(chatAction.ChatID, 10, 64); err == nil {
			config.ChatID = id
		} else {
			config.ChannelUsername = chatAction.ChatID
		}
	}
	if config.ChatID == 0 && config.ChannelUsername == "" {
		config.ChatID = defaultChatID
	}
	return config
}
