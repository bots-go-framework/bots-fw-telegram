package telegram

import (
	"github.com/bots-go-framework/bots-fw/botinput"
)

var (
	_ botinput.WebhookInput                 = (*tgWebhookNewChatMembersMessage)(nil)
	_ botinput.WebhookMessage               = (*tgWebhookNewChatMembersMessage)(nil)
	_ botinput.WebhookNewChatMembersMessage = (*tgWebhookNewChatMembersMessage)(nil)
)

type tgWebhookNewChatMembersMessage struct {
	tgWebhookMessage
}

func newTgWebhookNewChatMembersMessage(input tgWebhookInput) tgWebhookNewChatMembersMessage {
	return tgWebhookNewChatMembersMessage{tgWebhookMessage: newTelegramWebhookMessage(input, input.update.Message)}
}

func (m tgWebhookNewChatMembersMessage) NewChatMembers() []botinput.WebhookActor {
	members := make([]botinput.WebhookActor, len(m.message.NewChatMembers))
	for i, m := range m.message.NewChatMembers {
		members[i] = m
	}
	return members
}
