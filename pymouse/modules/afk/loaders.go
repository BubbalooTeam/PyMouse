package afk

import (
	"pymouse/pymouse/client"
	"pymouse/pymouse/middlewares"

	th "github.com/mymmrac/telego/telegohandler"
)

func LoadModule(bS *client.BotStruct) {
	middlewares.Help.RegisterHelp("AFK")
	bS.Handler.Handle(SetAway, th.CommandEqual("afk"))
	bS.Handler.Use(CheckAway)
}
