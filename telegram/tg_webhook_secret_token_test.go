package telegram

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/bots-go-framework/bots-fw/botsfw"
)

// SEC-4: anyone who knows/guesses a bot's webhook URL (the `?id=` value is a public,
// non-secret bot username, not a credential) could POST forged Telegram updates and be
// treated as any Telegram user, because nothing verified the request actually came from
// Telegram. These tests cover verifyWebhookSecretToken, the guard added to close that gap.

func newWebhookRequest(t *testing.T, headerValue string) *http.Request {
	t.Helper()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://example.com/bot/tg/hook?id=somebot", nil)
	if err != nil {
		t.Fatalf("http.NewRequestWithContext failed: %v", err)
	}
	if headerValue != "" {
		req.Header.Set(TelegramWebhookSecretTokenHeader, headerValue)
	}
	return req
}

func TestVerifyWebhookSecretToken_ValidSecretPasses(t *testing.T) {
	settings := &botsfw.BotSettings{Code: "somebot", WebhookSecretToken: "s3cr3t"}
	req := newWebhookRequest(t, "s3cr3t")
	if err := verifyWebhookSecretToken(context.Background(), req, settings); err != nil {
		t.Errorf("verifyWebhookSecretToken() = %v, want nil for matching secret", err)
	}
}

func TestVerifyWebhookSecretToken_WrongSecretRejected(t *testing.T) {
	settings := &botsfw.BotSettings{Code: "somebot", WebhookSecretToken: "s3cr3t"}
	req := newWebhookRequest(t, "wrong-value")
	err := verifyWebhookSecretToken(context.Background(), req, settings)
	assertAuthFailed(t, err)
}

func TestVerifyWebhookSecretToken_MissingHeaderRejectedWhenSecretConfigured(t *testing.T) {
	settings := &botsfw.BotSettings{Code: "somebot", WebhookSecretToken: "s3cr3t"}
	req := newWebhookRequest(t, "")
	err := verifyWebhookSecretToken(context.Background(), req, settings)
	assertAuthFailed(t, err)
}

func TestVerifyWebhookSecretToken_NoSecretConfigured_CompatAllowsByDefault(t *testing.T) {
	settings := &botsfw.BotSettings{Code: "somebot"} // WebhookSecretToken left empty on purpose
	req := newWebhookRequest(t, "")
	if err := verifyWebhookSecretToken(context.Background(), req, settings); err != nil {
		t.Errorf("verifyWebhookSecretToken() = %v, want nil (backward-compat: no secret configured yet)", err)
	}
}

func TestVerifyWebhookSecretToken_NoSecretConfigured_RequireWebhookSecretRejects(t *testing.T) {
	settings := &botsfw.BotSettings{Code: "somebot", RequireWebhookSecret: true} // no secret set, but required
	req := newWebhookRequest(t, "")
	err := verifyWebhookSecretToken(context.Background(), req, settings)
	assertAuthFailed(t, err)
}

func TestVerifyWebhookSecretToken_NilSettingsDoesNotPanic(t *testing.T) {
	req := newWebhookRequest(t, "")
	if err := verifyWebhookSecretToken(context.Background(), req, nil); err != nil {
		t.Errorf("verifyWebhookSecretToken() = %v, want nil for nil settings (defensive fallback)", err)
	}
}

func assertAuthFailed(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("verifyWebhookSecretToken() = nil, want a botsfw.ErrAuthFailed error")
	}
	var errAuthFailed botsfw.ErrAuthFailed
	if !errors.As(err, &errAuthFailed) {
		t.Errorf("verifyWebhookSecretToken() = %v (%T), want a botsfw.ErrAuthFailed", err, err)
	}
}
