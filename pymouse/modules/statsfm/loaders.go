package statsfm

import (
	"pymouse/pymouse/client"
	"pymouse/pymouse/helpers/telegram"
	"pymouse/pymouse/middlewares"
	"pymouse/pymouse/modules/statsfm/set"
	"pymouse/pymouse/modules/statsfm/sf"
	"regexp"

	"github.com/mymmrac/telego/telegohandler"
)

// LoadModule registers the stats.fm commands. Like ListenBrainz, reading a
// stats.fm profile needs no API key (see sf/utils.go) — it's an
// unauthenticated, undocumented-but-public API reverse-engineered from their
// own web app, so this module is always on.
func LoadModule(bS *client.BotStruct) {
	middlewares.Help.RegisterHelp("stats.fm")
	bS.Handler.Handle(sf.NowPlaying,
		telegohandler.Or(
			telegram.Command("sf", bS.Client.Username()),
			telegram.Command("statsfm", bS.Client.Username()),
			telegohandler.TextMatches(regexp.MustCompile(`^(?:sf)`)),
		))
	bS.Handler.Handle(set.SetUser, telegram.Command("sfset", bS.Client.Username()))
	bS.Handler.Handle(set.UnsetUser, telegram.Command("sfunset", bS.Client.Username()))
	// /sftop: ranked top tracks/artists/albums/genres, with buttons to switch
	// category and time range (today / 7 days / 4 weeks / 6 months / lifetime).
	bS.Handler.Handle(sf.Top, telegram.Command("sftop", bS.Client.Username()))
	bS.Handler.Handle(sf.TopSwitch, telegohandler.CallbackDataMatches(regexp.MustCompile(`^sftop\|(\d+)\|(\w+)\|(\w+)$`)))
	// "+" expands the now-playing card into the last 5 streams; "back"
	// restores it. Both are owner-only (validated inside the handlers).
	bS.Handler.Handle(sf.ExpandRecent, telegohandler.CallbackDataMatches(regexp.MustCompile(`^sf_recent\|(\d+)$`)))
	bS.Handler.Handle(sf.BackToNowPlaying, telegohandler.CallbackDataMatches(regexp.MustCompile(`^sf_back\|(\d+)$`)))
	// Inline mode: @PyMouseBOT sf (also "statsfm") shows the caller's
	// now-playing. Requires inline mode enabled in BotFather (/setinline).
	bS.Handler.HandleInlineQuery(sf.NowPlayingInline,
		telegohandler.InlineQueryMatches(regexp.MustCompile(`(?i)\b(sf|statsfm)\b`)))
}
