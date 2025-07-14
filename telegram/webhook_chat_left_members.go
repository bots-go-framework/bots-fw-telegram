package telegram

import (
	"github.com/bots-go-framework/bots-fw/botinput"
)

type tgWebhookLeftChatMembersMessage struct {
	tgWebhookMessage
}

func (*tgWebhookLeftChatMembersMessage) InputType() botinput.Type {
	return botinput.TypeLeftChatMembers
}

var _ botinput.LeftChatMembersMessage = (*tgWebhookLeftChatMembersMessage)(nil)

func newTgWebhookLeftChatMembersMessage(input tgWebhookInput) tgWebhookNewChatMembersMessage {
	return tgWebhookNewChatMembersMessage{tgWebhookMessage: newTelegramWebhookMessage(input, input.update.Message)}
}

func (m *tgWebhookLeftChatMembersMessage) LeftChatMembers() []botinput.Actor {
	return []botinput.Actor{m.message.LeftChatMember}
}
