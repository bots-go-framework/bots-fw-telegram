package telegram

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/pquerna/ffjson/ffjson"
	"time"
)

type tgWebhookInput struct {
	update     *tgbotapi.Update // TODO: Make a pointer?
	logRequest func()
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

var _ botsfw.WebhookInput = (*tgWebhookTextMessage)(nil)
var _ botsfw.WebhookInput = (*tgWebhookContactMessage)(nil)
var _ botsfw.WebhookInput = (*TgWebhookInlineQuery)(nil)
var _ botsfw.WebhookInput = (*tgWebhookChosenInlineResult)(nil)
var _ botsfw.WebhookInput = (*TgWebhookCallbackQuery)(nil)
var _ botsfw.WebhookInput = (*tgWebhookNewChatMembersMessage)(nil)

func (whi tgWebhookInput) GetID() interface{} {
	return whi.update.UpdateID
}

func message2input(input tgWebhookInput, tgMessageType TgMessageType, tgMessage *tgbotapi.Message) botsfw.WebhookInput {
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
	default:
		return nil
	}
}

// NewTelegramWebhookInput maps telegram update struct to bots framework interface
func NewTelegramWebhookInput(update *tgbotapi.Update, logRequest func()) (botsfw.WebhookInput, error) {
	input := tgWebhookInput{update: update, logRequest: logRequest}

	switch {

	case update.InlineQuery != nil:
		return newTelegramWebhookInlineQuery(input), nil

	case update.CallbackQuery != nil:
		return newTelegramWebhookCallbackQuery(input), nil

	case update.ChosenInlineResult != nil:
		return newTelegramWebhookChosenInlineResult(input), nil

	default:

		switch {

		case update.Message != nil:
			return message2input(input, TgMessageTypeRegular, update.Message), nil

		case update.EditedMessage != nil:
			return message2input(input, TgMessageTypeEdited, update.EditedMessage), nil

		case update.ChannelPost != nil:
			channelPost, _ := ffjson.MarshalFast(update.ChannelPost)
			return nil, fmt.Errorf("the ChannelPost is not supported at the moment: [%s]: %w", channelPost, botsfw.ErrNotImplemented)
			//return message2input(TgMessageTypeChannelPost, update.ChannelPost), nil

		case update.EditedChannelPost != nil:
			editedChannelPost, _ := ffjson.MarshalFast(update.EditedChannelPost)
			return nil, fmt.Errorf("the EditedChannelPost is not supported at the moment: [%s]: %w", string(editedChannelPost), botsfw.ErrNotImplemented)
			//	return message2input(TgMessageTypeEditedChannelPost, update.EditedChannelPost), nil

		default:
			return nil, botsfw.ErrNotImplemented

		}
	}
}

func (whi tgWebhookInput) GetSender() botsfw.WebhookSender {
	switch {
	case whi.update.Message != nil:
		return tgSender{tgUser: whi.update.Message.From}
	case whi.update.EditedMessage != nil:
		return tgSender{tgUser: whi.update.EditedMessage.From}
	case whi.update.CallbackQuery != nil:
		return tgSender{tgUser: whi.update.CallbackQuery.From}
	case whi.update.InlineQuery != nil:
		return tgSender{tgUser: whi.update.InlineQuery.From}
	case whi.update.ChosenInlineResult != nil:
		return tgSender{tgUser: whi.update.ChosenInlineResult.From}
	//case whi.update.ChannelPost != nil:
	//	chat := whi.update.ChannelPost.Chat
	//	return tgSender{  // TODO: Seems to be dirty hack.
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

func (whi tgWebhookInput) GetRecipient() botsfw.WebhookRecipient {
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

func (whi tgWebhookInput) Chat() botsfw.WebhookChat {
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
