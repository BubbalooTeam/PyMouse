package utilsdb

import (
	"log"
	"pymouse/pymouse/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersInformations struct {
	UserID    int64  `bson:"user_id,omitempty" json:"user_id,omitempty"`
	UserName  string `bson:"username" json:"username" default:"nil"`
	FirstName string `bson:"first_name" json:"first_name" default:"nil"`
	Language  string `bson:"language" json:"language" default:"en_us"`
}

var UsersCollection = database.NewMongoCollection("users")

func FindUser(UserID int64) (UI *UsersInformations) {
	dftUser := &UsersInformations{
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
	UI := FindUser(UserID)

	if UI != nil {
		if UI.FirstName == FirstName && UI.UserName == UserName {
			return
		}
		UI.FirstName = FirstName
		UI.UserName = UserName
	} else {
		UI = &UsersInformations{
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
