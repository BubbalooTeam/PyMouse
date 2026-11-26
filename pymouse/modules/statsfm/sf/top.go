package sf

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"pymouse/pymouse/database/repositories"
	"pymouse/pymouse/helpers/i18n"
	"pymouse/pymouse/helpers/rapidhttp"
	"pymouse/pymouse/helpers/scrobblecard"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
)

// topLimit is how many ranked entries /sftop shows.
const topLimit = 10

// topCategory/topRange are the only values the /sftop keyboard cycles
// through; both are embedded verbatim in callback_data, so they're kept
// short and must stay URL/callback-safe (no "|").
type topCategory string
type topRange string

const (
	categoryTracks  topCategory = "trk"
	categoryArtists topCategory = "art"
	categoryAlbums  topCategory = "alb"
	categoryGenres  topCategory = "gen"
)

const (
	rangeToday  topRange = "today"
	rangeDays   topRange = "days"   // stats.fm's own "last 7 days" window
	rangeWeeks  topRange = "weeks"  // stats.fm's own "last 4 weeks" window
	rangeMonths topRange = "months" // stats.fm's own "last 6 months" window
	rangeLife   topRange = "life"
)

var topCategoryOrder = []topCategory{categoryTracks, categoryArtists, categoryAlbums, categoryGenres}
var topRangeOrder = []topRange{rangeToday, rangeDays, rangeWeeks, rangeMonths, rangeLife}

// apiRange maps our compact range code to stats.fm's query value: identical
// for every range except "life", kept distinct from "lifetime" so
// callback_data stays a couple bytes shorter.
func (r topRange) apiRange() string {
	if r == rangeLife {
		return "lifetime"
	}
	return string(r)
}

func parseCategory(s string) (topCategory, bool) {
	for _, c := range topCategoryOrder {
		if string(c) == s {
			return c, true
		}
	}
	return "", false
}

func parseRange(s string) (topRange, bool) {
	for _, r := range topRangeOrder {
		if string(r) == s {
			return r, true
		}
	}
	return "", false
}

// categoryEmoji is UI sugar for the keyboard buttons only — the card header
// uses the plain localized label (loadFont's bundled fonts don't carry emoji
// glyphs).
func categoryEmoji(c topCategory) string {
	switch c {
	case categoryTracks:
		return "🎵"
	case categoryArtists:
		return "🎤"
	case categoryAlbums:
		return "💿"
	case categoryGenres:
		return "🏷️"
	}
	return ""
}

func categoryLabelKey(c topCategory) string {
	switch c {
	case categoryTracks:
		return "statsfm.top.category.tracks"
	case categoryArtists:
		return "statsfm.top.category.artists"
	case categoryAlbums:
		return "statsfm.top.category.albums"
	case categoryGenres:
		return "statsfm.top.category.genres"
	}
	return ""
}

func rangeLabelKey(r topRange) string {
	switch r {
	case rangeToday:
		return "statsfm.top.range.today"
	case rangeDays:
		return "statsfm.top.range.days"
	case rangeWeeks:
		return "statsfm.top.range.weeks"
	case rangeMonths:
		return "statsfm.top.range.months"
	case rangeLife:
		return "statsfm.top.range.life"
	}
	return ""
}

// rangeButtonLabel is the short, English, locale-independent label shown on
// the range buttons themselves — the full localized phrase (used in the card
// header instead, via rangeLabelKey) is too long for the button width and
// renders truncated by Telegram's client ("Último...", "Últimas...").
func rangeButtonLabel(r topRange) string {
	switch r {
	case rangeToday:
		return "1d"
	case rangeDays:
		return "7d"
	case rangeWeeks:
		return "4w"
	case rangeMonths:
		return "6m"
	case rangeLife:
		return "all"
	}
	return string(r)
}

func joinArtistNames(artists []sfmArtist) string {
	names := make([]string, 0, len(artists))
	for _, a := range artists {
		names = append(names, a.Name)
	}
	return strings.Join(names, ", ")
}

func getTopTracks(httpClient *http.Client, username string, rng topRange, limit int) ([]scrobblecard.TopEntry, error) {
	var out struct {
		Items []struct {
			Streams int64    `json:"streams"`
			Track   sfmTrack `json:"track"`
		} `json:"items"`
	}
	status, err := fetchJSON(httpClient, "users/"+url.PathEscape(username)+"/top/tracks",
		map[string]string{"range": rng.apiRange(), "limit": strconv.Itoa(limit)}, &out)
	if err != nil {
		return nil, err
	}
	if status == 404 {
		return nil, fmt.Errorf("invalid username.")
	}
	if status != 200 {
		return nil, fmt.Errorf("failed to fetch top tracks.")
	}
	entries := make([]scrobblecard.TopEntry, 0, len(out.Items))
	for _, it := range out.Items {
		entries = append(entries, scrobblecard.TopEntry{
			Label:     it.Track.Name,
			Secondary: artistNames(it.Track),
			Streams:   it.Streams,
		})
	}
	return entries, nil
}

func getTopArtists(httpClient *http.Client, username string, rng topRange, limit int) ([]scrobblecard.TopEntry, error) {
	var out struct {
		Items []struct {
			Streams int64 `json:"streams"`
			Artist  struct {
				Name string `json:"name"`
			} `json:"artist"`
		} `json:"items"`
	}
	status, err := fetchJSON(httpClient, "users/"+url.PathEscape(username)+"/top/artists",
		map[string]string{"range": rng.apiRange(), "limit": strconv.Itoa(limit)}, &out)
	if err != nil {
		return nil, err
	}
	if status == 404 {
		return nil, fmt.Errorf("invalid username.")
	}
	if status != 200 {
		return nil, fmt.Errorf("failed to fetch top artists.")
	}
	entries := make([]scrobblecard.TopEntry, 0, len(out.Items))
	for _, it := range out.Items {
		entries = append(entries, scrobblecard.TopEntry{Label: it.Artist.Name, Streams: it.Streams})
	}
	return entries, nil
}

func getTopAlbums(httpClient *http.Client, username string, rng topRange, limit int) ([]scrobblecard.TopEntry, error) {
	var out struct {
		Items []struct {
			Streams int64 `json:"streams"`
			Album   struct {
				Name    string      `json:"name"`
				Artists []sfmArtist `json:"artists"`
			} `json:"album"`
		} `json:"items"`
	}
	status, err := fetchJSON(httpClient, "users/"+url.PathEscape(username)+"/top/albums",
		map[string]string{"range": rng.apiRange(), "limit": strconv.Itoa(limit)}, &out)
	if err != nil {
		return nil, err
	}
	if status == 404 {
		return nil, fmt.Errorf("invalid username.")
	}
	if status != 200 {
		return nil, fmt.Errorf("failed to fetch top albums.")
	}
	entries := make([]scrobblecard.TopEntry, 0, len(out.Items))
	for _, it := range out.Items {
		entries = append(entries, scrobblecard.TopEntry{
			Label:     it.Album.Name,
			Secondary: joinArtistNames(it.Album.Artists),
			Streams:   it.Streams,
		})
	}
	return entries, nil
}

func getTopGenres(httpClient *http.Client, username string, rng topRange, limit int) ([]scrobblecard.TopEntry, error) {
	var out struct {
		Items []struct {
			Streams int64 `json:"streams"`
			Genre   struct {
				Tag string `json:"tag"`
			} `json:"genre"`
		} `json:"items"`
	}
	status, err := fetchJSON(httpClient, "users/"+url.PathEscape(username)+"/top/genres",
		map[string]string{"range": rng.apiRange(), "limit": strconv.Itoa(limit)}, &out)
	if err != nil {
		return nil, err
	}
	if status == 404 {
		return nil, fmt.Errorf("invalid username.")
	}
	if status != 200 {
		return nil, fmt.Errorf("failed to fetch top genres.")
	}
	entries := make([]scrobblecard.TopEntry, 0, len(out.Items))
	for _, it := range out.Items {
		entries = append(entries, scrobblecard.TopEntry{Label: it.Genre.Tag, Streams: it.Streams})
	}
	return entries, nil
}

// getTopEntries dispatches to the right top/* endpoint and truncates the
// result to limit client-side: /top/genres ignores the "limit" query param
// (observed returning 100+ items for limit=5) while the others respect it,
// so this truncation is required for genres and a harmless no-op elsewhere.
func getTopEntries(httpClient *http.Client, username string, cat topCategory, rng topRange, limit int) ([]scrobblecard.TopEntry, error) {
	var entries []scrobblecard.TopEntry
	var err error
	switch cat {
	case categoryArtists:
		entries, err = getTopArtists(httpClient, username, rng, limit)
	case categoryAlbums:
		entries, err = getTopAlbums(httpClient, username, rng, limit)
	case categoryGenres:
		entries, err = getTopGenres(httpClient, username, rng, limit)
	default:
		entries, err = getTopTracks(httpClient, username, rng, limit)
	}
	if err != nil {
		return nil, err
	}
	if len(entries) > limit {
		entries = entries[:limit]
	}
	return entries, nil
}

// topKeyboard builds the category row + range row, prefixing the active
// choice in each with "• " so the current view is visible on the card itself.
func topKeyboard(ownerID int64, cat topCategory, rng topRange) [][]telego.InlineKeyboardButton {
	catRow := make([]telego.InlineKeyboardButton, 0, len(topCategoryOrder))
	for _, c := range topCategoryOrder {
		text := categoryEmoji(c)
		if c == cat {
			text = "• " + text
		}
		catRow = append(catRow, telego.InlineKeyboardButton{
			Text:         text,
			CallbackData: fmt.Sprintf("sftop|%d|%s|%s", ownerID, c, rng),
		})
	}

	rngRow := make([]telego.InlineKeyboardButton, 0, len(topRangeOrder))
	for _, r := range topRangeOrder {
		text := rangeButtonLabel(r)
		if r == rng {
			text = "• " + text
		}
		rngRow = append(rngRow, telego.InlineKeyboardButton{
			Text:         text,
			CallbackData: fmt.Sprintf("sftop|%d|%s|%s", ownerID, cat, r),
		})
	}

	return [][]telego.InlineKeyboardButton{catRow, rngRow}
}

// renderTop fetches and draws the card for one category+range, returning the
// open file (caller must close+remove) and its keyboard.
func renderTop(
	httpClient *http.Client, username string, ownerID int64,
	cat topCategory, rng topRange, l func(string) string,
) (*os.File, [][]telego.InlineKeyboardButton, error) {
	entries, err := getTopEntries(httpClient, username, cat, rng, topLimit)
	if err != nil {
		return nil, nil, err
	}
	if len(entries) == 0 {
		return nil, nil, fmt.Errorf("no listens found.")
	}

	header := fmt.Sprintf("%s · %s", l(categoryLabelKey(cat)), l(rangeLabelKey(rng)))
	imgPath, err := scrobblecard.DrawTopList(scrobblecard.ScrobbleDir("statsfm"), username, entries, fontSet(), header)
	if err != nil {
		return nil, nil, err
	}
	defer os.Remove(imgPath)

	imF, err := os.Open(imgPath)
	if err != nil {
		return nil, nil, err
	}

	rows := topKeyboard(ownerID, cat, rng)

	return imF, rows, nil
}

func topErrorText(err error, l func(string) string) string {
	switch err.Error() {
	case "invalid username.":
		return l("statsfm.invalid-username")
	case "no listens found.":
		return l("statsfm.top.no-data")
	default:
		return l("statsfm.top.error")
	}
}

// Top handles /sftop and /top: shows the caller's top tracks for the last 7
// days, with buttons to switch between tracks/artists/albums/genres and
// today/7 days/4 weeks/6 months/lifetime.
func Top(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	l := i18n.Locale(update.Message.Chat)

	username := repositories.GetStatsFMUsername(update.Message.From.ID)
	if username == "" {
		bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      l("statsfm.top.no-username-set"),
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return nil
	}

	httpClient := rapidhttp.GetHTTPClient()
	ownerID := update.Message.From.ID

	imF, rows, err := renderTop(httpClient, username, ownerID, categoryTracks, rangeDays, l)
	if err != nil {
		bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      topErrorText(err, l),
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return nil
	}
	defer func() {
		imF.Close()
		os.Remove(imF.Name())
	}()

	bot.SendPhoto(ctx, &telego.SendPhotoParams{
		ChatID:          telegoutil.ID(update.Message.Chat.ID),
		Photo:           telego.InputFile{File: imF},
		ReplyParameters: &telego.ReplyParameters{MessageID: update.Message.MessageID},
		ReplyMarkup:     &telego.InlineKeyboardMarkup{InlineKeyboard: rows},
	})
	return nil
}

// TopSwitch handles taps on the category/range buttons: re-renders the card
// in place. Callback data is "sftop|ownerID|category|range"; non-owners get
// an alert instead (same owner-only pattern as the recent/back buttons).
func TopSwitch(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	callback := update.CallbackQuery
	l := i18n.Locale(callback.Message.GetChat())

	parts := strings.Split(callback.Data, "|")
	if len(parts) != 4 {
		return denyOthers(bot, ctx, callback.ID, l("statsfm.top.error"))
	}
	ownerID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || ownerID != callback.From.ID {
		return denyOthers(bot, ctx, callback.ID, l("statsfm.recent.not-yours"))
	}
	cat, ok := parseCategory(parts[2])
	if !ok {
		return denyOthers(bot, ctx, callback.ID, l("statsfm.top.error"))
	}
	rng, ok := parseRange(parts[3])
	if !ok {
		return denyOthers(bot, ctx, callback.ID, l("statsfm.top.error"))
	}

	username := repositories.GetStatsFMUsername(ownerID)
	httpClient := rapidhttp.GetHTTPClient()

	imF, rows, err := renderTop(httpClient, username, ownerID, cat, rng, l)
	if err != nil {
		return denyOthers(bot, ctx, callback.ID, topErrorText(err, l))
	}
	defer func() {
		imF.Close()
		os.Remove(imF.Name())
	}()

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
