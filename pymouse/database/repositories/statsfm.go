package repositories

import (
	"pymouse/pymouse/database"
	"pymouse/pymouse/database/models"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateStatsFMUsername(UserID int64, Username string) {
	UI := FindUser(UserID, "")
	if UI != nil {
		UI.StatsFM.Username = Username
	} else {
		UI = &models.UsersInformations{
			UserID: UserID,
			StatsFM: models.StatsFMInformations{
				Username: Username,
			},
		}
	}

	UsersCollection := database.NewMongoCollection("users")
	err := UsersCollection.UpdateOne(bson.M{"user_id": UserID}, UI)
	if err != nil {
		logrus.Errorf("%v - %d", err, UserID)
		return
	}

}

func SetStatsFMUsername(UserID int64, Username string) {
	UpdateStatsFMUsername(UserID, Username)
}

func UnSetStatsFMUsername(UserID int64) {
	UpdateStatsFMUsername(UserID, "")
}

func GetStatsFMUsername(UserID int64) string {
	UI := FindUser(UserID, "")
	if UI != nil && UI.StatsFM.Username != "" {
		return UI.StatsFM.Username
	}
	return ""
}
