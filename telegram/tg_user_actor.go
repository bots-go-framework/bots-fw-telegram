package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
)

var _ botinput.Actor = tgUserActor{}

// tgUserActor wraps tgbotapi.User to implement botinput.Actor
type tgUserActor struct {
	tgbotapi.User
}

func (u tgUserActor) IsBotUser() bool {
	return u.IsBot
}
