package telegram

import (
	"context"
	"testing"

	"github.com/bots-go-framework/bots-fw/botsfw"
)

// TELEGRAM_WEBHOOK_REQUIRE_SECRET makes a missing webhook secret fail CLOSED fleet-wide:
// when no secret resolves for a bot, the webhook is rejected instead of allowed-with-warning.

// withRequireWebhookSecret overrides the fleet-wide require flag for the duration of a test.
func withRequireWebhookSecret(t *testing.T, require bool) {
	t.Helper()
	prev := resolveRequireWebhookSecret
	resolveRequireWebhookSecret = func() bool { return require }
	t.Cleanup(func() { resolveRequireWebhookSecret = prev })
}

// require=true + no secret resolvable (no per-bot, no fleet) → reject (fail closed).
func TestVerifyWebhookSecretToken_FleetRequire_NoSecret_Rejects(t *testing.T) {
	withFleetWebhookSecret(t, "")
	withRequireWebhookSecret(t, true)
	settings := &botsfw.BotSettings{Code: "somebot"}
	req := newWebhookRequest(t, "")
	assertAuthFailed(t, verifyWebhookSecretToken(context.Background(), req, settings))
}

// require=false + no secret → legacy compat (allow with warning).
func TestVerifyWebhookSecretToken_FleetRequireOff_NoSecret_CompatAllows(t *testing.T) {
	withFleetWebhookSecret(t, "")
	withRequireWebhookSecret(t, false)
	settings := &botsfw.BotSettings{Code: "somebot"}
	req := newWebhookRequest(t, "")
	if err := verifyWebhookSecretToken(context.Background(), req, settings); err != nil {
		t.Errorf("verifyWebhookSecretToken() = %v, want nil (require off, no secret => compat allow)", err)
	}
}

// require=true but a fleet secret IS set → the require branch is not reached; normal strict
// verification applies (matching header passes, missing/wrong header rejects).
func TestVerifyWebhookSecretToken_FleetRequire_WithSecret_NormalStrict(t *testing.T) {
	withFleetWebhookSecret(t, "fleet-secret")
	withRequireWebhookSecret(t, true)
	settings := &botsfw.BotSettings{Code: "somebot"}

	if err := verifyWebhookSecretToken(context.Background(), newWebhookRequest(t, "fleet-secret"), settings); err != nil {
		t.Errorf("matching header should pass, got %v", err)
	}
	assertAuthFailed(t, verifyWebhookSecretToken(context.Background(), newWebhookRequest(t, ""), settings))
}
