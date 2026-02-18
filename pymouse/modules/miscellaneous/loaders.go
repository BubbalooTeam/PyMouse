package miscellaneous

import (
	"pymouse/pymouse/client"
	"pymouse/pymouse/middlewares"
	"pymouse/pymouse/modules/miscellaneous/gsmarena"
	"regexp"

	th "github.com/mymmrac/telego/telegohandler"
)

func LoadModule(bS *client.BotStruct) {
	middlewares.Help.RegisterHelp("Miscellaneous", middlewares.Help.WithSubmodules("Weather", "GSMArena"))
	bS.Handler.Handle(gsmarena.DeviceSearch, th.CommandEqual("d"))
	bS.Handler.Handle(gsmarena.DeviceSearchPagination, th.CallbackDataMatches(regexp.MustCompile(`^gsm_page\|(.*)$`)))
	bS.Handler.Handle(gsmarena.DeviceSearchSelect, th.CallbackDataMatches(regexp.MustCompile(`^d\|(.*)$`)))
}
