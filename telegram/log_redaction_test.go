package telegram

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-store/botsfwstore/botsfwstoretest"
	"github.com/bots-go-framework/bots-fw/botmsg"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/bots-go-framework/bots-fw/mocks/mock_botsfw"
	"github.com/bots-go-framework/bots-go-core/botkb"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"go.uber.org/mock/gomock"
)

var (
	redactionLogsOnce sync.Once
	redactionLogs     = &redactionLogCapture{entries: make(map[string][]string)}
	redactionTraceID  atomic.Uint64
)

type redactionLogCapture struct {
	mu      sync.Mutex
	entries map[string][]string
}

func (c *redactionLogCapture) Log(ctx context.Context, entry logus.LogEntry) error {
	if ctx == nil {
		return nil
	}
	traceID := logus.GetTraceID(ctx)
	if !strings.HasPrefix(traceID, "telegram-redaction-test-") {
		return nil
	}
	message := fmt.Sprintf(entry.MessageFormat, entry.MessageArgs...)
	c.mu.Lock()
	c.entries[traceID] = append(c.entries[traceID], message)
	c.mu.Unlock()
	return nil
}

func captureRedactionLogs(t *testing.T) (context.Context, func() string) {
	t.Helper()
	redactionLogsOnce.Do(func() {
		logus.AddLogEntryHandler(redactionLogs)
	})
	traceID := fmt.Sprintf("telegram-redaction-test-%d", redactionTraceID.Add(1))
	ctx := logus.WithTraceID(context.Background(), traceID)
	t.Cleanup(func() {
		redactionLogs.mu.Lock()
		delete(redactionLogs.entries, traceID)
		redactionLogs.mu.Unlock()
	})
	return ctx, func() string {
		redactionLogs.mu.Lock()
		defer redactionLogs.mu.Unlock()
		return strings.Join(redactionLogs.entries[traceID], "\n")
	}
}

func TestGetBotContextAndInputs_LogsRedactedEnvelope(t *testing.T) {
	const (
		privateText     = "PRIVATE-INBOUND-TEXT"
		privateName     = "PRIVATE-FIRST-NAME"
		privateUsername = "PRIVATE-USERNAME"
		privateQuery    = "PRIVATE-QUERY-VALUE"
	)
	ctx, logs := captureRedactionLogs(t)
	ctrl := gomock.NewController(t)
	provider := mock_botsfw.NewMockBotContextProvider(ctrl)
	provider.EXPECT().
		GetBotContext(ctx, PlatformID, "public-bot").
		Return(&botsfw.BotContext{BotSettings: &botsfw.BotSettings{Code: "public-bot"}}, nil)

	handler := NewTelegramWebhookHandler(
		provider,
		func(context.Context) i18n.Translator { return nil },
		testChatInstanceStore{},
	)
	body := fmt.Sprintf(`{
		"update_id": 42,
		"message": {
			"message_id": 7,
			"from": {"id": 123, "is_bot": false, "first_name": %q, "username": %q},
			"chat": {"id": 123, "type": "private"},
			"text": %q
		}
	}`, privateName, privateUsername, privateText)
	request := httptest.NewRequest(
		http.MethodPost,
		"https://example.test/tg/hook?id=public-bot&private="+privateQuery,
		bytes.NewBufferString(body),
	).WithContext(ctx)

	if _, _, err := handler.GetBotContextAndInputs(ctx, request); err != nil {
		t.Fatalf("GetBotContextAndInputs() error = %v", err)
	}
	logText := logs()
	if !strings.Contains(logText, "outcome=decoded") {
		t.Fatalf("logs do not contain decoded envelope: %s", logText)
	}
	assertLogOmits(t, logText, privateText, privateName, privateUsername, privateQuery, body, "update_id")
}

func TestGetBotContextAndInputs_InvalidJSONDoesNotLogBody(t *testing.T) {
	const privateBody = `{"text":"PRIVATE-MALFORMED-BODY"`
	ctx, logs := captureRedactionLogs(t)
	ctrl := gomock.NewController(t)
	provider := mock_botsfw.NewMockBotContextProvider(ctrl)
	provider.EXPECT().
		GetBotContext(ctx, PlatformID, "public-bot").
		Return(&botsfw.BotContext{BotSettings: &botsfw.BotSettings{Code: "public-bot"}}, nil)

	handler := NewTelegramWebhookHandler(
		provider,
		func(context.Context) i18n.Translator { return nil },
		testChatInstanceStore{},
	)
	request := httptest.NewRequest(
		http.MethodPost,
		"https://example.test/tg/hook?id=public-bot",
		bytes.NewBufferString(privateBody),
	).WithContext(ctx)

	if _, _, err := handler.GetBotContextAndInputs(ctx, request); err == nil {
		t.Fatal("GetBotContextAndInputs() error = nil, want malformed JSON error")
	}
	logText := logs()
	if !strings.Contains(logText, "outcome=decode_error") {
		t.Fatalf("logs do not contain decode_error envelope: %s", logText)
	}
	assertLogOmits(t, logText, privateBody, "PRIVATE-MALFORMED-BODY")
}

func TestSendMessage_LogsRedactedEnvelopeOverResponse(t *testing.T) {
	const (
		privateText     = "PRIVATE-OUTBOUND-TEXT"
		privateButton   = "PRIVATE-BUTTON-TEXT"
		privateCallback = "PRIVATE-CALLBACK-DATA"
	)
	ctx, logs := captureRedactionLogs(t)
	responder, recorder := newRedactionTestResponder(t, ctx, http.DefaultClient)
	keyboard := botkb.NewMessageKeyboard(
		botkb.KeyboardTypeInline,
		[]botkb.Button{botkb.NewDataButton(privateButton, privateCallback)},
	)
	message := botmsg.MessageFromBot{
		TextMessageFromBot: botmsg.TextMessageFromBot{
			Text:     privateText,
			Keyboard: keyboard,
		},
	}

	if _, err := responder.SendMessage(ctx, message, botsfw.BotAPISendMessageOverResponse); err != nil {
		t.Fatalf("SendMessage() error = %v", err)
	}
	responseBody := recorder.Body.String()
	for _, expected := range []string{privateText, privateButton, privateCallback} {
		if !strings.Contains(responseBody, expected) {
			t.Fatalf("Telegram response does not contain %q: %s", expected, responseBody)
		}
	}
	logText := logs()
	if !strings.Contains(logText, "Sending Telegram response") {
		t.Fatalf("logs do not contain outbound envelope: %s", logText)
	}
	assertLogOmits(t, logText, privateText, privateButton, privateCallback)
}

func TestSendMessage_HTTPSDoesNotEnableRawTelegramDebug(t *testing.T) {
	const (
		privateText     = "PRIVATE-HTTPS-TEXT"
		privateResponse = "PRIVATE-TELEGRAM-RESPONSE"
	)
	ctx, logs := captureRedactionLogs(t)
	client := &http.Client{Transport: redactionRoundTripFunc(func(*http.Request) (*http.Response, error) {
		body := fmt.Sprintf(
			`{"ok":true,"result":{"message_id":9,"chat":{"id":123,"type":"private"},"text":%q}}`,
			privateResponse,
		)
		return &http.Response{
			StatusCode:    http.StatusOK,
			Body:          io.NopCloser(strings.NewReader(body)),
			ContentLength: int64(len(body)),
			Header:        make(http.Header),
		}, nil
	})}
	responder, _ := newRedactionTestResponder(t, ctx, client)
	message := botmsg.MessageFromBot{
		TextMessageFromBot: botmsg.TextMessageFromBot{Text: privateText},
	}

	if _, err := responder.SendMessage(ctx, message, botsfw.BotAPISendMessageOverHTTPS); err != nil {
		t.Fatalf("SendMessage() error = %v", err)
	}
	assertLogOmits(
		t,
		logs(),
		privateText,
		privateResponse,
		"123456:PRIVATE-BOT-TOKEN",
		"Request to Telegram API",
	)
}

func newRedactionTestResponder(
	t *testing.T,
	ctx context.Context,
	client *http.Client,
) (tgWebhookResponder, *httptest.ResponseRecorder) {
	t.Helper()
	update := &tgbotapi.Update{
		UpdateID: 1,
		Message: &tgbotapi.Message{
			MessageID: 1,
			From:      &tgbotapi.User{ID: 123, FirstName: "PRIVATE-SENDER"},
			Chat:      &tgbotapi.Chat{ID: 123, Type: "private"},
			Text:      "PRIVATE-INBOUND",
		},
	}
	input, err := NewTelegramWebhookInput(update, nil)
	if err != nil {
		t.Fatalf("NewTelegramWebhookInput() error = %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/tg/hook?id=public-bot", nil).WithContext(ctx)
	host := redactionTestBotHost{client: client}
	appContext := redactionTestAppContext{}
	settings := &botsfw.BotSettings{
		Code:  "public-bot",
		Token: "123456:PRIVATE-BOT-TOKEN",
		Store: &botsfwstoretest.FakeStateStore{},
	}
	webhookContext, err := newTelegramWebhookContext(
		botsfw.NewCreateWebhookContextArgs(
			request,
			appContext,
			botsfw.BotContext{AppContext: appContext, BotHost: host, BotSettings: settings},
			input,
			settings.Store,
		),
		input.(TgWebhookInput),
		tgBotRecordsFieldsSetter{},
		testChatInstanceStore{},
	)
	if err != nil {
		t.Fatalf("newTelegramWebhookContext() error = %v", err)
	}
	recorder := httptest.NewRecorder()
	return newTgWebhookResponder(recorder, webhookContext), recorder
}

type redactionTestBotHost struct {
	client *http.Client
}

func (redactionTestBotHost) Context(r *http.Request) context.Context {
	return r.Context()
}

func (h redactionTestBotHost) GetHTTPClient(context.Context) *http.Client {
	return h.client
}

type redactionTestAppContext struct{}

func (redactionTestAppContext) SupportedLocales() []i18n.Locale {
	return nil
}

func (redactionTestAppContext) GetLocaleByCode5(string) (i18n.Locale, error) {
	return i18n.Locale{}, nil
}

func (redactionTestAppContext) GetTranslator(context.Context) i18n.Translator {
	return nil
}

type redactionRoundTripFunc func(*http.Request) (*http.Response, error)

func (f redactionRoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func assertLogOmits(t *testing.T, logs string, privateValues ...string) {
	t.Helper()
	for _, privateValue := range privateValues {
		if strings.Contains(logs, privateValue) {
			t.Fatalf("logs expose private value %q: %s", privateValue, logs)
		}
	}
}
