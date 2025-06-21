package telegram

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/bots-go-framework/bots-fw/botsfwconst"
)

// Platform is a bots platform descriptor (in this case - for Telegram)
var Platform botsfw.BotPlatform = platform{}

// platform describes Telegram platform
type platform struct {
}

// PlatformID is 'telegram'
const PlatformID botsfwconst.Platform = "telegram"

// ID returns 'telegram'
func (p platform) ID() string {
	return string(PlatformID)
}

// Version returns '2.0'
func (p platform) Version() string {
	return "2.0"
}
