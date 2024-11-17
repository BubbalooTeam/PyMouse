package main

import (
	"os"
	"os/signal"
	"pymouse/pymouse"
	"pymouse/pymouse/client"
	"pymouse/pymouse/config"
	"pymouse/pymouse/database"
	"syscall"

	th "github.com/mymmrac/telego/telegohandler"
	"github.com/sirupsen/logrus"
)

func main() {
	database.InitDB()

	logrus.Info("Creating Bot Client...")
	botClient, err := client.CreateBot(config.BotToken, config.TelegramAPIURL)
	if err != nil {
		logrus.Fatalf("Error in creating bot Client: %v", err)
	}

	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, syscall.SIGINT, syscall.SIGTERM)
	botSignal := make(chan struct{}, 1)

	logrus.Info("Bot Created, Starting Get Updates of Long Polling...")
	updates, err := client.GetUpdates(botClient)
	if err != nil {
		logrus.Fatalf("Error in get Updates of Telegram-Bot: %v", err)
	}
	logrus.Info("GetUpdates Started With Successfully, Creating Bot Handler...")

	botHandler, err := th.NewBotHandler(botClient, updates)
	if err != nil {
		logrus.Fatalf("Error in Create NewBotHandler: %v", err)
	}
	logrus.Info("Bot Handler Created, Registering Handlers...")

	handlerClass := pymouse.NewHandler(botClient, botHandler)
	handlerClass.Register()

	logrus.Info("Handler Registered, PyMouse is almost starting...")

	botUser, err := botClient.GetMe()
	if err != nil {
		logrus.Fatal(err)
	}
	go func() {
		<-chanSignal
		logrus.Info("Stopping PyMouse...")

		botClient.StopLongPolling()
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Info("Long polling stopped.")

		botHandler.Stop()
		logrus.Info("Bot handler stopped.")

		defer database.CloseDB()

		botSignal <- struct{}{}
	}()

	go botHandler.Start()
	logrus.Info("\U0001F680 Bot Started!")
	logrus.Infof("Bot Info: %v - @%v", botUser.FirstName, botUser.Username)

	<-botSignal
	logrus.Info("Done!")
}
