package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-store/botsfwdal"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/strongo/log"
	"io"
	"net/http"
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
	dataAccess botsfwdal.DataAccess,
	botsBy botsfw.SettingsProvider,
	recordsMaker botsfwmodels.BotRecordsMaker,
	translatorProvider botsfw.TranslatorProvider,
) botsfw.WebhookHandler {
	if translatorProvider == nil {
		panic("translatorProvider == nil")
	}
	return tgWebhookHandler{
		botsBy: botsBy,
		WebhookHandlerBase: botsfw.WebhookHandlerBase{
			DataAccess:         dataAccess,
			BotPlatform:        Platform{},
			RecordsMaker:       recordsMaker,
			TranslatorProvider: translatorProvider,
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
	router.Handle("GET", pathPrefix+"/tg/set-webhook", func(w http.ResponseWriter, r *http.Request) {
		h.SetWebhook(h.Context(r), w, r)
	})
	router.Handle("GET", pathPrefix+"/tg/test", func(w http.ResponseWriter, r *http.Request) {
		log.Debugf(h.Context(r), "Test request")
		if _, err := w.Write([]byte("Test response")); err != nil {
			log.Errorf(r.Context(), "Failed to write test response: %v", err)
		}
	})
}

func (h tgWebhookHandler) HandleWebhookRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			log.Criticalf(h.Context(r), "Unhandled panic in Telegram handler: %v", err)
		}
	}()

	h.HandleWebhook(w, r, h)
}

func (h tgWebhookHandler) SetWebhook(c context.Context, w http.ResponseWriter, r *http.Request) {
	log.Debugf(c, "tgWebhookHandler.SetWebhook()")
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
		log.Errorf(c, fmt.Sprintf("%v. All bots: %v", m, h.botsBy(c).ByCode))
		return
	}
	bot := tgbotapi.NewBotAPIWithClient(botSettings.Token, client)
	bot.EnableDebug(c)
	//bot.Debug = true

	webhookURL := fmt.Sprintf("https://%v/bot/tg/hook?id=%v&token=%v", r.Host, botCode, bot.Token)

	webhookConfig := tgbotapi.NewWebhook(webhookURL)
	webhookConfig.AllowedUpdates = []string{
		"message",
		"edited_message",
		"inline_query",
		"chosen_inline_result",
		"callback_query",
	}
	if response, err := bot.SetWebhook(*webhookConfig); err != nil {
		log.Errorf(c, "%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Errorf(c, "Failed to write error to response: %v", err)
		}
	} else {
		if _, err := w.Write([]byte(fmt.Sprintf("Webhook set\nErrorCode: %d\nDescription: %v\nContent: %v", response.ErrorCode, response.Description, string(response.Result)))); err != nil {
			log.Errorf(c, "Failed to write error to response: %v", err)
		}
	}
}

func (h tgWebhookHandler) GetBotContextAndInputs(c context.Context, r *http.Request) (botContext *botsfw.BotContext, entriesWithInputs []botsfw.EntryInputs, err error) {
	//log.Debugf(c, "tgWebhookHandler.GetBotContextAndInputs()")
	token := r.URL.Query().Get("token")
	botSettings, ok := h.botsBy(c).ByAPIToken[token]
	if !ok {
		errMess := fmt.Sprintf("Unknown token: [%v]", token)
		err = botsfw.ErrAuthFailed(errMess)
		return
	}
	botContext = botsfw.NewBotContext(h.BotHost, botSettings)
	var bodyBytes []byte
	defer func() {
		if r.Body != nil {
			if err := r.Body.Close(); err != nil {
				log.Errorf(c, "Failed to close request body: %v", err)
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
				log.Debugf(c, "Request body: %v", bodyStr)
			} else {
				log.Debugf(c, "Request len(body): %v", len(bodyBytes))
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
	log.Debugf(c, "Telegram input type: %T", input)
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
	appContext botsfw.BotAppContext,
	r *http.Request, botContext botsfw.BotContext,
	webhookInput botsfw.WebhookInput,
	botCoreStores botsfwdal.DataAccess,
	gaMeasurement botsfw.GaQueuer,
) botsfw.WebhookContext {
	return newTelegramWebhookContext(
		appContext, r, botContext, webhookInput.(TgWebhookInput), botCoreStores, gaMeasurement)
}

func (h tgWebhookHandler) GetResponder(w http.ResponseWriter, whc botsfw.WebhookContext) botsfw.WebhookResponder {
	if twhc, ok := whc.(*tgWebhookContext); ok {
		return newTgWebhookResponder(w, twhc)
	}
	panic(fmt.Sprintf("Expected tgWebhookContext, got: %T", whc))
}

func (h tgWebhookHandler) CreateBotCoreStores(appContext botsfw.BotAppContext, r *http.Request) botsfwdal.DataAccess {
	return h.WebhookHandlerBase.DataAccess
}
