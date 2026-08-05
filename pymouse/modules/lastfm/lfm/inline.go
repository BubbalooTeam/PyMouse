package lfm

import (
	"fmt"
	"net/url"
	"pymouse/pymouse/database/repositories"
	"pymouse/pymouse/helpers/i18n"
	"pymouse/pymouse/helpers/rapidhttp"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
)

// NowPlayingInline answers an inline query (@PyMouseBOT lfm) with the caller's
// current/last scrobbled track. Tapping the result sends a formatted text
// message (with a link to the user's Last.fm profile) into the chat — the card
// image is intentionally not used here because inline photo results require a
// publicly-hosted PhotoURL, which this bot doesn't have.
//
// The query text is ignored beyond triggering: the now-playing shown is always
// the caller's own (mirroring /lfm). Requires inline mode enabled in BotFather
// (/setinline) and "inline_query" in AllowedUpdates (see client/bot.go).
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

	// No username set -> offer one article that deep-links to /set.
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
					Description: empty,
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
		// Don't surface API errors as a visible result; answer empty so the
		// inline picker shows "no results" instead of a broken article.
		return bot.AnswerInlineQuery(ctx, &telego.AnswerInlineQueryParams{
			InlineQueryID: query.ID,
			CacheTime:     30,
			IsPersonal:    true,
		})
	}

	profileURL := "https://www.last.fm/user/" + url.PathEscape(username)
	messageText := fmt.Sprintf(
		"🎵 <b>%s</b> — <i>%s</i>\n👤 <a href=\"%s\">%s</a>",
		escapeHTML(trackInfo.Track),
		escapeHTML(trackInfo.Artist),
		profileURL,
		l("lastfm.inline.profile"),
	)

	return bot.AnswerInlineQuery(ctx, &telego.AnswerInlineQueryParams{
		InlineQueryID: query.ID,
		CacheTime:     30, // short: now-playing changes constantly
		IsPersonal:    true,
		Results: []telego.InlineQueryResult{
			&telego.InlineQueryResultArticle{
				Type:        "article",
				ID:          "lfm-nowplaying",
				Title:       "🎵 " + trackInfo.Track,
				Description: trackInfo.Artist,
				InputMessageContent: &telego.InputTextMessageContent{
					MessageText: messageText,
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
