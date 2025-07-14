package telegram

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"strconv"
	"time"
)

var _ botinput.InputMessage = (*tgWebhookInput)(nil)

type tgWebhookInput struct {
	update     *tgbotapi.Update
	logRequest func()
}

func (whi tgWebhookInput) BotChatID() (string, error) {
	tgChat := whi.update.Chat()
	if tgChat == nil {
		return "", nil
	}
	return strconv.FormatInt(tgChat.ID, 10), nil
}

func (whi tgWebhookInput) InputType() botinput.Type {
	switch {
	case whi.update.InlineQuery != nil:
		return botinput.TypeInlineQuery

	case whi.update.CallbackQuery != nil:
		return botinput.TypeCallbackQuery

	case whi.update.ChosenInlineResult != nil:
		return botinput.TypeChosenInlineResult

	case whi.update.ChannelPost != nil || whi.update.EditedChannelPost != nil:
		return botinput.TypeNotImplemented

	case whi.update.PreCheckoutQuery != nil:
		return botinput.TypePreCheckoutQuery

	case whi.update.Message.SuccessfulPayment != nil:
		return botinput.TypeSuccessfulPayment

	case whi.update.Message.RefundedPayment != nil:
		return botinput.TypeRefundedPayment

	case whi.update.Message != nil || whi.update.EditedMessage != nil:
		// This should be after any whi.update.Message.* checks
		return botinput.TypeText

	default:
		return botinput.TypeUnknown
	}
}

// TgWebhookInput is a wrapper of telegram update struct to bots framework interface
type TgWebhookInput interface {
	TgUpdate() *tgbotapi.Update
}

func (whi tgWebhookInput) LogRequest() {
	if whi.logRequest != nil {
		whi.logRequest()
	}
}

var _ TgWebhookInput = (*tgWebhookInput)(nil)

// tgWebhookUpdateProvider indicates that input can provide original Telegram update struct
type tgWebhookUpdateProvider interface {
	TgUpdate() *tgbotapi.Update
}

func (whi tgWebhookInput) TgUpdate() *tgbotapi.Update {
	return whi.update
}

var _ botinput.InputMessage = (*tgWebhookTextMessage)(nil)
var _ botinput.InputMessage = (*tgWebhookContactMessage)(nil)
var _ botinput.InputMessage = (*TgWebhookInlineQuery)(nil)
var _ botinput.InputMessage = (*tgWebhookChosenInlineResult)(nil)
var _ botinput.InputMessage = (*TgWebhookCallbackQuery)(nil)
var _ botinput.InputMessage = (*tgWebhookNewChatMembersMessage)(nil)

func (whi tgWebhookInput) GetID() interface{} {
	return whi.update.UpdateID
}

func message2input(input tgWebhookInput, tgMessageType TgMessageType, tgMessage *tgbotapi.Message) botinput.InputMessage {
	switch {
	case tgMessage.Text != "":
		return newTgWebhookTextMessage(input, tgMessageType, tgMessage)
	case tgMessage.Contact != nil:
		return newTgWebhookContact(input)
	case tgMessage.NewChatMembers != nil:
		return newTgWebhookNewChatMembersMessage(input)
	case tgMessage.LeftChatMember != nil:
		return newTgWebhookLeftChatMembersMessage(input)
	case tgMessage.Voice != nil:
		return newTgWebhookVoiceMessage(input, tgMessageType, tgMessage)
	case tgMessage.Photo != nil:
		return newTgWebhookPhotoMessage(input, tgMessageType, tgMessage)
	case tgMessage.Audio != nil:
		return newTgWebhookAudioMessage(input, tgMessageType, tgMessage)
	case tgMessage.Sticker != nil:
		return newTgWebhookStickerMessage(input, tgMessageType, tgMessage)
	case tgMessage.UsersShared != nil:
		return newTgWebhookUsersSharedMessage(input, tgMessageType, tgMessage)
	default:
		return nil
	}
}

// NewTelegramWebhookInput maps telegram update struct to bots framework interface
func NewTelegramWebhookInput(update *tgbotapi.Update, logRequest func()) (botinput.InputMessage, error) {
	input := tgWebhookInput{update: update, logRequest: logRequest}

	switch inputType := input.InputType(); inputType {
	case botinput.TypeInlineQuery:
		return newTelegramWebhookInlineQuery(input), nil
	case botinput.TypeCallbackQuery:
		return newTelegramWebhookCallbackQuery(input), nil
	case botinput.TypeChosenInlineResult:
		return newTelegramWebhookChosenInlineResult(input), nil
	case botinput.TypePreCheckoutQuery:
		return newTgWebhookPreCheckoutQuery(input), nil
	case botinput.TypeSuccessfulPayment:
		return newTgWebhookSuccessfulPayment(input), nil
	case botinput.TypeRefundedPayment:
		return newTgWebhookRefundedPayment(input), nil
	case botinput.TypeText:
		switch {
		case update.Message != nil:
			return message2input(input, TgMessageTypeRegular, update.Message), nil

		case update.EditedMessage != nil:
			return message2input(input, TgMessageTypeEdited, update.EditedMessage), nil

		}
	case botinput.TypeNotImplemented:
		switch {

		case update.ChannelPost != nil:
			channelPost, err := encodeToJsonString(update.ChannelPost)
			if err != nil {
				panic(err)
			}
			return nil, fmt.Errorf("the ChannelPost is not supported at the moment: [%s]: %w", channelPost, botsfw.ErrNotImplemented)

		case update.EditedChannelPost != nil:

			editedChannelPost, err := encodeToJsonString(update.EditedChannelPost)
			if err != nil {
				panic(err)
			}
			return nil, fmt.Errorf("the EditedChannelPost is not supported at the moment: [%s]: %w", editedChannelPost, botsfw.ErrNotImplemented)
		}
	default:
		return nil, fmt.Errorf("%w: %v", botsfw.ErrNotImplemented, inputType)
	}
	return nil, botsfw.ErrNotImplemented
}

func (whi tgWebhookInput) GetSender() botinput.User {
	switch {
	case whi.update.Message != nil:
		return tgWebhookUser{tgUser: whi.update.Message.From}
	case whi.update.EditedMessage != nil:
		return tgWebhookUser{tgUser: whi.update.EditedMessage.From}
	case whi.update.CallbackQuery != nil:
		return tgWebhookUser{tgUser: whi.update.CallbackQuery.From}
	case whi.update.InlineQuery != nil:
		return tgWebhookUser{tgUser: whi.update.InlineQuery.From}
	case whi.update.ChosenInlineResult != nil:
		return tgWebhookUser{tgUser: whi.update.ChosenInlineResult.From}
	case whi.update.PreCheckoutQuery != nil:
		return tgWebhookUser{tgUser: whi.update.PreCheckoutQuery.From}
	//case whi.update.ChannelPost != nil:
	//	chat := whi.update.ChannelPost.Chat
	//	return tgWebhookUser{  // TODO: Seems to be dirty hack.
	//		tgUser: &tgbotapi.User{
	//			ID: int(chat.ID),
	//			Name: chat.Name,
	//			FirstName: chat.FirstName,
	//			LastName: chat.LastName,
	//		},
	//	}
	default:
		panic("Unknown From sender")
	}
}

func (whi tgWebhookInput) GetRecipient() botinput.Recipient {
	panic("Not implemented")
}

func (whi tgWebhookInput) GetTime() time.Time {
	if whi.update.Message != nil {
		return whi.update.Message.Time()
	}
	if whi.update.EditedMessage != nil {
		return whi.update.EditedMessage.Time()
	}
	return time.Time{}
}

func (whi tgWebhookInput) StringID() string {
	return ""
}

func (whi tgWebhookInput) TelegramChatID() int64 {
	if whi.update.Message != nil {
		return whi.update.Message.Chat.ID
	}
	if whi.update.EditedMessage != nil {
		return whi.update.EditedMessage.Chat.ID
	}
	panic("Can't get Telegram chat ID from `update.Message` or `update.EditedMessage`.")
}

func (whi tgWebhookInput) Chat() botinput.Chat {
	update := whi.update
	if update.Message != nil {
		return TgWebhookChat{
			chat: update.Message.Chat,
		}
	} else if update.EditedMessage != nil {
		return TgWebhookChat{
			chat: update.EditedMessage.Chat,
		}
	} else if callbackQuery := update.CallbackQuery; callbackQuery != nil && callbackQuery.Message != nil {
		return TgWebhookChat{
			chat: callbackQuery.Message.Chat,
		}
	}
	return nil
}
