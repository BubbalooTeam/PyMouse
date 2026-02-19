package pm_menu

import (
	"pymouse/pymouse/client"
	"pymouse/pymouse/helpers/telegram"
	"pymouse/pymouse/modules/pm_menu/help"
	"pymouse/pymouse/modules/pm_menu/start"
	"regexp"

	th "github.com/mymmrac/telego/telegohandler"
)

func LoadModule(bS *client.BotStruct) {
	bS.Handler.Handle(help.HelpMenu, th.TextContains("/start help"))
	bS.Handler.Handle(help.HelpMenu, telegram.Command("help"))
	bS.Handler.Handle(help.HelpMenuCallback, th.CallbackDataMatches(regexp.MustCompile(`^HelpMenu$`)))
	bS.Handler.Handle(help.HelpModule, th.CallbackDataPrefix("help:"))
	bS.Handler.Handle(start.StartMessage, telegram.Command("start"))
	bS.Handler.Handle(start.StartBackCallback, th.CallbackDataMatches(regexp.MustCompile(`^StartBack$`)))
}
