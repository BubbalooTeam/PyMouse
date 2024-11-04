package client

import (
	"github.com/mymmrac/telego"
)

// CreateBot creates a new Telegram bot instance using telego.
func CreateBot(token string, telegram_api string) (*telego.Bot, error) {
	client, err := telego.NewBot(token)
	if telegram_api != "" {
		client, err = telego.NewBot(token, telego.WithAPIServer(telegram_api))
	}
	return client, err
}

// GetUpdates retrieves updates from the Telegram server.
func GetUpdates(bot *telego.Bot) (<-chan telego.Update, error) {
	updates, err := bot.UpdatesViaLongPolling(&telego.GetUpdatesParams{
		Timeout: 4,
	}, telego.WithLongPollingUpdateInterval(0))
	if err != nil {
		return nil, err
	}
	return updates, nil
}
