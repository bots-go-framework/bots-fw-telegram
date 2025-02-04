package telegram

import (
	"errors"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"strconv"
)

// GetMessageUID returns UID of the message to be edited
func GetMessageUID(whc botsfw.WebhookContext) (*ChatMessageUID, error) {
	if whc == nil {
		return nil, errors.New("whc is a required parameter for func GetMessageUID()")
	}
	input := whc.Input()
	chatID, err := input.BotChatID()
	if err != nil {
		return nil, err
	}
	var tgChatID int64
	if tgChatID, err = strconv.ParseInt(chatID, 10, 64); err != nil {
		return nil, err
	}
	messageID := whc.Input().(TgWebhookCallbackQuery).GetMessage().IntID()
	return NewChatMessageUID(tgChatID, int(messageID)), nil
}
