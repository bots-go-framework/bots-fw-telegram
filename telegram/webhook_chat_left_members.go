package telegram

import (
	"github.com/bots-go-framework/bots-fw/botinput"
)

type tgWebhookLeftChatMembersMessage struct {
	tgWebhookMessage
}

func (*tgWebhookLeftChatMembersMessage) InputType() botinput.WebhookInputType {
	return botinput.WebhookInputLeftChatMembers
}

var _ botinput.WebhookLeftChatMembersMessage = (*tgWebhookLeftChatMembersMessage)(nil)

func newTgWebhookLeftChatMembersMessage(input tgWebhookInput) tgWebhookNewChatMembersMessage {
	return tgWebhookNewChatMembersMessage{tgWebhookMessage: newTelegramWebhookMessage(input, input.update.Message)}
}

func (m *tgWebhookLeftChatMembersMessage) LeftChatMembers() []botinput.WebhookActor {
	return []botinput.WebhookActor{m.message.LeftChatMember}
}
