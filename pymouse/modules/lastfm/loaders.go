package lastfm

import (
	"pymouse/pymouse/client"
	"pymouse/pymouse/config"
	"pymouse/pymouse/helpers/telegram"
	"pymouse/pymouse/middlewares"
	"pymouse/pymouse/modules/lastfm/lfm"
	"pymouse/pymouse/modules/lastfm/set"
)

func LoadModule(bS *client.BotStruct) {
	if config.LastFMAPIKey != "" {
		middlewares.Help.RegisterHelp("LastFM")
		bS.Handler.Handle(lfm.NowPlaying, telegram.Command("lfm", bS.Client.Username()))
		bS.Handler.Handle(set.SetUser, telegram.Command("set", bS.Client.Username()))
		bS.Handler.Handle(set.UnsetUser, telegram.Command("unset", bS.Client.Username()))
	}
}
