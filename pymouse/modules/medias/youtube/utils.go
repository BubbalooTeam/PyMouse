package youtube

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"pymouse/pymouse/config"
	"pymouse/pymouse/helpers/rapidhttp"
	"pymouse/pymouse/helpers/utils"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
	yt_dl "github.com/kkdai/youtube/v2"
	"github.com/mymmrac/telego"
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
		response, err := rapidhttp.Request(
			rapidhttp.RequestOptions{
				URL:    ThumbLink,
				Method: "GET",
			},
		)
		if err != nil {
			log.Printf("[youtube/GetThumbURL]: An error occurred.\n\n- VideoID: %s\nQuality: %s\nError: %v", VideoID, Quality, err)
			continue
		}
		if response.StatusCode() == 200 {
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
			log.Println("[youtube/GetYouTubeClient][Error]: Error in parse Socks5Proxy...")
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
		log.Printf("[youtube/GetDownloadButtons]: Error in fetching video informations...")
		return nil
	}

	IKB = [][]telego.InlineKeyboardButton{
		{
			{
				Text:         "🥇 BEST - 🎥 MP4",
				CallbackData: fmt.Sprintf("yt_dl|%s|mp4+140|%d|v", videoID, userID),
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
			AudioBitrate := int(math.Round(float64(format.Bitrate) / 1000))
			audioDict[AudioBitrate] = fmt.Sprintf("📀 %dKbps (%s)",
				AudioBitrate,
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
					CallbackData: fmt.Sprintf("yt_dl|%s|%s+140|%d|v", videoID, maxItag, userID),
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
				CallbackData: fmt.Sprintf("yt_dl|%s|mp3|%d|a", videoID, userID),
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
			CallbackData: fmt.Sprintf("yt_dl|%s|%d|%d|a", videoID, bitrate, userID),
		})
	}

	if len(audioBtns) > 0 {
		IKB = append(IKB, utils.SplitIntoRows(audioBtns, 2)...)
	}

	return IKB
}
