package telegram

import (
	"github.com/bots-go-framework/bots-fw/botinput"
)

var (
	_ TgWebhookInput           = (*tgWebhookRefundedPayment)(nil)
	_ botinput.InputMessage    = (*tgWebhookRefundedPayment)(nil)
	_ botinput.RefundedPayment = (*tgWebhookRefundedPayment)(nil)
)

type tgWebhookRefundedPayment struct {
	tgWebhookMessage
}

func (t tgWebhookRefundedPayment) GetCurrency() string {
	return t.update.Message.RefundedPayment.Currency
}

func (t tgWebhookRefundedPayment) GetTotalAmount() int {
	return t.update.Message.RefundedPayment.TotalAmount
}

func (t tgWebhookRefundedPayment) GetInvoicePayload() string {
	return t.update.Message.RefundedPayment.InvoicePayload
}

func (t tgWebhookRefundedPayment) GetMessengerChargeID() string {
	return t.update.Message.RefundedPayment.TelegramPaymentChargeID
}

func (t tgWebhookRefundedPayment) GetPaymentProviderChargeID() string {
	return t.update.Message.RefundedPayment.ProviderPaymentChargeID
}

func newTgWebhookRefundedPayment(input tgWebhookInput) tgWebhookRefundedPayment {
	if input.update.Message.RefundedPayment == nil {
		panic("update.Message.RefundedPayment == nil")
	}
	return tgWebhookRefundedPayment{
		tgWebhookMessage: tgWebhookMessage{
			tgWebhookInput: input,
			message:        input.update.Message,
		},
	}
}
