package pymouse

import (
	"pymouse/pymouse/modules"
	"pymouse/pymouse/modules/afk"
	"pymouse/pymouse/modules/start"

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

func (bS *BotStruct) Register() {
	bS.Handler.Use(modules.SaveUsers)
	bS.Handler.Handle(start.Start, th.CommandEqual("start"))
	bS.Handler.Handle(afk.SetAFK, th.CommandEqual("afk"))
}
