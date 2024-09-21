package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var (
	BotToken     string
	DatabaseFile string
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	BotToken = os.Getenv("BOT_TOKEN")
	if BotToken == "" {
		log.Fatalf(`In order to initialize this bot, you must insert the "BOT_TOKEN" in the .env file.`)
	}

	DatabaseFile = os.Getenv("DATABASE_FILE")
	if DatabaseFile == "" {
		DatabaseFile = filepath.Join(".", "pymouse", "database", "files", "database.json")
	}
}
