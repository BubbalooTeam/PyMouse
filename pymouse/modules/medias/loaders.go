package medias

import (
	"pymouse/pymouse/client"
	"pymouse/pymouse/helpers/telegram"
	"pymouse/pymouse/middlewares"
	"pymouse/pymouse/modules/medias/youtube"
	"regexp"

	"github.com/mymmrac/telego/telegohandler"
)

func LoadModule(bS *client.BotStruct) {
	middlewares.Help.RegisterHelp("Medias", middlewares.Help.WithSubmodules("YouTube"))
	bS.Handler.Handle(youtube.YouTubeDLHandler, telegram.Command("ytdl", bS.Client.Username()))
	bS.Handler.Handle(youtube.YouTubeScrollCallbackHandler, telegohandler.CallbackDataMatches(regexp.MustCompile(`YouTubeScroll\|(.*)$`)))
	bS.Handler.Handle(youtube.YouTubeACallHandler, telegohandler.CallbackDataMatches(regexp.MustCompile(`yt\|(.*)$`)))
}
