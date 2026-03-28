package repositories

import (
	"pymouse/pymouse/database"
	"pymouse/pymouse/database/models"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindUser(UserID int64, UserName string) (UI *models.UsersInformations) {
	var Filter bson.M

	defaultUser := &models.UsersInformations{
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
		logrus.Error(err)
		UI = defaultUser
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
		UI = &models.UsersInformations{
			UserID:    UserID,
			UserName:  UserName,
			FirstName: FirstName,
		}
	}
	err := UsersCollection.UpdateOne(bson.M{"user_id": UserID}, UI)
	if err != nil {
		logrus.Errorf("%v - %d", err, UserID)
		return
	}
	logrus.Infof("%d - %s, Updated with successfully!", UserID, FirstName)
}
