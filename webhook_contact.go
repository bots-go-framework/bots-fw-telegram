package telegram

import "github.com/bots-go-framework/bots-fw/botsfw"

type tgWebhookContactMessage struct {
	tgWebhookMessage
}

func (tgWebhookContactMessage) InputType() botsfw.WebhookInputType {
	return botsfw.WebhookInputContact
}

var _ botsfw.WebhookContactMessage = (*tgWebhookContactMessage)(nil)

func newTgWebhookContact(input tgWebhookInput) tgWebhookContactMessage {
	return tgWebhookContactMessage{tgWebhookMessage: newTelegramWebhookMessage(input, input.update.Message)}
}

func (m tgWebhookContactMessage) FirstName() string {
	return m.update.Message.Contact.FirstName
}

func (m tgWebhookContactMessage) LastName() string {
	return m.update.Message.Contact.LastName
}

func (m tgWebhookContactMessage) PhoneNumber() string {
	return m.update.Message.Contact.PhoneNumber
}

func (m tgWebhookContactMessage) UserID() interface{} {
	return m.update.Message.Contact.UserID
}
