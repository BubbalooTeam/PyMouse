package lastfm

import (
	"pymouse/pymouse/client"
	"pymouse/pymouse/config"
	"pymouse/pymouse/helpers/telegram"
	"pymouse/pymouse/middlewares"
	"pymouse/pymouse/modules/lastfm/lfm"
	"pymouse/pymouse/modules/lastfm/set"
	"regexp"

	"github.com/mymmrac/telego/telegohandler"
)

func LoadModule(bS *client.BotStruct) {
	if config.LastFMAPIKey != "" {
		middlewares.Help.RegisterHelp("LastFM")
		bS.Handler.Handle(lfm.NowPlaying,
			telegohandler.Or(
				telegram.Command("lfm", bS.Client.Username()),
				telegram.Command("lt", bS.Client.Username()),
				telegram.Command("lastfm", bS.Client.Username()),
				telegohandler.TextMatches(regexp.MustCompile(`^(?:lt)`)),
			))
		bS.Handler.Handle(set.SetUser, telegram.Command("set", bS.Client.Username()))
		bS.Handler.Handle(set.UnsetUser, telegram.Command("unset", bS.Client.Username()))
	}
}
