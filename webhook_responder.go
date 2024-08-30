package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/strongo/logus"
	"net/http"
	"strconv"
	"time"
)

type tgWebhookResponder struct {
	w   http.ResponseWriter
	whc *tgWebhookContext
}

var _ botsfw.WebhookResponder = (*tgWebhookResponder)(nil)

func newTgWebhookResponder(w http.ResponseWriter, whc *tgWebhookContext) tgWebhookResponder {
	responder := tgWebhookResponder{w: w, whc: whc}
	whc.responder = responder
	return responder
}

func (r tgWebhookResponder) SendMessage(c context.Context, m botsfw.MessageFromBot, channel botsfw.BotAPISendMessageChannel) (resp botsfw.OnMessageSentResponse, err error) {
	logus.Debugf(c, "tgWebhookResponder.SendMessage(channel=%v, isEdit=%v)\nm: %+v", channel, m.IsEdit, m)
	switch channel {
	case botsfw.BotAPISendMessageOverHTTPS, botsfw.BotAPISendMessageOverResponse:
	// Known channels
	default:
		panic(fmt.Sprintf("Unknown channel: [%v]. Expected either 'https' or 'response'.", channel))
	}
	//ctx := tc.Context()

	var cancel context.CancelFunc
	c, cancel = context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var chattable tgbotapi.Chattable

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
		logus.Debugf(c, botsfw.NoMessageToSend)
		return
	} else if botMessage = m.BotMessage; botMessage != nil {
		logus.Debugf(c, "m.BotMessage != nil")
		switch m.BotMessage.BotMessageType() {
		case botsfw.BotMessageTypeInlineResults:
			chattable = tgbotapi.InlineConfig(m.BotMessage.(InlineBotMessage))
		case botsfw.BotMessageTypeCallbackAnswer:
			callbackAnswer := tgbotapi.AnswerCallbackQueryConfig(m.BotMessage.(CallbackAnswer))
			if callbackAnswer.CallbackQueryID == "" && tgUpdate.CallbackQuery != nil {
				callbackAnswer.CallbackQueryID = tgUpdate.CallbackQuery.ID
			}
			chattable = callbackAnswer
		case botsfw.BotMessageTypeLeaveChat:
			leaveChat := tgbotapi.LeaveChatConfig(m.BotMessage.(LeaveChat))
			if leaveChat.ChatID == "" {
				leaveChat.ChatID = strconv.FormatInt(tgUpdate.Chat().ID, 10)
			}
			chattable = leaveChat
		case botsfw.BotMessageTypeExportChatInviteLink:
			exportChatInviteLink := tgbotapi.ExportChatInviteLink(m.BotMessage.(ExportChatInviteLink))
			if exportChatInviteLink.ChatID == "" {
				exportChatInviteLink.ChatID = strconv.FormatInt(tgUpdate.Chat().ID, 10)
			}
			chattable = exportChatInviteLink
		case botsfw.BotMessageTypeUndefined:
			err = fmt.Errorf("bot message type %v==undefined", m.BotMessage.BotMessageType())
			return
		default:
			err = fmt.Errorf("unknown bot message type %v==%T", m.BotMessage.BotMessageType(), botMessage)
			return
		}
	} else if m.IsEdit || (tgUpdate.CallbackQuery != nil && tgUpdate.CallbackQuery.InlineMessageID != "" && m.ToChat == nil) {
		if m.IsEdit {
			logus.Debugf(c, "m.IsEdit")
		} else if tgUpdate.CallbackQuery != nil {
			logus.Debugf(c, "tgUpdate.CallbackQuery != nil")
		}

		// Edit message
		inlineMessageID, chatID, messageID := getTgMessageIDs(tgUpdate)
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
		logus.Debugf(c, "Edit message => inlineMessageID: %v, chatID: %d, messageID: %d", inlineMessageID, chatID, messageID)
		if inlineMessageID == "" && chatID == 0 && messageID == 0 {
			err = errors.New("Can't edit Telegram message as inlineMessageID is empty && chatID == 0 && messageID == 0")
			return
		}
		if m.Text == "" && m.Keyboard != nil {
			chattable = tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, inlineMessageID, m.Keyboard.(*tgbotapi.InlineKeyboardMarkup))
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
			chattable = editMessageTextConfig
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

		chattable = messageConfig
	} else {
		switch inputType := r.whc.InputType(); inputType {
		case botinput.WebhookInputInlineQuery: // pass
			logus.Debugf(c, "No response to WebhookInputInlineQuery")
		case botinput.WebhookInputChosenInlineResult: // pass
		default:
			mBytes, err := ffjson.Marshal(m)
			if err != nil {
				logus.Errorf(c, "Failed to marshal MessageFromBot to JSON: %v", err)
			}
			inputTypeName := botinput.GetWebhookInputTypeIdNameString(inputType)
			logus.Debugf(c, "Not inline answer, Not inline, Not edit inline, Text is empty. r.whc.InputType(): %v\nMessageFromBot:\n%v", inputTypeName, string(mBytes))
			ffjson.Pool(mBytes)
		}
		return
	}

	jsonStr, err := ffjson.Marshal(chattable)
	if err != nil {
		logus.Errorf(c, "Failed to marshal message config to json: %v\n\tJSON: %v\n\tchattable: %v", err, jsonStr, chattable)
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
	logus.Debugf(c, "Sending to Telegram, Text: %v\n------------------------\nAs JSON: %v", m.Text, indentedJSONStr)

	//if values, err := chattable.Values(); err != nil {
	//	logus.Errorf(c, "Failed to marshal message config to url.Values: %v", err)
	//	return resp, err
	//} else {
	//	logus.Debugf(c, "Message for sending to Telegram as URL values: %v", values)
	//}

	switch channel {
	case botsfw.BotAPISendMessageOverResponse:
		if _, err := tgbotapi.ReplyToResponse(chattable, r.w); err != nil {
			logus.Errorf(c, "Failed to send message to Telegram throw HTTP response: %v", err)
		}
		return resp, err
	case botsfw.BotAPISendMessageOverHTTPS:
		botContext := r.whc.BotContext()
		botAPI := tgbotapi.NewBotAPIWithClient(
			botContext.BotSettings.Token,
			botContext.BotHost.GetHTTPClient(c),
		)
		botAPI.EnableDebug(c)
		message, err := botAPI.Send(chattable)
		if err != nil {
			return resp, err
		} else if message.MessageID != 0 {
			logus.Debugf(c, "Telegram API: MessageID=%v", message.MessageID)
		} else {
			messageJSON, err := ffjson.Marshal(message)
			if err != nil {
				logus.Warningf(c, "Telegram API response as raw: %v", message)
			} else {
				logus.Debugf(c, "Telegram API response as JSON: %v", string(messageJSON))
			}
			ffjson.Pool(messageJSON)
		}
		return botsfw.OnMessageSentResponse{TelegramMessage: message}, nil
	default:
		panic(fmt.Sprintf("Unknown channel: %v", channel))
	}
}
