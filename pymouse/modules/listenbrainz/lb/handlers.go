package lb

import (
	"fmt"
	"net/url"
	"os"
	"pymouse/pymouse/database/repositories"
	"pymouse/pymouse/helpers/i18n"
	"pymouse/pymouse/helpers/rapidhttp"
	"pymouse/pymouse/helpers/scrobblecard"

	"github.com/cavaliergopher/grab/v3"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
)

func NowPlaying(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	l := i18n.Locale(update.Message.Chat)

	glabClient := grab.NewClient()

	username := repositories.GetListenBrainzUsername(update.Message.From.ID)
	if username == "" {
		bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      l("listenbrainz.nowplaying.no-username-set"),
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
			text = l("listenbrainz.invalid-username")

		case "no listens found.":
			text = l("listenbrainz.nowplaying.no-scrobbles")

		default:
			text = l("listenbrainz.nowplaying.error")
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

	im, _ := scrobblecard.DrawScrobble(
		glabClient,
		scrobblecard.ScrobbleDir("listenbrainz"),
		imageURL,
		trackInfo.Track,
		trackInfo.Artist,
		username,
		listeningText,
		false, // ListenBrainz exposes no per-track loved/feedback state via a keyless GET
		"",
		"",
		fontSet(),
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
						{
							Text: "🎵 " + l("listenbrainz.recent.profile-btn"),
							URL:  "https://listenbrainz.org/user/" + url.PathEscape(username),
						},
						{
							Text:         "📋 " + l("listenbrainz.recent.expand-btn"),
							CallbackData: fmt.Sprintf("lb_recent|%d", update.Message.From.ID),
						},
					},
				},
			},
		},
	)

	return nil
}
