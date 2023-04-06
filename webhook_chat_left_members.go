package telegram

import "github.com/bots-go-framework/bots-fw/botsfw"

type tgWebhookLeftChatMembersMessage struct {
	tgWebhookMessage
}

func (tgWebhookLeftChatMembersMessage) InputType() botsfw.WebhookInputType {
	return botsfw.WebhookInputLeftChatMembers
}

var _ botsfw.WebhookLeftChatMembersMessage = (*tgWebhookLeftChatMembersMessage)(nil)

func newTgWebhookLeftChatMembersMessage(input tgWebhookInput) tgWebhookNewChatMembersMessage {
	return tgWebhookNewChatMembersMessage{tgWebhookMessage: newTelegramWebhookMessage(input, input.update.Message)}
}

func (m *tgWebhookLeftChatMembersMessage) LeftChatMembers() []botsfw.WebhookActor {
	return []botsfw.WebhookActor{m.message.LeftChatMember}
}
