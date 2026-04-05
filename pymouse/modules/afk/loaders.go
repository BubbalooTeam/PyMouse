package afk

import (
	"pymouse/pymouse/client"
	"pymouse/pymouse/helpers/telegram"
	"pymouse/pymouse/middlewares"
	"regexp"

	"github.com/mymmrac/telego/telegohandler"
)

func LoadModule(bS *client.BotStruct) {
	middlewares.Help.RegisterHelp("AFK")
	bS.Handler.Handle(SetAway, telegohandler.Or(
		telegram.Command("afk", bS.Client.Username()),
		telegohandler.TextMatches(regexp.MustCompile(`^(?:brb)(\s.+)?`)),
	),
	)
	bS.Handler.Use(CheckAway)
}
