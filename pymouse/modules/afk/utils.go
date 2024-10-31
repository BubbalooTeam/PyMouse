package afk

import (
	"fmt"
	"pymouse/pymouse/database/utilitiesdb"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

func StopAway(bot *telego.Bot, update telego.Update) {
	User := update.Message.From

	// Stop Away From Keyboard (AFK), updating database informations.
	utilitiesdb.UnSetAway(User.ID)

	bot.SendChatAction(
		&telego.SendChatActionParams{
			ChatID: telegoutil.ID(update.Message.Chat.ID),
			Action: "typing",
		},
	)
	bot.SendMessage(
		&telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      fmt.Sprintf("<b>%s is back!</b>", User.FirstName),
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		},
	)
}
