package repositories

import (
	"pymouse/pymouse/database"
	"pymouse/pymouse/database/models"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateListenBrainzUsername(UserID int64, Username string) {
	UI := FindUser(UserID, "")
	if UI != nil {
		UI.ListenBrainz.Username = Username
	} else {
		UI = &models.UsersInformations{
			UserID: UserID,
			ListenBrainz: models.ListenBrainzInformations{
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

func SetListenBrainzUsername(UserID int64, Username string) {
	UpdateListenBrainzUsername(UserID, Username)
}

func UnSetListenBrainzUsername(UserID int64) {
	UpdateListenBrainzUsername(UserID, "")
}

func GetListenBrainzUsername(UserID int64) string {
	UI := FindUser(UserID, "")
	if UI != nil && UI.ListenBrainz.Username != "" {
		return UI.ListenBrainz.Username
	}
	return ""
}
