package set

import (
	"fmt"
	"pymouse/pymouse/database/repositories"
	"pymouse/pymouse/helpers/i18n"
	"pymouse/pymouse/helpers/rapidhttp"
	"pymouse/pymouse/helpers/utils"
	"strings"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
)

func SetUser(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	l := i18n.Locale(update.Message.Chat)

	query := strings.ToLower(utils.GetArgs(update))
	if query == "" {
		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				Text:      l("statsfm.set.no-username-provided"),
				ParseMode: "HTML",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
				ReplyMarkup: &telego.InlineKeyboardMarkup{
					InlineKeyboard: [][]telego.InlineKeyboardButton{
						{
							telego.InlineKeyboardButton{
								Text: l("buttons.create-statsfm"),
								URL:  "https://stats.fm/",
							},
						},
					},
				},
			},
		)
		return nil
	}

	httpClient := rapidhttp.GetHTTPClient()
	username := normalizeUsername(query)
	if !CheckUsername(httpClient, username) {
		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				Text:      l("statsfm.invalid-username"),
				ParseMode: "HTML",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			},
		)
		return nil
	}

	repositories.SetStatsFMUsername(update.Message.From.ID, username)
	bot.SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      fmt.Sprintf(l("statsfm.set.username-set"), username),
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		},
	)
	return nil
}

func UnsetUser(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	l := i18n.Locale(update.Message.Chat)

	username := repositories.GetStatsFMUsername(update.Message.From.ID)
	if username == "" {
		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				Text:      l("statsfm.set.no-username-set"),
				ParseMode: "HTML",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			},
		)
		return nil
	}

	repositories.UnSetStatsFMUsername(update.Message.From.ID)
	bot.SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      l("statsfm.set.username-unset"),
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		},
	)
	return nil
}
