package repositories

import (
	"pymouse/pymouse/database"
	"pymouse/pymouse/database/models"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateAway(
	UserID int64,
	AwayState bool,
	AwayTime time.Time,
	AwayReason string,
) {
	// Get Mongo Collection of "users".
	UsersCollection := database.NewMongoCollection("users")

	// Get user informations.
	UI := FindUser(UserID, "")

	if UI != nil {
		UI.Away.IsAway = AwayState
		UI.Away.AwayTime = AwayTime
		UI.Away.AwayReason = AwayReason
	} else {
		UI = &models.UsersInformations{
			UserID: UserID,
			Away: models.AwayInformations{
				IsAway:     AwayState,
				AwayTime:   AwayTime,
				AwayReason: AwayReason,
			},
		}
	}
	err := UsersCollection.UpdateOne(bson.M{"user_id": UserID}, UI)
	if err != nil {
		logrus.Errorf("%v - %d", err, UserID)
		return
	}
}

func SetAway(
	UserID int64,
	AwayTime time.Time,
	AwayReason string,
) {
	UpdateAway(UserID, true, AwayTime, AwayReason)
}

func UnSetAway(
	UserID int64,
) {
	UpdateAway(UserID, false, time.Time{}, "")
}

func GetAway(
	UserID int64,
) (AI models.AwayInformations) {
	UI := FindUser(UserID, "")
	AI = UI.Away
	return AI
}
