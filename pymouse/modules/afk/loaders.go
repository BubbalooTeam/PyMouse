package afk

import (
	"pymouse/pymouse/client"
	"pymouse/pymouse/helpers/telegram"
	"pymouse/pymouse/middlewares"
)

func LoadModule(bS *client.BotStruct) {
	middlewares.Help.RegisterHelp("AFK")
	bS.Handler.Handle(SetAway, telegram.Command("afk"))
	bS.Handler.Use(CheckAway)
}
