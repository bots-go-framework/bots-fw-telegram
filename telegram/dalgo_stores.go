package telegram

//func NewDalgoStores(db dalgo4botsfw.DbProvider) (botsfw.BotChatStore, botsfw.BotUserStore) {
//	return newDalgoBotChatStore(db), newDalgoBotUserStore(db)
//}
//
//func newDalgoBotChatStore(db dalgo4botsfw.DbProvider) botsfw.BotChatStore {
//	newChatData := func() botsfw.BotChat {
//		return new(botsfwtgmodels.TgChatBase)
//	}
//	return dalgo4botsfw.NewBotChatStore(botsfwtgmodels.TgChatCollection, db, newChatData)
//}
//
//func newDalgoBotUserStore(db dalgo4botsfw.DbProvider) botsfw.BotUserStore {
//
//	newUserData := func() botsfw.BotUser {
//		return new(botsfwtgmodels.PlatformBotUserBaseData)
//	}
//
//	createBotUser := func(c context.Context, botID string, apiUser botsfw.WebhookActor) (botsfw.BotUser, error) {
//		if apiUser == nil {
//			return &botsfwtgmodels.PlatformBotUserBaseData{}, nil
//		}
//		return &botsfwtgmodels.PlatformBotUserBaseData{
//			PlatformUserData: botsfw.PlatformUserData{
//				BotEntity: botsfw.BotEntity{
//					OwnedByUserWithID: user.NewOwnedByUserWithIntID(0, time.Now()),
//				},
//				FirstName: apiUser.GetFirstName(),
//				LastName:  apiUser.GetLastName(),
//				UserName:  apiUser.GetUserName(),
//			},
//		}, nil
//	}
//
//	return dalgo4botsfw.NewBotUserStore(botsfwtgmodels.BotUserCollection, db, newUserData, createBotUser)
//}
