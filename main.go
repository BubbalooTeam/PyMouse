package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"pymouse/pymouse"
	"pymouse/pymouse/client"
	"pymouse/pymouse/config"
	"pymouse/pymouse/database"
	"pymouse/pymouse/helpers/i18n"
	"pymouse/pymouse/middlewares"
	"runtime"
	"strings"
	"syscall"

	"github.com/lrstanley/go-ytdlp"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/sirupsen/logrus"
)

func main() {
	database.InitDB()

	logrus.Info("[CompileLocales]: Loading localization files...")
	err := i18n.CompileLocales()
	if err != nil {
		logrus.Errorf("[CompileLocales]: %v", err)
		return
	}
	logrus.Info("[CompileLocales]: All location files have been loaded.")

	logrus.Info("Creating Bot Client...")
	botClient, err := client.CreateBot(config.BotToken, config.TelegramAPIURL)
	if err != nil {
		logrus.Fatalf("Error in creating bot Client: %v", err)
	}

	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, syscall.SIGINT, syscall.SIGTERM)
	botSignal := make(chan struct{}, 1)

	ctx := context.Background()

	logrus.Info("Bot Created, Starting Get Updates of Long Polling...")

	ytdlp.MustInstall(ctx, &ytdlp.InstallOptions{})

	updates, err := client.GetUpdates(ctx, botClient, config.WebhookURL)
	if err != nil {
		logrus.Fatalf("Error in get Updates of Telegram-Bot: %v", err)
	}
	logrus.Info("GetUpdates Started With Successfully, Creating Bot Handler...")

	botHandler, err := th.NewBotHandler(botClient, updates)
	if err != nil {
		logrus.Fatalf("Error in Create NewBotHandler: %v", err)
	}
	logrus.Info("Bot Handler Created, Registering Handlers...")

	middlewares.NewHelp()
	handlerClass := client.NewHandler(botClient, botHandler)
	pymouse.Register(handlerClass)

	logrus.Info("Handler Registered, PyMouse is almost starting...")

	botUser, err := botClient.GetMe(ctx)
	if err != nil {
		logrus.Fatal(err)
	}

	systemName, err := exec.Command("uname", "-sr").Output()
	if err != nil {
		logrus.Errorf("Error getting system name: %v", err)
	}

	go func() {
		<-chanSignal
		logrus.Info("Stopping PyMouse...")

		if err != nil {
			logrus.Fatal(err)
		}

		if config.WebhookURL != "" {
			botClient.DeleteWebhook(
				ctx,
				&telego.DeleteWebhookParams{
					DropPendingUpdates: true,
				},
			)
		}

		botHandler.Stop()
		logrus.Info("Bot handler stopped.")

		_, err = botClient.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID: telegoutil.ID(config.LogChannelID),
				Text: fmt.Sprintf(
					"<b>🚀 PyMouse stopped!</b>\n\n<b><i>System:</i></b> <code>%s</code>\n<b><i>GoLang:</i></b> <code>%s</code>",
					strings.ReplaceAll(strings.TrimSpace(string(systemName)), "\n", ""),
					runtime.Version(),
				),
				ParseMode: "HTML",
			},
		)
		if err != nil {
			logrus.Fatalf("The 'LOG_CHANNEL_ID' parameter in the .env file is invalid.")
		}

		defer database.CloseDB()

		botSignal <- struct{}{}
	}()

	go botHandler.Start()

	_, err = botClient.SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID: telegoutil.ID(config.LogChannelID),
			Text: fmt.Sprintf(
				"<b>🚀 PyMouse started!</b>\n\n<b><i>System:</i></b> <code>%s</code>\n<b><i>GoLang:</i></b> <code>%s</code>",
				strings.ReplaceAll(strings.TrimSpace(string(systemName)), "\n", ""),
				runtime.Version(),
			),
			ParseMode: "HTML",
		},
	)
	if err != nil {
		logrus.Fatalf("The 'LOG_CHANNEL_ID' parameter in the .env file is invalid.")
	}
	logrus.Info("\U0001F680 Bot Started!")
	logrus.Infof("Bot Info: %v - @%v", botUser.FirstName, botUser.Username)

	<-botSignal
	logrus.Info("Done!")
}
