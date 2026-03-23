package youtube

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"pymouse/pymouse/config"
	"pymouse/pymouse/helpers/rapidhttp"
	"pymouse/pymouse/helpers/utils"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lrstanley/go-ytdlp"
	"github.com/mymmrac/telego"
	"github.com/patrickmn/go-cache"
	"github.com/raitonoberu/ytsearch"
	"github.com/sirupsen/logrus"
)

type VidCacheInformations struct {
	UUID            string
	VidInformations []*ytsearch.VideoItem
}

type VidExtractCacheInformations struct {
	UUID           string
	VideoChannel   string
	VideoDuration  int64
	VideoID        string
	VideoTitle     string
	VideoThumbnail string
	Formats        map[string]map[string]int64
}

var (
	vidCache         = cache.New(15*time.Minute, 30*time.Minute)
	vidCallbackCache = cache.New(15*time.Minute, 30*time.Minute)
	YouTubeRegex_URL = regexp.MustCompile(`(?m)http(?:s?):\/\/(?:www\.)?(?:music\.)?youtu(?:be\.com\/(watch\?v=|shorts\/|embed\/)|\.be\/|)(?P<id>[\w\-\_]{11})(?:&(amp;)?[\w\?=]*)?`)
)

func SetYouTubeSearchCache(uuid string, videos []*ytsearch.VideoItem) {
	vidCache.Set(uuid, VidCacheInformations{
		UUID:            uuid,
		VidInformations: videos,
	}, cache.DefaultExpiration)
}

func GetYouTubeSearchCache(uuid string) (*VidCacheInformations, bool) {
	if v, found := vidCache.Get(uuid); found {
		val, ok := v.(VidCacheInformations)
		if !ok {
			return nil, false
		}
		if uuid == val.UUID {
			return &val, true
		}
	}
	return nil, false
}

func SetYouTubeVideoCallbackCache(
	uuid string,
	videoChannel string,
	videoDuration int64,
	videoID string,
	videoTitle string,
	videoThumbnail string,
	formats map[string]map[string]int64,
) {
	vidCallbackCache.Set(uuid, VidExtractCacheInformations{
		UUID:           uuid,
		VideoChannel:   videoChannel,
		VideoDuration:  videoDuration,
		VideoID:        videoID,
		VideoTitle:     videoTitle,
		VideoThumbnail: videoThumbnail,
		Formats:        formats,
	}, cache.DefaultExpiration)
}

func GetYouTubeVideoCallbackCache(uuid string) (*VidExtractCacheInformations, bool) {
	if v, found := vidCallbackCache.Get(uuid); found {
		val, ok := v.(VidExtractCacheInformations)
		if !ok {
			return nil, false
		}
		if uuid == val.UUID {
			return &val, true
		}
	}
	return nil, false
}

func GetThumbURL(VideoID string) (ThumbURL string) {
	client := rapidhttp.GetHTTPClient()
	ThumbQuality := []string{
		"maxresdefault.jpg", // Best quality
		"hqdefault.jpg",
		"sddefault.jpg",
		"mqdefault.jpg",
		"default.jpg", // Worst quality
	}
	ThumbURL = "https://imgur.com/4LwPLai"
	for _, Quality := range ThumbQuality {
		ThumbLink := fmt.Sprintf("https://i.ytimg.com/vi/%s/%s", VideoID, Quality)

		// Getting Thumbnail Informations
		response, err := rapidhttp.Request(
			client,
			rapidhttp.HTTPStruct{
				Method: "GET",
				URL:    ThumbLink,
			},
		)
		if err != nil {
			logrus.Errorf("An error occurred.\n\n- VideoID: %s\nQuality: %s\nError: %v", VideoID, Quality, err)
			continue
		}
		if response.StatusCode == 200 {
			ThumbURL = ThumbLink
			break
		}
	}
	return ThumbURL
}

func GetVideoInformations(
	videoInformations []*ytsearch.VideoItem,
	page int,
) *ytsearch.VideoItem {
	pageIndex := page - 1
	if pageIndex < 0 || videoInformations == nil || pageIndex >= len(videoInformations) {
		logrus.Errorf("failed to extract information from YouTube scrolling video.")
		return nil
	}

	return videoInformations[pageIndex]

}

func GetYouTubeScrollText(
	video *ytsearch.VideoItem,
	l func(string) string,
) string {
	videoText := fmt.Sprintf("<a href=\"%s\">%s</a>", video.URL, video.Title)
	videoText += fmt.Sprintf(l("youtube-dl.formatter.published-time"), video.PublishedTime)
	videoText += fmt.Sprintf(l("youtube-dl.formatter.duration-time"), utils.TimeFormatter(float64(video.Duration)))
	videoText += fmt.Sprintf(l("youtube-dl.formatter.views-count"), utils.FormatInteger(video.ViewCount))
	videoText += fmt.Sprintf(l("youtube-dl.formatter.creator-of-content"), video.Channel.Title)
	return videoText
}

func getSize(f *ytdlp.ExtractedFormat, duration *float64) int64 {
	if f.FileSize != nil {
		return int64(*f.FileSize)
	}

	if f.FileSizeApprox != nil {
		return int64(*f.FileSizeApprox)
	}

	if f.TBR != nil && duration != nil {
		bits := (*f.TBR * *duration * 1000) / 8
		return int64(bits)
	}

	return 0
}

func MP3getSize(duration float64) int64 {
	audioBitrate := 245.9
	bits := (audioBitrate * duration * 1000) / 8
	return int64(bits)
}

func ExtractVideoInfo(videoID string) (*ytdlp.ExtractedInfo, error) {
	dl := ytdlp.New().
		FormatSort("res,ext:mp4").
		SkipDownload().
		PrintJSON().
		JsRuntimes("node")

	if _, err := os.Stat("youtubeCookies.txt"); err == nil {
		dl = dl.Cookies("youtubeCookies.txt")
	}

	run, err := dl.Run(
		context.TODO(),
		fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to process video informations: %v", err)
	}

	video, err := run.GetExtractedInfo()
	if err != nil || len(video) == 0 {
		return nil, fmt.Errorf("failed to extract video informations.")
	}

	return video[0], nil

}

func GetQualMap(videoInfo *ytdlp.ExtractedInfo, qualList []string) (map[string]map[string]int64, bool) {
	qualDict := map[string]map[string]int64{}

	for _, f := range videoInfo.Formats {
		ext := ""
		if f.Extension != nil {
			ext = *f.Extension
		}

		formatID := ""
		if f.FormatID != nil {
			formatID = *f.FormatID
		}

		size := getSize(f, videoInfo.Duration)

		if ext == "mp4" && f.Height != nil {
			res := fmt.Sprintf("%dp", int(*f.Height))

			for _, q := range qualList {
				if res == q {
					if qualDict[q] == nil {
						qualDict[q] = map[string]int64{}
					}

					qualDict[q][formatID] = size
				}
			}
		}
	}

	if len(qualDict) == 0 {
		return nil, false
	}

	return qualDict, true
}

func GetDownloadButton(callbackKey string, videoInfo *ytdlp.ExtractedInfo, userID int64) (*telego.InlineKeyboardMarkup, error) {
	qualList := []string{"1440p", "1080p", "720p", "480p", "360p", "240p", "144p"}
	callbackQualDict := map[string]map[string]int64{}
	videoID := ""
	if videoInfo.ID != "" {
		videoID = videoInfo.ID
	}

	qualDict, found := GetQualMap(videoInfo, qualList)
	if !found {
		return nil, fmt.Errorf("failed to extract quality list.")
	}

	var rows [][]telego.InlineKeyboardButton
	var currentRow []telego.InlineKeyboardButton

	for _, q := range qualList {
		frmtDict := qualDict[q]

		if len(frmtDict) == 0 {
			continue
		}

		var bestID string
		var bestSize int64
		var bestInt int

		for id, s := range frmtDict {
			idInt, _ := strconv.Atoi(id)

			if idInt > bestInt {
				bestInt = idInt
				bestID = id
				bestSize = s
			}
		}

		if bestID == "" && bestSize == 0 {
			continue
		}
		if callbackQualDict[bestID] == nil {
			callbackQualDict[bestID] = map[string]int64{}
		}

		callbackQualDict[bestID][bestID] = bestSize
		currentRow = append(currentRow, telego.InlineKeyboardButton{
			Text: fmt.Sprintf("📹 %s (%s)", q, utils.HumanBytes(bestSize)),
			CallbackData: fmt.Sprintf(
				"yt|dl|%s|%s|%d|v",
				callbackKey,
				bestID,
				userID,
			),
		})

		if len(currentRow) == 2 {
			rows = append(rows, currentRow)
			currentRow = []telego.InlineKeyboardButton{}
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}
	SetYouTubeVideoCallbackCache(
		callbackKey,
		*videoInfo.Channel,
		int64(*videoInfo.Duration),
		videoID,
		*videoInfo.Title,
		GetThumbURL(videoInfo.ID),
		callbackQualDict,
	)

	markup := telego.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}

	return &markup, nil
}

func DownloadByFormatID(videoID string, formatID int) (string, error) {
	dir := fmt.Sprintf("%s/medias/youtubedl", config.DownloadPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)

	var dl *ytdlp.Command
	var template string
	var ext string
	var suffix string

	if formatID == 140 {
		template = fmt.Sprintf("%s/%%(title)s_a.%%(ext)s", dir)
		ext = "mp3"
		suffix = "_a"

		dl = ytdlp.New().
			ExtractAudio().
			AudioFormat("mp3").
			AudioQuality("0").
			EmbedMetadata().
			Output(template)
	} else {
		template = fmt.Sprintf("%s/%%(title)s_v.%%(ext)s", dir)
		ext = "mp4"
		suffix = "_v"

		dl = ytdlp.New().
			Format(fmt.Sprintf("%d+bestaudio", formatID)).
			MergeOutputFormat("mp4").
			Output(template)
	}

	dl.JsRuntimes("node")
	dl.PrintJSON()

	if _, err := os.Stat("youtubeCookies.txt"); err == nil {
		dl.Cookies("youtubeCookies.txt")
	}

	run, err := dl.Run(context.TODO(), url)
	if err != nil {
		return "", err
	}

	video, err := run.GetExtractedInfo()
	if err != nil {
		return "", fmt.Errorf("failed to extract informations!")
	}

	videoInfo := video[0]

	var base string
	if videoInfo.AltFilename != nil {
		base = filepath.Base(*videoInfo.AltFilename)
		base = strings.TrimSuffix(base, filepath.Ext(base))
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var filename string

	for _, f := range files {
		name := f.Name()

		if strings.HasPrefix(name, base) &&
			strings.Contains(name, suffix) &&
			strings.HasSuffix(name, "."+ext) {

			filename = filepath.Join(dir, name)
			break
		}
	}

	if filename == "" {
		return "", fmt.Errorf("final file not found")
	}

	dirname := filepath.Dir(filename)
	file := filepath.Base(filename)

	fileExt := filepath.Ext(file)
	name := strings.TrimSuffix(file, fileExt)

	sanitized := utils.SanitizeFilename(name)

	if sanitized != name {
		newPath := filepath.Join(dirname, sanitized+fileExt)

		if _, err := os.Stat(newPath); err == nil {
			for i := 1; i <= 10; i++ {
				try := filepath.Join(dirname, fmt.Sprintf("%s_%d%s", sanitized, i, fileExt))
				if _, err := os.Stat(try); os.IsNotExist(err) {
					newPath = try
					break
				}
			}
		}

		if err := os.Rename(filename, newPath); err != nil {
			return "", err
		}

		filename = newPath
	}

	return filename, nil
}
