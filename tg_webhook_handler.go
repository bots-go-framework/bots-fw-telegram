package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/strongo/logus"
	"io"
	"net/http"
	"runtime/debug"

	//"github.com/kylelemons/go-gypsy/yaml"
	//"bytes"
	"bytes"
	"github.com/pquerna/ffjson/ffjson"
	"strings"
	"time"
)

var _ botsfw.WebhookHandler = (*tgWebhookHandler)(nil)

type tgWebhookHandler struct {
	botsfw.WebhookHandlerBase
	botsBy botsfw.SettingsProvider
}

// NewTelegramWebhookHandler creates new Telegram webhooks handler
func NewTelegramWebhookHandler(
	//dataAccess botsfwdal.DataAccess,
	botsBy botsfw.SettingsProvider,
	translatorProvider botsfw.TranslatorProvider,
	recordsMaker botsfwmodels.BotRecordsMaker,
	setAppUserFields func(botsfwmodels.AppUserData, botsfw.WebhookSender) error,
) botsfw.WebhookHandler {
	if translatorProvider == nil {
		panic("translatorProvider == nil")
	}
	return tgWebhookHandler{
		botsBy: botsBy,
		WebhookHandlerBase: botsfw.WebhookHandlerBase{
			//DataAccess:         dataAccess,
			BotPlatform:        platform{},
			RecordsMaker:       recordsMaker,
			TranslatorProvider: translatorProvider,
			RecordsFieldsSetter: tgBotRecordsFieldsSetter{
				setAppUserFields: setAppUserFields,
			},
		},
	}
}

func (h tgWebhookHandler) HandleUnmatched(whc botsfw.WebhookContext) (m botsfw.MessageFromBot) {
	switch whc.InputType() {
	case botsfw.WebhookInputCallbackQuery:
		m.BotMessage = CallbackAnswer(tgbotapi.AnswerCallbackQueryConfig{
			Text:      "⚠️ Error: Not matched to any command",
			ShowAlert: true,
		})
	}
	return
}

func (h tgWebhookHandler) RegisterHttpHandlers(driver botsfw.WebhookDriver, host botsfw.BotHost, router botsfw.HttpRouter, pathPrefix string) {
	if router == nil {
		panic("router == nil")
	}
	h.WebhookHandlerBase.Register(driver, host)

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
	c := h.Context(r)
	logus.Debugf(c, "tgWebhookHandler.SetWebhook()")
	ctxWithDeadline, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()
	client := h.GetHTTPClient(ctxWithDeadline)
	botCode := r.URL.Query().Get("code")
	if botCode == "" {
		http.Error(w, "tgWebhookHandler: Missing required parameter: code", http.StatusBadRequest)
		return
	}
	botSettings, ok := h.botsBy(c).ByCode[botCode]
	if !ok {
		m := fmt.Sprintf("Bot not found by code: %v", botCode)
		http.Error(w, m, http.StatusBadRequest)
		logus.Errorf(c, fmt.Sprintf("%v. All bots: %v", m, h.botsBy(c).ByCode))
		return
	}
	bot := tgbotapi.NewBotAPIWithClient(botSettings.Token, client)
	bot.EnableDebug(c)
	//bot.Debug = true

	webhookURL := fmt.Sprintf("https://%v/bot/tg/hook?id=%v", r.Host, botCode)

	webhookConfig := tgbotapi.NewWebhook(webhookURL)
	webhookConfig.AllowedUpdates = []string{
		"message",
		"edited_message",
		"inline_query",
		"chosen_inline_result",
		"callback_query",
	}
	if response, err := bot.SetWebhook(*webhookConfig); err != nil {
		logus.Errorf(c, "%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			logus.Errorf(c, "Failed to write error to response: %v", err)
		}
	} else {
		if _, err := w.Write([]byte(fmt.Sprintf("Webhook set\nErrorCode: %d\nDescription: %v\nContent: %v", response.ErrorCode, response.Description, string(response.Result)))); err != nil {
			logus.Errorf(c, "Failed to write error to response: %v", err)
		}
	}
}

func (h tgWebhookHandler) GetBotContextAndInputs(c context.Context, r *http.Request) (botContext *botsfw.BotContext, entriesWithInputs []botsfw.EntryInputs, err error) {
	logus.Debugf(c, "tgWebhookHandler.GetBotContextAndInputs(): %s", r.URL.RequestURI())
	botID := r.URL.Query().Get("id")
	botSettings, ok := h.botsBy(c).ByCode[botID]
	if !ok {
		errMess := fmt.Sprintf("Unknown bot ID (username): [%v]", botID)
		err = botsfw.ErrAuthFailed(errMess)
		return
	}
	botContext = botsfw.NewBotContext(h.BotHost, botSettings)
	var bodyBytes []byte
	defer func() {
		if r.Body != nil {
			if err := r.Body.Close(); err != nil {
				logus.Errorf(c, "Failed to close request body: %v", err)
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
				logus.Debugf(c, "Request body: %v", bodyStr)
			} else {
				logus.Debugf(c, "Request len(body): %v", len(bodyBytes))
			}
		}
	}

	var update *tgbotapi.Update
	if update, err = h.unmarshalUpdate(c, bodyBytes); err != nil {
		logRequestBody()
		return
	}

	var input botsfw.WebhookInput
	if input, err = NewTelegramWebhookInput(update, logRequestBody); err != nil {
		logRequestBody()
		return
	}
	logRequestBody()

	entriesWithInputs = []botsfw.EntryInputs{
		{
			Entry:  tgWebhookEntry{update: update},
			Inputs: []botsfw.WebhookInput{input},
		},
	}

	if input == nil {
		logRequestBody()
		err = fmt.Errorf("telegram input is <nil>: %w", botsfw.ErrNotImplemented)
		return
	}
	logus.Debugf(c, "Telegram input type: %T", input)
	return
}

func (h tgWebhookHandler) unmarshalUpdate(_ context.Context, content []byte) (update *tgbotapi.Update, err error) {
	update = new(tgbotapi.Update)
	if err = ffjson.UnmarshalFast(content, update); err != nil {
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
