package client

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type BotStruct struct {
	Client  *telego.Bot
	Handler *th.BotHandler
}

func NewHandler(bot *telego.Bot, botHandler *th.BotHandler) *BotStruct {
	return &BotStruct{
		Client:  bot,
		Handler: botHandler,
	}
}
