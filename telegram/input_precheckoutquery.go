package telegram

import (
	"github.com/bots-go-framework/bots-fw/botinput"
)

var (
	_ TgWebhookInput            = (*tgWebhookPreCheckoutQuery)(nil)
	_ botinput.InputMessage     = (*tgWebhookPreCheckoutQuery)(nil)
	_ botinput.PreCheckoutQuery = (*tgWebhookPreCheckoutQuery)(nil)
)

type tgWebhookPreCheckoutQuery struct {
	tgInput
}

func (q tgWebhookPreCheckoutQuery) GetPreCheckoutQueryID() string {
	return q.update.PreCheckoutQuery.ID
}

func (q tgWebhookPreCheckoutQuery) GetCurrency() string {
	return q.update.PreCheckoutQuery.Currency
}

func (q tgWebhookPreCheckoutQuery) GetTotalAmount() int {
	return q.update.PreCheckoutQuery.TotalAmount
}

func (q tgWebhookPreCheckoutQuery) GetInvoicePayload() string {
	return q.update.PreCheckoutQuery.InvoicePayload
}

func (q tgWebhookPreCheckoutQuery) GetFrom() botinput.Sender {
	return tgWebhookUser{tgUser: q.update.ChosenInlineResult.From}
}

func (q tgWebhookPreCheckoutQuery) GetShippingOptionID() string {
	return q.update.PreCheckoutQuery.ShippingOptionID
}

func (q tgWebhookPreCheckoutQuery) GetOrderInfo() botinput.OrderInfo {
	return (*tgOrderInfo)(q.update.PreCheckoutQuery.OrderInfo)
}

func newTgWebhookPreCheckoutQuery(input tgInput) tgWebhookPreCheckoutQuery {
	if input.update.PreCheckoutQuery == nil {
		panic("update.PreCheckoutQuery == nil")
	}
	return tgWebhookPreCheckoutQuery{
		tgInput: input,
	}
}
