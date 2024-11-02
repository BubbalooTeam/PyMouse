package afk

import (
	"fmt"
	"pymouse/pymouse/database/utilitiesdb"
	"pymouse/pymouse/helpers/utils"
	"regexp"
	"strings"
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
	message := update.Message
	re := regexp.MustCompile(`(?i)^\/?(afk|away|brb)\b`)

	if message == nil ||
		message.From == nil ||
		re.MatchString(message.Text) {
		return
	}

	if strings.Contains(message.Chat.Type, "private") {
		return
	}

	if utilitiesdb.GetAway(message.From.ID).IsAway {
		StopAway(bot, update)
		return
	}

	CaSAway(bot, update)

	next(bot, update)
}
