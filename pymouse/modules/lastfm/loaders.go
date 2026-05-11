package lastfm

import (
	"pymouse/pymouse/client"
	"pymouse/pymouse/helpers/telegram"
	"pymouse/pymouse/modules/lastfm/lfm"
	"pymouse/pymouse/modules/lastfm/set"
)

func LoadModule(bS *client.BotStruct) {
	bS.Handler.Handle(lfm.NowPlaying, telegram.Command("lfm", bS.Client.Username()))
	bS.Handler.Handle(set.SetUser, telegram.Command("set", bS.Client.Username()))
	bS.Handler.Handle(set.UnsetUser, telegram.Command("unset", bS.Client.Username()))
}
