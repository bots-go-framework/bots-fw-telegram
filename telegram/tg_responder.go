package telegram

import (
	"context"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botmsg"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/bots-go-framework/bots-go-core/botkb"
	"github.com/strongo/logus"
	"net/http"
	"strconv"
	"time"
)

type tgWebhookResponder struct {
	w   http.ResponseWriter
	whc *tgWebhookContext
}

func telegramParseMode(format botmsg.Format) string {
	switch format {
	case botmsg.FormatHTML:
		return "HTML"
	case botmsg.FormatMarkdown:
		return "MarkdownV2"
	case botmsg.FormatText:
		return ""
	default:
		panic(fmt.Sprintf("Unknown message parse_mode value: %d", format))
	}
}

func configureTelegramTextMessage(messageConfig *tgbotapi.MessageConfig, m botmsg.MessageFromBot) {
	messageConfig.ParseMode = telegramParseMode(m.Format)
	messageConfig.DisableWebPagePreview = m.DisableWebPagePreview
	messageConfig.DisableNotification = m.DisableNotification
}

func (r tgWebhookResponder) DeleteMessage(ctx context.Context, messageID string) (err error) {
	var msgID int
	if msgID, err = strconv.Atoi(messageID); err != nil {
		err = errors.New("failed to parse Telegram message ID as integer")
		return
	}
	chatID := r.whc.chatID
	if chatID == "" {
		input := r.whc.Input()
		var chat botinput.Chat
		if inputWithChat, ok := input.(interface{ Chat() botinput.Chat }); ok {
			chat = inputWithChat.Chat()
		}
		if chat != nil {
			chatID = chat.GetID()
		} else {
			var message *tgbotapi.Message
			switch tgInput := input.(type) {
			case tgWebhookTextMessage:
				chat = tgInput.Chat()
				message = tgInput.update.Message
			case callbackQueryInput:
				chat = tgInput.Chat()
				message = tgInput.update.Message
			}
			if message != nil && message.Chat != nil && message.Chat.ID != 0 {
				chatID = strconv.FormatInt(message.Chat.ID, 10)
			}
		}
		if chatID == "" && chat != nil {
			chatID = chat.GetID()
		}
	}
	if chatID == "" {
		return errors.New("can not determine chatID from current WebhookContext")
	}
	botContext := r.whc.BotContext()
	httpClient := botContext.BotHost.GetHTTPClient(ctx)
	botAPI := tgbotapi.NewBotAPIWithClient(botContext.BotSettings.Token, httpClient)
	_, err = botAPI.DeleteMessage(chatID, msgID)
	return
}

var _ botsfw.WebhookResponder = (*tgWebhookResponder)(nil)

func newTgWebhookResponder(w http.ResponseWriter, whc *tgWebhookContext) tgWebhookResponder {
	responder := tgWebhookResponder{w: w, whc: whc}
	whc.responder = responder
	return responder
}

func (r tgWebhookResponder) SendMessage(ctx context.Context, m botmsg.MessageFromBot, channel botmsg.BotAPISendMessageChannel) (resp botsfw.OnMessageSentResponse, err error) {
	logus.Debugf(
		ctx,
		"Telegram response prepared: channel=%q, is_edit=%t, text_bytes=%d, has_keyboard=%t, bot_message_type=%T",
		channel,
		m.IsEdit,
		len(m.Text),
		m.Keyboard != nil,
		m.BotMessage,
	)
	switch channel {
	case botsfw.BotAPISendMessageOverHTTPS, botsfw.BotAPISendMessageOverResponse:
	// Known channels
	default:
		panic(fmt.Sprintf("Unknown channel: [%v]. Expected either 'https' or 'response'.", channel))
	}
	var sendable tgbotapi.Sendable

	tgUpdate := r.whc.Input().(tgWebhookUpdateProvider).TgUpdate()

	var botMessage botmsg.BotMessage

	if m.Text == botmsg.NoMessageToSend {
		logus.Debugf(ctx, botmsg.NoMessageToSend)
		return
	} else if botMessage = m.BotMessage; botMessage != nil {
		logus.Debugf(ctx, "m.BotMessage != nil")
		switch m.BotMessage.BotMessageType() {
		case botmsg.BotMessageTypeInlineResults:
			sendable = tgbotapi.InlineConfig(m.BotMessage.(InlineBotMessage))
		case botmsg.TypeCallbackAnswer:
			var callbackAnswer tgbotapi.AnswerCallbackQueryConfig
			switch botMsg := botMessage.(type) {
			case CallbackAnswer:
				callbackAnswer = tgbotapi.AnswerCallbackQueryConfig(botMsg)
			case botmsg.AnswerCallbackQuery:
				callbackAnswer = tgbotapi.AnswerCallbackQueryConfig{
					CallbackQueryID: botMsg.CallbackQueryID,
					Text:            botMsg.Text,
					ShowAlert:       botMsg.ShowAlert,
					URL:             botMsg.URL,
					CacheTime:       botMsg.CacheTime,
				}
			}
			if callbackAnswer.CallbackQueryID == "" && tgUpdate.CallbackQuery != nil {
				callbackAnswer.CallbackQueryID = tgUpdate.CallbackQuery.ID
			}
			sendable = callbackAnswer
		case botmsg.TypeLeaveChat:
			leaveChat := tgbotapi.LeaveChatConfig(m.BotMessage.(LeaveChat))
			if leaveChat.ChatID == "" {
				leaveChat.ChatID = strconv.FormatInt(tgUpdate.Chat().ID, 10)
			}
			sendable = leaveChat
		case botmsg.TypeExportChatInviteLink:
			exportChatInviteLink := tgbotapi.ExportChatInviteLink(m.BotMessage.(ExportChatInviteLink))
			if exportChatInviteLink.ChatID == "" {
				exportChatInviteLink.ChatID = strconv.FormatInt(tgUpdate.Chat().ID, 10)
			}
			sendable = exportChatInviteLink
		case botmsg.TypeUndefined:
			err = fmt.Errorf("bot message type %v==undefined", m.BotMessage.BotMessageType())
			return
		case botmsg.TypeSendInvoice:
			invoiceConfig := tgbotapi.InvoiceConfig(m.BotMessage.(Invoice))
			if invoiceConfig.ChatID == 0 {
				invoiceConfig.ChatID = tgUpdate.Chat().ID
			}
			sendable = &invoiceConfig
		case botmsg.TypeSetDescription:
			setBotDescription := m.BotMessage.(SetBotDescription)
			sendable = (tgbotapi.SetMyDescription)(setBotDescription)
		case botmsg.TypeSetShortDescription:
			setBotDescription := m.BotMessage.(SetBotShortDescription)
			sendable = (tgbotapi.SetMyShortDescription)(setBotDescription)
		case botmsg.TypeSetCommands:
			setBotDescription := m.BotMessage.(SetBotCommands)
			sendable = (tgbotapi.SetMyCommandsConfig)(setBotDescription)
		case botmsg.TypeAnswerPreCheckoutQuery:
			answerPreCheckoutQuery := m.BotMessage.(PreCheckoutQueryAnswer)
			sendable = (tgbotapi.AnswerPreCheckoutQueryConfig)(answerPreCheckoutQuery)
		case botmsg.TypeSendPhoto:
			photoConfig := m.BotMessage.(SendPhoto)
			if photoConfig.ChatID == 0 {
				photoConfig.ChatID = tgUpdate.Chat().ID
			}
			if photoConfig.Caption != "" {
				photoConfig.ParseMode = telegramParseMode(m.Format)
			}
			sendable = (tgbotapi.PhotoConfig)(photoConfig)
		case botmsg.TypeChatAction:
			sendable = chatActionSendable(m.BotMessage.(botmsg.ChatAction), tgUpdate.Chat().ID)
		default:
			//var ok bool
			//sendable, ok = m.BotMessage.(tgbotapi.Sendable)
			//if !ok {
			err = fmt.Errorf("unknown bot message type %v==%T", m.BotMessage.BotMessageType(), botMessage)
			return
			//}
		}
	} else if m.IsEdit || m.EditMessageIntID != 0 || (tgUpdate.CallbackQuery != nil && tgUpdate.CallbackQuery.InlineMessageID != "" && m.ToChat == nil) {
		// Edit message
		inlineMessageID, chatID, messageID := getTgMessageIDs(tgUpdate)
		if m.EditMessageIntID != 0 {
			messageID = m.EditMessageIntID
			inlineMessageID = ""
		}
		if m.EditMessageUID != nil {
			switch messageUID := m.EditMessageUID.(type) { // TODO: How do we remove duplicates for value & pointer cases?
			case callbackCurrent:
				// do nothing
			case InlineMessageUID, *InlineMessageUID:
				inlineMessageID = messageUID.UID()
				chatID = 0
				messageID = 0
			default:
				err = fmt.Errorf("unknown EditMessageUID type %T", m.EditMessageUID)
				return
			case ChatMessageUID, *ChatMessageUID:
				inlineMessageID = ""
				if uid, ok := messageUID.(interface {
					GetChatID() int64
					GetMessageID() int
				}); ok {
					chatID = uid.GetChatID()
					messageID = uid.GetMessageID()
				}
			}
		}
		logus.Debugf(
			ctx,
			"Telegram edit target resolved: has_inline_id=%t, has_chat_id=%t, has_message_id=%t",
			inlineMessageID != "",
			chatID != 0,
			messageID != 0,
		)
		if inlineMessageID == "" && chatID == 0 && messageID == 0 {
			err = errors.New("can't edit Telegram message as inlineMessageID is empty && chatID == 0 && messageID == 0")
			return
		}
		if m.Text == "" && m.Keyboard != nil {
			keyboard := getTelegramKeyboard(m.Keyboard)
			switch kb := keyboard.(type) {
			case *tgbotapi.InlineKeyboardMarkup:
				sendable = tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, inlineMessageID, kb)
			case *tgbotapi.ReplyKeyboardMarkup, *tgbotapi.ReplyKeyboardHide:
				msg := tgbotapi.NewMessage(chatID, "")
				msg.ReplyMarkup = kb
				sendable = msg
			default:
				err = fmt.Errorf("unknown Telegram keyboard type %T", keyboard)
				return
			}
		} else if m.Text != "" {
			kb := getTelegramKeyboard(m.Keyboard)

			createEditMessage := func() *tgbotapi.EditMessageTextConfig {
				editMessageTextConfig := tgbotapi.NewEditMessageText(chatID, messageID, inlineMessageID, m.Text)
				editMessageTextConfig.ParseMode = telegramParseMode(m.Format)
				editMessageTextConfig.DisableWebPagePreview = m.DisableWebPagePreview
				sendable = editMessageTextConfig
				return editMessageTextConfig
			}

			if kb == nil {
				createEditMessage()
			} else {
				switch kb := kb.(type) {
				case *tgbotapi.InlineKeyboardMarkup:
					editMessageTextConfig := createEditMessage()
					editMessageTextConfig.ReplyMarkup = kb
					sendable = editMessageTextConfig
				case *tgbotapi.ReplyKeyboardMarkup, *tgbotapi.ReplyKeyboardHide:
					messageConfig := tgbotapi.NewMessage(chatID, m.Text)
					messageConfig.ReplyMarkup = kb
					configureTelegramTextMessage(messageConfig, m)
					sendable = messageConfig
				default:
					err = fmt.Errorf("unknown Telegram keyboard type %T", kb)
					return
				}
			}
		} else {
			err = errors.New("can't edit Telegram message without text or keyboard")
			panic(err)
			// return
		}
	} else if m.Text != "" {
		messageConfig := r.whc.NewTgMessage(m.Text)
		if m.ToChat != nil {
			messageConfig.ChatID = int64(m.ToChat.(botmsg.ChatIntID))
		}
		configureTelegramTextMessage(messageConfig, m)
		if m.Keyboard != nil {
			messageConfig.ReplyMarkup = getTelegramKeyboard(m.Keyboard)
		}

		sendable = messageConfig
	} else {
		switch inputType := r.whc.InputType(); inputType {
		case botinput.TypeInlineQuery: // pass
			logus.Debugf(ctx, "No response to WebhookInputInlineQuery")
		case botinput.TypeChosenInlineResult: // pass
		default:
			logus.Debugf(
				ctx,
				"Telegram response omitted: input_type=%v, is_edit=%t, has_keyboard=%t, bot_message_type=%T",
				inputType,
				m.IsEdit,
				m.Keyboard != nil,
				m.BotMessage,
			)
		}
		return
	}

	logus.Debugf(
		ctx,
		"Sending Telegram response: channel=%q, sendable_type=%T, text_bytes=%d",
		channel,
		sendable,
		len(m.Text),
	)

	switch channel {
	case botsfw.BotAPISendMessageOverResponse:
		if _, err = tgbotapi.ReplyToResponse(sendable, r.w); err != nil {
			logus.Errorf(ctx, "Failed to send message to Telegram via HTTP response: error_type=%T", err)
		}
		return resp, err
	case botsfw.BotAPISendMessageOverHTTPS:
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		var message tgbotapi.Message
		message, err = r.sendOverHttps(ctx, sendable)
		return botsfw.OnMessageSentResponse{Message: &message}, err
	default:
		panic(fmt.Sprintf("Unknown channel: %v", channel))
	}
}

func (r tgWebhookResponder) sendOverHttps(ctx context.Context, chattable tgbotapi.Sendable) (message tgbotapi.Message, err error) {
	botContext := r.whc.BotContext()
	botAPI := tgbotapi.NewBotAPIWithClient(
		botContext.BotSettings.Token,
		botContext.BotHost.GetHTTPClient(ctx),
	)

	if message, err = botAPI.Send(chattable); err != nil {
		return
	} else if message.MessageID != 0 {
		logus.Debugf(ctx, "Telegram API response received: has_message_id=true")
	} else {
		logus.Debugf(ctx, "Telegram API response received without message ID: sendable_type=%T", chattable)
	}
	return
}

func getTelegramKeyboard(keyboard botkb.Keyboard) tgbotapi.Keyboard {
	// A text-only edit has no reply markup. The responder intentionally calls
	// this helper for both text-only and keyboard-bearing edits, so nil must be
	// represented as a missing Telegram keyboard rather than a programmer error.
	if keyboard == nil {
		return nil
	}
	if kb, ok := keyboard.(tgbotapi.Keyboard); ok {
		return kb
	}
	switch kb := keyboard.(type) {
	case *botkb.MessageKeyboard:
		switch keyboard.KeyboardType() {
		case botkb.KeyboardTypeInline:
			return getInlineKeyboard(kb)
		case botkb.KeyboardTypeBottom:
			return getReplyKeyboard(kb)
		case botkb.KeyboardTypeHide:
			return getHideKeyboard(kb)
		default:
			panic(fmt.Sprintf("keyboard.KeyboardType() returns unsupported type %v", kb.KeyboardType()))
		}
	default:
		panic(fmt.Sprintf("keyboard is of unsupported type %T", keyboard))
	}
}

func getHideKeyboard(_ *botkb.MessageKeyboard) *tgbotapi.ReplyKeyboardHide {
	return &tgbotapi.ReplyKeyboardHide{HideKeyboard: true}
}

func getReplyKeyboard(kb *botkb.MessageKeyboard) *tgbotapi.ReplyKeyboardMarkup {
	tgButtons := make([][]tgbotapi.KeyboardButton, len(kb.Buttons))
	for i, buttons := range kb.Buttons {
		tgButtons[i] = make([]tgbotapi.KeyboardButton, len(buttons))
		for j, button := range buttons {
			tgButtons[i][j] = tgbotapi.KeyboardButton{Text: button.GetText()}
		}
	}
	replyKb := tgbotapi.NewReplyKeyboard(tgButtons...)
	replyKb.OneTimeKeyboard = kb.IsOneTime()
	replyKb.ResizeKeyboard = true
	return replyKb
}

func getInlineKeyboard(kb *botkb.MessageKeyboard) *tgbotapi.InlineKeyboardMarkup {
	tgButtons := make([][]tgbotapi.InlineKeyboardButton, len(kb.Buttons))
	for i, buttons := range kb.Buttons {
		tgButtons[i] = make([]tgbotapi.InlineKeyboardButton, len(buttons))
		for j, button := range buttons {
			switch btn := button.(type) {
			case botkb.DataButton:
				tgButtons[i][j] = tgbotapi.NewInlineKeyboardButtonData(btn.Text, btn.Data)
			case *botkb.DataButton:
				tgButtons[i][j] = tgbotapi.NewInlineKeyboardButtonData(btn.Text, btn.Data)
			case botkb.UrlButton:
				tgButtons[i][j] = tgbotapi.NewInlineKeyboardButtonURL(btn.Text, btn.URL)
			case *botkb.UrlButton:
				tgButtons[i][j] = tgbotapi.NewInlineKeyboardButtonURL(btn.Text, btn.URL)
			case botkb.SwitchInlineQueryButton:
				tgButtons[i][j] = tgbotapi.NewInlineKeyboardButtonSwitchInlineQuery(btn.Text, btn.Query)
			case *botkb.SwitchInlineQueryButton:
				tgButtons[i][j] = tgbotapi.NewInlineKeyboardButtonSwitchInlineQuery(btn.Text, btn.Query)
			case botkb.SwitchInlineQueryCurrentChatButton:
				tgButtons[i][j] = tgbotapi.NewInlineKeyboardButtonSwitchInlineQueryCurrentChat(btn.Text, btn.Query)
			case *botkb.SwitchInlineQueryCurrentChatButton:
				tgButtons[i][j] = tgbotapi.NewInlineKeyboardButtonSwitchInlineQueryCurrentChat(btn.Text, btn.Query)
			default:
				panic(fmt.Sprintf("Unknown button type at [%d][%d]: %T", i, j, btn))
			}
		}
	}
	return tgbotapi.NewInlineKeyboardMarkup(tgButtons...)
}

func GetTelegramBotAPIClient(ctx context.Context, botContext botsfw.BotContext) *tgbotapi.BotAPI {
	return tgbotapi.NewBotAPIWithClient(
		botContext.BotSettings.Token,
		botContext.BotHost.GetHTTPClient(ctx),
	)
}
