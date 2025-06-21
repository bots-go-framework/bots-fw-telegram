package telegram

import (
	"testing"
)

func Test_botRecordsFieldsSetter_Platform(t *testing.T) {
	actual := tgBotRecordsFieldsSetter{}.Platform()
	if actual != string(PlatformID) {
		t.Errorf("Expected %s, got %s", PlatformID, actual)
	}
}
