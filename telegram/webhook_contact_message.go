package telegram

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"strconv"
)

var (
	_ botinput.InputMessage   = (*tgWebhookContactMessage)(nil)
	_ botinput.Message        = (*tgWebhookContactMessage)(nil)
	_ botinput.ContactMessage = (*tgWebhookContactMessage)(nil)
)

type tgWebhookContactMessage struct {
	tgWebhookMessage
}

func (m tgWebhookContactMessage) GetVCard() string {
	return m.update.Message.Contact.VCard
}

func newTgWebhookContact(input tgWebhookInput) tgWebhookContactMessage {
	return tgWebhookContactMessage{tgWebhookMessage: newTelegramWebhookMessage(input, input.update.Message)}
}

func (m tgWebhookContactMessage) GetFirstName() string {
	return m.update.Message.Contact.FirstName
}

func (m tgWebhookContactMessage) GetLastName() string {
	return m.update.Message.Contact.LastName
}

func (m tgWebhookContactMessage) GetPhoneNumber() string {
	return m.update.Message.Contact.PhoneNumber
}

func (m tgWebhookContactMessage) GetBotUserID() string {
	return strconv.FormatInt(m.update.Message.Contact.UserID, 10)
}
