package afk

import (
	"fmt"
	"pymouse/pymouse/database/repositories"
	"pymouse/pymouse/helpers/i18n"
	"pymouse/pymouse/helpers/utils"
	"regexp"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
)

func SetAway(ctx *th.Context, update telego.Update) error {
	bot := ctx.Bot()
	l := i18n.Locale(update.Message.Chat)

	User := update.Message.From
	Away := repositories.GetAway(User.ID)
	if Away.IsAway {
		StopAway(ctx, bot, update, l)
		return nil
	}
	AwayReason := utils.GetArgs(update)

	// Set User to Away From Keyboard (AFK)
	repositories.SetAway(User.ID, time.Now().UTC(), AwayReason)

	// Send ChatAction via Telegram
	bot.SendChatAction(
		ctx,
		&telego.SendChatActionParams{
			ChatID: telegoutil.ID(update.Message.Chat.ID),
			Action: "typing",
		},
	)
	afkMessage := fmt.Sprintf(l("afk.afk-set"), User.FirstName)

	// Send a notification to notify AFK
	if AwayReason != "" {
		afkMessage += fmt.Sprintf(l("generic-strings.reason"), AwayReason)
	}
	bot.SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      afkMessage,
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
	l := i18n.Locale(message.Chat)

	if message.Chat.Type != "private" {
		UserAway := repositories.GetAway(message.From.ID)

		if UserAway == nil || UserAway.IsAway {
			StopAway(ctx, bot, update, l)
			return ctx.Next(update)
		}

		// Get mentioned user and notify user who mentioned
		CaSAway(ctx, bot, update, l)
	}
	return ctx.Next(update)
}
