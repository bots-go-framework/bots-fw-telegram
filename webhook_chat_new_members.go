package telegram

import "github.com/strongo/bots-framework/botsfw"

type tgWebhookNewChatMembersMessage struct {
	tgWebhookMessage
}

func (tgWebhookNewChatMembersMessage) InputType() botsfw.WebhookInputType {
	return botsfw.WebhookInputNewChatMembers
}

var _ botsfw.WebhookNewChatMembersMessage = (*tgWebhookNewChatMembersMessage)(nil)

func newTgWebhookNewChatMembersMessage(input tgWebhookInput) tgWebhookNewChatMembersMessage {
	return tgWebhookNewChatMembersMessage{tgWebhookMessage: newTelegramWebhookMessage(input, input.update.Message)}
}

func (m tgWebhookNewChatMembersMessage) NewChatMembers() []botsfw.WebhookActor {
	members := make([]botsfw.WebhookActor, len(m.message.NewChatMembers))
	for i, m := range m.message.NewChatMembers {
		members[i] = m
	}
	return members
}
