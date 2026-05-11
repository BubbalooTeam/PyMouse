package repositories

import (
	"pymouse/pymouse/database"
	"pymouse/pymouse/database/models"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateLastFMUsername(UserID int64, Username string) {
	UI := FindUser(UserID, "")
	if UI != nil {
		UI.LastFM.Username = Username
	} else {
		UI = &models.UsersInformations{
			UserID: UserID,
			LastFM: models.LastFMInformations{
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

func SetLastFMUsername(UserID int64, Username string) {
	UpdateLastFMUsername(UserID, Username)
}

func UnSetLastFMUsername(UserID int64) {
	UpdateLastFMUsername(UserID, "")
}

func GetLastFMUsername(UserID int64) string {
	UI := FindUser(UserID, "")
	if UI != nil && UI.LastFM.Username != "" {
		return UI.LastFM.Username
	}
	return ""
}
