package afk

import (
	"encoding/json"
	"log"
	"pymouse/pymouse/database"
	"time"
)

var AFKMapperFormated map[string]interface{}

type AFKInformation struct {
	AFKbool   bool      `json:"AFKbool"`
	AFKtime   time.Time `json:"AFKtime"`
	AFKreason string    `json:"AFKreason"`
}

type AFKResponse struct {
	UserID string         `json:"UserID"`
	AFK    AFKInformation `json:"AFK"`
}

func updateAFK(
	UserID string,
	AFKbool bool,
	AFKtime time.Time,
	AFKreason string,
) {
	usersCollection := database.NewCollection("users")
	afkFilter := map[string]interface{}{"UserID": UserID}

	afkMapper := AFKResponse{
		UserID: UserID,
		AFK: AFKInformation{
			AFKbool:   AFKbool,
			AFKtime:   AFKtime,
			AFKreason: AFKreason,
		},
	}

	AFKMapperByte, err := json.Marshal(afkMapper)
	if err != nil {
		log.Fatalf("Error in convert AFKMapper to JSON format: %v", err)
	}
	err = json.Unmarshal(AFKMapperByte, &AFKMapperFormated)
	if err != nil {
		log.Fatalf("Error in format JSON in AFKmap: %v", err)
	}
	usersCollection.InsertOrUpdate(afkFilter, AFKMapperFormated)
}

func SetupAFK(
	UserID string,
	AFKtime time.Time,
	AFKreason string,
) {
	updateAFK(UserID, true, AFKtime, AFKreason)
}

func UnsetupAFK(
	UserID string,
) {
	updateAFK(UserID, false, time.Time{}, "")
}
