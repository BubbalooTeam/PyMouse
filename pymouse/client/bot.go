package client

import (
	"context"

	"github.com/mymmrac/telego"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// CreateBot creates a new Telegram bot instance using telego.
func CreateBot(token string, telegram_api string) (*telego.Bot, error) {
	client, err := telego.NewBot(token, telego.WithLogger(logrus.StandardLogger()))
	if telegram_api != "" {
		client, err = telego.NewBot(token, telego.WithAPIServer(telegram_api), telego.WithLogger(logrus.StandardLogger()))
	}
	return client, err
}

// GetUpdates retrieves updates from the Telegram server.
func GetUpdates(ctx context.Context, bot *telego.Bot, webhookURL string) (<-chan telego.Update, error) {
	var updates <-chan telego.Update
	var err error

	bot.DeleteWebhook(
		ctx,
		&telego.DeleteWebhookParams{
			DropPendingUpdates: true,
		},
	)
	if webhookURL == "" {
		updates, err = bot.UpdatesViaLongPolling(ctx, &telego.GetUpdatesParams{
			Timeout: 4,
		}, telego.WithLongPollingUpdateInterval(0))
		if err != nil {
			return nil, err
		}
	} else {
		err := bot.SetWebhook(
			ctx,
			&telego.SetWebhookParams{
				URL: webhookURL + bot.Token(),
			},
		)
		if err != nil {
			return nil, err
		}
		updates, err = bot.UpdatesViaWebhook(
			ctx,
			telego.WebhookFastHTTP(
				&fasthttp.Server{},
				"/bot"+bot.Token(),
			),
		)
	}

	return updates, nil
}
