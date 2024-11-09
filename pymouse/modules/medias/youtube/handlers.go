package youtube

import (
	"fmt"
	"log"
	"pymouse/pymouse/helpers/telegram"
	"pymouse/pymouse/helpers/utils"
	"regexp"
	"strconv"
	"strings"

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

	if query == "" {
		bot.SendChatAction(
			&telego.SendChatActionParams{
				ChatID: telegoutil.ID(update.Message.Chat.ID),
				Action: "typing",
			},
		)
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

			bot.SendChatAction(
				&telego.SendChatActionParams{
					ChatID: telegoutil.ID(update.Message.Chat.ID),
					Action: "typing",
				},
			)
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
			bot.SendChatAction(
				&telego.SendChatActionParams{
					ChatID: telegoutil.ID(update.Message.Chat.ID),
					Action: "typing",
				},
			)
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

		// Get YouTube Buttons
		YouTubeKeyboard := telegram.KeyboardPaginate(
			len(VideoCache.VidInformations),
			1,
			fmt.Sprintf(
				"YouTubeScroll|%s|{number}|%d",
				SearchKey,
				update.Message.From.ID,
			),
		)

		// Send YouTube informations
		bot.SendChatAction(
			&telego.SendChatActionParams{
				ChatID: telegoutil.ID(update.Message.Chat.ID),
				Action: "upload_photo",
			},
		)
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
				ReplyMarkup: YouTubeKeyboard,
			},
		)
		return
	}
	YouTubeClient := GetYouTubeClient()

	VideoID := utils.MatchByGroup(YouTubeRegex_URL, query, "id")
	YouTubeVideo, _ := YouTubeClient.GetVideo(VideoID)

	VideoQualKeyboard := GetDownloadButtons(YouTubeVideo.ID, update.Message.From.ID)
	ThumbnailURL := GetThumbURL(YouTubeVideo.ID)

	out := fmt.Sprintf(
		"<b><a href=\"%s\">%s</a></b>\n\n",
		fmt.Sprintf("https://www.youtube.com/watch?v=%s", YouTubeVideo.ID),
		YouTubeVideo.Title,
	)
	bot.SendChatAction(
		&telego.SendChatActionParams{
			ChatID: telegoutil.ID(update.Message.Chat.ID),
			Action: "upload_photo",
		},
	)
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
			ReplyMarkup: &telego.InlineKeyboardMarkup{InlineKeyboard: VideoQualKeyboard},
		},
	)
}

func YouTubeScrollCallback(bot *telego.Bot, update telego.Update) {
	message := update.CallbackQuery.Message.(*telego.Message)
	CallbackData := strings.Split(update.CallbackQuery.Data, "|")
	// Get the utility information
	SearchKey := CallbackData[1]
	VidPageNumber, err := strconv.Atoi(CallbackData[2])
	if err != nil {
		log.Printf("[youtube/YouTubeScrollCallback][Error]: Error in get the utility information (VidPageNumber): %v", err)
		return
	}
	UserID, err := strconv.Atoi(CallbackData[3])
	if err != nil {
		log.Printf("[youtube/YouTubeScrollCallback][Error]: Error in get the utility information (UserID): %v", err)
		return
	}

	if update.CallbackQuery.From.ID != int64(UserID) {
		bot.AnswerCallbackQuery(
			&telego.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "This YouTube Downloader/Scroll button is not directed at you!",
				ShowAlert:       true,
				CacheTime:       3,
			},
		)
		return
	}
	VideoCache := GetVideoByUUID(SearchKey)
	if VideoCache == nil {
		bot.AnswerCallbackQuery(
			&telego.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "This Search is not avalaible in my variables, please, perform a new search.",
				ShowAlert:       true,
				CacheTime:       3,
			},
		)
		return
	}

	VideoPage := VideoCache.VidInformations[VidPageNumber-1]

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

	// Get YouTube Buttons
	YouTubeKeyboard := telegram.KeyboardPaginate(
		len(VideoCache.VidInformations),
		VidPageNumber,
		fmt.Sprintf(
			"YouTubeScroll|%s|{number}|%d",
			SearchKey,
			update.CallbackQuery.From.ID,
		),
	)
	if VidPageNumber == 0 {
		if len(VideoCache.VidInformations) == 1 {
			bot.AnswerCallbackQuery(
				&telego.AnswerCallbackQueryParams{
					CallbackQueryID: update.CallbackQuery.ID,
					Text:            "That's the end of video list.",
					ShowAlert:       true,
					CacheTime:       3,
				},
			)
			return
		}
	}
	_, err = bot.EditMessageMedia(
		&telego.EditMessageMediaParams{
			ChatID:    telegoutil.ID(message.Chat.ID),
			MessageID: update.CallbackQuery.Message.GetMessageID(),
			Media: &telego.InputMediaPhoto{
				Type: "photo",
				Media: telego.InputFile{
					URL: ThumbnailURL,
				},
				Caption:   out,
				ParseMode: "HTML",
			},
			ReplyMarkup: YouTubeKeyboard,
		},
	)
	if err != nil {
		log.Println(err)
	}

}
