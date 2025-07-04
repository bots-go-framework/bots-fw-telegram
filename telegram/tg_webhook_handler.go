package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/strongo/logus"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

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
	setAppUserFields func(botsfwmodels.AppUserData, botinput.WebhookSender) error, // TODO: Move to botsfwdal.AppUserDal ?
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

func (h tgWebhookHandler) HandleUnmatched(whc botsfw.WebhookContext) (m botsfw.MessageFromBot) {
	switch whc.Input().InputType() {
	case botinput.WebhookInputCallbackQuery:
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

func (h tgWebhookHandler) GetBotContextAndInputs(ctx context.Context, r *http.Request) (botContext *botsfw.BotContext, entriesWithInputs []botsfw.EntryInputs, err error) {
	logus.Debugf(ctx, "tgWebhookHandler.GetBotContextAndInputs(): %s", r.URL.RequestURI())
	botID := r.URL.Query().Get("id")
	if botContext, err = h.botContextProvider.GetBotContext(ctx, PlatformID, botID); err != nil {
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

	var input botinput.WebhookInput
	if input, err = NewTelegramWebhookInput(update, logRequestBody); err != nil {
		logRequestBody()
		return
	}
	logRequestBody()

	entriesWithInputs = []botsfw.EntryInputs{
		{
			Entry:  tgWebhookEntry{update: update},
			Inputs: []botinput.WebhookInput{input},
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
