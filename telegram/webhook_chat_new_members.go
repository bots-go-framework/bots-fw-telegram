package telegram

import (
	"github.com/bots-go-framework/bots-fw/botinput"
)

var (
	_ botinput.InputMessage          = (*tgWebhookNewChatMembersMessage)(nil)
	_ botinput.Message               = (*tgWebhookNewChatMembersMessage)(nil)
	_ botinput.NewChatMembersMessage = (*tgWebhookNewChatMembersMessage)(nil)
)

type tgWebhookNewChatMembersMessage struct {
	tgWebhookMessage
}

func newTgWebhookNewChatMembersMessage(input tgWebhookInput) tgWebhookNewChatMembersMessage {
	return tgWebhookNewChatMembersMessage{tgWebhookMessage: newTelegramWebhookMessage(input, input.update.Message)}
}

func (m tgWebhookNewChatMembersMessage) NewChatMembers() []botinput.Actor {
	members := make([]botinput.Actor, len(m.message.NewChatMembers))
	for i, m := range m.message.NewChatMembers {
		members[i] = m
	}
	return members
}
