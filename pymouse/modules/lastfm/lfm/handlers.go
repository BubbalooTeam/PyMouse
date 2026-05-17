package lfm

import (
	"fmt"
	"os"
	"pymouse/pymouse/database/repositories"
	"pymouse/pymouse/helpers/i18n"
	"pymouse/pymouse/helpers/rapidhttp"
	"strings"

	"github.com/cavaliergopher/grab/v3"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
)

func NowPlaying(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	l := i18n.Locale(update.Message.Chat)

	glabClient := grab.NewClient()

	username := repositories.GetLastFMUsername(update.Message.From.ID)
	if username == "" {
		bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      l("lastfm.nowplaying.no-username-set"),
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
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
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      text,
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})

		return nil
	}

	listeningText := getListeningText(trackInfo, l)
	imageURL := trackInfo.Image
	if imageURL == "" {
		imageURL = "https://telegra.ph/file/bdcf492162713ea5633a1.jpg"
	}

	youtubeURL := fmt.Sprintf("https://www.youtube.com/results?search_query=%s+%s", strings.ReplaceAll(trackInfo.Artist, " ", "+"), strings.ReplaceAll(trackInfo.Track, " ", "+"))

	im, _ := DrawScrobble(
		glabClient,
		imageURL,
		trackInfo.Track,
		trackInfo.Artist,
		username,
		listeningText,
		trackInfo.Loved,
		Fonts{
			OpenSans: "pymouse/assets/fonts/opensans.ttf",
			Poppins:  "pymouse/assets/fonts/poppins-semibolditalic.ttf",
			Arial:    "pymouse/assets/fonts/arial.ttf",
		},
		l,
	)

	imF, _ := os.Open(im)
	defer func() {
		imF.Close()
		os.Remove(im)
	}()

	bot.SendPhoto(
		ctx,
		&telego.SendPhotoParams{
			ChatID: telegoutil.ID(update.Message.Chat.ID),
			Photo: telego.InputFile{
				File: imF,
			},
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
