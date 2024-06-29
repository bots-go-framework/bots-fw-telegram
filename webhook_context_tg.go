package telegram

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
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
	//whi          tgWebhookInput

	// This 3 props are cache for getLocalAndChatIDByChatInstance()
	isInGroup func() bool
	locale    string
	chatID    string
}

func (twhc *tgWebhookContext) NewEditMessage(text string, format botsfw.MessageFormat) (m botsfw.MessageFromBot, err error) {
	m.Text = text
	m.Format = format
	m.IsEdit = true
	return
}

func (twhc *tgWebhookContext) CreateOrUpdateTgChatInstance() (err error) {
	c := twhc.Context()
	logus.Debugf(c, "*tgWebhookContext.CreateOrUpdateTgChatInstance()")
	tgUpdate := twhc.tgInput.TgUpdate()
	if tgUpdate.CallbackQuery == nil {
		logus.Debugf(c, "CreateOrUpdateTgChatInstance() => tgUpdate.CallbackQuery == nil")
		return
	}
	if chatInstanceID := tgUpdate.CallbackQuery.ChatInstance; chatInstanceID == "" {
		logus.Debugf(c, "CreateOrUpdateTgChatInstance() => no chatInstanceID")
	} else {
		chatID := tgUpdate.CallbackQuery.Message.Chat.ID
		logus.Debugf(c, "CreateOrUpdateTgChatInstance() => chatID: %v, chatInstanceID: %v", chatID, chatInstanceID)
		if chatID == 0 {
			return
		}
		tgChatData := twhc.ChatData().(botsfwtgmodels.TgChatData)
		if tgChatData.GetTgChatInstanceID() != chatInstanceID {
			tgChatData.SetTgChatInstanceID(chatInstanceID)
			//if err = twhc.SaveBotChat(c, twhc.GetBotCode(), twhc.MustBotChatID(), tgChatEntity.(botsfw.BotChat)); err != nil {
			//	return
			//}
		}

		var chatInstance botsfwtgmodels.TgChatInstanceData
		preferredLanguage := tgChatData.GetPreferredLanguage()
		logus.Debugf(c, "CreateOrUpdateTgChatInstance() => checking tg chat instance within tx")
		changed := false
		if chatInstance, err = tgChatInstanceDal.GetTelegramChatInstanceByID(c, chatInstanceID); err != nil {
			if !dal.IsNotFound(err) {
				return
			}
			logus.Debugf(c, "CreateOrUpdateTgChatInstance() => new tg chat instance")
			chatInstance = tgChatInstanceDal.NewTelegramChatInstance(chatInstanceID, chatID, preferredLanguage)
			changed = true
		} else { // Update if needed
			logus.Debugf(c, "CreateOrUpdateTgChatInstance() => existing tg chat instance")
			if tgChatInstanceId := chatInstance.GetTgChatID(); tgChatInstanceId != chatID {
				err = fmt.Errorf("chatInstance.GetTgChatID():%d != chatID:%d", tgChatInstanceId, chatID)
			} else if prefLang := chatInstance.GetPreferredLanguage(); prefLang != preferredLanguage {
				chatInstance.SetPreferredLanguage(preferredLanguage)
				changed = true
			}
		}
		if changed {
			logus.Debugf(c, "Saving tg chat instance...")
			if err = tgChatInstanceDal.SaveTelegramChatInstance(c, chatInstanceID, chatInstance); err != nil {
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
) (*tgWebhookContext, error) {
	twhc := &tgWebhookContext{
		tgInput: input,
	}
	chat := twhc.tgInput.TgUpdate().Chat()

	isInGroup := func() bool { // Checks if current chat is a group chat
		if chat != nil && chat.IsGroup() {
			return true
		}

		if callbackQuery := twhc.tgInput.TgUpdate().CallbackQuery; callbackQuery != nil && callbackQuery.ChatInstance != "" {
			c := args.BotContext.BotHost.Context(args.HttpRequest)
			var isGroupChat bool
			chatInstance, err := tgChatInstanceDal.GetTelegramChatInstanceByID(c, callbackQuery.ChatInstance)
			if err != nil {
				if !dal.IsNotFound(err) {
					logus.Errorf(c, "failed to get tg chat instance: %v", err)
				}
				return isGroupChat
			} else if chatInstance != nil {
				isGroupChat = chatInstance.GetTgChatID() < 0
			}
			return isGroupChat
		}

		return false
	}

	whcb, err := botsfw.NewWebhookContextBase(
		args,
		Platform,
		recordsFieldsSetter,
		isInGroup,
		twhc.getLocalAndChatIDByChatInstance,
	)
	twhc.WebhookContextBase = whcb
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

func (twhc *tgWebhookContext) AppUserData() (botsfwmodels.AppUserData, error) {
	appUserID := twhc.AppUserID()
	//appUser := twhc.BotAppContext().NewBotAppUserEntity()
	ctx := twhc.Context()
	tx := twhc.Tx()
	appUser, err := twhc.BotContext().BotSettings.GetAppUserByID(ctx, tx, appUserID)
	return appUser.Data, err
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
		panic(err)
	}
	if botChatID == "" {
		panic(fmt.Sprintf("Not able to send message as BotChatID() returned empty string. text: %v", text))
	}
	botChatIntID, err := strconv.ParseInt(botChatID, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Not able to parse BotChatID(%v) as int: %v", botChatID, err))
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

func (twhc *tgWebhookContext) getLocalAndChatIDByChatInstance(c context.Context) (locale, chatID string, err error) {
	logus.Debugf(c, "*tgWebhookContext.getLocalAndChatIDByChatInstance()")
	if chatID == "" && locale == "" { // we need to cache to make sure not called within transaction
		if cbq := twhc.tgInput.TgUpdate().CallbackQuery; cbq != nil && cbq.ChatInstance != "" {
			if cbq.Message != nil && cbq.Message.Chat != nil && cbq.Message.Chat.ID != 0 {
				logus.Errorf(c, "getLocalAndChatIDByChatInstance() => should not be here")
			} else {
				if chatInstance, err := tgChatInstanceDal.GetTelegramChatInstanceByID(c, cbq.ChatInstance); err != nil {
					if !dal.IsNotFound(err) {
						return "", "", err
					}
				} else if tgChatID := chatInstance.GetTgChatID(); tgChatID != 0 {
					twhc.chatID = strconv.FormatInt(tgChatID, 10)
					twhc.locale = chatInstance.GetPreferredLanguage()
					isInGroup := tgChatID < 0
					twhc.isInGroup = func() bool {
						return isInGroup
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
