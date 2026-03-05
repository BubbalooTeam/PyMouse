package weather

import (
	"os"
	"pymouse/pymouse/helpers/i18n"
	"pymouse/pymouse/helpers/utils"
	"strings"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/sirupsen/logrus"
)

func WeatherHandler(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	l := i18n.Locale(update.Message.Chat)

	weatherQuery := utils.GetArgs(update)
	if weatherQuery == "" {
		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				Text:      l("weather.checkers.location-not-provided"),
				ParseMode: "HTML",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			},
		)
		return nil
	}

	response, err := GetWeatherResponse(weatherQuery, l)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "failed to extract location informations: location not found"):
			bot.SendMessage(
				ctx,
				&telego.SendMessageParams{
					ChatID:    telegoutil.ID(update.Message.Chat.ID),
					Text:      l("weather.checkers.location-not-found"),
					ParseMode: "HTML",
					ReplyParameters: &telego.ReplyParameters{
						MessageID: update.Message.MessageID,
					},
				},
			)
			return nil
		default:
			logrus.Error(err)
			return nil
		}
	}
	filename, err := MakeWeatherInterface(response, l)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		file.Close()
		os.Remove(filename)
	}()

	bot.SendPhoto(
		ctx,
		&telego.SendPhotoParams{
			ChatID: telegoutil.ID(update.Message.Chat.ID),
			Photo:  telegoutil.File(file),
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		},
	)
	return nil
}
