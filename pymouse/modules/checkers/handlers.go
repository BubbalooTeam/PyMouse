package checkers

import (
	"fmt"
	utilsDB "pymouse/pymouse/database/utilsdb"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func SaveUsers(bot *telego.Bot, update telego.Update, next th.Handler) {
	var UserName string

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
	// Telegram user informations
	UserID := message.From.ID
	FirstName := message.From.FirstName
	if message.From.Username != "" {
		UserName = fmt.Sprintf("@%s", message.From.Username)
	}

	// Update or Insert User Informations
	utilsDB.UpdateUser(UserID, UserName, FirstName)

	// Pass to next handler
	next(bot, update)
}
