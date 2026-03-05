package miscellaneous

import (
	"pymouse/pymouse/client"
	"pymouse/pymouse/helpers/telegram"
	"pymouse/pymouse/middlewares"
	"pymouse/pymouse/modules/miscellaneous/gsmarena"
	"pymouse/pymouse/modules/miscellaneous/weather"
	"regexp"

	th "github.com/mymmrac/telego/telegohandler"
)

func LoadModule(bS *client.BotStruct) {
	middlewares.Help.RegisterHelp("Miscellaneous", middlewares.Help.WithSubmodules("Weather", "GSMArena"))
	bS.Handler.Handle(gsmarena.DeviceSearch, telegram.Command("d"))
	bS.Handler.Handle(gsmarena.DeviceSearchPagination, th.CallbackDataMatches(regexp.MustCompile(`^gsm_page\|(.*)$`)))
	bS.Handler.Handle(gsmarena.DeviceSearchSelect, th.CallbackDataMatches(regexp.MustCompile(`^d\|(.*)$`)))
	bS.Handler.Handle(weather.WeatherHandler, telegram.Command("weather"))
}
