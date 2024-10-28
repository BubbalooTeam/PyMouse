package main

import (
	"log"
	"os"
	"os/signal"
	"pymouse/pymouse"
	"pymouse/pymouse/client"
	"pymouse/pymouse/config"
	"pymouse/pymouse/database"
	"syscall"

	th "github.com/mymmrac/telego/telegohandler"
)

func main() {
	database.InitDB()

	log.Println("\033[0;33mCreating Bot Client...\033[0m")
	botClient, err := client.CreateBot(config.BotToken)
	if err != nil {
		log.Fatalf("\033[0;31mError in creating bot Client: %v\033[0m", err)
	}

	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, syscall.SIGINT, syscall.SIGTERM)
	botSignal := make(chan struct{}, 1)

	log.Println("\033[0;32mBot Created, Starting Get Updates of Long Polling...\033[0m")
	updates, err := client.GetUpdates(botClient)
	if err != nil {
		log.Fatalf("\033[0;31mError in get Updates of TelegramBot: %v\033[0m", err)
	}
	log.Println("\033[0;32mGetUpdates Started With Successfully, Creating Bot Handler...\033[0m")

	botHandler, err := th.NewBotHandler(botClient, updates)
	if err != nil {
		log.Fatalf("\033[0;31mError in Create NewBotHandler: %v\033[0m", err)
	}
	log.Println("\033[0;32mBot Handler Created, Registering Handlers...\033[0m")

	handlerClass := pymouse.NewHandler(botClient, botHandler)
	handlerClass.Register()

	log.Println("\033[0;33mHandler Registered, PyMouse is almost starting...\033[0m")

	botUser, err := botClient.GetMe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		<-chanSignal
		log.Println("\033[0;31mStopping PyMouse...\033[0m")

		botClient.StopLongPolling()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("\033[0;32mLong polling stopped.\033[0m")

		botHandler.Stop()
		log.Println("\033[0;32mBot handler stopped.\033[0m")

		defer database.CloseDB()

		botSignal <- struct{}{}
	}()

	go botHandler.Start()
	log.Println("\033[0;32m\U0001F680 Bot Started\033[0m")
	log.Printf("\033[0;36mBot Info:\033[0m %v - @%v\n", botUser.FirstName, botUser.Username)

	<-botSignal
	log.Println("\033[0;34mDone!\033[0m")
}
