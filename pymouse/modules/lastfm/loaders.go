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
		// Inline mode: @PyMouseBOT lfm (also "lt", "lastfm") shows the caller's
		// now-playing. Requires inline mode enabled in BotFather (/setinline).
		bS.Handler.HandleInlineQuery(lfm.NowPlayingInline,
			telegohandler.InlineQueryMatches(regexp.MustCompile(`(?i)\b(lfm|lt|lastfm)\b`)))
	}
}
