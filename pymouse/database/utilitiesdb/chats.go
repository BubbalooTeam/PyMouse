package utilitiesdb

import (
	"log"
	"pymouse/pymouse/database"
	"pymouse/pymouse/database/modeldb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindChat(ChatID int64) (CI *modeldb.ChatsInformations) {
	ChatsCollection := database.NewMongoCollection("chats")
	dftChat := &modeldb.ChatsInformations{
		ChatID: ChatID,
	}
	err := ChatsCollection.FindOne(bson.M{"chat_id": ChatID}).Decode(&CI)
	if err == mongo.ErrNoDocuments {
		CI = nil
	} else if err != nil {
		log.Printf("[MongoDB][Chats/FindChat][Error]: %v", err)
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
		CI = &modeldb.ChatsInformations{
			ChatID:    ChatID,
			Username:  Username,
			ChatTitle: ChatTitle,
		}
	}
	err := ChatsCollection.UpdateOne(bson.M{"chat_id": ChatID}, CI)
	if err != nil {
		log.Printf("[MongoDB][Chats/UpdateChat][Error]: %v - %d", err, ChatID)
		return
	}
	log.Printf("[MongoDB][Chats/UpdateChat]: %d - %s, Updated with successfully!", ChatID, ChatTitle)
}
