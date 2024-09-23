package modules

import (
	"encoding/json"
	"fmt"
	"log"
	"pymouse/pymouse/database"
	"strconv"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

var UserInfoJSON map[string]interface{}

func SaveUsers(bot *telego.Bot, update telego.Update, next th.Handler) {
	var UserUserName string

	message := update.Message

	if message == nil {
		if update.CallbackQuery == nil {
			return
		}
		message = update.CallbackQuery.Message.(*telego.Message)
	}
	if message.SenderChat != nil {
		return
	}
	UserID := strconv.FormatInt(message.From.ID, 10)
	UserFirstName := message.From.FirstName
	userLanguage := message.From.LanguageCode
	if message.From.Username != "" {
		UserUserName = fmt.Sprintf("@%s", message.From.Username)
	}

	usersCollection := database.NewCollection("users")
	userFilter := map[string]interface{}{"UserID": UserID}
	findUser := func() bool {
		UserInfoList := usersCollection.FindMatches(userFilter)
		fmt.Println(UserInfoList)
		if len(UserInfoList) > 0 {
			UserInfoByte, err := json.Marshal(UserInfoList[0])
			if err != nil {
				log.Fatalf("Error in marshaling UserInfo: %v", err)
			}
			err = json.Unmarshal(UserInfoByte, &UserInfoJSON)
			if err != nil {
				log.Fatalf("Error in unmarshaling UserInfo in JSON format: %v", err)
			}
			if UserInfoJSON != nil {
				DBFirstName := UserInfoJSON["FirstName"]
				DBUserName := UserInfoJSON["UserName"]
				DBUserTGLanguage := UserInfoJSON["UserLanguage"]
				// Check if user already is updated
				if UserFirstName == DBFirstName && UserUserName == DBUserName && userLanguage == DBUserTGLanguage {
					return false
				} else {
					return true
				}
			} else {
				return false
			}
		} else {
			return true
		}
	}
	updateUser := func() {
		UserNeedUpdate := findUser()
		if UserNeedUpdate {
			userData := map[string]interface{}{
				"UserID":       UserID,
				"FirstName":    UserFirstName,
				"UserName":     UserUserName,
				"UserLanguage": userLanguage,
			}
			usersCollection.InsertOrUpdate(userFilter, userData)
			fmt.Printf("\033[0;36mInserted/Updated User %s[%s] successfully!\033[0m", UserFirstName, UserID)
		}
	}
	updateUser()
	next(bot, update)
}
