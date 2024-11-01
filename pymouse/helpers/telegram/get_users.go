package telegram

import (
	"pymouse/pymouse/database/modeldb"
	"pymouse/pymouse/database/utilitiesdb"

	"github.com/mymmrac/telego"
)

func GetUserViaEntities(bot *telego.Bot, update telego.Update, entities telego.MessageEntity) (UI *modeldb.UsersInformations) {
	message := update.Message

	EntOffset := entities.Offset
	EntLenght := entities.Length

	// Gets the message entity, and gets a user.
	UserEntity := message.Text[EntOffset : EntOffset+EntLenght]
	UI = utilitiesdb.FindUser(0, UserEntity)
	return UI
}

func GetUserMentioned(bot *telego.Bot, update telego.Update) (UI *modeldb.UsersInformations) {
	message := update.Message

	if message.Entities != nil {
		for _, y := range message.Entities {
			if y.Type == "mention" {
				UI = GetUserViaEntities(bot, update, y)
				return UI
			} else if y.Type == "text_mention" {
				UserID := y.User.ID

				// Get Absolute User from Database
				UI = utilitiesdb.FindUser(UserID, "")
				return UI
			}
		}
	}
	UserID := message.From.ID

	// Get Absolute User from Database
	UI = utilitiesdb.FindUser(UserID, "")
	return UI
}

func GetUserReplied(bot *telego.Bot, update telego.Update) (UI *modeldb.UsersInformations) {
	UserID := update.Message.ReplyToMessage.From.ID

	// Get Absolute User from Database
	UI = utilitiesdb.FindUser(UserID, "")
	return UI
}
