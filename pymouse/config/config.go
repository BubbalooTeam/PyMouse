package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	BotToken       string
	DatabaseURI    string
	Socks5Proxy    string
	TelegramAPIURL string
)

func init() {
	err := ConfigureLogging()
	if err != nil {
		logrus.Errorf("Error in configuring your logger: %v", err)
	}
	if err = godotenv.Load(); err != nil {
		logrus.Errorf("Error loading .env file: %v", err)
	}
	BotToken = os.Getenv("BOT_TOKEN")
	if BotToken == "" {
		logrus.Errorf(`In order to initialize this bot, you must insert the "BOT_TOKEN" in the .env file.`)
	}

	DatabaseURI = os.Getenv("DATABASE_URI")
	if DatabaseURI == "" {
		logrus.Errorf(`In order to initialize this bot, you must insert the "DATABASE_URI" in the .env file.`)
	}

	Socks5Proxy = os.Getenv("SOCKS5_PROXY")

	TelegramAPIURL = os.Getenv("TELEGRAM_API_URL")
}
