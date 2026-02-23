package afk

import (
	"fmt"
	"pymouse/pymouse/database/modeldb"
	"pymouse/pymouse/database/utilitiesdb"
	"pymouse/pymouse/helpers/telegram"
	"pymouse/pymouse/helpers/utils"
	"time"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
)

func StopAway(ctx *telegohandler.Context, bot *telego.Bot, update telego.Update, l func(string) string) {
	User := update.Message.From

	// Stop Away From Keyboard (AFK), updating database informations.
	utilitiesdb.UnSetAway(User.ID)

	bot.SendChatAction(
		ctx,
		&telego.SendChatActionParams{
			ChatID: telegoutil.ID(update.Message.Chat.ID),
			Action: "typing",
		},
	)
	bot.SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      fmt.Sprintf(l("afk.afk-back"), User.FirstName),
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		},
	)
}

func SenderAway(ctx *telegohandler.Context, bot *telego.Bot, update telego.Update, UI *modeldb.UsersInformations, l func(string) string) {
	var AwayText string

	Away := utilitiesdb.GetAway(UI.UserID)

	if Away.IsAway {
		bot.SendChatAction(
			ctx,
			&telego.SendChatActionParams{
				ChatID: telegoutil.ID(update.Message.Chat.ID),
				Action: "typing",
			},
		)
		AwayText += fmt.Sprintf(l("afk.afk-response"), UI.FirstName)
		if Away.AwayReason != "" {
			AwayText += fmt.Sprintf(l("generic-strings.reason"), Away.AwayReason)
		}
		if !Away.AwayTime.IsZero() {
			AwayText += fmt.Sprintf(l("generic-strings.last-seen"), utils.TimeFormatter(time.Now().UTC().Sub(Away.AwayTime).Seconds()))
		}
		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				Text:      AwayText,
				ParseMode: "HTML",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			},
		)
	}
}

func CaSAway(ctx *telegohandler.Context, bot *telego.Bot, update telego.Update, l func(string) string) {
	var UI *modeldb.UsersInformations
	// Get mentioned user and notify user who mentioned
	UI = telegram.GetUserMentioned(bot, update)

	if UI != nil {
		SenderAway(ctx, bot, update, UI, l)
		return
	}

	UI = telegram.GetUserReplied(bot, update)
	if UI != nil {
		SenderAway(ctx, bot, update, UI, l)
		return
	}
}
