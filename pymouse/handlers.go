package pymouse

import (
	"pymouse/pymouse/modules/afk"
	"pymouse/pymouse/modules/checkers"
	"pymouse/pymouse/modules/medias/youtube"
	"pymouse/pymouse/modules/start"
	"regexp"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type BotStruct struct {
	Client  *telego.Bot
	Handler *th.BotHandler
}

func NewHandler(bot *telego.Bot, botHandler *th.BotHandler) *BotStruct {
	return &BotStruct{
		Client:  bot,
		Handler: botHandler,
	}
}

func (bS *BotStruct) Register() {
	// Checkers of DataBase and Utilities of BOT on bot_incoming.
	bS.Handler.Use(checkers.SaveUsers)
	bS.Handler.Use(checkers.SaveChats)
	// Checkers of Database of Utilities
	bS.Handler.Use(afk.CheckAway)

	// Bot Commands, comming soon, add a dinamic commands loader.
	bS.Handler.Handle(start.Start, th.CommandEqual("start"))
	bS.Handler.Handle(afk.SetAway, th.CommandEqual("afk"))
	bS.Handler.Handle(youtube.GetYoutubeMedias, th.CommandEqual("ytdl"))
	bS.Handler.Handle(youtube.YouTubeScrollCallback, th.CallbackDataMatches(regexp.MustCompile(`YouTubeScroll\|(.*)$`)))
	bS.Handler.Handle(youtube.YouTubeACallHandler, th.CallbackDataMatches(regexp.MustCompile(`yt\|(.*)$`)))
}
