package repositories

import (
	"pymouse/pymouse/database"
	"pymouse/pymouse/database/models"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindChat(ChatID int64) (CI *models.ChatsInformations) {
	ChatsCollection := database.NewMongoCollection("chats")
	dftChat := &models.ChatsInformations{
		ChatID: ChatID,
	}
	err := ChatsCollection.FindOne(bson.M{"chat_id": ChatID}).Decode(&CI)
	if err == mongo.ErrNoDocuments {
		CI = nil
	} else if err != nil {
		logrus.Error(err)
		CI = dftChat
	}
	return CI
}

func UpdateChat(ChatID int64, Username string, ChatTitle string) {
	ChatsCollection := database.NewMongoCollection("chats")
	CI := FindChat(ChatID)

	if CI != nil {
		if CI.Username == Username && CI.ChatTitle == ChatTitle {
			return
		}
		CI.Username = Username
		CI.ChatTitle = ChatTitle
	} else {
		CI = &models.ChatsInformations{
			ChatID:    ChatID,
			Username:  Username,
			ChatTitle: ChatTitle,
		}
	}
	err := ChatsCollection.UpdateOne(bson.M{"chat_id": ChatID}, CI)
	if err != nil {
		logrus.Errorf("%v - %d", err, ChatID)
		return
	}
	logrus.Infof("%d - %s, Updated with successfully!", ChatID, ChatTitle)
}
