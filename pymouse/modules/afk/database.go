package afk

import (
	"encoding/json"
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

func GetAFK(UserID int64) (map[string]interface{}, error) {
	usersCollection := database.NewCollection("users")
	AFKch := make(chan map[string]interface{})
	AFKerrCh := make(chan error)

	go func() {
		AFKresponse := usersCollection.FindMatches(map[string]interface{}{"UserID": strconv.FormatInt(UserID, 10)})
		if len(AFKresponse) > 0 {
			AFKByte, err := json.Marshal(AFKresponse[0])
			if err != nil {
				AFKerrCh <- err
				return
			}
			err = json.Unmarshal(AFKByte, &GetAFKJSON)
			if err != nil {
				AFKerrCh <- err
				return
			}

			afkData := GetAFKJSON["AFK"].(map[string]interface{})
			AFKch <- afkData
		} else {
			AFKch <- nil
		}
	}()

	select {
	case afkData := <-AFKch:
		return afkData, nil
	case err := <-AFKerrCh:
		return nil, err
	}
}
