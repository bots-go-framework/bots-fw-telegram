package telegram

import (
	"testing"
)

func Test_botRecordsFieldsSetter_Platform(t *testing.T) {
	actual := tgBotRecordsFieldsSetter{}.Platform()
	expected := PlatformID
	if actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
