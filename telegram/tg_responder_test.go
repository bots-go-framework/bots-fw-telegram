package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-go-core/botkb"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetInlineKeyboard(t *testing.T) {
	t.Run("EmptyKeyboard", func(t *testing.T) {
		kb := botkb.NewMessageKeyboard(botkb.KeyboardTypeInline)
		inlineKb := getInlineKeyboard(kb)

		if len(inlineKb.InlineKeyboard) != 0 {
			t.Errorf("Expected empty keyboard, got %v rows", len(inlineKb.InlineKeyboard))
		}
	})

	t.Run("SingleDataButton", func(t *testing.T) {
		kb := botkb.NewMessageKeyboard(
			botkb.KeyboardTypeInline,
			[]botkb.Button{botkb.NewDataButton("Test Button", "test-data")},
		)
		inlineKb := getInlineKeyboard(kb)

		if len(inlineKb.InlineKeyboard) != 1 {
			t.Errorf("Expected 1 row, got %v", len(inlineKb.InlineKeyboard))
		}

		if len(inlineKb.InlineKeyboard[0]) != 1 {
			t.Errorf("Expected 1 button in row, got %v", len(inlineKb.InlineKeyboard[0]))
		}

		button := inlineKb.InlineKeyboard[0][0]
		if button.Text != "Test Button" {
			t.Errorf("Expected button text 'Test Button', got '%v'", button.Text)
		}

		if button.CallbackData != "test-data" {
			t.Errorf("Expected button data 'test-data', got '%v'", button.CallbackData)
		}
	})

	t.Run("SingleUrlButton", func(t *testing.T) {
		kb := botkb.NewMessageKeyboard(
			botkb.KeyboardTypeInline,
			[]botkb.Button{botkb.NewUrlButton("URL Button", "https://example.com")},
		)
		inlineKb := getInlineKeyboard(kb)

		if len(inlineKb.InlineKeyboard) != 1 {
			t.Errorf("Expected 1 row, got %v", len(inlineKb.InlineKeyboard))
		}

		if len(inlineKb.InlineKeyboard[0]) != 1 {
			t.Errorf("Expected 1 button in row, got %v", len(inlineKb.InlineKeyboard[0]))
		}

		button := inlineKb.InlineKeyboard[0][0]
		if button.Text != "URL Button" {
			t.Errorf("Expected button text 'URL Button', got '%v'", button.Text)
		}

		if button.URL != "https://example.com" {
			t.Errorf("Expected button URL 'https://example.com', got '%v'", button.URL)
		}
	})

	t.Run("SingleSwitchInlineQueryButton", func(t *testing.T) {
		kb := botkb.NewMessageKeyboard(
			botkb.KeyboardTypeInline,
			[]botkb.Button{botkb.NewSwitchInlineQueryButton("Switch Query", "query")},
		)
		inlineKb := getInlineKeyboard(kb)

		if len(inlineKb.InlineKeyboard) != 1 {
			t.Errorf("Expected 1 row, got %v", len(inlineKb.InlineKeyboard))
		}

		if len(inlineKb.InlineKeyboard[0]) != 1 {
			t.Errorf("Expected 1 button in row, got %v", len(inlineKb.InlineKeyboard[0]))
		}

		button := inlineKb.InlineKeyboard[0][0]
		if button.Text != "Switch Query" {
			t.Errorf("Expected button text 'Switch Query', got '%v'", button.Text)
		}

		if button.SwitchInlineQuery == nil || *button.SwitchInlineQuery != "query" {
			t.Errorf("Expected button switch inline query 'query', got '%v'", *button.SwitchInlineQuery)
		}
	})

	t.Run("SingleSwitchInlineQueryCurrentChatButton", func(t *testing.T) {
		kb := botkb.NewMessageKeyboard(
			botkb.KeyboardTypeInline,
			[]botkb.Button{botkb.NewSwitchInlineQueryCurrentChatButton("Switch Current", "current-query")},
		)
		inlineKb := getInlineKeyboard(kb)

		if len(inlineKb.InlineKeyboard) != 1 {
			t.Errorf("Expected 1 row, got %v", len(inlineKb.InlineKeyboard))
		}

		if len(inlineKb.InlineKeyboard[0]) != 1 {
			t.Errorf("Expected 1 button in row, got %v", len(inlineKb.InlineKeyboard[0]))
		}

		button := inlineKb.InlineKeyboard[0][0]
		if button.Text != "Switch Current" {
			t.Errorf("Expected button text 'Switch Current', got '%v'", button.Text)
		}

		if button.SwitchInlineQueryCurrentChat == nil || *button.SwitchInlineQueryCurrentChat != "current-query" {
			t.Errorf("Expected button switch inline query current chat 'current-query', got '%v'", *button.SwitchInlineQueryCurrentChat)
		}
	})

	t.Run("MultipleButtonsAndRows", func(t *testing.T) {
		kb := botkb.NewMessageKeyboard(
			botkb.KeyboardTypeInline,
			[]botkb.Button{botkb.NewDataButton("Button 1", "data-1"), botkb.NewUrlButton("Button 2", "https://example.com")},
			[]botkb.Button{botkb.NewSwitchInlineQueryButton("Button 3", "query")},
		)
		inlineKb := getInlineKeyboard(kb)

		if len(inlineKb.InlineKeyboard) != 2 {
			t.Errorf("Expected 2 rows, got %v", len(inlineKb.InlineKeyboard))
		}

		if len(inlineKb.InlineKeyboard[0]) != 2 {
			t.Errorf("Expected 2 buttons in first row, got %v", len(inlineKb.InlineKeyboard[0]))
		}

		if len(inlineKb.InlineKeyboard[1]) != 1 {
			t.Errorf("Expected 1 button in second row, got %v", len(inlineKb.InlineKeyboard[1]))
		}

		button1 := inlineKb.InlineKeyboard[0][0]
		if button1.Text != "Button 1" || button1.CallbackData != "data-1" {
			t.Errorf("Button 1 has incorrect properties: %+v", button1)
		}

		button2 := inlineKb.InlineKeyboard[0][1]
		if button2.Text != "Button 2" || button2.URL != "https://example.com" {
			t.Errorf("Button 2 has incorrect properties: %+v", button2)
		}

		button3 := inlineKb.InlineKeyboard[1][0]
		if button3.Text != "Button 3" || *button3.SwitchInlineQuery != "query" {
			t.Errorf("Button 3 has incorrect properties: %+v", button3)
		}
	})
}

func TestGetReplyKeyboard(t *testing.T) {
	t.Run("EmptyKeyboard", func(t *testing.T) {
		kb := botkb.NewMessageKeyboard(botkb.KeyboardTypeBottom)
		replyKb := getReplyKeyboard(kb)

		if len(replyKb.Keyboard) != 0 {
			t.Errorf("Expected empty keyboard, got %v rows", len(replyKb.Keyboard))
		}
	})

	t.Run("SingleTextButton", func(t *testing.T) {
		kb := botkb.NewMessageKeyboard(
			botkb.KeyboardTypeBottom,
			[]botkb.Button{&botkb.DataButton{Text: "Test Button"}},
		)
		replyKb := getReplyKeyboard(kb)

		if len(replyKb.Keyboard) != 1 {
			t.Errorf("Expected 1 row, got %v", len(replyKb.Keyboard))
		}

		if len(replyKb.Keyboard[0]) != 1 {
			t.Errorf("Expected 1 button in row, got %v", len(replyKb.Keyboard[0]))
		}

		button := replyKb.Keyboard[0][0]
		if button.Text != "Test Button" {
			t.Errorf("Expected button text 'Test Button', got '%v'", button.Text)
		}
	})

	t.Run("MultipleButtonsAndRows", func(t *testing.T) {
		kb := botkb.NewMessageKeyboard(
			botkb.KeyboardTypeBottom,
			[]botkb.Button{&botkb.DataButton{Text: "Button 1"}, &botkb.DataButton{Text: "Button 2"}},
			[]botkb.Button{&botkb.DataButton{Text: "Button 3"}},
		)
		replyKb := getReplyKeyboard(kb)

		if len(replyKb.Keyboard) != 2 {
			t.Errorf("Expected 2 rows, got %v", len(replyKb.Keyboard))
		}

		if len(replyKb.Keyboard[0]) != 2 {
			t.Errorf("Expected 2 buttons in first row, got %v", len(replyKb.Keyboard[0]))
		}

		if len(replyKb.Keyboard[1]) != 1 {
			t.Errorf("Expected 1 button in second row, got %v", len(replyKb.Keyboard[1]))
		}

		button1 := replyKb.Keyboard[0][0]
		if button1.Text != "Button 1" {
			t.Errorf("Button 1 has incorrect text: %v", button1.Text)
		}

		button2 := replyKb.Keyboard[0][1]
		if button2.Text != "Button 2" {
			t.Errorf("Button 2 has incorrect text: %v", button2.Text)
		}

		button3 := replyKb.Keyboard[1][0]
		if button3.Text != "Button 3" {
			t.Errorf("Button 3 has incorrect text: %v", button3.Text)
		}
	})
}

func TestGetTelegramKeyboard(t *testing.T) {
	t.Run("DirectTelegramKeyboard", func(t *testing.T) {
		// Test when a Telegram keyboard is passed directly
		tgKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Test", "test"),
			),
		)

		result := getTelegramKeyboard(tgKeyboard)

		// Should return the same keyboard
		if result != tgKeyboard {
			t.Errorf("Expected the same keyboard to be returned")
		}
	})

	t.Run("InlineKeyboard", func(t *testing.T) {
		// Test with a botkb.MessageKeyboard with KeyboardTypeInline
		kb := botkb.NewMessageKeyboard(
			botkb.KeyboardTypeInline,
			[]botkb.Button{botkb.NewDataButton("Test Button", "test-data")},
		)

		result := getTelegramKeyboard(kb)

		// Check that it's an InlineKeyboardMarkup
		inlineKb, ok := result.(*tgbotapi.InlineKeyboardMarkup)
		if !ok {
			t.Fatalf("Expected InlineKeyboardMarkup, got %T", result)
		}

		if len(inlineKb.InlineKeyboard) != 1 || len(inlineKb.InlineKeyboard[0]) != 1 {
			t.Errorf("Keyboard has incorrect structure")
		}

		button := inlineKb.InlineKeyboard[0][0]
		if button.Text != "Test Button" || button.CallbackData != "test-data" {
			t.Errorf("Button has incorrect properties: %+v", button)
		}
	})

	t.Run("ReplyKeyboard", func(t *testing.T) {
		// Test with a botkb.MessageKeyboard with KeyboardTypeBottom
		kb := botkb.NewMessageKeyboard(
			botkb.KeyboardTypeBottom,
			[]botkb.Button{&botkb.DataButton{Text: "Test Button"}},
		)

		result := getTelegramKeyboard(kb)

		// Check that it's a ReplyKeyboardMarkup
		replyKb, ok := result.(*tgbotapi.ReplyKeyboardMarkup)
		if !ok {
			t.Fatalf("Expected ReplyKeyboardMarkup, got %T", result)
		}

		if len(replyKb.Keyboard) != 1 || len(replyKb.Keyboard[0]) != 1 {
			t.Errorf("Keyboard has incorrect structure")
		}

		button := replyKb.Keyboard[0][0]
		if button.Text != "Test Button" {
			t.Errorf("Button has incorrect text: %v", button.Text)
		}
	})

	t.Run("UnsupportedKeyboardType", func(t *testing.T) {
		// Test with an unsupported keyboard type
		kb := botkb.NewMessageKeyboard(botkb.KeyboardTypeForceReply)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for unsupported keyboard type")
			}
		}()

		getTelegramKeyboard(kb)
	})

	t.Run("UnsupportedKeyboardImplementation", func(t *testing.T) {
		// Test with an unsupported keyboard implementation
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for unsupported keyboard implementation")
			}
		}()

		// Create a mock keyboard that's not a MessageKeyboard
		mockKeyboard := &mockKeyboard{}
		getTelegramKeyboard(mockKeyboard)
	})
}

func TestGetHideKeyboard(t *testing.T) {
	t.Run("BasicTest", func(t *testing.T) {
		// Create a keyboard with KeyboardTypeHide
		kb := botkb.NewMessageKeyboard(botkb.KeyboardTypeHide)

		// Call the function being tested
		hideKb := getHideKeyboard(kb)

		// Verify the result is a ReplyKeyboardHide with HideKeyboard set to true
		if hideKb == nil {
			t.Fatal("Expected non-nil ReplyKeyboardHide")
		}

		if !hideKb.HideKeyboard {
			t.Errorf("Expected HideKeyboard to be true, got false")
		}
	})

	t.Run("IgnoresButtons", func(t *testing.T) {
		// Create a keyboard with KeyboardTypeHide and some buttons (which should be ignored)
		kb := botkb.NewMessageKeyboard(
			botkb.KeyboardTypeHide,
			[]botkb.Button{botkb.NewDataButton("Button 1", "data-1")},
		)

		// Call the function being tested
		hideKb := getHideKeyboard(kb)

		// Verify the result is a ReplyKeyboardHide with HideKeyboard set to true
		if hideKb == nil {
			t.Fatal("Expected non-nil ReplyKeyboardHide")
		}

		if !hideKb.HideKeyboard {
			t.Errorf("Expected HideKeyboard to be true, got false")
		}
	})
}

// Mock keyboard for testing
type mockKeyboard struct{}

func (m *mockKeyboard) KeyboardType() botkb.KeyboardType {
	return botkb.KeyboardTypeInline
}

func TestNewTgWebhookResponder(t *testing.T) {
	// Create a mock http.ResponseWriter
	w := httptest.NewRecorder()

	// Create a mock tgWebhookContext
	whc := &tgWebhookContext{
		responder: nil, // Initially nil, should be set by newTgWebhookResponder
	}

	// Call the function being tested
	responder := newTgWebhookResponder(w, whc)

	// Verify that the responder has the correct fields
	if responder.w != w {
		t.Errorf("Expected responder.w to be %v, got %v", w, responder.w)
	}

	if responder.whc != whc {
		t.Errorf("Expected responder.whc to be %v, got %v", whc, responder.whc)
	}

	// Verify that the responder was set in the webhook context
	if !reflect.DeepEqual(whc.responder, responder) {
		t.Errorf("Expected whc.responder to be %v, got %v", responder, whc.responder)
	}
}
