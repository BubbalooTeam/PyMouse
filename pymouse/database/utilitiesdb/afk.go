package utilitiesdb

import (
	"log"
	"pymouse/pymouse/database"
	"pymouse/pymouse/database/modeldb"
	"time"

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
		UI = &modeldb.UsersInformations{
			UserID: UserID,
			Away: modeldb.AwayInformations{
				IsAway:     AwayState,
				AwayTime:   AwayTime,
				AwayReason: AwayReason,
			},
		}
	}
	err := UsersCollection.UpdateOne(bson.M{"user_id": UserID}, UI)
	if err != nil {
		log.Printf("[MongoDB][Afk/UpdateAway][Error]: %v - %d", err, UserID)
		return
	}
	log.Printf("[MongoDB][Afk/UpdateAway]: %d, Updated with successfully to state: %t!", UserID, AwayState)
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
) (AI modeldb.AwayInformations) {
	UI := FindUser(UserID, "")
	AI = UI.Away
	return AI
}
