package lfm

import (
	"context"
	"fmt"
	"os"
	"pymouse/pymouse/database/repositories"
	"pymouse/pymouse/helpers/i18n"
	"pymouse/pymouse/helpers/rapidhttp"
	"strconv"
	"strings"

	"github.com/cavaliergopher/grab/v3"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
)

// recentLimit is how many scrobbles the "+" button expands to.
const recentLimit = 5

// fontSet is the same Fonts literal used by NowPlaying; kept here so the
// callback handlers can regenerate cards without duplicating the paths.
func fontSet() Fonts {
	return Fonts{
		OpenSans: "pymouse/assets/fonts/opensans.ttf",
		Poppins:  "pymouse/assets/fonts/poppins-semibolditalic.ttf",
		Arial:    "pymouse/assets/fonts/arial.ttf",
		Unicode:  "pymouse/assets/fonts/notosans-unicode.ttf",
		CJK:      "pymouse/assets/fonts/notosanscjk-sc.otf",
	}
}

// parseOwner extracts the owner userID from a "prefix|userID" callback data
// and reports whether the caller is that owner.
func parseOwner(data string, fromID int64) (int64, bool) {
	parts := strings.Split(data, "|")
	if len(parts) != 2 {
		return 0, false
	}
	uid, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, false
	}
	return uid, uid == fromID
}

// denyOthers answers the callback with an alert when someone other than the
// card's owner taps the button.
func denyOthers(bot *telego.Bot, ctx context.Context, queryID, msg string) error {
	return bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: queryID,
		Text:            msg,
		ShowAlert:       true,
	})
}

// ExpandRecent handles the "+" button: replaces the now-playing photo with a
// list of the owner's last `recentLimit` scrobbles.
func ExpandRecent(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	callback := update.CallbackQuery
	l := i18n.Locale(callback.Message.GetChat())

	ownerID, ok := parseOwner(callback.Data, callback.From.ID)
	if !ok {
		return denyOthers(bot, ctx, callback.ID, l("lastfm.recent.not-yours"))
	}

	username := repositories.GetLastFMUsername(ownerID)
	httpClient := rapidhttp.GetHTTPClient()

	tracks, err := getRecentTracks(httpClient, username, recentLimit)
	if err != nil {
		return denyOthers(bot, ctx, callback.ID, l("lastfm.nowplaying.error"))
	}

	imgPath, err := DrawRecentScrobble(username, tracks, fontSet(), l)
	if err != nil {
		return denyOthers(bot, ctx, callback.ID, l("lastfm.nowplaying.error"))
	}
	imF, err := os.Open(imgPath)
	if err != nil {
		return denyOthers(bot, ctx, callback.ID, l("lastfm.nowplaying.error"))
	}
	defer func() {
		imF.Close()
		os.Remove(imgPath)
	}()

	profileBtn := []telego.InlineKeyboardButton{
		{Text: "👤 " + l("lastfm.recent.profile-btn"), URL: "https://www.last.fm/user/" + username},
		{Text: "↩️ " + l("lastfm.recent.back-btn"), CallbackData: fmt.Sprintf("lfm_back|%d", ownerID)},
	}

	if _, err := bot.EditMessageMedia(ctx, &telego.EditMessageMediaParams{
		ChatID:    telegoutil.ID(callback.Message.GetChat().ID),
		MessageID: callback.Message.GetMessageID(),
		Media: &telego.InputMediaPhoto{
			Type:  "photo",
			Media: telego.InputFile{File: imF},
		},
		ReplyMarkup: &telego.InlineKeyboardMarkup{InlineKeyboard: [][]telego.InlineKeyboardButton{profileBtn}},
	}); err != nil {
		// editing can fail if the message is unchanged; not fatal
	}

	return bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: callback.ID,
	})
}

// BackToNowPlaying handles the "← back" button: regenerates the original
// now-playing card and swaps it back in.
func BackToNowPlaying(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	callback := update.CallbackQuery
	l := i18n.Locale(callback.Message.GetChat())

	ownerID, ok := parseOwner(callback.Data, callback.From.ID)
	if !ok {
		return denyOthers(bot, ctx, callback.ID, l("lastfm.recent.not-yours"))
	}

	username := repositories.GetLastFMUsername(ownerID)
	httpClient := rapidhttp.GetHTTPClient()
	glabClient := grab.NewClient()

	trackInfo, err := getTrack(httpClient, username)
	if err != nil {
		return denyOthers(bot, ctx, callback.ID, l("lastfm.nowplaying.error"))
	}

	imageURL := trackInfo.Image
	if imageURL == "" {
		imageURL = "https://telegra.ph/file/bdcf492162713ea5633a1.jpg"
	}
	listeningText := getListeningText(trackInfo, l)

	imgPath, err := DrawScrobble(
		glabClient, imageURL, trackInfo.Track, trackInfo.Artist,
		username, listeningText, trackInfo.Loved, fontSet(), l,
	)
	if err != nil {
		return denyOthers(bot, ctx, callback.ID, l("lastfm.nowplaying.error"))
	}
	imF, err := os.Open(imgPath)
	if err != nil {
		return denyOthers(bot, ctx, callback.ID, l("lastfm.nowplaying.error"))
	}
	defer func() {
		imF.Close()
		os.Remove(imgPath)
	}()

	rows := [][]telego.InlineKeyboardButton{{
		{Text: "▶️ YouTube", URL: trackInfo.YouTubeURL},
		{Text: "👤 " + l("lastfm.recent.profile-btn"), URL: "https://www.last.fm/user/" + username},
		{Text: "📋 " + l("lastfm.recent.expand-btn"), CallbackData: fmt.Sprintf("lfm_recent|%d", ownerID)},
	}}

	if _, err := bot.EditMessageMedia(ctx, &telego.EditMessageMediaParams{
		ChatID:    telegoutil.ID(callback.Message.GetChat().ID),
		MessageID: callback.Message.GetMessageID(),
		Media: &telego.InputMediaPhoto{
			Type:  "photo",
			Media: telego.InputFile{File: imF},
		},
		ReplyMarkup: &telego.InlineKeyboardMarkup{InlineKeyboard: rows},
	}); err != nil {
		// editing can fail if the message is unchanged; not fatal
	}

	return bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: callback.ID,
	})
}
