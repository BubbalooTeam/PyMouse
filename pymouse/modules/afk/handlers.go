package afk

import (
	"fmt"
	"pymouse/pymouse/database/modeldb"
	"pymouse/pymouse/database/utilitiesdb"
	"pymouse/pymouse/helpers/telegram"
	"pymouse/pymouse/helpers/utils"
	"regexp"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
)

func SetAway(bot *telego.Bot, update telego.Update) {
	User := update.Message.From
	Away := utilitiesdb.GetAway(User.ID)
	if Away.IsAway {
		StopAway(bot, update)
		return
	}
	AwayReason := utils.GetArgs(update)

	// Set User to Away From Keyboard (AFK)
	utilitiesdb.SetAway(User.ID, time.Now().UTC(), AwayReason)

	// Send ChatAction via Telegram
	bot.SendChatAction(&telego.SendChatActionParams{
		ChatID: telegoutil.ID(update.Message.Chat.ID),
		Action: "typing",
	})

	// Send a notification to notify AFK
	if AwayReason != "" {
		bot.SendMessage(&telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      fmt.Sprintf("<b>%s is now unavailable!</b>\n<b>Reason:</b> %s", update.Message.From.FirstName, AwayReason),
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return
	}
	bot.SendMessage(&telego.SendMessageParams{
		ChatID:    telegoutil.ID(update.Message.Chat.ID),
		Text:      fmt.Sprintf("<b>%s is now unavailable!</b>", update.Message.From.FirstName),
		ParseMode: "HTML",
		ReplyParameters: &telego.ReplyParameters{
			MessageID: update.Message.MessageID,
		},
	})
}

func CheckAway(bot *telego.Bot, update telego.Update, next th.Handler) {
	var UI *modeldb.UsersInformations

	message := update.Message
	re := regexp.MustCompile(`(?i)^\/?(afk|away|brb)\b`)

	if message == nil ||
		message.SenderChat != nil ||
		message.Chat.Type == "private" ||
		re.MatchString(message.Text) {
		next(bot, update)
		return
	}

	UserAway := utilitiesdb.GetAway(message.From.ID)

	if UserAway.IsAway {
		StopAway(bot, update)
		next(bot, update)
		return
	}
	// Get mentioned user and notify user who mentioned
	UI = telegram.GetUserMentioned(bot, update)
	switch {
	case func() bool {
		UI = telegram.GetUserMentioned(bot, update)
		return UI != nil
	}():
		SenderAway(bot, update, UI)
	case func() bool {
		UI = telegram.GetUserReplied(bot, update)
		return UI != nil
	}():
		SenderAway(bot, update, UI)
	default:
		next(bot, update)
	}
}
