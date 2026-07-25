package telegram

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botmsg"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/strongo/logus"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// TelegramWebhookSecretTokenHeader is the HTTP header Telegram sends the configured
// `secret_token` back on with every webhook call, so a handler can verify the request
// actually came from Telegram (and not from anyone who guessed/discovered the webhook
// URL, which is not itself a secret - see SEC-4).
// https://core.telegram.org/bots/api#setwebhook
const TelegramWebhookSecretTokenHeader = "X-Telegram-Bot-Api-Secret-Token"

// ForwardedHostHeader is the de-facto standard reverse-proxy header carrying the
// original Host the client requested, before any Host-rewriting hop (e.g. a proxy
// that swaps Host to match a backend's own routing requirements). SetWebhook prefers
// this over r.Host so the webhook URL registered with Telegram reflects the public
// host callers actually want, even when reached through such a proxy.
const ForwardedHostHeader = "X-Forwarded-Host"

var _ botsfw.WebhookHandler = (*tgWebhookHandler)(nil)

type tgWebhookHandler struct {
	botsfw.WebhookHandlerBase
	botContextProvider botsfw.BotContextProvider
	//botsBy botsfw.BotSettingsProvider
	pathPrefix string
	// adminPathPrefix, when non-empty, is the prefix the admin-only endpoints
	// (set-webhook, test) mount under instead of pathPrefix — see WithAdminPathPrefix.
	adminPathPrefix    string
	chatInstances      ChatInstanceStore
	responderTransform WebhookResponderTransform
}

// TgWebhookHandlerOption configures a Telegram webhook handler at construction
// time (functional-options pattern).
type TgWebhookHandlerOption func(*tgWebhookHandler)

// WebhookResponderTransform wraps the responder used for a single webhook.
// The transformed responder is installed into WebhookContext as well as passed
// to router dispatch, so direct responder sends and router-returned messages
// share one host-enforced delivery boundary.
type WebhookResponderTransform func(botsfw.WebhookContext, botsfw.WebhookResponder) botsfw.WebhookResponder

// WithResponderTransform installs a host-owned responder transform. The
// transform is applied once per webhook after the Telegram responder has been
// created and before any router action can run.
func WithResponderTransform(transform WebhookResponderTransform) TgWebhookHandlerOption {
	if transform == nil {
		panic("transform == nil")
	}
	return func(h *tgWebhookHandler) { h.responderTransform = transform }
}

// WithAdminPathPrefix mounts the admin-only endpoints (set-webhook and the test
// route) under adminPrefix instead of the public webhook pathPrefix, so a host can
// gate that prefix with admin auth. set-webhook (re-)points a bot's live webhook
// using the bot's own token, so on a publicly reachable origin it must NOT be open.
//
// The public /tg/hook endpoint (which Telegram POSTs to) and the hook URL that
// set-webhook registers with Telegram are unaffected — both keep using the public
// pathPrefix. An empty adminPrefix (the default) preserves the legacy behaviour of
// mounting every endpoint under the public pathPrefix.
func WithAdminPathPrefix(adminPrefix string) TgWebhookHandlerOption {
	return func(h *tgWebhookHandler) { h.adminPathPrefix = adminPrefix }
}

// NewTelegramWebhookHandler creates new Telegram webhooks handler
func NewTelegramWebhookHandler(
	botContextProvider botsfw.BotContextProvider,
	translatorProvider botsfw.TranslatorProvider,
	chatInstances ChatInstanceStore,
	opts ...TgWebhookHandlerOption,
) botsfw.WebhookHandler {
	if botContextProvider == nil {
		panic("botContextProvider == nil")
	}
	if translatorProvider == nil {
		panic("translatorProvider == nil")
	}
	if chatInstances == nil {
		panic("chatInstances == nil")
	}
	h := tgWebhookHandler{
		botContextProvider: botContextProvider,
		chatInstances:      chatInstances,
		WebhookHandlerBase: botsfw.WebhookHandlerBase{
			BotPlatform:         platform{},
			TranslatorProvider:  translatorProvider,
			RecordsFieldsSetter: tgBotRecordsFieldsSetter{},
		},
	}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

func (h tgWebhookHandler) HandleUnmatched(whc botsfw.WebhookContext) (m botmsg.MessageFromBot) {
	switch whc.Input().InputType() {
	case botinput.TypeCallbackQuery:
		m.BotMessage = CallbackAnswer(tgbotapi.AnswerCallbackQueryConfig{
			Text:      "⚠️ Error: Not matched to any command",
			ShowAlert: true,
		})
	default:
		// TODO: Do nothing?
	}
	return
}

func (h tgWebhookHandler) RegisterHttpHandlers(driver botsfw.WebhookDriver, host botsfw.BotHost, router botsfw.HttpRouter, pathPrefix string) {
	if router == nil {
		panic("router == nil")
	}
	h.Register(driver, host)

	pathPrefix = strings.TrimSuffix(pathPrefix, "/")
	h.pathPrefix = pathPrefix

	// Admin-only endpoints (set-webhook, test) mount under adminPathPrefix when set,
	// so a host can gate that prefix with admin auth; otherwise they fall back to the
	// public pathPrefix (legacy behaviour). The public /tg/hook endpoint and the hook
	// URL that set-webhook registers always use the public pathPrefix.
	adminPrefix := strings.TrimSuffix(h.adminPathPrefix, "/")
	if adminPrefix == "" {
		adminPrefix = pathPrefix
	}

	//router.POST(pathPrefix+"/telegram/webhook", h.HandleWebhookRequest) // TODO: Remove obsolete
	router.Handle("POST", pathPrefix+"/tg/hook", h.HandleWebhookRequest)
	router.Handle("GET", adminPrefix+"/tg/set-webhook", h.SetWebhook)
	router.Handle("GET", adminPrefix+"/tg/test/time-now", httpHandlerTestTimeNow)
}

func httpHandlerTestTimeNow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logus.Debugf(ctx, "Test request")
	now := time.Now().Format(time.RFC3339Nano)
	if _, err := w.Write([]byte("Test: " + now)); err != nil {
		logus.Errorf(ctx, "Failed to write test response: %v", err)
	}
}

func (h tgWebhookHandler) HandleWebhookRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			stack := string(debug.Stack())
			logus.Criticalf(h.Context(r), "Unhandled panic in Telegram handler: panic_type=%T\nStack trace: %s", recovered, stack)
		}
	}()

	h.HandleWebhook(w, r, h)
}

// publicHost resolves the host that should appear in URLs handed back to third
// parties (e.g. the Telegram webhook URL), preferring ForwardedHostHeader - set by
// a reverse proxy that rewrites the wire Host to something the origin needs - over
// r.Host, which reflects that rewritten value rather than what the caller reached.
func publicHost(r *http.Request) string {
	if host := r.Header.Get(ForwardedHostHeader); host != "" {
		return host
	}
	return r.Host
}

func (h tgWebhookHandler) SetWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := h.Context(r)
	logus.Debugf(ctx, "tgWebhookHandler.SetWebhook()")
	ctxWithDeadline, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	client := h.GetHTTPClient(ctxWithDeadline)
	botCode := r.URL.Query().Get("code")
	if botCode == "" {
		http.Error(w, "tgWebhookHandler: Missing required parameter: code", http.StatusBadRequest)
		return
	}
	botContext, err := h.botContextProvider.GetBotContext(ctx, PlatformID, botCode)
	if err != nil {
		logus.Errorf(ctx, "Failed to get bot context for Telegram webhook setup: error_type=%T", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bot := tgbotapi.NewBotAPIWithClient(botContext.BotSettings.Token, client)

	webhookURL := fmt.Sprintf("https://%s%s/tg/hook?id=%s", publicHost(r), h.pathPrefix, botCode)

	webhookConfig := tgbotapi.NewWebhook(webhookURL)
	webhookConfig.AllowedUpdates = []string{
		"message",
		"edited_message",
		"inline_query",
		"chosen_inline_result",
		"callback_query",
		"pre_checkout_query",
		"successful_payment",
		"refunded_payment",
		"purchased_paid_media",
	}
	if webhookConfig.SecretToken = effectiveWebhookSecret(botContext.BotSettings); webhookConfig.SecretToken == "" {
		// SEC-4: neither a per-bot BotSettings.WebhookSecretToken nor the fleet-wide
		// TELEGRAM_WEBHOOK_SECRET (see effectiveWebhookSecret) is set, so this webhook is
		// registered without a secret_token and stays unauthenticated - verifyWebhookSecretToken
		// will log a warning (or reject, if RequireWebhookSecret is set) on every incoming
		// request until a secret is configured and the webhook is re-registered.
		logus.Warningf(ctx, "SEC-4 WARNING: registering webhook for bot %q WITHOUT a secret_token - its webhook will be unauthenticated", botCode)
	}
	var response tgbotapi.APIResponse
	if response, err = bot.SetWebhook(*webhookConfig); err != nil {
		logus.Errorf(ctx, "Failed to set Telegram webhook: error_type=%T", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	} else if _, err = fmt.Fprintf(w, `Webhook set
ErrorCode: %d
Description: %v
Content: %v

Parametes:
	allowed_updates: %s
`, response.ErrorCode, response.Description, string(response.Result), strings.Join(webhookConfig.AllowedUpdates, ",")); err != nil {
		logus.Errorf(ctx, "Failed to write error to response: %v", err)
	}
}

// EnvTelegramWebhookSecret is the environment variable holding the single, fleet-wide
// Telegram webhook secret_token. It is applied to EVERY bot that does not set its own
// BotSettings.WebhookSecretToken, so one provisioned secret authenticates the whole
// fleet's webhooks with no per-bot configuration (SEC-4). A bot that needs an isolated
// secret can still override it via BotSettings.WebhookSecretToken.
const EnvTelegramWebhookSecret = "TELEGRAM_WEBHOOK_SECRET"

// resolveFleetWebhookSecret reads the fleet-wide webhook secret. Indirected through a
// package var so tests can substitute a value without mutating the process environment.
var resolveFleetWebhookSecret = func() string { return os.Getenv(EnvTelegramWebhookSecret) }

// EnvTelegramRequireWebhookSecret, when truthy ("true"/"1"/"yes"/"on"), makes a missing
// webhook secret a HARD FAILURE fleet-wide: if no secret resolves for a bot (neither a
// per-bot WebhookSecretToken nor the fleet-wide TELEGRAM_WEBHOOK_SECRET), the webhook is
// rejected rather than allowed with a warning. It is the fleet-wide equivalent of the
// per-bot BotSettings.RequireWebhookSecret — set it in production so an accidentally-unset
// secret fails CLOSED (bots reject) instead of silently serving unauthenticated webhooks.
const EnvTelegramRequireWebhookSecret = "TELEGRAM_WEBHOOK_REQUIRE_SECRET"

// resolveRequireWebhookSecret reports whether a webhook secret is required fleet-wide.
// Indirected through a package var so tests can override without touching the environment.
var resolveRequireWebhookSecret = func() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(EnvTelegramRequireWebhookSecret))) {
	case "true", "1", "yes", "on":
		return true
	}
	return false
}

// effectiveWebhookSecret returns the secret_token to register (SetWebhook) and to verify
// for a bot: the bot's own BotSettings.WebhookSecretToken when set, otherwise the
// fleet-wide TELEGRAM_WEBHOOK_SECRET. It returns "" only when neither is configured, in
// which case the webhook stays in the unauthenticated compat posture (see
// verifyWebhookSecretToken). This is the single, framework-owned standard: hosts provision
// one secret and every bot is authenticated uniformly.
func effectiveWebhookSecret(settings *botsfw.BotSettings) string {
	if settings != nil && settings.WebhookSecretToken != "" {
		return settings.WebhookSecretToken
	}
	return resolveFleetWebhookSecret()
}

// verifyWebhookSecretToken checks the X-Telegram-Bot-Api-Secret-Token header Telegram sends
// on every webhook call against the secret resolved for this bot (see SEC-4: without this
// check, anyone who knows/guesses a bot's webhook URL can POST forged updates and be treated
// as any Telegram user, since botID/`?id=` is a public, non-secret value - the bot's own
// @username). The expected secret is effectiveWebhookSecret(settings): the per-bot
// WebhookSecretToken if set, else the fleet-wide TELEGRAM_WEBHOOK_SECRET.
//
// Compat posture: if no secret resolves for the bot - neither a per-bot WebhookSecretToken
// nor the fleet-wide TELEGRAM_WEBHOOK_SECRET - this does NOT block the request (existing
// deployments that haven't rolled out a secret yet keep working) but logs a high-visibility
// warning on every single request so the gap is impossible to miss, unless a secret is
// required (per-bot settings.RequireWebhookSecret or fleet-wide TELEGRAM_WEBHOOK_REQUIRE_SECRET),
// in which case an unconfigured secret is a hard misconfiguration and the request is
// rejected (fail closed). Once a secret IS configured (per-bot or fleet-wide), verification
// is always strictly enforced.
func verifyWebhookSecretToken(ctx context.Context, r *http.Request, settings *botsfw.BotSettings) error {
	if settings == nil { // defensive: should not happen for a bot resolved via BotContextProvider
		logus.Warningf(ctx, "SEC-4 WARNING: BotSettings is nil, cannot verify %s header - treating as unauthenticated", TelegramWebhookSecretTokenHeader)
		return nil
	}
	expected := effectiveWebhookSecret(settings)
	if expected == "" {
		if settings.RequireWebhookSecret || resolveRequireWebhookSecret() {
			logus.Criticalf(ctx,
				"SEC-4: rejecting Telegram webhook request for bot %q: a webhook secret is required (per-bot RequireWebhookSecret or fleet-wide %s) but none is configured (no per-bot WebhookSecretToken, no %s)",
				settings.Code, EnvTelegramRequireWebhookSecret, EnvTelegramWebhookSecret)
			return botsfw.ErrAuthFailed(fmt.Sprintf("webhook secret required but not configured for bot %q", settings.Code))
		}
		logus.Warningf(ctx,
			"SEC-4 WARNING: bot %q has NO webhook secret configured - its webhook is UNAUTHENTICATED and anyone who knows the webhook URL can forge Telegram updates and impersonate any user. Configure BotSettings.WebhookSecretToken (and re-register the webhook with a matching secret_token) to close this gap.",
			settings.Code)
		return nil
	}

	got := r.Header.Get(TelegramWebhookSecretTokenHeader)
	if got == "" || subtle.ConstantTimeCompare([]byte(got), []byte(expected)) != 1 {
		logus.Warningf(ctx, "SEC-4: rejecting Telegram webhook request for bot %q: missing/invalid %s header", settings.Code, TelegramWebhookSecretTokenHeader)
		return botsfw.ErrAuthFailed(fmt.Sprintf("invalid or missing %s header", TelegramWebhookSecretTokenHeader))
	}
	return nil
}

func (h tgWebhookHandler) GetBotContextAndInputs(ctx context.Context, r *http.Request) (botContext *botsfw.BotContext, entriesWithInputs []botinput.EntryInputs, err error) {
	logus.Debugf(ctx, "tgWebhookHandler.GetBotContextAndInputs(): path=%q", r.URL.EscapedPath())
	botID := r.URL.Query().Get("id")
	if botContext, err = h.botContextProvider.GetBotContext(ctx, PlatformID, botID); err != nil {
		return
	}

	if err = verifyWebhookSecretToken(ctx, r, botContext.BotSettings); err != nil {
		botContext = nil
		return
	}

	var bodyBytes []byte
	defer func() {
		if r.Body != nil {
			if err := r.Body.Close(); err != nil {
				logus.Errorf(ctx, "Failed to close request body: %v", err)
			}
		}
	}()
	if bodyBytes, err = io.ReadAll(r.Body); err != nil {
		err = fmt.Errorf("failed to read request body: %w", err)
		return
	}

	var (
		requestLogged bool
		update        *tgbotapi.Update
		input         botinput.InputMessage
		logOutcome    = "decode_error"
	)
	logRequestMetadata := func() {
		if !requestLogged {
			requestLogged = true
			inputType := botinput.TypeUnknown
			if input != nil {
				inputType = input.InputType()
			}
			logus.Debugf(
				ctx,
				"Telegram webhook received: bot=%q, body_bytes=%d, input_type=%v, outcome=%s",
				botID,
				len(bodyBytes),
				inputType,
				logOutcome,
			)
		}
	}

	if update, err = h.unmarshalUpdate(ctx, bodyBytes); err != nil {
		logRequestMetadata()
		return
	}

	logOutcome = "unsupported"
	if input, err = NewTelegramWebhookInput(update, logRequestMetadata); err != nil {
		logRequestMetadata()
		return
	}
	if input == nil {
		logRequestMetadata()
		err = fmt.Errorf("telegram input is <nil>: %w", botsfw.ErrNotImplemented)
		return
	}
	logOutcome = "decoded"
	logRequestMetadata()

	entriesWithInputs = []botinput.EntryInputs{
		{
			Entry:  tgWebhookEntry{update: update},
			Inputs: []botinput.InputMessage{input},
		},
	}

	logus.Debugf(ctx, "Telegram input type: %T", input)
	return
}

func (h tgWebhookHandler) unmarshalUpdate(_ context.Context, content []byte) (update *tgbotapi.Update, err error) {
	update = new(tgbotapi.Update)
	if err = json.Unmarshal(content, update); err != nil {
		return
	}
	return
}

func (h tgWebhookHandler) CreateWebhookContext(
	args botsfw.CreateWebhookContextArgs,
) (botsfw.WebhookContext, error) {
	return newTelegramWebhookContext(args, args.WebhookInput.(TgWebhookInput), h.RecordsFieldsSetter, h.chatInstances)
}

func (h tgWebhookHandler) GetResponder(w http.ResponseWriter, whc botsfw.WebhookContext) botsfw.WebhookResponder {
	if twhc, ok := whc.(*tgWebhookContext); ok {
		return h.installResponder(twhc, newTgWebhookResponder(w, twhc))
	}
	panic(fmt.Sprintf("Expected tgWebhookContext, got: %T", whc))
}

func (h tgWebhookHandler) installResponder(whc *tgWebhookContext, responder botsfw.WebhookResponder) botsfw.WebhookResponder {
	if h.responderTransform != nil {
		responder = h.responderTransform(whc, responder)
		if responder == nil {
			panic("WebhookResponderTransform returned nil")
		}
	}
	whc.responder = responder
	return responder
}

//func (h tgWebhookHandler) CreateBotCoreStores(appContext botsfw.BotAppContext, r *http.Request) botsfwdal.DataAccess {
//	return h.WebhookHandlerBase.DataAccess
//}
