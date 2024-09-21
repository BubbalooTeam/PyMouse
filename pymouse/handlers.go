package pymouse

import (
	"pymouse/pymouse/modules"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type BotStruct struct {
	Client  *telego.Bot
	Handler *th.BotHandler
}

func NewHandler(bot *telego.Bot, bh *th.BotHandler) *BotStruct {
	return &BotStruct{
		Client:  bot,
		Handler: bh,
	}
}

func (h *BotStruct) Register() {
	h.Handler.Handle(modules.Start, th.CommandEqual("start"))
}
