package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
)

var _ botinput.OrderInfo = (*tgOrderInfo)(nil)

type tgOrderInfo tgbotapi.OrderInfo

func (t *tgOrderInfo) GetUserName() string {
	return t.Name
}

func (t *tgOrderInfo) GetPhoneNumber() string {
	return t.PhoneNumber
}

func (t *tgOrderInfo) GetEmailAddress() string {
	return t.Email
}

func (t *tgOrderInfo) GetShippingAddress() botinput.ShippingAddress {
	return (*tgShippingAddress)(t.ShippingAddress)
}

var _ botinput.ShippingAddress = (*tgShippingAddress)(nil)

type tgShippingAddress tgbotapi.ShippingAddress

func (t *tgShippingAddress) GetCountryCode() string {
	return t.CountryCode
}

func (t *tgShippingAddress) GetState() string {
	return t.State
}

func (t *tgShippingAddress) GetCity() string {
	return t.City
}

func (t *tgShippingAddress) GetStreetLine1() string {
	return t.StreetLine1
}

func (t *tgShippingAddress) GetStreetLine2() string {
	return t.StreetLine2
}

func (t *tgShippingAddress) GetPostCode() string {
	return t.PostCode
}
