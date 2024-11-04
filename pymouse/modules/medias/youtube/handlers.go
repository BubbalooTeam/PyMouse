package youtube

import (
	"fmt"
	"log"
	"pymouse/pymouse/helpers/utils"
	"regexp"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/raitonoberu/ytsearch"
)

type VidCache struct {
	UUID            string
	VidInformations []*ytsearch.VideoItem
}

var (
	YouTubeRegex_URL = regexp.MustCompile(`(?m)http(?:s?):\/\/(?:www\.)?(?:music\.)?youtu(?:be\.com\/(watch\?v=|shorts\/|embed\/)|\.be\/|)(?P<id>[\w\-\_]{11})(?:&(amp;)?[\w\?=]*)?`)
	YouTubeCaches    = []VidCache{}
)

func GetYoutubeMedias(bot *telego.Bot, update telego.Update) {
	query := utils.GetArgs(update)

	bot.SendChatAction(
		&telego.SendChatActionParams{
			ChatID: telegoutil.ID(update.Message.Chat.ID),
			Action: "typing",
		},
	)
	if query == "" {
		bot.SendMessage(
			&telego.SendMessageParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				Text:      "<b>What is the video title, channel name or video URL to download?</b>",
				ParseMode: "HTML",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			},
		)
		return
	}
	if !YouTubeRegex_URL.MatchString(query) {
		BaseSearch := ytsearch.VideoSearch(query)
		search, err := BaseSearch.Next()
		if err != nil {
			log.Printf("There was an error when searching on YouTube: %v", err)
			bot.SendMessage(
				&telego.SendMessageParams{
					ChatID:    telegoutil.ID(update.Message.Chat.ID),
					Text:      "<b>There was an error when searching on YouTube.</b>",
					ParseMode: "HTML",
					ReplyParameters: &telego.ReplyParameters{
						MessageID: update.Message.MessageID,
					},
				},
			)
			return
		}
		SearchKey := RandYouTubeKey()
		if len(search.Videos) == 0 {
			bot.SendMessage(
				&telego.SendMessageParams{
					ChatID:    telegoutil.ID(update.Message.Chat.ID),
					Text:      fmt.Sprintf("<i>No results found for:</i> <b>%s</b>", query),
					ParseMode: "HTML",
					ReplyParameters: &telego.ReplyParameters{
						MessageID: update.Message.MessageID,
					},
				},
			)
			return
		}
		YouTubeCaches = append(
			YouTubeCaches,
			VidCache{
				UUID:            SearchKey,
				VidInformations: search.Videos,
			},
		)
		VideoCache := GetVideoByUUID(SearchKey)
		VideoPage := VideoCache.VidInformations[0]
		out := YouTubeMakeTextWithInfos(
			VideoPage.URL,
			VideoPage.Title,
			utils.TimeFormatter(float64(VideoPage.Duration)),
			utils.FormatInteger(VideoPage.ViewCount),
			VideoPage.PublishedTime,
			VideoPage.Channel.URL,
			VideoPage.Channel.Title,
		)
		ThumbnailURL := GetThumbURL(VideoPage.ID)
		bot.SendPhoto(
			&telego.SendPhotoParams{
				ChatID: telegoutil.ID(update.Message.Chat.ID),
				Photo: telego.InputFile{
					URL: ThumbnailURL,
				},
				Caption:   out,
				ParseMode: "HTML",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			},
		)
		return
	}
	bot.SendMessage(
		&telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      query,
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		},
	)
}
