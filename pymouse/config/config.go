package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	BotToken       string
	DatabaseURI    string
	Socks5Proxy    string
	TelegramAPIURL string
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	BotToken = os.Getenv("BOT_TOKEN")
	if BotToken == "" {
		log.Fatalf(`In order to initialize this bot, you must insert the "BOT_TOKEN" in the .env file.`)
	}

	DatabaseURI = os.Getenv("DATABASE_URI")
	if DatabaseURI == "" {
		log.Fatalf(`In order to initialize this bot, you must insert the "DATABASE_URI" in the .env file.`)
	}

	Socks5Proxy = os.Getenv("SOCKS5_PROXY")

	TelegramAPIURL = os.Getenv("TELEGRAM_API_URL")
}
