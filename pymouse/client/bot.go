package client

import (
	"context"

	"github.com/mymmrac/telego"
	"github.com/sirupsen/logrus"
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
func GetUpdates(ctx context.Context, bot *telego.Bot) (<-chan telego.Update, error) {
	updates, err := bot.UpdatesViaLongPolling(ctx, &telego.GetUpdatesParams{
		Timeout: 4,
	}, telego.WithLongPollingUpdateInterval(0))
	if err != nil {
		return nil, err
	}
	return updates, nil
}
