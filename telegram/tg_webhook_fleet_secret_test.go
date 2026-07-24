package telegram

import (
	"context"
	"testing"

	"github.com/bots-go-framework/bots-fw/botsfw"
)

// SEC-4 fleet-wide secret: a single TELEGRAM_WEBHOOK_SECRET (see EnvTelegramWebhookSecret)
// authenticates every bot that sets no per-bot WebhookSecretToken, so one provisioned secret
// secures the whole fleet's webhooks with no per-bot configuration. These tests cover
// effectiveWebhookSecret (the resolver) and its effect on verifyWebhookSecretToken.

// withFleetWebhookSecret overrides the fleet-wide TELEGRAM_WEBHOOK_SECRET resolver for the
// duration of a test, without touching the real process environment.
func withFleetWebhookSecret(t *testing.T, secret string) {
	t.Helper()
	prev := resolveFleetWebhookSecret
	resolveFleetWebhookSecret = func() string { return secret }
	t.Cleanup(func() { resolveFleetWebhookSecret = prev })
}

func TestEffectiveWebhookSecret_PerBotOverridesFleet(t *testing.T) {
	withFleetWebhookSecret(t, "fleet-secret")
	settings := &botsfw.BotSettings{Code: "somebot", WebhookSecretToken: "per-bot-secret"}
	if got := effectiveWebhookSecret(settings); got != "per-bot-secret" {
		t.Errorf("effectiveWebhookSecret() = %q, want the per-bot secret to win over the fleet secret", got)
	}
}

func TestEffectiveWebhookSecret_FleetUsedWhenNoPerBot(t *testing.T) {
	withFleetWebhookSecret(t, "fleet-secret")
	settings := &botsfw.BotSettings{Code: "somebot"}
	if got := effectiveWebhookSecret(settings); got != "fleet-secret" {
		t.Errorf("effectiveWebhookSecret() = %q, want the fleet secret as fallback", got)
	}
}

func TestEffectiveWebhookSecret_EmptyWhenNeitherSet(t *testing.T) {
	withFleetWebhookSecret(t, "")
	settings := &botsfw.BotSettings{Code: "somebot"}
	if got := effectiveWebhookSecret(settings); got != "" {
		t.Errorf("effectiveWebhookSecret() = %q, want empty when neither per-bot nor fleet secret is set", got)
	}
}

func TestEffectiveWebhookSecret_NilSettingsFallsBackToFleet(t *testing.T) {
	withFleetWebhookSecret(t, "fleet-secret")
	if got := effectiveWebhookSecret(nil); got != "fleet-secret" {
		t.Errorf("effectiveWebhookSecret(nil) = %q, want the fleet secret", got)
	}
}

// Fleet secret authenticates a bot that sets no per-bot secret: a matching header passes;
// a wrong or missing header is rejected. This is the fleet-wide "enforce" posture that takes
// effect automatically once TELEGRAM_WEBHOOK_SECRET is provisioned and webhooks re-registered.
func TestVerifyWebhookSecretToken_FleetSecretMatchingHeaderPasses(t *testing.T) {
	withFleetWebhookSecret(t, "fleet-secret")
	settings := &botsfw.BotSettings{Code: "somebot"} // no per-bot secret on purpose
	req := newWebhookRequest(t, "fleet-secret")
	if err := verifyWebhookSecretToken(context.Background(), req, settings); err != nil {
		t.Errorf("verifyWebhookSecretToken() = %v, want nil for a header matching the fleet secret", err)
	}
}

func TestVerifyWebhookSecretToken_FleetSecretWrongHeaderRejected(t *testing.T) {
	withFleetWebhookSecret(t, "fleet-secret")
	settings := &botsfw.BotSettings{Code: "somebot"}
	req := newWebhookRequest(t, "wrong-value")
	assertAuthFailed(t, verifyWebhookSecretToken(context.Background(), req, settings))
}

func TestVerifyWebhookSecretToken_FleetSecretMissingHeaderRejected(t *testing.T) {
	withFleetWebhookSecret(t, "fleet-secret")
	settings := &botsfw.BotSettings{Code: "somebot"}
	req := newWebhookRequest(t, "")
	assertAuthFailed(t, verifyWebhookSecretToken(context.Background(), req, settings))
}

// A per-bot secret overrides the fleet secret end-to-end: the header must match the per-bot
// value, and the fleet value is NOT accepted for that bot.
func TestVerifyWebhookSecretToken_PerBotOverridesFleet_FleetValueRejected(t *testing.T) {
	withFleetWebhookSecret(t, "fleet-secret")
	settings := &botsfw.BotSettings{Code: "somebot", WebhookSecretToken: "per-bot-secret"}
	req := newWebhookRequest(t, "fleet-secret") // sending the fleet value, not the per-bot one
	assertAuthFailed(t, verifyWebhookSecretToken(context.Background(), req, settings))
}
