package twrp

import (
	"fmt"
	"pymouse/pymouse/helpers/i18n"
	"pymouse/pymouse/helpers/telegram"
	"pymouse/pymouse/helpers/utils"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
)

func TWRPHandler(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	l := i18n.Locale(update.Message.Chat)

	query := utils.GetArgs(update)
	if query == "" {
		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				Text:      l("android.twrp.checkers.device-not-provided"),
				ParseMode: "HTML",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			},
		)
		return nil
	}

	results, found := getTWRP(query)
	if !found || len(results) == 0 {
		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				Text:      l("android.twrp.checkers.device-not-found"),
				ParseMode: "HTML",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			},
		)
		return nil
	}

	callbackKey := utils.RandKey()
	setTWRPResults(callbackKey, query, update.Message.From.ID, results)
	bot.SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      twrpText(query, results[0], l),
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
			ReplyMarkup: telegoutil.InlineKeyboardGrid(twrpKeyboard(
				callbackKey,
				results,
				1,
				update.Message.From.ID,
				l,
			)),
		},
	)
	return nil
}

func twrpText(device string, result TWRPBaseResult, l func(string) string) string {
	return fmt.Sprintf(
		l("android.twrp.info"),
		device,
		result.Updated,
		result.FileName,
		result.FileSize,
	)
}

func twrpKeyboard(key string, results []TWRPBaseResult, page int, ownerID int64, l func(string) string) [][]telego.InlineKeyboardButton {
	keyboard := telegram.KeyboardPaginate(
		len(results),
		page,
		fmt.Sprintf("twrp_page|%s|{number}|%d", key, ownerID),
	)
	keyboard = append(keyboard, []telego.InlineKeyboardButton{{
		Text: l("android.twrp.download"),
		URL:  results[page-1].DownloadURL,
	}})
	return keyboard
}

func TWRPPageCallbackHandler(ctx *telegohandler.Context, update telego.Update) error {
	callback := update.CallbackQuery
	l := i18n.Locale(callback.Message.GetChat())
	data := strings.Split(callback.Data, "|")
	if len(data) != 4 {
		return nil
	}

	page, err := strconv.Atoi(data[2])
	ownerID, ownerErr := strconv.ParseInt(data[3], 10, 64)
	if err != nil || ownerErr != nil {
		return nil
	}

	results, found := getTWRPResults(data[1])
	if !found || callback.From.ID != ownerID || results.OwnerID != callback.From.ID {
		return twrpUnauthorized(ctx, update, l)
	}
	if page < 1 || page > len(results.Results) {
		return nil
	}

	bot := ctx.Bot()
	bot.EditMessageText(
		ctx,
		&telego.EditMessageTextParams{
			ChatID:    telegoutil.ID(callback.Message.GetChat().ID),
			MessageID: callback.Message.GetMessageID(),
			Text:      twrpText(results.Device, results.Results[page-1], l),
			ParseMode: "HTML",
			ReplyMarkup: telegoutil.InlineKeyboardGrid(twrpKeyboard(
				data[1],
				results.Results,
				page,
				results.OwnerID,
				l,
			)),
		},
	)
	return nil
}

func twrpUnauthorized(ctx *telegohandler.Context, update telego.Update, l func(string) string) error {
	ctx.Bot().AnswerCallbackQuery(
		ctx,
		&telego.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            l("android.twrp.checkers.not-for-you"),
			ShowAlert:       true,
			CacheTime:       3,
		},
	)
	return nil
}
