package afk

import (
	"fmt"
	"pymouse/pymouse/database/utilitiesdb"
	"pymouse/pymouse/helpers/utils"
	"time"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

func SetAway(bot *telego.Bot, update telego.Update) {
	User := update.Message.From
	IsAway := utilitiesdb.GetAway(User.ID).IsAway
	if IsAway {
		StopAway(bot, update)
		return
	}
	AwayReason := utils.GetArgs(update)

	// Set User to Away From Keyboard (AFK)
	utilitiesdb.SetAway(
		User.ID,
		time.Now().UTC(),
		AwayReason,
	)
	// Send ChatAction via Telegram
	bot.SendChatAction(
		&telego.SendChatActionParams{
			ChatID: telegoutil.ID(update.Message.Chat.ID),
			Action: "typing",
		},
	)

	// Send a notification to notify AFK
	if AwayReason != "" {
		bot.SendMessage(
			&telego.SendMessageParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				Text:      fmt.Sprintf("<b>%s is now AFK!</b>\n<b>Reason:</b> %s", update.Message.From.FirstName, AwayReason),
				ParseMode: "HTML",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			},
		)
		return
	}
	bot.SendMessage(
		&telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      fmt.Sprintf("<b>%s is now AFK!</b>", update.Message.From.FirstName),
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		},
	)
}
