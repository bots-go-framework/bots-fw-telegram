package telegram

import (
	"bytes"
	"context"
	"encoding/json"
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

func (r tgWebhookResponder) DeleteMessage(ctx context.Context, messageID string) (err error) {
	var msgID int
	if msgID, err = strconv.Atoi(messageID); err != nil {
		err = fmt.Errorf("failed to parse messageID='%s' as int: %w", messageID, err)
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
			case TgWebhookCallbackQuery:
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
	botAPI.EnableDebug(ctx)
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
	logus.Debugf(ctx, "tgWebhookResponder.SendMessage(channel=%v, isEdit=%v)\nm: %+v", channel, m.IsEdit, m)
	channel = botsfw.BotAPISendMessageOverHTTPS
	switch channel {
	case botsfw.BotAPISendMessageOverHTTPS, botsfw.BotAPISendMessageOverResponse:
	// Known channels
	default:
		panic(fmt.Sprintf("Unknown channel: [%v]. Expected either 'https' or 'response'.", channel))
	}
	//ctx := tc.Context()

	var sendable tgbotapi.Sendable

	parseMode := func() string {
		switch m.Format {
		case botmsg.FormatHTML:
			return "HTML"
		case botmsg.FormatMarkdown:
			return "MarkdownV2"
		case botmsg.FormatText:
			return ""
		default:
			panic(fmt.Sprintf("Unknown message parse_mode value: %d", m.Format))
		}
	}

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
			callbackAnswer := tgbotapi.AnswerCallbackQueryConfig(m.BotMessage.(CallbackAnswer))
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
				photoConfig.ParseMode = parseMode()
			}
			sendable = (tgbotapi.PhotoConfig)(photoConfig)
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
			switch m.EditMessageUID.(type) { // TODO: How do we remove duplicates for value & pointer cases?
			case callbackCurrent:
				// do nothing
			case InlineMessageUID:
				inlineMessageID = m.EditMessageUID.(InlineMessageUID).InlineMessageID
				chatID = 0
				messageID = 0
			case *InlineMessageUID:
				inlineMessageID = m.EditMessageUID.(*InlineMessageUID).InlineMessageID
				chatID = 0
				messageID = 0
			case ChatMessageUID:
				chatMessageUID := m.EditMessageUID.(ChatMessageUID)
				inlineMessageID = ""
				if chatMessageUID.ChatID != 0 {
					chatID = chatMessageUID.ChatID
				}
				if chatMessageUID.MessageID != 0 {
					messageID = chatMessageUID.MessageID
				}
			case *ChatMessageUID:
				chatMessageUID := m.EditMessageUID.(*ChatMessageUID)
				inlineMessageID = ""
				if chatMessageUID.ChatID != 0 {
					chatID = chatMessageUID.ChatID
				}
				if chatMessageUID.MessageID != 0 {
					messageID = chatMessageUID.MessageID
				}
			default:
				err = fmt.Errorf("unknown EditMessageUID type %T(%v)", m.EditMessageUID, m.EditMessageUID)
				return
			}
		}
		logus.Debugf(ctx, "Edit message => inlineMessageID: %v, chatID: %d, messageID: %d", inlineMessageID, chatID, messageID)
		if inlineMessageID == "" && chatID == 0 && messageID == 0 {
			err = errors.New("can't edit Telegram message as inlineMessageID is empty && chatID == 0 && messageID == 0")
			return
		}
		if m.Text == "" && m.Keyboard != nil {
			switch keyboard := m.Keyboard.(type) {
			case *tgbotapi.InlineKeyboardMarkup:
				sendable = tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, inlineMessageID, keyboard)
			case *botkb.MessageKeyboard:
				switch keyboard.KeyboardType() {
				case botkb.KeyboardTypeInline:
					kb := getTelegramInlineKeyboard(m.Keyboard)
					sendable = tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, inlineMessageID, tgbotapi.NewInlineKeyboardMarkup(kb.InlineKeyboard...))
				default:
					err = fmt.Errorf("unknown keyboard type %T(%v)", keyboard.KeyboardType(), keyboard)
					return
				}
			}
		} else if m.Text != "" {
			editMessageTextConfig := tgbotapi.NewEditMessageText(chatID, messageID, inlineMessageID, m.Text)
			editMessageTextConfig.ParseMode = parseMode()
			editMessageTextConfig.DisableWebPagePreview = m.DisableWebPagePreview
			if m.Keyboard != nil {
				editMessageTextConfig.ReplyMarkup = getTelegramInlineKeyboard(m.Keyboard)
			}
			sendable = editMessageTextConfig
		} else {
			err = fmt.Errorf("can't edit telegram message as got unknown output: %v", m)
			panic(err)
			// return
		}
	} else if m.Text != "" {
		messageConfig := r.whc.NewTgMessage(m.Text)
		if m.ToChat != nil {
			messageConfig.ChatID = int64(m.ToChat.(botmsg.ChatIntID))
		}
		messageConfig.DisableWebPagePreview = m.DisableWebPagePreview
		messageConfig.DisableNotification = m.DisableNotification
		if m.Keyboard != nil {
			messageConfig.ReplyMarkup = getTelegramInlineKeyboard(m.Keyboard)
		}

		messageConfig.ParseMode = parseMode()

		sendable = messageConfig
	} else {
		switch inputType := r.whc.InputType(); inputType {
		case botinput.TypeInlineQuery: // pass
			logus.Debugf(ctx, "No response to WebhookInputInlineQuery")
		case botinput.TypeChosenInlineResult: // pass
		default:
			var mJson string
			if mJson, err = encodeToJsonString(m); err != nil {
				logus.Errorf(ctx, "Failed to marshal MessageFromBot to JSON: %v", err)
			} else {
				inputTypeName := inputType.String()
				logus.Debugf(ctx, "Not inline answer, Not inline, Not edit inline, Text is empty. r.whc.InputType(): %v\nMessageFromBot:\n%v", inputTypeName, mJson)
			}
		}
		return
	}

	var sendableStr string

	if sendableStr, err = encodeToJsonString(sendable); err != nil {
		logus.Errorf(ctx, "Failed to marshal message config to json: %v\n\tsendable: %v", err, sendable)
		return resp, err
	}
	logus.Debugf(ctx, "Sending to Telegram, Text: %v\n------------------------\nAs JSON: %v", m.Text, sendableStr)

	//if values, err := sendable.Values(); err != nil {
	//	logus.Errorf(ctx, "Failed to marshal message config to url.Values: %v", err)
	//	return resp, err
	//} else {
	//	logus.Debugf(ctx, "Message for sending to Telegram as URL values: %v", values)
	//}

	switch channel {
	case botsfw.BotAPISendMessageOverResponse:
		if _, err = tgbotapi.ReplyToResponse(sendable, r.w); err != nil {
			logus.Errorf(ctx, "Failed to send message to Telegram via HTTP response: %v", err)
		}
		return resp, err
	case botsfw.BotAPISendMessageOverHTTPS:
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		var message tgbotapi.Message
		message, err = r.sendOverHttps(ctx, sendable)
		return botsfw.OnMessageSentResponse{TelegramMessage: message}, err
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
	botAPI.EnableDebug(ctx)

	if message, err = botAPI.Send(chattable); err != nil {
		return
	} else if message.MessageID != 0 {
		logus.Debugf(ctx, "Telegram API: MessageID=%v", message.MessageID)
	} else {
		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		encoder.SetIndent("", "\t")
		encoder.SetEscapeHTML(false)
		if err = encoder.Encode(chattable); err != nil {
			logus.Warningf(ctx, "Telegram API response as raw: %v", message)
		} else {
			logus.Debugf(ctx, "Telegram API response as JSON: %v", string(buf.String()))
		}
	}
	return
}

func getTelegramInlineKeyboard(keyboard botkb.Keyboard) *tgbotapi.InlineKeyboardMarkup {
	switch kb := keyboard.(type) {
	case *tgbotapi.InlineKeyboardMarkup:
		return kb
	case *botkb.MessageKeyboard:
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
	default:
		panic(fmt.Sprintf("m.Keyboard has unsupported type %v", kb.KeyboardType()))
	}
}

func GetTelegramBotAPIClient(ctx context.Context, botContext botsfw.BotContext) *tgbotapi.BotAPI {
	return tgbotapi.NewBotAPIWithClient(
		botContext.BotSettings.Token,
		botContext.BotHost.GetHTTPClient(ctx),
	)
}
