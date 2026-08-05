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
	DownloadPath   string
	LastFMAPIKey   string
	LogChannelID   int64
	OwnerID        int64
	Socks5Proxy    string
	TelegramAPIURL string
	WebhookURL     string

	// Cloudflare R2 (optional). When all four are set, the last.fm inline
	// result uploads its card image to R2 and serves it via R2PublicURL.
	// When unset, the inline result falls back to a plain text article.
	R2AccountID    string
	R2Bucket       string
	R2AccessKeyID  string
	R2SecretKey    string
	R2PublicURL    string
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

	DownloadPath = os.Getenv("DOWNLOAD_PATH")
	if DownloadPath == "" {
		DownloadPath = "pymouse/downloads"
	}

	LogChannelID, err = strconv.ParseInt(os.Getenv("LOG_CHANNEL_ID"), 10, 64)
	if err != nil || LogChannelID == 0 {
		logrus.Fatal(`In order to initialize this bot, you must set the "LOG_CHANNEL_ID" in the .env file.`)
	}
	OwnerID, err = strconv.ParseInt(os.Getenv("OWNER_ID"), 10, 64)
	if err != nil || OwnerID == 0 {
		logrus.Fatal(`In order to initialize this bot, you must set the "OWNER_ID" in the .env file.`)
	}

	LastFMAPIKey = os.Getenv("LASTFM_API_KEY")
	Socks5Proxy = os.Getenv("SOCKS5_PROXY")
	TelegramAPIURL = os.Getenv("TELEGRAM_API_URL")
	WebhookURL = os.Getenv("WEBHOOK_URL")

	// R2 is optional — only read if present.
	R2AccountID = os.Getenv("R2_ACCOUNT_ID")
	R2Bucket = os.Getenv("R2_BUCKET")
	R2AccessKeyID = os.Getenv("R2_ACCESS_KEY_ID")
	R2SecretKey = os.Getenv("R2_SECRET_ACCESS_KEY")
	R2PublicURL = os.Getenv("R2_PUBLIC_URL")
}
