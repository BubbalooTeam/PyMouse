package utilitiesdb

import (
	"log"
	"pymouse/pymouse/database"
	"pymouse/pymouse/database/modeldb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindUser(UserID int64) (UI *modeldb.UsersInformations) {
	UsersCollection := database.NewMongoCollection("users")
	dftUser := &modeldb.UsersInformations{
		UserID: UserID,
	}
	err := UsersCollection.FindOne(bson.M{"user_id": UserID}).Decode(&UI)
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
	UI := FindUser(UserID)

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
		log.Printf("[MongoDB][users/UpdateUser][Error]: %v - %d", err, UserID)
		return
	}
	log.Printf("[MongoDB][users/UpdateUser]: %d - %s, Updated with successfully!", UserID, FirstName)
}
