package youtube

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"pymouse/pymouse/config"
	"pymouse/pymouse/helpers/rapidhttp"
	"pymouse/pymouse/helpers/utils"
	"pymouse/pymouse/modules/medias"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	yt_dl "github.com/kkdai/youtube/v2"
	"github.com/mymmrac/telego"
	"github.com/sirupsen/logrus"
)

func RandYouTubeKey() (RYK string) {
	NewID := uuid.New()
	RYK = NewID.String()[:8]
	return RYK
}

func YouTubeMakeTextWithInfos(
	Vidurl string,
	Vidtitle string,
	Vidduration string,
	Vidviews string,
	publishedTime string,
	CreatorofContentUrl string,
	CreatorofContentName string,
) string {
	headInfo := fmt.Sprintf("<b><a href=\"%s\">%s</a></b>\n\n", Vidurl, Vidtitle)

	out := fmt.Sprintf("<b>❯ Published:</b> <code>%s</code>\n", publishedTime)
	out += fmt.Sprintf("<b>❯ Duration:</b> <code>%s</code>\n", Vidduration)
	out += fmt.Sprintf("<b>❯ Views:</b> <code>%s</code>\n", Vidviews)

	CreatorOfContent := fmt.Sprintf("<a href=\"%s\">%s</a>", CreatorofContentUrl, CreatorofContentName)
	out += fmt.Sprintf("<b>❯ Creator:</b> <i>%s</i>\n", CreatorOfContent)

	return headInfo + out
}

func GetVideoByUUID(UUID string) *VidCache {
	for _, CachedVideo := range YouTubeCaches {
		if CachedVideo.UUID == UUID {
			return &CachedVideo
		}
	}
	return nil
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

func GetYouTubeClient() yt_dl.Client {
	if config.Socks5Proxy != "" {
		ParsedProxyURL, err := url.Parse(config.Socks5Proxy)
		if err != nil {
			logrus.Errorf("Error in parse Socks5Proxy, please check the bot .env file...")
			return yt_dl.Client{}
		}
		HTTPClient := http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(ParsedProxyURL),
			},
		}
		return yt_dl.Client{HTTPClient: &HTTPClient}
	}
	return yt_dl.Client{}
}

func GetDownloadButtons(videoID string, userID int64) (IKB [][]telego.InlineKeyboardButton) {
	client := GetYouTubeClient()
	video, err := client.GetVideo(videoID)
	if err != nil {
		logrus.Errorf("Error in fetching video informations...")
		return nil
	}

	IKB = [][]telego.InlineKeyboardButton{
		{
			{
				Text:         "🥇 BEST - 🎥 MP4",
				CallbackData: fmt.Sprintf("yt|dl|%s|mp4|%d|v", videoID, userID),
			},
		},
	}

	VidQList := []string{"1080p60", "720p60", "1080p", "720p", "480p", "360p", "240p", "144p"}
	vidQDict := make(map[string]map[string]int64)

	for _, VidQuality := range VidQList {
		vidQDict[VidQuality] = make(map[string]int64)
	}

	audioDict := make(map[int]string)

	for _, format := range video.Formats {
		if strings.Contains(format.MimeType, "video/mp4") {
			if format.ContentLength == 0 {
				continue
			}

			itagStr := strconv.Itoa(format.ItagNo)
			VidQuality := format.QualityLabel

			if _, exists := vidQDict[VidQuality]; exists {
				vidQDict[VidQuality][itagStr] = format.ContentLength
			}
		} else if strings.Contains(format.MimeType, "audio/mp4") && format.AudioChannels > 0 {
			AudioBitrate := format.ItagNo
			audioDict[AudioBitrate] = fmt.Sprintf("📀 %dKbps (%s)",
				int(math.Round(float64(format.Bitrate)/1000.0)),
				utils.HumanBytes(format.ContentLength),
			)
		}
	}

	var videoButtons []telego.InlineKeyboardButton
	for _, VidQuality := range VidQList {
		formats := vidQDict[VidQuality]
		if len(formats) > 0 {
			var maxItag string
			var maxSize int64
			for itag, size := range formats {
				if size > maxSize {
					maxSize = size
					maxItag = itag
				}
			}

			if maxSize > 0 {
				videoButtons = append(videoButtons, telego.InlineKeyboardButton{
					Text:         fmt.Sprintf("🎥 %s (%s)", VidQuality, utils.HumanBytes(maxSize)),
					CallbackData: fmt.Sprintf("yt|dl|%s|%s|%d|v", videoID, maxItag, userID),
				})
			}
		}
	}

	if len(videoButtons) > 0 {
		IKB = append(IKB, utils.SplitIntoRows(videoButtons, 2)...)
	}

	IKB = append(
		IKB,
		[]telego.InlineKeyboardButton{
			{
				Text:         "🥇 BEST - 📀 320Kbps - MP3",
				CallbackData: fmt.Sprintf("yt|dl|%s|mp3|%d|a", videoID, userID),
			},
		},
	)
	var audioBtns []telego.InlineKeyboardButton
	var bitrates []int
	for bitrate := range audioDict {
		bitrates = append(bitrates, bitrate)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(bitrates)))

	for _, bitrate := range bitrates {
		audioBtns = append(audioBtns, telego.InlineKeyboardButton{
			Text:         audioDict[bitrate],
			CallbackData: fmt.Sprintf("yt|dl|%s|%d|%d|a", videoID, bitrate, userID),
		})
	}

	if len(audioBtns) > 0 {
		IKB = append(IKB, utils.SplitIntoRows(audioBtns, 2)...)
	}

	return IKB
}

func GetYouTubeFormat(Video *yt_dl.Video, Itag int) *yt_dl.Format {
	YouTubeFormat := Video.Formats.Itag(Itag)
	if len(YouTubeFormat) == 0 {
		logrus.Error("This YouTube Video Itag is Invalid!")
		return nil
	}
	return &YouTubeFormat[0]
}

func GetBestQuality(formats []yt_dl.Format, mediaType string) yt_dl.Format {
	var bestQuality yt_dl.Format
	var maxBitrate int

	isDesiredQuality := func(qualityLabel string) bool {
		supportedQualities := []string{"1080p60", "720p60", "1080p", "720p", "480p", "360p", "240p", "144p"}
		for _, supported := range supportedQualities {
			if strings.Contains(qualityLabel, supported) {
				return true
			}
		}
		return false
	}

	for _, format := range formats {
		switch mediaType {
		case "video":
			if format.Bitrate > maxBitrate && isDesiredQuality(format.QualityLabel) {
				maxBitrate = format.Bitrate
				bestQuality = format
			}
		case "audio":
			if format.AudioChannels > 0 && format.QualityLabel == "" && format.Bitrate > maxBitrate {
				maxBitrate = format.Bitrate
				bestQuality = format
			}
		}
	}
	return bestQuality
}

func copyStreamWithRetries(ytClient *yt_dl.Client, ytVideo *yt_dl.Video, ytFormat *yt_dl.Format, mediaFile *os.File) error {
	for attempt := 1; attempt <= 5; attempt++ {
		stream, _, err := ytClient.GetStream(ytVideo, ytFormat)
		if err != nil {
			logrus.Errorf("Failed to get stream from YouTube: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		_, err = io.Copy(mediaFile, stream)
		stream.Close()

		logrus.Error(err)
		if err == nil {
			return nil
		}

		mediaFile.Seek(0, 0)
		mediaFile.Truncate(0)
		time.Sleep(2 * time.Second)
	}

	os.Remove(mediaFile.Name())
	return fmt.Errorf("youtube — Failed to copy stream, after 5 attempts")
}

func downloadAndMergeAudio(ytClient *yt_dl.Client, ytVideo *yt_dl.Video, mediaFile *os.File) error {
	audioFormat := GetYouTubeFormat(ytVideo, 140)

	audioFile, err := os.CreateTemp("", "PyMouse_YouTube_*.m4a")
	if err != nil {
		return err
	}
	defer audioFile.Close()

	err = copyStreamWithRetries(ytClient, ytVideo, audioFormat, audioFile)
	if err != nil {
		return err
	}

	return medias.MergeAudioVideo(mediaFile, audioFile)
}

func DownloadYouTubeVideo(
	VideoID string,
	MediaType string,
	VideoSItag string,
) (*os.File, string, error) {
	var MediaFile *os.File
	var VideoFormat *yt_dl.Format

	YouTubeClient := GetYouTubeClient()
	YouTubeVideo, err := YouTubeClient.GetVideo(VideoID)
	if err != nil {
		logrus.Errorf("Failed to Get Video in YouTube, please check your Proxy or YouTube-Downloader.")
		return nil, "", err
	}

	ytFilename := fmt.Sprintf("PyMouse_YouTube_%s", RandYouTubeKey())

	formatType := "audio/mp4"
	if strings.Contains(MediaType, "video") {
		formatType = "video/mp4"
	}

	switch VideoSItag {
	case "mp3", "mp4":
		VideoQuality := GetBestQuality(YouTubeVideo.Formats.Type(formatType), MediaType)
		logrus.Info(VideoQuality.ItagNo)
		VideoFormat = GetYouTubeFormat(YouTubeVideo, VideoQuality.ItagNo)
		if VideoFormat == nil {
			return nil, "", fmt.Errorf("youtube stream with this format is not avalaible")
		}
	default:
		VideoItag, err := strconv.Atoi(VideoSItag)
		if err != nil {
			logrus.Errorf("Error in get Download information (VideoItag): %v", err)
			return nil, "", err
		}
		VideoFormat = GetYouTubeFormat(YouTubeVideo, VideoItag)
		if VideoFormat == nil {
			return nil, "", fmt.Errorf("youtube stream with this format is not avalaible")
		}
	}

	// Switch Download Method, According to MediaType
	var fileExtension string
	switch MediaType {
	case "audio":
		fileExtension = ".mp3"
	case "video":
		fileExtension = ".mp4"
	default:
		fileExtension = ".mp4"
	}

	MediaFile, err = os.Create(filepath.Join(os.TempDir(), ytFilename+fileExtension))
	if err != nil {
		logrus.Error("Failed to create a temporary file.")
		return nil, "", err
	}

	defer func() {
		if err != nil {
			MediaFile.Close()
			os.Remove(MediaFile.Name())
		}
	}()

	// Downloader of YouTube.
	logrus.Info("Trying to download YouTube stream...")
	err = copyStreamWithRetries(&YouTubeClient, YouTubeVideo, VideoFormat, MediaFile)
	if err != nil {
		logrus.Errorf("Failed to download a YouTube stream: %v", err)
		return nil, "", err
	}

	if strings.Contains(MediaType, "video") {
		err = downloadAndMergeAudio(&YouTubeClient, YouTubeVideo, MediaFile)
		if err != nil {
			logrus.Errorf("Failed to merge audio in video: %v", err)
			return nil, "", err
		}
	}
	MediaFile.Seek(0, 0)

	return MediaFile, YouTubeMakeTextWithInfos(
		fmt.Sprintf("https://www.youtube.com/watch?v=%s", YouTubeVideo.ID),
		YouTubeVideo.Title,
		utils.TimeFormatter(YouTubeVideo.Duration.Seconds()),
		utils.FormatInteger(YouTubeVideo.Views),
		YouTubeVideo.PublishDate.Format(time.RFC822),
		fmt.Sprintf("www.youtube.com/channel/%s", YouTubeVideo.ChannelID),
		YouTubeVideo.Author,
	), nil
}
