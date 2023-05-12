package telegram

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_botRecordsFieldsSetter_Platform(t *testing.T) {
	assert.Equal(t, PlatformID, tgBotRecordsFieldsSetter{}.Platform())
}
