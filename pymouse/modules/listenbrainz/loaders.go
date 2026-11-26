package listenbrainz

import (
	"pymouse/pymouse/client"
	"pymouse/pymouse/helpers/telegram"
	"pymouse/pymouse/middlewares"
	"pymouse/pymouse/modules/listenbrainz/lb"
	"pymouse/pymouse/modules/listenbrainz/set"
	"regexp"

	"github.com/mymmrac/telego/telegohandler"
)

// LoadModule registers the ListenBrainz commands. Unlike Last.fm, reading
// listens from ListenBrainz needs no API key (see lb/utils.go), so this
// module is always on — there's no key to blank out to disable it.
func LoadModule(bS *client.BotStruct) {
	middlewares.Help.RegisterHelp("ListenBrainz")
	bS.Handler.Handle(lb.NowPlaying,
		telegohandler.Or(
			telegram.Command("lb", bS.Client.Username()),
			telegram.Command("listenbrainz", bS.Client.Username()),
			telegohandler.TextMatches(regexp.MustCompile(`^(?:lb)`)),
		))
	bS.Handler.Handle(set.SetUser, telegram.Command("lbset", bS.Client.Username()))
	bS.Handler.Handle(set.UnsetUser, telegram.Command("lbunset", bS.Client.Username()))
	// "+" expands the now-playing card into the last 5 listens; "back"
	// restores it. Both are owner-only (validated inside the handlers).
	bS.Handler.Handle(lb.ExpandRecent, telegohandler.CallbackDataMatches(regexp.MustCompile(`^lb_recent\|(\d+)$`)))
	bS.Handler.Handle(lb.BackToNowPlaying, telegohandler.CallbackDataMatches(regexp.MustCompile(`^lb_back\|(\d+)$`)))
	// Inline mode: @PyMouseBOT lb (also "listenbrainz") shows the caller's
	// now-playing. Requires inline mode enabled in BotFather (/setinline).
	bS.Handler.HandleInlineQuery(lb.NowPlayingInline,
		telegohandler.InlineQueryMatches(regexp.MustCompile(`(?i)\b(lb|listenbrainz)\b`)))
}
