package telegram

import (
	"pymouse/pymouse/database/models"
	"pymouse/pymouse/database/repositories"
	"strings"

	"github.com/mymmrac/telego"
)

func GetUserViaEntities(bot *telego.Bot, update telego.Update, entities telego.MessageEntity) (UI *models.UsersInformations) {
	message := update.Message

	EntOffset := entities.Offset
	EntLenght := entities.Length

	// Gets the message entity, and gets a user.
	UserEntity := message.Text[EntOffset : EntOffset+EntLenght]
	// Remove @ from the username and get the user informations
	UserEntity = strings.Replace(UserEntity, "@", "", -1)
	UI = repositories.FindUser(0, UserEntity)
	return UI
}

func GetUserMentioned(bot *telego.Bot, update telego.Update) (UI *models.UsersInformations) {
	var UserID int64
	message := update.Message

	if message.Entities != nil {
		for _, y := range message.Entities {
			switch y.Type {
			case "mention":
				UI = GetUserViaEntities(bot, update, y)
				return UI
			case "text_mention":
				if y.User != nil {
					UserID = y.User.ID

					// Get Absolute User from Database
					UI = repositories.FindUser(UserID, "")
					return UI
				}
			}
		}
	}
	return nil
}

func GetUserReplied(bot *telego.Bot, update telego.Update) (UI *models.UsersInformations) {
	if update.Message != nil && update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From != nil {
		UserID := update.Message.ReplyToMessage.From.ID

		// Get Absolute User from Database
		UI = repositories.FindUser(UserID, "")
		return UI
	}
	return nil
}
