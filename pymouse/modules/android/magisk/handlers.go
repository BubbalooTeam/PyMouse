package magisk

import (
	"fmt"
	"pymouse/pymouse/helpers/i18n"
	"pymouse/pymouse/helpers/utils"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
)

func MagiskHandler(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	l := i18n.Locale(update.Message.Chat)

	releases, found := getMagisk()
	if !found {
		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				Text:      l("android.magisk.checkers.github-error"),
				ParseMode: "HTML",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			},
		)
		return nil
	}

	callbackKey := utils.RandKey()
	setMagiskResults(callbackKey, update.Message.From.ID, releases)

	bot.SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      magiskText(releases[0], l),
			ParseMode: "HTML",
			LinkPreviewOptions: &telego.LinkPreviewOptions{
				IsDisabled: true,
			},
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
			ReplyMarkup: telegoutil.InlineKeyboardGrid(magiskKeyboard(callbackKey, releases[0], update.Message.From.ID, l)),
		},
	)
	return nil
}

func magiskText(release MagiskRelease, l func(string) string) string {
	versionPrefix := ""
	if len(release.Version) > 0 && release.Version[0] >= '0' && release.Version[0] <= '9' {
		versionPrefix = "v"
	}
	return fmt.Sprintf(
		l("android.magisk.info"),
		l("android.magisk."+release.Variant),
		versionPrefix,
		release.Version,
		release.VersionCode,
	)
}

func magiskKeyboard(key string, release MagiskRelease, ownerID int64, l func(string) string) [][]telego.InlineKeyboardButton {
	keyboard := [][]telego.InlineKeyboardButton{{
		{Text: l("android.magisk.stable"), CallbackData: fmt.Sprintf("magisk_variant|%s|stable|%d", key, ownerID)},
		{Text: l("android.magisk.beta"), CallbackData: fmt.Sprintf("magisk_variant|%s|beta|%d", key, ownerID)},
		{Text: l("android.magisk.canary"), CallbackData: fmt.Sprintf("magisk_variant|%s|canary|%d", key, ownerID)},
	}}
	keyboard = append(keyboard, []telego.InlineKeyboardButton{
		{Text: l("android.magisk.download"), URL: release.DownloadURL},
		{Text: l("android.magisk.changelog"), URL: release.ChangelogURL},
	})
	return keyboard
}

func MagiskVariantCallbackHandler(ctx *telegohandler.Context, update telego.Update) error {
	callback := update.CallbackQuery
	l := i18n.Locale(callback.Message.GetChat())
	data := strings.Split(callback.Data, "|")
	if len(data) != 4 {
		return nil
	}

	ownerID, ownerErr := strconv.ParseInt(data[3], 10, 64)
	if ownerErr != nil {
		return nil
	}

	results, found := getMagiskResults(data[1])
	if !found || callback.From.ID != ownerID || results.OwnerID != callback.From.ID {
		return magiskUnauthorized(ctx, update, l)
	}
	var release MagiskRelease
	for _, result := range results.Results {
		if result.Variant == data[2] {
			release = result
			break
		}
	}
	if release.Variant == "" {
		return nil
	}

	ctx.Bot().EditMessageText(
		ctx,
		&telego.EditMessageTextParams{
			ChatID:    telegoutil.ID(callback.Message.GetChat().ID),
			MessageID: callback.Message.GetMessageID(),
			Text:      magiskText(release, l),
			ParseMode: "HTML",
			ReplyMarkup: telegoutil.InlineKeyboardGrid(magiskKeyboard(
				data[1],
				release,
				results.OwnerID,
				l,
			)),
		},
	)
	return nil
}

func magiskUnauthorized(ctx *telegohandler.Context, update telego.Update, l func(string) string) error {
	ctx.Bot().AnswerCallbackQuery(
		ctx,
		&telego.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            l("android.magisk.checkers.not-for-you"),
			ShowAlert:       true,
			CacheTime:       3,
		},
	)
	return nil
}
