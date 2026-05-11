package lfm

import (
	"fmt"
	"pymouse/pymouse/database/repositories"
	"pymouse/pymouse/helpers/i18n"
	"pymouse/pymouse/helpers/rapidhttp"
	"strings"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
)

func NowPlaying(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	l := i18n.Locale(update.Message.Chat)

	reply := &telego.ReplyParameters{
		MessageID: update.Message.MessageID,
	}

	username := repositories.GetLastFMUsername(update.Message.From.ID)
	if username == "" {
		bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID:          telegoutil.ID(update.Message.Chat.ID),
			Text:            l("lastfm.nowplaying.no-username-set"),
			ParseMode:       "HTML",
			ReplyParameters: reply,
		})
		return nil
	}

	httpClient := rapidhttp.GetHTTPClient()

	trackInfo, err := getTrack(httpClient, username)
	if err != nil {
		var text string

		switch err.Error() {
		case "invalid username.":
			text = l("lastfm.invalid-username")

		case "no scrobbles found.":
			text = l("lastfm.nowplaying.no-scrobbles")

		default:
			text = l("lastfm.nowplaying.error")
		}

		bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID:          telegoutil.ID(update.Message.Chat.ID),
			Text:            text,
			ParseMode:       "HTML",
			ReplyParameters: reply,
		})

		return nil
	}

	textKey := "lastfm.nowplaying.was-listening"
	if trackInfo.Now {
		textKey = "lastfm.nowplaying.is-listening"
	}
	imageURL := trackInfo.Image
	if imageURL == "" {
		imageURL = "https://telegra.ph/file/3ad207681d56059a7d90d.jpg"
	}

	nowText := fmt.Sprintf(
		l(textKey),
		update.Message.From.FirstName,
		trackInfo.Artist,
		trackInfo.Track,
		trackInfo.Playcount,
	)
	if trackInfo.Loved {
		nowText += l("lastfm.nowplaying.loved")
	}
	youtubeURL := fmt.Sprintf("https://www.youtube.com/results?search_query=%s+%s", strings.ReplaceAll(trackInfo.Artist, " ", "+"), strings.ReplaceAll(trackInfo.Track, " ", "+"))

	bot.SendPhoto(
		ctx,
		&telego.SendPhotoParams{
			ChatID: telegoutil.ID(update.Message.Chat.ID),
			Photo: telego.InputFile{
				URL: imageURL,
			},
			Caption:   nowText,
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
			ReplyMarkup: &telego.InlineKeyboardMarkup{
				InlineKeyboard: [][]telego.InlineKeyboardButton{
					{
						telego.InlineKeyboardButton{
							Text: "📽️",
							URL:  youtubeURL,
						},
					},
				},
			},
		},
	)

	return nil
}
