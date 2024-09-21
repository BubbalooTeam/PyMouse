package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"pymouse/pymouse"
	"pymouse/pymouse/client"
	"pymouse/pymouse/config"
	"syscall"

	th "github.com/mymmrac/telego/telegohandler"
)

func main() {
	fmt.Println("\033[0;34mCreating Bot Client...\033")
	botClient, err := client.CreateBot(config.BotToken)
	if err != nil {
		log.Fatalf("Error in creating bot Client: %v", err)
	}

	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, syscall.SIGINT, syscall.SIGTERM)
	botSignal := make(chan struct{}, 1)

	fmt.Println("\033[0;34mBot Created, Starting Get Updates of Long Polling...\033")
	updates, err := client.GetUpdates(botClient)
	if err != nil {
		log.Fatalf("Error in get Updates of TelegramBot: %v", err)
	}
	fmt.Println("\033[0;34mGetUpdates Started With Successfully, Creating Bot Handler...\033")
	botHandler, err := th.NewBotHandler(botClient, updates)
	if err != nil {
		log.Fatalf("Error in Create NewBotHandler: %v", err)
	}
	fmt.Println("\033[0;33mBot Handler Created, Registering Handlers...\033")
	handlerClass := pymouse.NewHandler(botClient, botHandler)
	handlerClass.Register()
	fmt.Println("\033[0;33mHandler Registered, PyMouse is almost starting...\033")

	botUser, err := botClient.GetMe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		<-chanSignal
		fmt.Println("\033[0;31mStopping PyMouse...\033[0m")

		botClient.StopLongPolling()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Long polling stopped.")

		botHandler.Stop()
		fmt.Println("Bot handler stopped.")

		botSignal <- struct{}{}
	}()

	go botHandler.Start()
	fmt.Println("\033[0;32m\U0001F680 Bot Started\033[0m")
	fmt.Printf("\033[0;36mBot Info:\033[0m %v - @%v\n", botUser.FirstName, botUser.Username)

	<-botSignal
	fmt.Println("Done!")
}
