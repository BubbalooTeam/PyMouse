package lineageos

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

func LineageOSHandler(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	l := i18n.Locale(update.Message.Chat)
	device := strings.ToLower(strings.TrimSpace(utils.GetArgs(update)))
	if device == "" {
		bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      l("android.lineageos.checkers.device-not-provided"),
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return nil
	}

	builds, found := getLineageBuilds(device)
	if !found {
		bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      l("android.lineageos.checkers.device-not-found"),
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return nil
	}

	key := utils.RandKey()
	setLineageResults(key, update.Message.From.ID, device, builds)
	bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:    telegoutil.ID(update.Message.Chat.ID),
		Text:      lineageText(builds[0], device, l),
		ParseMode: "HTML",
		LinkPreviewOptions: &telego.LinkPreviewOptions{
			IsDisabled: true,
		},
		ReplyParameters: &telego.ReplyParameters{
			MessageID: update.Message.MessageID,
		},
		ReplyMarkup: telegoutil.InlineKeyboardGrid(lineageKeyboard(key, builds, 1, update.Message.From.ID, l)),
	})
	return nil
}

func lineageText(build LineageBuild, device string, l func(string) string) string {
	return fmt.Sprintf(
		l("android.lineageos.info"),
		device,
		build.Filename,
		build.ROMType,
		build.Size,
		build.Version,
		build.Datetime,
	)
}

func lineageKeyboard(key string, builds []LineageBuild, page int, ownerID int64, l func(string) string) [][]telego.InlineKeyboardButton {
	keyboard := telegram.KeyboardPaginate(
		len(builds),
		page,
		fmt.Sprintf("lineageos_page|%s|{number}|%d", key, ownerID),
	)
	keyboard = append(keyboard, []telego.InlineKeyboardButton{{
		Text: l("android.lineageos.download"),
		URL:  builds[page-1].URL,
	}})
	return keyboard
}

func LineageOSPageCallbackHandler(ctx *telegohandler.Context, update telego.Update) error {
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

	results, found := getLineageResults(data[1])
	if !found || callback.From.ID != ownerID || results.OwnerID != callback.From.ID {
		return lineageUnauthorized(ctx, update, l)
	}
	if page < 1 || page > len(results.Builds) {
		return nil
	}

	ctx.Bot().EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:    telegoutil.ID(callback.Message.GetChat().ID),
		MessageID: callback.Message.GetMessageID(),
		Text:      lineageText(results.Builds[page-1], results.Device, l),
		ParseMode: "HTML",
		LinkPreviewOptions: &telego.LinkPreviewOptions{
			IsDisabled: true,
		},
		ReplyMarkup: telegoutil.InlineKeyboardGrid(lineageKeyboard(
			data[1], results.Builds, page, results.OwnerID, l,
		)),
	})
	return nil
}

func lineageUnauthorized(ctx *telegohandler.Context, update telego.Update, l func(string) string) error {
	ctx.Bot().AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		Text:            l("android.lineageos.checkers.not-for-you"),
		ShowAlert:       true,
		CacheTime:       3,
	})
	return nil
}
