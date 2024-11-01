package utilitiesdb

import (
	"log"
	"pymouse/pymouse/database"
	"pymouse/pymouse/database/modeldb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindUser(UserID int64, UserName string) (UI *modeldb.UsersInformations) {
	var Filter bson.M

	dftUser := &modeldb.UsersInformations{
		UserID:   UserID,
		UserName: UserName,
	}

	UsersCollection := database.NewMongoCollection("users")

	if UserID != 0 {
		Filter = bson.M{"user_id": UserID}
	} else if UserName != "" {
		Filter = bson.M{"username": UserName}
	} else {
		return nil
	}

	err := UsersCollection.FindOne(Filter).Decode(&UI)
	if err == mongo.ErrNoDocuments {
		UI = nil
	} else if err != nil {
		log.Printf("[MongoDB][Users/FindUser][Error]: %v", err)
		UI = dftUser
	}

	return UI
}

func UpdateUser(UserID int64, UserName string, FirstName string) {
	UsersCollection := database.NewMongoCollection("users")
	UI := FindUser(UserID, "")

	if UI != nil {
		if UI.FirstName == FirstName && UI.UserName == UserName {
			return
		}
		UI.FirstName = FirstName
		UI.UserName = UserName
	} else {
		UI = &modeldb.UsersInformations{
			UserID:    UserID,
			UserName:  UserName,
			FirstName: FirstName,
		}
	}
	err := UsersCollection.UpdateOne(bson.M{"user_id": UserID}, UI)
	if err != nil {
		log.Printf("[MongoDB][Users/UpdateUser][Error]: %v - %d", err, UserID)
		return
	}
	log.Printf("[MongoDB][Users/UpdateUser]: %d - %s, Updated with successfully!", UserID, FirstName)
}
