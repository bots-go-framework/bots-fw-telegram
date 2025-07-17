package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

var _ botinput.LocationMessage = (*locationInputMessage)(nil)

var _ botinput.LocationMessage = (*locationInputMessage)(nil)
var _ botsfw.InputMessage = (*locationInputMessage)(nil)

type locationInputMessage struct {
	tgInputMessage
}

func (v locationInputMessage) Text() string {
	return ""
}

func (v locationInputMessage) GetLatitude() float64 {
	return v.update.Message.Location.Latitude
}

func (v locationInputMessage) GetLongitude() float64 {
	return v.update.Message.Location.Longitude
}

func newLocationMessage(input tgInput, tgMessage *tgbotapi.Message) locationInputMessage {
	return locationInputMessage{
		tgInputMessage: newTelegramWebhookMessage(input, tgMessage),
	}
}
