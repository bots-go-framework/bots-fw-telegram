package store

import "testing"

func TestNewTgChat(t *testing.T) {
	t.Run("should panic if data is nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("should panic")
			}
		}()
		NewTgChat("id", nil)
	})
	t.Run("should not panic if data is not nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatal("should not panic")
			}
		}()
		NewTgChat("id", &TgChatBase{})
	})
}
