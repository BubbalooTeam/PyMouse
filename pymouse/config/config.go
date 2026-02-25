package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	BotToken       string
	DatabaseURI    string
	LogChannelID   int64
	OwnerID        int64
	Socks5Proxy    string
	TelegramAPIURL string
	WebhookURL     string
)

func init() {
	err := ConfigureLogging()
	if err != nil {
		logrus.Errorf("Error in configuring your logger: %v", err)
	}
	if err = godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}
	BotToken = os.Getenv("BOT_TOKEN")
	if BotToken == "" {
		logrus.Fatal(`In order to initialize this bot, you must insert the "BOT_TOKEN" in the .env file.`)
	}

	DatabaseURI = os.Getenv("DATABASE_URI")
	if DatabaseURI == "" {
		logrus.Fatal(`In order to initialize this bot, you must insert the "DATABASE_URI" in the .env file.`)
	}
	LogChannelID, err = strconv.ParseInt(os.Getenv("LOG_CHANNEL_ID"), 10, 64)
	if err != nil || LogChannelID == 0 {
		logrus.Fatal(`In order to initialize this bot, you must set the "LOG_CHANNEL_ID" in the .env file.`)
	}
	OwnerID, err = strconv.ParseInt(os.Getenv("OWNER_ID"), 10, 64)
	if err != nil || OwnerID == 0 {
		logrus.Fatal(`In order to initialize this bot, you must set the "OWNER_ID" in the .env file.`)
	}

	Socks5Proxy = os.Getenv("SOCKS5_PROXY")
	TelegramAPIURL = os.Getenv("TELEGRAM_API_URL")
	WebhookURL = os.Getenv("WEBHOOK_URL")
}
