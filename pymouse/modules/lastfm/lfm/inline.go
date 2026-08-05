package lfm

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"pymouse/pymouse/database/repositories"
	"pymouse/pymouse/helpers/i18n"
	"pymouse/pymouse/helpers/rapidhttp"

	"github.com/cavaliergopher/grab/v3"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
)

// NowPlayingInline answers an inline query (@PyMouseBOT lfm) with the caller's
// current/last scrobbled track. The result is the same card image the /lfm
// command sends, hosted on telegra.ph so Telegram can serve it inline, with the
// YouTube / Profile / Share buttons attached.
//
// The query text is ignored beyond triggering: the now-playing shown is always
// the caller's own (mirroring /lfm). Requires inline mode enabled in BotFather
// (/setinline) and "inline_query" in AllowedUpdates (see client/bot.go).
//
// If anything fails (no username set, last.fm error, image draw, telegraph
// upload), the handler falls back to a plain text article so the inline picker
// never appears broken.
func NowPlayingInline(ctx *telegohandler.Context, query telego.InlineQuery) error {
	bot := ctx.Bot()
	// Inline queries have no chat, but the caller's language preference is
	// stored per-user (same path PM chats use). Pass the user id as a private
	// chat so i18n.Locale -> GetChatLanguage resolves it; falls back to en_us.
	l := i18n.Locale(telego.Chat{
		ID:   query.From.ID,
		Type: telego.ChatTypePrivate,
	})

	username := repositories.GetLastFMUsername(query.From.ID)

	// No username set -> offer one article that tells them to /set.
	if username == "" {
		empty := l("lastfm.nowplaying.no-username-set")
		return bot.AnswerInlineQuery(ctx, &telego.AnswerInlineQueryParams{
			InlineQueryID: query.ID,
			CacheTime:     300,
			IsPersonal:    true,
			Results: []telego.InlineQueryResult{
				&telego.InlineQueryResultArticle{
					Type:        "article",
					ID:          "lfm-noset",
					Title:       "Last.fm",
					Description: stripHTML(empty),
					InputMessageContent: &telego.InputTextMessageContent{
						MessageText: empty,
						ParseMode:   "HTML",
					},
				},
			},
		})
	}

	httpClient := rapidhttp.GetHTTPClient()
	trackInfo, err := getTrack(httpClient, username)
	if err != nil {
		return answerEmpty(bot, ctx, query.ID)
	}

	imageURL := trackInfo.Image
	if imageURL == "" {
		imageURL = "https://telegra.ph/file/bdcf492162713ea5633a1.jpg"
	}
	listeningText := getListeningText(trackInfo, l)

	// Generate the same card image the /lfm command sends, then host it on
	// telegra.ph so it can be used as an inline photo result.
	imgPath, derr := DrawScrobble(
		grab.NewClient(), imageURL, trackInfo.Track, trackInfo.Artist,
		username, listeningText, trackInfo.Loved, fontSet(), l,
	)
	if derr != nil {
		return answerTextFallback(bot, ctx, query.ID, trackInfo, username, l)
	}
	defer os.Remove(imgPath)

	photoURL, uerr := uploadToTelegraph(imgPath)
	if uerr != nil {
		return answerTextFallback(bot, ctx, query.ID, trackInfo, username, l)
	}

	caption := fmt.Sprintf(
		"🎵 <b>%s</b> — <i>%s</i>\n👤 <a href=\"%s\">%s</a>",
		escapeHTML(trackInfo.Track),
		escapeHTML(trackInfo.Artist),
		"https://www.last.fm/user/"+url.PathEscape(username),
		l("lastfm.inline.profile"),
	)

	shareQuery := "lfm"
	rows := [][]telego.InlineKeyboardButton{{
		{Text: "▶️ YouTube", URL: trackInfo.YouTubeURL},
		{Text: "👤 " + l("lastfm.recent.profile-btn"), URL: "https://www.last.fm/user/" + url.PathEscape(username)},
		{Text: "🔗 " + l("lastfm.recent.share-btn"), SwitchInlineQuery: &shareQuery},
	}}

	return bot.AnswerInlineQuery(ctx, &telego.AnswerInlineQueryParams{
		InlineQueryID: query.ID,
		CacheTime:     30, // short: now-playing changes constantly
		IsPersonal:    true,
		Results: []telego.InlineQueryResult{
			&telego.InlineQueryResultPhoto{
				Type:        "photo",
				ID:          "lfm-nowplaying",
				PhotoURL:    photoURL,
				ThumbnailURL: photoURL,
				PhotoWidth:  600,
				PhotoHeight: 250,
				Title:       "🎵 " + trackInfo.Track,
				Description: trackInfo.Artist,
				Caption:     caption,
				ParseMode:   "HTML",
				ReplyMarkup: &telego.InlineKeyboardMarkup{InlineKeyboard: rows},
			},
		},
	})
}

// answerEmpty replies with no results so the picker shows "nothing found"
// instead of a broken article.
func answerEmpty(bot *telego.Bot, ctx context.Context, queryID string) error {
	return bot.AnswerInlineQuery(ctx, &telego.AnswerInlineQueryParams{
		InlineQueryID: queryID,
		CacheTime:     30,
		IsPersonal:    true,
	})
}

// answerTextFallback is used when the image cannot be produced/hosted: it sends
// the now-playing as a plain text article instead, so the share still works.
func answerTextFallback(
	bot *telego.Bot, ctx context.Context, queryID string,
	trackInfo LastFMTrackInformations, username string, l func(string) string,
) error {
	text := fmt.Sprintf(
		"🎵 <b>%s</b> — <i>%s</i>\n👤 <a href=\"%s\">%s</a>",
		escapeHTML(trackInfo.Track),
		escapeHTML(trackInfo.Artist),
		"https://www.last.fm/user/"+url.PathEscape(username),
		l("lastfm.inline.profile"),
	)
	return bot.AnswerInlineQuery(ctx, &telego.AnswerInlineQueryParams{
		InlineQueryID: queryID,
		CacheTime:     30,
		IsPersonal:    true,
		Results: []telego.InlineQueryResult{
			&telego.InlineQueryResultArticle{
				Type:        "article",
				ID:          "lfm-nowplaying",
				Title:       "🎵 " + trackInfo.Track,
				Description: trackInfo.Artist,
				InputMessageContent: &telego.InputTextMessageContent{
					MessageText: text,
					ParseMode:   "HTML",
				},
			},
		},
	})
}

// escapeHTML replaces the few characters that can break an HTML-formatted
// Telegram message. Used because track/artist names are user-controlled.
func escapeHTML(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		switch r {
		case '<':
			out = append(out, []rune("&lt;")...)
		case '>':
			out = append(out, []rune("&gt;")...)
		case '&':
			out = append(out, []rune("&amp;")...)
		default:
			out = append(out, r)
		}
	}
	return string(out)
}

// stripHTML removes <b>/<i>/<code> tags from a localized string so it can be
// used as a plain-text description (inline result descriptions don't render
// HTML). Good enough for the few tags our locales use.
func stripHTML(s string) string {
	var out []rune
	inTag := false
	for _, r := range s {
		switch r {
		case '<':
			inTag = true
		case '>':
			inTag = false
		default:
			if !inTag {
				out = append(out, r)
			}
		}
	}
	return string(out)
}
