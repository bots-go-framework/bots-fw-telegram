package telegram

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botmsg"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/strongo/logus"
	"io"
	"net/http"
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

var _ botsfw.WebhookHandler = (*tgWebhookHandler)(nil)

type tgWebhookHandler struct {
	botsfw.WebhookHandlerBase
	botContextProvider botsfw.BotContextProvider
	//botsBy botsfw.BotSettingsProvider
}

// NewTelegramWebhookHandler creates new Telegram webhooks handler
func NewTelegramWebhookHandler(
	botContextProvider botsfw.BotContextProvider,
	translatorProvider botsfw.TranslatorProvider,
	setAppUserFields func(botsfwmodels.AppUserData, botinput.Sender) error, // TODO: Move to botsfwdal.AppUserDal ?
) botsfw.WebhookHandler {
	if botContextProvider == nil {
		panic("botContextProvider == nil")
	}
	if translatorProvider == nil {
		panic("translatorProvider == nil")
	}
	if setAppUserFields == nil {
		panic("setAppUserFields == nil")
	}
	return tgWebhookHandler{
		botContextProvider: botContextProvider,
		WebhookHandlerBase: botsfw.WebhookHandlerBase{
			BotPlatform:        platform{},
			TranslatorProvider: translatorProvider,
			RecordsFieldsSetter: tgBotRecordsFieldsSetter{
				setAppUserFields: setAppUserFields,
			},
		},
	}
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
	//router.POST(pathPrefix+"/telegram/webhook", h.HandleWebhookRequest) // TODO: Remove obsolete
	router.Handle("POST", pathPrefix+"/tg/hook", h.HandleWebhookRequest)
	router.Handle("GET", pathPrefix+"/tg/set-webhook", h.SetWebhook)
	router.Handle("GET", pathPrefix+"/tg/test/time-now", httpHandlerTestTimeNow)
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
		if err := recover(); err != nil {
			stack := string(debug.Stack())
			logus.Criticalf(h.Context(r), "Unhandled panic in Telegram handler: %v\nStack trace: %s", err, stack)
		}
	}()

	h.HandleWebhook(w, r, h)
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
		err = fmt.Errorf("failed to get bot context: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bot := tgbotapi.NewBotAPIWithClient(botContext.BotSettings.Token, client)
	bot.EnableDebug(ctx)
	//bot.Debug = true

	webhookURL := fmt.Sprintf("https://%s/bot/tg/hook?id=%s", r.Host, botCode)

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
	if webhookConfig.SecretToken = botContext.BotSettings.WebhookSecretToken; webhookConfig.SecretToken == "" {
		// SEC-4: registering a webhook without a secret_token leaves it unauthenticated -
		// verifyWebhookSecretToken will log a warning (or reject, if RequireWebhookSecret is
		// set) on every incoming request until BotSettings.WebhookSecretToken is configured
		// for this bot and the webhook is re-registered.
		logus.Warningf(ctx, "SEC-4 WARNING: registering webhook for bot %q WITHOUT a secret_token - its webhook will be unauthenticated", botCode)
	}
	var response tgbotapi.APIResponse
	if response, err = bot.SetWebhook(*webhookConfig); err != nil {
		logus.Errorf(ctx, "%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			logus.Errorf(ctx, "Failed to write error to response: %v", err)
		}
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

// verifyWebhookSecretToken checks the X-Telegram-Bot-Api-Secret-Token header Telegram sends
// on every webhook call against the secret configured for this bot (see SEC-4: without this
// check, anyone who knows/guesses a bot's webhook URL can POST forged updates and be treated
// as any Telegram user, since botID/`?id=` is a public, non-secret value - the bot's own
// @username).
//
// Compat posture: if no secret is configured for the bot (settings.WebhookSecretToken == ""),
// this does NOT block the request - existing deployments that haven't rolled out a secret yet
// keep working - but it logs a high-visibility warning on every single request so the gap is
// impossible to miss in logs, unless settings.RequireWebhookSecret is set, in which case an
// unconfigured secret is treated as a hard misconfiguration and the request is rejected.
// Once a secret IS configured, verification is always strictly enforced.
func verifyWebhookSecretToken(ctx context.Context, r *http.Request, settings *botsfw.BotSettings) error {
	if settings == nil { // defensive: should not happen for a bot resolved via BotContextProvider
		logus.Warningf(ctx, "SEC-4 WARNING: BotSettings is nil, cannot verify %s header - treating as unauthenticated", TelegramWebhookSecretTokenHeader)
		return nil
	}
	expected := settings.WebhookSecretToken
	if expected == "" {
		if settings.RequireWebhookSecret {
			logus.Criticalf(ctx,
				"SEC-4: rejecting Telegram webhook request for bot %q: RequireWebhookSecret is true but no WebhookSecretToken is configured",
				settings.Code)
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
	logus.Debugf(ctx, "tgWebhookHandler.GetBotContextAndInputs(): %s", r.URL.RequestURI())
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

	var requestLogged bool
	logRequestBody := func() {
		if !requestLogged {
			requestLogged = true
			if len(bodyBytes) < 1024*10 {
				var bodyToLog bytes.Buffer
				var bodyStr string
				if indentErr := json.Indent(&bodyToLog, bodyBytes, "", "\t"); indentErr == nil {
					bodyStr = bodyToLog.String()
				} else {
					bodyStr = string(bodyBytes)
				}
				logus.Debugf(ctx, "Request body (%s): %s", r.URL.String(), bodyStr)
			} else {
				logus.Debugf(ctx, "Request len(body): %v", len(bodyBytes))
			}
		}
	}

	var update *tgbotapi.Update
	if update, err = h.unmarshalUpdate(ctx, bodyBytes); err != nil {
		logRequestBody()
		return
	}

	var input botinput.InputMessage
	if input, err = NewTelegramWebhookInput(update, logRequestBody); err != nil {
		logRequestBody()
		return
	}
	logRequestBody()

	entriesWithInputs = []botinput.EntryInputs{
		{
			Entry:  tgWebhookEntry{update: update},
			Inputs: []botinput.InputMessage{input},
		},
	}

	if input == nil {
		logRequestBody()
		err = fmt.Errorf("telegram input is <nil>: %w", botsfw.ErrNotImplemented)
		return
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
	return newTelegramWebhookContext(args, args.WebhookInput.(TgWebhookInput), h.RecordsFieldsSetter)
}

func (h tgWebhookHandler) GetResponder(w http.ResponseWriter, whc botsfw.WebhookContext) botsfw.WebhookResponder {
	if twhc, ok := whc.(*tgWebhookContext); ok {
		return newTgWebhookResponder(w, twhc)
	}
	panic(fmt.Sprintf("Expected tgWebhookContext, got: %T", whc))
}

//func (h tgWebhookHandler) CreateBotCoreStores(appContext botsfw.BotAppContext, r *http.Request) botsfwdal.DataAccess {
//	return h.WebhookHandlerBase.DataAccess
//}
