package repositories

import (
	"pymouse/pymouse/database"
	"pymouse/pymouse/database/models"
	"strings"

	"github.com/mymmrac/telego"
	"go.mongodb.org/mongo-driver/bson"
)

func updateUserLanguage(chatID int64, language string) bool {
	UI := FindUser(chatID, "")
	if UI != nil {
		if UI.Language == language {
			return false
		}
	}
	UI = &models.UsersInformations{
		UserID:   chatID,
		Language: language,
	}

	UsersCollection := database.NewMongoCollection("users")
	err := UsersCollection.UpdateOne(bson.M{"user_id": chatID}, UI)
	if err == nil {
		return true
	}
	return false
}

func updateChatLanguage(chatID int64, language string) bool {
	CI := FindChat(chatID)
	if CI != nil {
		if CI.Language == language {
			return false
		}
	}
	CI = &models.ChatsInformations{
		ChatID:   chatID,
		Language: language,
	}

	ChatsCollection := database.NewMongoCollection("chats")
	err := ChatsCollection.UpdateOne(bson.M{"chat_id": chatID}, CI)
	if err == nil {
		return true
	}
	return false
}

func GetChatLanguage(chat telego.Chat) string {
	var chatLanguage string
	if strings.Contains(chat.Type, telego.ChatTypePrivate) {
		chatLanguage = FindUser(chat.ID, "").Language
	} else {
		chatLanguage = FindChat(chat.ID).Language
	}
	return chatLanguage
}

func SetChatLanguage(chat telego.Chat, language string) bool {
	if strings.Contains(chat.Type, telego.ChatTypePrivate) {
		return updateUserLanguage(chat.ID, language)
	}
	return updateChatLanguage(chat.ID, language)
}
