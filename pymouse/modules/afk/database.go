package afk

import (
	"encoding/json"
	"fmt"
	"log"
	"pymouse/pymouse/database"
	"strconv"
	"time"
)

var AFKMapperFormated map[string]interface{}
var GetAFKJSON map[string]interface{}

type AFKInformation struct {
	AFKbool   string    `json:"AFKbool"`
	AFKtime   time.Time `json:"AFKtime"`
	AFKreason string    `json:"AFKreason"`
}

type AFKResponse struct {
	UserID string         `json:"UserID"`
	AFK    AFKInformation `json:"AFK"`
}

func updateAFK(
	UserID int64,
	AFKbool bool,
	AFKtime time.Time,
	AFKreason string,
) {
	usersCollection := database.NewCollection("users")
	afkFilter := map[string]interface{}{"UserID": strconv.FormatInt(UserID, 10)}

	afkMapper := AFKResponse{
		UserID: strconv.FormatInt(UserID, 10),
		AFK: AFKInformation{
			AFKbool:   strconv.FormatBool(AFKbool),
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
	UserID int64,
	AFKtime time.Time,
	AFKreason string,
) {
	updateAFK(UserID, true, AFKtime, AFKreason)
}

func UnsetupAFK(
	UserID int64,
) {
	updateAFK(UserID, false, time.Time{}, "")
}

func GetAFK(UserID int64) map[string]interface{} {
	usersCollection := database.NewCollection("users")

	AFKresponse := usersCollection.FindMatches(map[string]interface{}{"UserID": strconv.FormatInt(UserID, 10)})
	if len(AFKresponse) > 0 {
		AFKByte, err := json.Marshal(AFKresponse[0])
		if err != nil {
			log.Fatalf("Error marshalling AFK response: %v", err)
			return nil // Return an empty map if there's an error
		}

		err = json.Unmarshal(AFKByte, &GetAFKJSON)
		if err != nil {
			fmt.Println("Error unmarshalling AFK response:", err)
			return nil // Return an empty map if there's an error
		}

		afkData := GetAFKJSON["AFK"]
		if afkData != nil {
			return afkData.(map[string]interface{})
		} else {
			return nil
		}
	} else {
		fmt.Println("User Not Found for AFK function.")
		return nil
	}
}
