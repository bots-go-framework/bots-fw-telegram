package telegram

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"time"
)

var (
	_ TgWebhookInput             = (*tgWebhookSuccessfulPayment)(nil)
	_ botinput.InputMessage      = (*tgWebhookSuccessfulPayment)(nil)
	_ botinput.SuccessfulPayment = (*tgWebhookSuccessfulPayment)(nil)
)

type tgWebhookSuccessfulPayment struct {
	tgWebhookMessage
}

func (t tgWebhookSuccessfulPayment) GetSubscriptionExpirationDate() time.Time {
	return time.Unix(t.update.Message.SuccessfulPayment.SubscriptionExpirationDate, 0)
}

func (t tgWebhookSuccessfulPayment) GetIsRecurring() bool {
	return t.update.Message.SuccessfulPayment.IsRecurring
}

func (t tgWebhookSuccessfulPayment) GetIsFirstRecurring() bool {
	return t.update.Message.SuccessfulPayment.IsFirstRecurring
}

func (t tgWebhookSuccessfulPayment) GetShippingOptionID() string {
	return t.update.Message.SuccessfulPayment.ShippingOptionID
}

func (t tgWebhookSuccessfulPayment) GetOrderInfo() botinput.OrderInfo {
	if t.update.Message.SuccessfulPayment.OrderInfo == nil {
		return nil
	}
	oi := tgOrderInfo(*t.update.Message.SuccessfulPayment.OrderInfo)
	return &oi
}

func (t tgWebhookSuccessfulPayment) GetCurrency() string {
	return t.update.Message.SuccessfulPayment.Currency
}

func (t tgWebhookSuccessfulPayment) GetTotalAmount() int {
	return t.update.Message.SuccessfulPayment.TotalAmount
}

func (t tgWebhookSuccessfulPayment) GetInvoicePayload() string {
	return t.update.Message.SuccessfulPayment.InvoicePayload
}

func (t tgWebhookSuccessfulPayment) GetMessengerChargeID() string {
	return t.update.Message.SuccessfulPayment.TelegramPaymentChargeID
}

func (t tgWebhookSuccessfulPayment) GetPaymentProviderChargeID() string {
	return t.update.Message.SuccessfulPayment.ProviderPaymentChargeID
}

func newTgWebhookSuccessfulPayment(input tgWebhookInput) tgWebhookSuccessfulPayment {
	if input.update.Message.SuccessfulPayment == nil {
		panic("update.Message.SuccessfulPayment == nil")
	}
	return tgWebhookSuccessfulPayment{
		tgWebhookMessage: tgWebhookMessage{
			tgWebhookInput: input,
			message:        input.update.Message,
		},
	}
}
