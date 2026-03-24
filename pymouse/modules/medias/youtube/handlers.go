package youtube

import (
	"fmt"
	"os"
	"pymouse/pymouse/helpers/i18n"
	"pymouse/pymouse/helpers/telegram"
	"pymouse/pymouse/helpers/utils"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/raitonoberu/ytsearch"
	"github.com/sirupsen/logrus"
)

func YouTubeDLHandler(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	l := i18n.Locale(update.Message.Chat)

	query := utils.GetArgs(update)
	callbackKey := utils.RandKey()
	if query == "" {
		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				Text:      l("youtube-dl.checkers.no-args"),
				ParseMode: "HTML",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			},
		)
		return nil
	}
	if !YouTubeRegex_URL.MatchString(query) {
		youtubeSearch, err := ytsearch.VideoSearch(query).Next()
		if err != nil {
			logrus.Errorf("There was an error when searching on YouTube: %v", err)
			return nil
		}
		youtubeVideos := youtubeSearch.Videos
		if len(youtubeVideos) == 0 {
			bot.SendMessage(
				ctx,
				&telego.SendMessageParams{
					ChatID:    telegoutil.ID(update.Message.Chat.ID),
					Text:      l("youtube-dl.checkers.video-not-found"),
					ParseMode: "HTML",
					ReplyParameters: &telego.ReplyParameters{
						MessageID: update.Message.MessageID,
					},
				},
			)
			return nil
		}
		searchKey := utils.RandKey()
		SetYouTubeSearchCache(searchKey, youtubeVideos)

		video := GetVideoInformations(youtubeVideos, 1)
		thumbnailURL := GetThumbURL(video.ID)
		outText := GetYouTubeScrollText(video, l)
		SetYouTubeVideoCallbackCache(
			callbackKey,
			video.Channel.Title,
			int64(video.Duration),
			video.ID,
			video.Title,
			thumbnailURL,
			map[string]map[string]int64{
				"140": {
					"140": MP3getSize(float64(video.Duration)),
				},
			},
		)

		// Get YouTube Buttons
		YouTubeKeyboard := telegram.KeyboardPaginate(
			len(youtubeVideos),
			1,
			fmt.Sprintf(
				"YouTubeScroll|%s|{number}|%d",
				searchKey,
				update.Message.From.ID,
			),
		)
		DownloadYouTubeButton := []telego.InlineKeyboardButton{
			{
				Text:         "📹️ Vídeo",
				CallbackData: fmt.Sprintf("yt|gen|%s|%v|%d", callbackKey, nil, update.Message.From.ID),
			},
			{
				Text:         "💿️ Áudio",
				CallbackData: fmt.Sprintf("yt|dl|%s|140|%d|a", callbackKey, update.Message.From.ID),
			},
		}
		YouTubeKeyboard = append(YouTubeKeyboard, DownloadYouTubeButton)
		// Send YouTube informations
		bot.SendPhoto(
			ctx,
			&telego.SendPhotoParams{
				ChatID: telegoutil.ID(update.Message.Chat.ID),
				Photo: telego.InputFile{
					URL: thumbnailURL,
				},
				Caption:   outText,
				ParseMode: "HTML",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
				ReplyMarkup: telegoutil.InlineKeyboardGrid(YouTubeKeyboard),
			},
		)
		return nil
	}

	videoID := utils.MatchByGroup(YouTubeRegex_URL, query, "id")
	videoInfo, err := ExtractVideoInfo(videoID)
	if err != nil {
		logrus.Errorf("failed to extract video informations: %v", err)
		return nil
	}
	thumbnailURL := GetThumbURL(videoInfo.ID)
	outText := fmt.Sprintf("📽️ <b>%s</b> - %s\n⏳️ <code>%s</code>", *videoInfo.Title, *videoInfo.Channel, utils.TimeFormatter(*videoInfo.Duration))
	SetYouTubeVideoCallbackCache(
		callbackKey,
		*videoInfo.Channel,
		int64(*videoInfo.Duration),
		videoInfo.ID,
		*videoInfo.Title,
		thumbnailURL,
		map[string]map[string]int64{
			"140": {
				"140": MP3getSize(*videoInfo.Duration),
			},
		},
	)
	thumbnailURL = GetThumbURL(videoInfo.ID)
	DownloadYouTubeButton := [][]telego.InlineKeyboardButton{
		{
			{
				Text:         "📹️ Video",
				CallbackData: fmt.Sprintf("yt|gen|%s|%v|%d", callbackKey, nil, update.Message.From.ID),
			},
			{
				Text:         "💿️ Audio",
				CallbackData: fmt.Sprintf("yt|dl|%s|140|%d|a", callbackKey, update.Message.From.ID),
			},
		},
	}
	bot.SendPhoto(
		ctx,
		&telego.SendPhotoParams{
			ChatID: telegoutil.ID(update.Message.Chat.ID),
			Photo: telego.InputFile{
				URL: thumbnailURL,
			},
			Caption:   outText,
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
			ReplyMarkup: telegoutil.InlineKeyboardGrid(DownloadYouTubeButton),
		},
	)
	return nil
}

func YouTubeScrollCallbackHandler(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	l := i18n.Locale(update.CallbackQuery.Message.GetChat())
	callbackData := strings.Split(update.CallbackQuery.Data, "|")

	if len(callbackData) < 2 {
		logrus.Error("failed to extract [searchKey] information from callbackData.")
		return nil
	}
	searchKey := callbackData[1]

	vidPageNumber, err := strconv.Atoi(callbackData[2])
	if err != nil {
		logrus.Error("failed to extract [vidPageNumber] information from callbackData.")
		return nil
	}

	userID, err := strconv.Atoi(callbackData[3])
	if err != nil {
		logrus.Error("failed to extract [userID] information from callbackData.")
		return nil
	}

	if update.CallbackQuery.From.ID != int64(userID) {
		bot.AnswerCallbackQuery(
			ctx,
			&telego.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            l("youtube-dl.checkers.not-for-you"),
				ShowAlert:       true,
				CacheTime:       3,
			},
		)
		return nil
	}

	youtubeSearch, found := GetYouTubeSearchCache(searchKey)
	if !found {
		bot.AnswerCallbackQuery(
			ctx,
			&telego.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            l("youtube-dl.checkers.video-cache-not-found"),
				ShowAlert:       true,
				CacheTime:       3,
			},
		)
		return nil
	}
	youtubeVideos := youtubeSearch.VidInformations

	callbackKey := utils.RandKey()
	video := GetVideoInformations(youtubeVideos, vidPageNumber)
	thumbnailURL := GetThumbURL(video.ID)
	outText := GetYouTubeScrollText(video, l)
	SetYouTubeVideoCallbackCache(
		callbackKey,
		video.Channel.Title,
		int64(video.Duration),
		video.ID,
		video.Title,
		thumbnailURL,
		map[string]map[string]int64{
			"140": {
				"140": MP3getSize(float64(video.Duration)),
			},
		},
	)

	// Get YouTube Buttons
	YouTubeKeyboard := telegram.KeyboardPaginate(
		len(youtubeVideos),
		vidPageNumber,
		fmt.Sprintf(
			"YouTubeScroll|%s|{number}|%d",
			searchKey,
			update.CallbackQuery.From.ID,
		),
	)
	DownloadYouTubeButton := []telego.InlineKeyboardButton{
		{
			Text:         "📹️ Vídeo",
			CallbackData: fmt.Sprintf("yt|gen|%s|%v|%d", callbackKey, nil, update.CallbackQuery.From.ID),
		},
		{
			Text:         "💿️ Áudio",
			CallbackData: fmt.Sprintf("yt|dl|%s|140|%d|a", callbackKey, update.CallbackQuery.From.ID),
		},
	}
	YouTubeKeyboard = append(YouTubeKeyboard, DownloadYouTubeButton)

	bot.EditMessageMedia(
		ctx,
		&telego.EditMessageMediaParams{
			ChatID:    telegoutil.ID(update.CallbackQuery.Message.GetChat().ID),
			MessageID: update.CallbackQuery.Message.GetMessageID(),
			Media: &telego.InputMediaPhoto{
				Type: "photo",
				Media: telego.InputFile{
					URL: thumbnailURL,
				},
				Caption:   outText,
				ParseMode: "HTML",
			},
			ReplyMarkup: telegoutil.InlineKeyboardGrid(YouTubeKeyboard),
		},
	)
	return nil
}

func YouTubeACallHandler(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	l := i18n.Locale(update.CallbackQuery.Message.GetChat())
	callbackData := strings.Split(update.CallbackQuery.Data, "|")

	if len(callbackData) < 2 {
		logrus.Error("failed to extract [searchKey] information from callbackData.")
		return nil
	}
	actionType := callbackData[1]
	callbackKey := callbackData[2]
	userID, err := strconv.Atoi(callbackData[4])
	if err != nil {
		logrus.Error("failed to extract [userID] information from callbackData.")
		return nil
	}

	if update.CallbackQuery.From.ID != int64(userID) {
		bot.AnswerCallbackQuery(
			ctx,
			&telego.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            l("youtube-dl.checkers.not-for-you"),
				ShowAlert:       true,
				CacheTime:       3,
			},
		)
		return nil
	}
	callbackInfo, ok := GetYouTubeVideoCallbackCache(callbackKey)
	if !ok {
		bot.AnswerCallbackQuery(
			ctx,
			&telego.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            l("youtube-dl.checkers.video-cache-not-found"),
				ShowAlert:       true,
				CacheTime:       3,
			},
		)
		return nil
	}
	switch actionType {
	case "gen":
		bot.EditMessageCaption(
			ctx,
			&telego.EditMessageCaptionParams{
				ChatID:    telegoutil.ID(update.CallbackQuery.Message.GetChat().ID),
				MessageID: update.CallbackQuery.Message.GetMessageID(),
				Caption:   l("youtube-dl.getting-formats"),
				ParseMode: "HTML",
			},
		)
		videoInfo, err := ExtractVideoInfo(callbackInfo.VideoID)
		if err != nil {
			bot.EditMessageCaption(
				ctx,
				&telego.EditMessageCaptionParams{
					ChatID:    telegoutil.ID(update.CallbackQuery.Message.GetChat().ID),
					MessageID: update.CallbackQuery.Message.GetMessageID(),
					Caption:   l("youtube-dl.checkers.extract-info-failed"),
					ParseMode: "HTML",
				},
			)
			logrus.Errorf("Failed to extract video informations: %v", err)
			return nil
		}
		ThumbnailURL := GetThumbURL(videoInfo.ID)
		keyboard, err := GetDownloadButton(
			callbackKey,
			videoInfo,
			int64(userID),
		)
		if err != nil {
			bot.EditMessageMedia(
				ctx,
				&telego.EditMessageMediaParams{
					ChatID:    telegoutil.ID(update.CallbackQuery.Message.GetChat().ID),
					MessageID: update.CallbackQuery.Message.GetMessageID(),
					Media: &telego.InputMediaPhoto{
						Type: "photo",
						Media: telego.InputFile{
							URL: ThumbnailURL,
						},
						Caption:   "<b>There was an error getting the video quality buttons!\nThis occurs due to several factors such as VPS blocked by YouTube, problematic Proxy or one that stopped working...</b>",
						ParseMode: "HTML",
					},
				},
			)
			return nil
		}
		out := fmt.Sprintf(
			"<b><a href=\"%s\">%s</a></b>\n\n",
			fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoInfo.ID),
			*videoInfo.Title,
		)
		bot.EditMessageMedia(
			ctx,
			&telego.EditMessageMediaParams{
				ChatID:    telegoutil.ID(update.CallbackQuery.Message.GetChat().ID),
				MessageID: update.CallbackQuery.Message.GetMessageID(),
				Media: &telego.InputMediaPhoto{
					Type: "photo",
					Media: telego.InputFile{
						URL: ThumbnailURL,
					},
					Caption:   out,
					ParseMode: "HTML",
				},
				ReplyMarkup: keyboard,
			},
		)
		return nil
	case "dl":
		formatID, err := strconv.Atoi(callbackData[3])
		if err != nil {
			logrus.Error("failed to extract [formatID] information from callbackData.")
			return nil
		}
		mediaType := callbackData[5]
		callbackInfo, found := GetYouTubeVideoCallbackCache(callbackKey)
		if !found {
			bot.AnswerCallbackQuery(
				ctx,
				&telego.AnswerCallbackQueryParams{
					CallbackQueryID: update.CallbackQuery.ID,
					Text:            l("youtube-dl.checkers.video-cache-not-found"),
					ShowAlert:       true,
					CacheTime:       3,
				},
			)
			return nil
		}

		bot.EditMessageCaption(
			ctx,
			&telego.EditMessageCaptionParams{
				ChatID:    telegoutil.ID(update.CallbackQuery.Message.GetChat().ID),
				MessageID: update.CallbackQuery.Message.GetMessageID(),
				Caption:   l("youtube-dl.downloading"),
				ParseMode: "HTML",
			},
		)
		filename, err := DownloadByFormatID(callbackInfo.VideoID, formatID)
		if err != nil {
			bot.EditMessageCaption(
				ctx,
				&telego.EditMessageCaptionParams{
					ChatID:    telegoutil.ID(update.CallbackQuery.Message.GetChat().ID),
					MessageID: update.CallbackQuery.Message.GetMessageID(),
					Caption:   l("youtube-dl.checkers.download-failed"),
					ParseMode: "HTML",
				},
			)
			logrus.Errorf("Failed to download video/song: %v", err)
			return nil
		}

		bot.EditMessageCaption(
			ctx,
			&telego.EditMessageCaptionParams{
				ChatID:    telegoutil.ID(update.CallbackQuery.Message.GetChat().ID),
				MessageID: update.CallbackQuery.Message.GetMessageID(),
				Caption:   l("youtube-dl.uploading"),
				ParseMode: "HTML",
			},
		)

		file, err := os.Open(filename)
		if err != nil {
			logrus.Error(err)
			return nil
		}
		defer func() {
			file.Close()
			os.Remove(filename)
		}()

		keyboard := [][]telego.InlineKeyboardButton{
			{
				{
					Text: fmt.Sprintf(l("buttons.open-in"), "YouTube"),
					URL:  fmt.Sprintf("https://youtube.com/watch?v=%s", callbackInfo.VideoID),
				},
			},
		}
		outText := fmt.Sprintf("<b><i>%s</i></b>\n\n", callbackInfo.VideoTitle)
		outText += fmt.Sprintf(
			l("youtube-dl.formatter.creator-of-content"),
			callbackInfo.VideoChannel,
		)
		outText += fmt.Sprintf(
			l("youtube-dl.formatter.duration-time"),
			utils.TimeFormatter(float64(callbackInfo.VideoDuration)),
		)

		if strings.Contains(mediaType, "a") {
			bot.EditMessageMedia(
				ctx,
				&telego.EditMessageMediaParams{
					ChatID:    telegoutil.ID(update.CallbackQuery.Message.GetChat().ID),
					MessageID: update.CallbackQuery.Message.GetMessageID(),
					Media: &telego.InputMediaAudio{
						Type: "audio",
						Media: telego.InputFile{
							File: file,
						},
						Caption:   outText,
						Duration:  int(callbackInfo.VideoDuration),
						Performer: callbackInfo.VideoChannel,
						Thumbnail: &telego.InputFile{
							URL: callbackInfo.VideoThumbnail,
						},
					},
					ReplyMarkup: &telego.InlineKeyboardMarkup{
						InlineKeyboard: keyboard,
					},
				},
			)
		} else {
			bot.EditMessageMedia(
				ctx,
				&telego.EditMessageMediaParams{
					ChatID:    telegoutil.ID(update.CallbackQuery.Message.GetChat().ID),
					MessageID: update.CallbackQuery.Message.GetMessageID(),
					Media: &telego.InputMediaVideo{
						Type: "video",
						Media: telego.InputFile{
							File: file,
						},
						Caption:  outText,
						Duration: int(callbackInfo.VideoDuration),
						Thumbnail: &telego.InputFile{
							URL: callbackInfo.VideoThumbnail,
						},
					},
					ReplyMarkup: &telego.InlineKeyboardMarkup{
						InlineKeyboard: keyboard,
					},
				},
			)
		}
	}
	return nil
}
