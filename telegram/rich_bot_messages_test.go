package telegram

import (
	"encoding/json"
	"testing"

	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botmsg"
	"github.com/bots-go-framework/bots-go-core/botkb"
)

func nativeTableRichMessage() tgbotapi.InputRichMessage {
	heading := tgbotapi.RichText{PlainText: "Player"}
	value := tgbotapi.RichText{PlainText: "Score"}
	return tgbotapi.InputRichMessage{
		Blocks: []tgbotapi.InputRichBlock{{
			Type:       tgbotapi.RichBlockTypeTable,
			IsBordered: true,
			IsStriped:  true,
			Cells: [][]tgbotapi.RichBlockTableCell{{
				{Text: &heading, IsHeader: true},
				{Text: &value, IsHeader: true},
			}},
		}},
	}
}

func TestSendRichMessageCarriesNativeTableAndKeyboard(t *testing.T) {
	message := NewSendRichMessage(42, nativeTableRichMessage())
	if got := message.BotMessageType(); got != botmsg.TypeSendRichMessage {
		t.Fatalf("BotMessageType() = %v", got)
	}

	config := tgbotapi.RichMessageConfig(message)
	config.ReplyMarkup = getTelegramKeyboard(botkb.NewMessageKeyboard(
		botkb.KeyboardTypeInline,
		[]botkb.Button{botkb.NewDataButton("Play", "pref:play")},
	))
	values, err := config.Values()
	if err != nil {
		t.Fatal(err)
	}
	if got := values.Get("chat_id"); got != "42" {
		t.Errorf("chat_id = %q, want 42", got)
	}

	var rich map[string]any
	if err = json.Unmarshal([]byte(values.Get("rich_message")), &rich); err != nil {
		t.Fatal(err)
	}
	blocks, ok := rich["blocks"].([]any)
	if !ok || len(blocks) != 1 {
		t.Fatalf("rich blocks = %#v", rich["blocks"])
	}
	table, ok := blocks[0].(map[string]any)
	if !ok || table["type"] != tgbotapi.RichBlockTypeTable {
		t.Fatalf("first block = %#v", blocks[0])
	}
	if values.Get("reply_markup") == "" {
		t.Fatal("reply_markup was not serialized")
	}
}

func TestEditRichMessageTargetsArbitraryPlayerCard(t *testing.T) {
	message := NewEditRichMessage(-100123, 77, nativeTableRichMessage())
	if got := message.BotMessageType(); got != botmsg.TypeEditRichMessage {
		t.Fatalf("BotMessageType() = %v", got)
	}

	values, err := tgbotapi.EditMessageTextConfig(message).Values()
	if err != nil {
		t.Fatal(err)
	}
	if got := values.Get("chat_id"); got != "-100123" {
		t.Errorf("chat_id = %q", got)
	}
	if got := values.Get("message_id"); got != "77" {
		t.Errorf("message_id = %q", got)
	}
	if values.Get("rich_message") == "" {
		t.Fatal("rich_message was not serialized")
	}
}

func TestSendRichMessageDraftType(t *testing.T) {
	message := NewSendRichMessageDraft(42, 9, tgbotapi.InputRichMessage{
		Blocks: []tgbotapi.InputRichBlock{{
			Type: tgbotapi.RichBlockTypeThinking,
			Text: &tgbotapi.RichText{PlainText: "Choosing a legal move…"},
		}},
	})
	if got := message.BotMessageType(); got != botmsg.TypeSendRichMessageDraft {
		t.Fatalf("BotMessageType() = %v", got)
	}
	values, err := tgbotapi.RichMessageDraftConfig(message).Values()
	if err != nil {
		t.Fatal(err)
	}
	if got := values.Get("draft_id"); got != "9" {
		t.Errorf("draft_id = %q", got)
	}
}

func TestTelegramBooleanRequestMakesEphemeralMethodsHandlerReachable(t *testing.T) {
	request := NewTelegramBooleanRequest(tgbotapi.NewDeleteEphemeralMessage(-1001, 42, 7))
	if request.ReturnsMessage {
		t.Fatal("ephemeral deletion must use a boolean result")
	}
	if got := request.TelegramMethod(); got != "deleteEphemeralMessage" {
		t.Errorf("TelegramMethod() = %q", got)
	}
	values, err := request.Values()
	if err != nil {
		t.Fatal(err)
	}
	if values.Get("receiver_user_id") != "42" || values.Get("ephemeral_message_id") != "7" {
		t.Errorf("values = %v", values)
	}

	// The generic wrapper remains a normal BotMessage and can therefore travel
	// through MessageFromBot.BotMessage from any handler.
	var _ botmsg.BotMessage = request
}
