package telegram

import (
	"github.com/bots-go-framework/bots-fw/botinput"
)

type tgWebhookNewChatMembersMessage struct {
	tgWebhookMessage
}

func (tgWebhookNewChatMembersMessage) InputType() botinput.WebhookInputType {
	return botinput.WebhookInputNewChatMembers
}

var _ botinput.WebhookNewChatMembersMessage = (*tgWebhookNewChatMembersMessage)(nil)

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
