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
	botClient, err := client.CreateBot(config.BotToken)
	if err != nil {
		log.Fatalf("Error in creating bot Client: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{}, 1)

	updates, err := client.GetUpdates(botClient)
	if err != nil {
		log.Fatalf("Error in get Updates of TelegramBot: %v", err)
	}
	botHandler, err := th.NewBotHandler(botClient, updates)
	if err != nil {
		log.Fatalf("Error in create NewBotHandler: %v", err)
	}
	handlerClass := pymouse.NewHandler(botClient, botHandler)
	handlerClass.Register()

	botUser, err := botClient.GetMe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		// Wait for stop signal
		<-sigs
		fmt.Println("\033[0;31mStopping...\033[0m")

		botClient.StopLongPolling()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Long polling stopped")

		botHandler.Stop()
		fmt.Println("Bot handler stopped")

		done <- struct{}{}
	}()

	go botHandler.Start()
	fmt.Println("\033[0;32m\U0001F680 Bot Started\033[0m")
	fmt.Printf("\033[0;36mBot Info:\033[0m %v - @%v\n", botUser.FirstName, botUser.Username)

	<-done
	fmt.Println("Done!")
}
