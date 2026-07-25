package telegram

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/bots-go-framework/bots-fw/botmsg"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/strongo/logus"
	"net/http"
	"strconv"
)

var _ botsfw.WebhookContext = (*tgWebhookContext)(nil)

type tgWebhookContext struct {
	*botsfw.WebhookContextBase
	tgInput TgWebhookInput
	//update         tgbotapi.Update // TODO: Consider removing?
	//responseWriter http.ResponseWriter
	responder botsfw.WebhookResponder
	//whi          tgInput

	// This 3 props are cache for getLocalAndChatIDByChatInstance()
	isInGroup     func() bool
	locale        string
	chatID        string
	chatInstances ChatInstanceStore
}

func (twhc *tgWebhookContext) NewEditMessage(text string, format botmsg.Format) (m botmsg.MessageFromBot, err error) {
	m.Text = text
	m.Format = format
	m.IsEdit = true
	return
}

func (twhc *tgWebhookContext) CreateOrUpdateTgChatInstance() (err error) {
	ctx := twhc.Context()
	logus.Debugf(ctx, "*tgWebhookContext.CreateOrUpdateTgChatInstance()")
	tgUpdate := twhc.tgInput.TgUpdate()
	if tgUpdate.CallbackQuery == nil {
		logus.Debugf(ctx, "CreateOrUpdateTgChatInstance() => tgUpdate.CallbackQuery == nil")
		return
	}
	if chatInstanceID := tgUpdate.CallbackQuery.ChatInstance; chatInstanceID == "" {
		logus.Debugf(ctx, "CreateOrUpdateTgChatInstance() => no chatInstanceID")
	} else {
		chatID := tgUpdate.CallbackQuery.Message.Chat.ID
		logus.Debugf(ctx, "CreateOrUpdateTgChatInstance() => callback has chat and chat-instance identifiers")
		if chatID == 0 {
			return
		}
		tgChatData := twhc.ChatData().(botsfwtgmodels.TgChatData)
		if tgChatData.GetTgChatInstanceID() != chatInstanceID {
			tgChatData.SetTgChatInstanceID(chatInstanceID)
			//if err = twhc.SaveBotChat(ctx, twhc.GetBotCode(), twhc.MustBotChatID(), tgChatEntity.(botsfw.BotChat)); err != nil {
			//	return
			//}
		}

		var chatInstanceData botsfwtgmodels.TgChatInstanceData
		preferredLanguage := tgChatData.GetPreferredLanguage()
		logus.Debugf(ctx, "CreateOrUpdateTgChatInstance() => checking chat-instance store")
		changed := false
		botCode := twhc.GetBotCode()
		var found bool
		if chatInstanceData, found, err = twhc.chatInstances.Get(ctx, botCode, chatInstanceID); err != nil {
			return
		} else if !found {
			logus.Debugf(ctx, "CreateOrUpdateTgChatInstance() => new tg chat instance")
			chatInstanceData = NewTelegramChatInstance(chatInstanceID, chatID, preferredLanguage)
			changed = true
		} else { // Update if needed
			logus.Debugf(ctx, "CreateOrUpdateTgChatInstance() => existing tg chat instance")
			if tgChatInstanceId := chatInstanceData.GetTgChatID(); tgChatInstanceId != chatID {
				err = fmt.Errorf("stored Telegram chat instance does not match callback chat")
			} else if prefLang := chatInstanceData.GetPreferredLanguage(); prefLang != preferredLanguage {
				chatInstanceData.SetPreferredLanguage(preferredLanguage)
				changed = true
			}
		}
		if changed {
			logus.Debugf(ctx, "Saving tg chat instance...")
			if err = twhc.chatInstances.Save(ctx, botCode, chatInstanceID, chatInstanceData); err != nil {
				return
			}
		}
		return
	}
	return
}

func getTgMessageIDs(update *tgbotapi.Update) (inlineMessageID string, chatID int64, messageID int) {
	if update.CallbackQuery != nil {
		if update.CallbackQuery.InlineMessageID != "" {
			inlineMessageID = update.CallbackQuery.InlineMessageID
		} else if update.CallbackQuery.Message != nil {
			messageID = update.CallbackQuery.Message.MessageID
			chatID = update.CallbackQuery.Message.Chat.ID
		}
	} else if update.Message != nil {
		messageID = update.Message.MessageID
		chatID = update.Message.Chat.ID
	} else if update.EditedMessage != nil {
		messageID = update.EditedMessage.MessageID
		chatID = update.EditedMessage.Chat.ID
	} else if update.ChannelPost != nil {
		messageID = update.ChannelPost.MessageID
		chatID = update.ChannelPost.Chat.ID
	} else if update.ChosenInlineResult != nil {
		if update.ChosenInlineResult.InlineMessageID != "" {
			inlineMessageID = update.ChosenInlineResult.InlineMessageID
		}
	} else if update.EditedChannelPost != nil {
		messageID = update.EditedChannelPost.MessageID
		chatID = update.EditedChannelPost.Chat.ID
	}

	return
}

func newTelegramWebhookContext(
	args botsfw.CreateWebhookContextArgs,
	input TgWebhookInput,
	recordsFieldsSetter botsfw.BotRecordsFieldsSetter,
	chatInstances ChatInstanceStore,
) (twhc *tgWebhookContext, err error) {
	if chatInstances == nil {
		return nil, fmt.Errorf("chat instance store is required")
	}
	twhc = &tgWebhookContext{tgInput: input, chatInstances: chatInstances}

	chat := twhc.tgInput.TgUpdate().Chat()

	getIsInGroup := func() (isInGroup bool, err error) { // Checks if current chat is a group chat
		if chat != nil && chat.IsGroup() {
			return true, nil
		}

		if callbackQuery := twhc.tgInput.TgUpdate().CallbackQuery; callbackQuery != nil && callbackQuery.ChatInstance != "" {
			c := args.BotContext.BotHost.Context(args.HttpRequest)
			var isGroupChat bool
			var chatInstance botsfwtgmodels.TgChatInstanceData
			botCode := twhc.GetBotCode()
			var found bool
			if chatInstance, found, err = chatInstances.Get(c, botCode, callbackQuery.ChatInstance); err != nil {
				logus.Errorf(c, "failed to get tg chat instance: %v", err)
				return isGroupChat, err
			} else if found && chatInstance != nil {
				isGroupChat = chatInstance.GetTgChatID() < 0
			}
			return isGroupChat, err
		}

		return false, err
	}

	twhc.WebhookContextBase, err = botsfw.NewWebhookContextBase(
		args,
		Platform,
		recordsFieldsSetter,
		getIsInGroup,
		twhc.getLocalAndChatIDByChatInstance,
	)
	return twhc, err
}

func (twhc *tgWebhookContext) Close(context.Context) error {
	return nil
}

func (twhc *tgWebhookContext) Responder() botsfw.WebhookResponder {
	return twhc.responder
}

//type tgBotAPIUser struct {
//	user tgbotapi.User
//}
//
//func (tc tgBotAPIUser) FirstName() string {
//	return tc.user.FirstName
//}
//
//func (tc tgBotAPIUser) LastName() string {
//	return tc.user.LastName
//}

//func (tc tgBotAPIUser) IdAsString() string {
//	return ""
//}

//func (tc tgBotAPIUser) IdAsInt64() int64 {
//	return int64(tc.user.ID)
//}

func (twhc *tgWebhookContext) Init(http.ResponseWriter, *http.Request) error {
	return nil
}

func (twhc *tgWebhookContext) BotAPI() *tgbotapi.BotAPI {
	botContext := twhc.BotContext()
	return tgbotapi.NewBotAPIWithClient(botContext.BotSettings.Token, botContext.BotHost.GetHTTPClient(twhc.Context()))
}

func (twhc *tgWebhookContext) IsNewerThen( /*chatEntity*/ data botsfwmodels.BotChatData) bool {
	return true
	//if telegramChat, ok := whc.Data().(*TgChatBaseData); ok && telegramChat != nil {
	//	return whc.Input().whi.update.UpdateID > telegramChat.LastProcessedUpdateID
	//}
	//return false
}

//func (twhc *tgWebhookContext) getTelegramSenderID() int {
//	senderID := twhc.Input().GetSender().GetID()
//	if tgUserID, ok := senderID.(int); ok {
//		return tgUserID
//	}
//	panic("int expected")
//}

func (twhc *tgWebhookContext) NewTgMessage(text string) *tgbotapi.MessageConfig {
	//inputMessage := tc.InputMessage()
	//if inputMessage != nil {
	//ctx := tc.Context()
	//Data := inputMessage.TgChat()
	//chatID := Data.GetID()
	//logus.Infof(ctx, "NewTgMessage(): tc.update.Message.TgChat.ID: %v", chatID)
	botChatID, err := twhc.BotChatID()
	if err != nil {
		panic(fmt.Sprintf("failed to resolve Telegram BotChatID: %T", err))
	}
	if botChatID == "" {
		panic("not able to send Telegram message because BotChatID() returned an empty string")
	}
	botChatIntID, err := strconv.ParseInt(botChatID, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("not able to parse Telegram BotChatID as integer: %T", err))
	}
	//tgbotapi.NewEditMessageText()
	return tgbotapi.NewMessage(botChatIntID, text)
}

func (twhc *tgWebhookContext) UpdateLastProcessed( /*chatEntity*/ data botsfwmodels.BotChatData) error {
	return nil
	//if telegramChat, ok := chatEntity.(*TgChatBaseData); ok {
	//	telegramChat.LastProcessedUpdateID = tc.whi.update.UpdateID
	//	return nil
	//}
	//return fmt.Errorf("Expected *TgChatBaseData, got: %T", chatEntity)
}

func (twhc *tgWebhookContext) getLocalAndChatIDByChatInstance(ctx context.Context) (locale, chatID string, err error) {
	logus.Debugf(ctx, "*tgWebhookContext.getLocalAndChatIDByChatInstance()")
	if twhc.chatID == "" && twhc.locale == "" { // we need to cache to make sure not called within transaction
		if cbq := twhc.tgInput.TgUpdate().CallbackQuery; cbq != nil {
			if chatInstanceID := cbq.ChatInstance; chatInstanceID != "" {
				if cbq.Message != nil && cbq.Message.Chat != nil && cbq.Message.Chat.ID != 0 {
					logus.Errorf(ctx, "getLocalAndChatIDByChatInstance() => should not be here")
				} else {
					var chatInstanceData botsfwtgmodels.TgChatInstanceData
					botCode := twhc.GetBotCode()
					var found bool
					if chatInstanceData, found, err = twhc.chatInstances.Get(ctx, botCode, chatInstanceID); err != nil {
						return "", "", err
					} else if found && chatInstanceData != nil && chatInstanceData.GetTgChatID() != 0 {
						tgChatID := chatInstanceData.GetTgChatID()
						twhc.chatID = strconv.FormatInt(tgChatID, 10)
						twhc.locale = chatInstanceData.GetPreferredLanguage()
						isInGroup := tgChatID < 0
						twhc.isInGroup = func() bool {
							return isInGroup
						}
					}
				}
			}
		}
	}
	return twhc.locale, twhc.chatID, nil
}

//func (twhc *tgWebhookContext) ChatData() botsfwmodels.ChatData {
//	if _, err := twhc.BotChatID(); err != nil {
//		logus.Errorf(twhc.Context(), fmt.Errorf("whc.BotChatID(): %w", err).Error())
//		return nil
//	}
//	//tgUpdate := twhc.tgInput.TgUpdate()
//	//if tgUpdate.CallbackQuery != nil {
//	//
//	//}
//
//	return twhc.WebhookContextBase.ChatData()
//}
