package telegram

import (
	"context"
	"net/http"
	"testing"

	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/bots-go-framework/bots-fw/mocks/mock_botsfw"
	"github.com/strongo/i18n"
	"go.uber.org/mock/gomock"
)

// fakeWebhookDriver is a no-op botsfw.WebhookDriver — RegisterHttpHandlers only needs
// a non-nil driver to pass Register()'s guard; it never calls back into it here.
type fakeWebhookDriver struct{}

func (fakeWebhookDriver) RegisterWebhookHandlers(botsfw.HttpRouter, string, ...botsfw.WebhookHandler) {
}
func (fakeWebhookDriver) HandleWebhook(http.ResponseWriter, *http.Request, botsfw.WebhookHandler) {}

// recordingRouter records the (method, path) of every route registered on it.
type recordingRouter struct{ routes []string }

func (rr *recordingRouter) Handle(method, path string, _ http.HandlerFunc) {
	rr.routes = append(rr.routes, method+" "+path)
}

func (rr *recordingRouter) has(route string) bool {
	for _, r := range rr.routes {
		if r == route {
			return true
		}
	}
	return false
}

func registerAndRecord(t *testing.T, opts ...TgWebhookHandlerOption) *recordingRouter {
	t.Helper()
	ctrl := gomock.NewController(t)
	var tp botsfw.TranslatorProvider = func(context.Context) i18n.Translator { return nil }
	h := NewTelegramWebhookHandler(
		mock_botsfw.NewMockBotContextProvider(ctrl), tp, testChatInstanceStore{}, opts...)
	rr := &recordingRouter{}
	h.RegisterHttpHandlers(fakeWebhookDriver{}, mock_botsfw.NewMockBotHost(ctrl), rr, "/bot")
	return rr
}

// WithAdminPathPrefix moves the admin-only endpoints (set-webhook, test) under the
// admin prefix, while the public /tg/hook endpoint stays under the public prefix.
func TestRegisterHttpHandlers_WithAdminPathPrefix(t *testing.T) {
	rr := registerAndRecord(t, WithAdminPathPrefix("/admin/bot"))

	for _, want := range []string{
		"POST /bot/tg/hook",               // public — Telegram POSTs here
		"GET /admin/bot/tg/set-webhook",   // admin — moved under the admin prefix
		"GET /admin/bot/tg/test/time-now", // admin — moved too
	} {
		if !rr.has(want) {
			t.Errorf("route %q not registered; got %v", want, rr.routes)
		}
	}
	if rr.has("GET /bot/tg/set-webhook") {
		t.Errorf("set-webhook must NOT be under the public prefix when an admin prefix is set; got %v", rr.routes)
	}
}

// Without an admin prefix, every endpoint stays under the public prefix (legacy behaviour).
func TestRegisterHttpHandlers_NoAdminPathPrefix_LegacyLayout(t *testing.T) {
	rr := registerAndRecord(t)

	for _, want := range []string{
		"POST /bot/tg/hook",
		"GET /bot/tg/set-webhook",
		"GET /bot/tg/test/time-now",
	} {
		if !rr.has(want) {
			t.Errorf("route %q not registered; got %v", want, rr.routes)
		}
	}
}
