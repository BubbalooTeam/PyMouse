package afk

import (
	"fmt"
	"pymouse/pymouse/database/utilitiesdb"
	"pymouse/pymouse/helpers/utils"
	"regexp"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
)

func SetAway(ctx *th.Context, update telego.Update) error {
	bot := ctx.Bot()

	User := update.Message.From
	Away := utilitiesdb.GetAway(User.ID)
	if Away.IsAway {
		StopAway(ctx, bot, update)
		return nil
	}
	AwayReason := utils.GetArgs(update)

	// Set User to Away From Keyboard (AFK)
	utilitiesdb.SetAway(User.ID, time.Now().UTC(), AwayReason)

	// Send ChatAction via Telegram
	bot.SendChatAction(
		ctx,
		&telego.SendChatActionParams{
			ChatID: telegoutil.ID(update.Message.Chat.ID),
			Action: "typing",
		})

	// Send a notification to notify AFK
	if AwayReason != "" {
		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				Text:      fmt.Sprintf("<b>%s is now unavailable!</b>\n<b>Reason:</b> %s", update.Message.From.FirstName, AwayReason),
				ParseMode: "HTML",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			},
		)
		return nil
	}
	bot.SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      fmt.Sprintf("<b>%s is now unavailable!</b>", update.Message.From.FirstName),
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		},
	)
	return nil
}

func CheckAway(ctx *th.Context, update telego.Update) error {
	bot := ctx.Bot()
	message := update.Message
	re := regexp.MustCompile(`(?i)^\/?(afk|away|brb)\b`)

	if message == nil || message.SenderChat != nil || re.MatchString(message.Text) {
		return ctx.Next(update)
	}

	if message.Chat.Type != "private" {
		UserAway := utilitiesdb.GetAway(message.From.ID)

		if UserAway.IsAway {
			StopAway(ctx, bot, update)
			return ctx.Next(update)
		}

		// Get mentioned user and notify user who mentioned
		CaSAway(ctx, bot, update)
	}
	return ctx.Next(update)
}
