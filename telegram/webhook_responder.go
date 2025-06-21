package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/strongo/logus"
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
		var chat botinput.WebhookChat
		if inputWithChat, ok := input.(interface{ Chat() botinput.WebhookChat }); ok {
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

func (r tgWebhookResponder) SendMessage(ctx context.Context, m botsfw.MessageFromBot, channel botsfw.BotAPISendMessageChannel) (resp botsfw.OnMessageSentResponse, err error) {
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
		case botsfw.MessageFormatHTML:
			return "html"
		case botsfw.MessageFormatMarkdown:
			return "markdown"
		case botsfw.MessageFormatText:
			return ""
		default:
			panic(fmt.Sprintf("Unknown message parse_mode value: %d", m.Format))
		}
	}

	tgUpdate := r.whc.Input().(tgWebhookUpdateProvider).TgUpdate()

	var botMessage botsfw.BotMessage

	if m.Text == botsfw.NoMessageToSend {
		logus.Debugf(ctx, botsfw.NoMessageToSend)
		return
	} else if botMessage = m.BotMessage; botMessage != nil {
		logus.Debugf(ctx, "m.BotMessage != nil")
		switch m.BotMessage.BotMessageType() {
		case botsfw.BotMessageTypeInlineResults:
			sendable = tgbotapi.InlineConfig(m.BotMessage.(InlineBotMessage))
		case botsfw.BotMessageTypeCallbackAnswer:
			callbackAnswer := tgbotapi.AnswerCallbackQueryConfig(m.BotMessage.(CallbackAnswer))
			if callbackAnswer.CallbackQueryID == "" && tgUpdate.CallbackQuery != nil {
				callbackAnswer.CallbackQueryID = tgUpdate.CallbackQuery.ID
			}
			sendable = callbackAnswer
		case botsfw.BotMessageTypeLeaveChat:
			leaveChat := tgbotapi.LeaveChatConfig(m.BotMessage.(LeaveChat))
			if leaveChat.ChatID == "" {
				leaveChat.ChatID = strconv.FormatInt(tgUpdate.Chat().ID, 10)
			}
			sendable = leaveChat
		case botsfw.BotMessageTypeExportChatInviteLink:
			exportChatInviteLink := tgbotapi.ExportChatInviteLink(m.BotMessage.(ExportChatInviteLink))
			if exportChatInviteLink.ChatID == "" {
				exportChatInviteLink.ChatID = strconv.FormatInt(tgUpdate.Chat().ID, 10)
			}
			sendable = exportChatInviteLink
		case botsfw.BotMessageTypeUndefined:
			err = fmt.Errorf("bot message type %v==undefined", m.BotMessage.BotMessageType())
			return
		case botsfw.BotMessageTypeSendInvoice:
			invoiceConfig := tgbotapi.InvoiceConfig(m.BotMessage.(Invoice))
			if invoiceConfig.ChatID == 0 {
				invoiceConfig.ChatID = tgUpdate.Chat().ID
			}
			sendable = &invoiceConfig
		case botsfw.BotMessageTypeSetDescription:
			setBotDescription := m.BotMessage.(SetBotDescription)
			sendable = (tgbotapi.SetMyDescription)(setBotDescription)
		case botsfw.BotMessageTypeSetShortDescription:
			setBotDescription := m.BotMessage.(SetBotShortDescription)
			sendable = (tgbotapi.SetMyShortDescription)(setBotDescription)
		case botsfw.BotMessageTypeSetCommands:
			setBotDescription := m.BotMessage.(SetBotCommands)
			sendable = (tgbotapi.SetMyCommandsConfig)(setBotDescription)
		case botsfw.BotMessageTypeAnswerPreCheckoutQuery:
			answerPreCheckoutQuery := m.BotMessage.(PreCheckoutQueryAnswer)
			sendable = (tgbotapi.AnswerPreCheckoutQueryConfig)(answerPreCheckoutQuery)
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
			sendable = tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, inlineMessageID, m.Keyboard.(*tgbotapi.InlineKeyboardMarkup))
		} else if m.Text != "" {
			editMessageTextConfig := tgbotapi.NewEditMessageText(chatID, messageID, inlineMessageID, m.Text)
			editMessageTextConfig.ParseMode = parseMode()
			editMessageTextConfig.DisableWebPagePreview = m.DisableWebPagePreview
			if m.Keyboard != nil {
				switch keyboard := m.Keyboard.(type) {
				case *tgbotapi.InlineKeyboardMarkup:
					editMessageTextConfig.ReplyMarkup = keyboard
				//case tgbotapi.ForceReply:
				//	editMessageTextConfig.ReplyMarkup = keyboard
				default:
					panic(fmt.Sprintf("m.Keyboard has unsupported type %T", m.Keyboard))
				}
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
			messageConfig.ChatID = int64(m.ToChat.(botsfw.ChatIntID))
		}
		messageConfig.DisableWebPagePreview = m.DisableWebPagePreview
		messageConfig.DisableNotification = m.DisableNotification
		if m.Keyboard != nil {
			messageConfig.ReplyMarkup = m.Keyboard
		}

		messageConfig.ParseMode = parseMode()

		sendable = messageConfig
	} else {
		switch inputType := r.whc.InputType(); inputType {
		case botinput.WebhookInputInlineQuery: // pass
			logus.Debugf(ctx, "No response to WebhookInputInlineQuery")
		case botinput.WebhookInputChosenInlineResult: // pass
		default:
			var mBytes []byte
			mBytes, err = ffjson.Marshal(m)
			if err != nil {
				logus.Errorf(ctx, "Failed to marshal MessageFromBot to JSON: %v", err)
			}
			inputTypeName := botinput.GetWebhookInputTypeIdNameString(inputType)
			logus.Debugf(ctx, "Not inline answer, Not inline, Not edit inline, Text is empty. r.whc.InputType(): %v\nMessageFromBot:\n%v", inputTypeName, string(mBytes))
			ffjson.Pool(mBytes)
		}
		return
	}

	var jsonStr []byte
	jsonStr, err = ffjson.Marshal(sendable)
	if err != nil {
		logus.Errorf(ctx, "Failed to marshal message config to json: %v\n\tJSON: %v\n\tsendable: %v", err, jsonStr, sendable)
		ffjson.Pool(jsonStr)
		return resp, err
	}
	var indentedJSON bytes.Buffer
	var indentedJSONStr string
	if indentedErr := json.Indent(&indentedJSON, jsonStr, "", "\t"); indentedErr == nil {
		indentedJSONStr = indentedJSON.String()
	} else {
		indentedJSONStr = string(jsonStr)
	}
	ffjson.Pool(jsonStr)
	logus.Debugf(ctx, "Sending to Telegram, Text: %v\n------------------------\nAs JSON: %v", m.Text, indentedJSONStr)

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
		var messageJSON []byte
		if messageJSON, err = ffjson.Marshal(message); err != nil {
			logus.Warningf(ctx, "Telegram API response as raw: %v", message)
		} else {
			logus.Debugf(ctx, "Telegram API response as JSON: %v", string(messageJSON))
		}
		ffjson.Pool(messageJSON)
	}
	return
}

func GetTelegramBotAPIClient(ctx context.Context, botContext botsfw.BotContext) *tgbotapi.BotAPI {
	return tgbotapi.NewBotAPIWithClient(
		botContext.BotSettings.Token,
		botContext.BotHost.GetHTTPClient(ctx),
	)
}
