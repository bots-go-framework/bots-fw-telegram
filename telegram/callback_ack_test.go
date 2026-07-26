package telegram

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botmsg"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

type callbackAckResponder struct {
	messages []botmsg.MessageFromBot
	err      error
}

type concurrentAckResponder struct {
	count atomic.Int32
}

func (r *concurrentAckResponder) SendMessage(_ context.Context, _ botmsg.MessageFromBot, _ botmsg.BotAPISendMessageChannel) (botsfw.OnMessageSentResponse, error) {
	r.count.Add(1)
	return botsfw.OnMessageSentResponse{}, nil
}

func (*concurrentAckResponder) DeleteMessage(context.Context, string) error { return nil }

func (r *callbackAckResponder) SendMessage(_ context.Context, m botmsg.MessageFromBot, channel botmsg.BotAPISendMessageChannel) (botsfw.OnMessageSentResponse, error) {
	if channel != botsfw.BotAPISendMessageOverHTTPS {
		return botsfw.OnMessageSentResponse{}, errors.New("callback acknowledgement must use HTTPS")
	}
	r.messages = append(r.messages, m)
	return botsfw.OnMessageSentResponse{}, r.err
}

func (*callbackAckResponder) DeleteMessage(context.Context, string) error {
	return nil
}

func TestAcknowledgeCallbackQueryImmediatelyAndOnlyOnce(t *testing.T) {
	responder := &callbackAckResponder{}
	update := &tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{ID: "callback-123"},
	}
	whc := &tgWebhookContext{
		tgInput:   tgInput{update: update},
		responder: responder,
	}

	if err := whc.AcknowledgeCallbackQuery("Thinking…", false); err != nil {
		t.Fatal(err)
	}
	if !whc.WasCallbackQueryAcknowledged() {
		t.Fatal("acknowledgement marker was not set")
	}
	if len(responder.messages) != 1 {
		t.Fatalf("sent %d acknowledgements, want 1", len(responder.messages))
	}
	answer, ok := responder.messages[0].BotMessage.(botmsg.AnswerCallbackQuery)
	if !ok {
		t.Fatalf("BotMessage = %T, want botmsg.AnswerCallbackQuery", responder.messages[0].BotMessage)
	}
	if answer.CallbackQueryID != "callback-123" {
		t.Errorf("CallbackQueryID = %q", answer.CallbackQueryID)
	}
	if answer.Text != "Thinking…" {
		t.Errorf("Text = %q", answer.Text)
	}

	if err := whc.AcknowledgeCallbackQuery("", false); err != nil {
		t.Fatal(err)
	}
	if len(responder.messages) != 1 {
		t.Fatalf("duplicate call sent %d acknowledgements, want 1", len(responder.messages))
	}
}

func TestAcknowledgeCallbackQueryMarksOnlySuccessfulSend(t *testing.T) {
	responder := &callbackAckResponder{err: errors.New("network down")}
	update := &tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{ID: "callback-123"},
	}
	whc := &tgWebhookContext{
		tgInput:   tgInput{update: update},
		responder: responder,
	}

	if err := whc.AcknowledgeCallbackQuery("", false); err == nil {
		t.Fatal("expected error")
	}
	if whc.WasCallbackQueryAcknowledged() {
		t.Fatal("failed send must not set acknowledgement marker")
	}
}

func TestAcknowledgeCallbackQueryConcurrentCallsSendOnce(t *testing.T) {
	responder := &concurrentAckResponder{}
	whc := &tgWebhookContext{
		tgInput: tgInput{update: &tgbotapi.Update{
			CallbackQuery: &tgbotapi.CallbackQuery{ID: "callback-123"},
		}},
		responder: responder,
	}
	var waitGroup sync.WaitGroup
	for range 20 {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			if err := whc.AcknowledgeCallbackQuery("", false); err != nil {
				t.Errorf("AcknowledgeCallbackQuery: %v", err)
			}
		}()
	}
	waitGroup.Wait()
	if got := responder.count.Load(); got != 1 {
		t.Fatalf("sent %d acknowledgements, want 1", got)
	}
}
