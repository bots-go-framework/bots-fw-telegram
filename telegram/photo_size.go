package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
)

var _ botinput.PhotoMessageItem = (*photoSize)(nil)

type photoSize struct {
	tgbotapi.PhotoSize
}

func (v photoSize) GetFileID() string {
	return v.FileID
}

func (v photoSize) GetWidth() int {
	return v.Width
}

func (v photoSize) GetHeight() int {
	return v.Height
}

func (v photoSize) GetFileSize() int {
	return v.FileSize
}
