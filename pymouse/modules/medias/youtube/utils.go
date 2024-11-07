package youtube

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"pymouse/pymouse/config"
	"pymouse/pymouse/helpers/rapidhttp"
	"sort"
	"strconv"

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

	out := fmt.Sprintf("<b>❯ Published:</b> %s\n", publishedTime)
	out += fmt.Sprintf("<b>❯ Duration:</b> %s\n", Vidduration)
	out += fmt.Sprintf("<b>❯ Views:</b> %s\n", Vidviews)

	CreatorOfContent := fmt.Sprintf("<a href=\"%s\">%s</a>", CreatorofContentUrl, CreatorofContentName)
	out += fmt.Sprintf("<b>❯ Creator:</b> %s\n", CreatorOfContent)

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

// HumanBytes converte bytes para formato legível
func HumanBytes(size int64) string {
	if size == 0 {
		return ""
	}

	power := float64(1024)
	units := []string{" ", "Ki", "Mi", "Gi", "Ti"}
	i := 0
	s := float64(size)

	for s >= power && i < len(units)-1 {
		s /= power
		i++
	}

	return fmt.Sprintf("%.2f %sB", s, units[i])
}

// SplitIntoRows divide uma slice em múltiplas slices de tamanho fixo
func SplitIntoRows(items []telego.InlineKeyboardButton, width int) [][]telego.InlineKeyboardButton {
	var rows [][]telego.InlineKeyboardButton
	for i := 0; i < len(items); i += width {
		end := i + width
		if end > len(items) {
			end = len(items)
		}
		rows = append(rows, items[i:end])
	}
	return rows
}

func GetDownloadButtons(videoID string, userID int64) ([][]telego.InlineKeyboardButton, error) {
	client := GetYouTubeClient()
	video, err := client.GetVideo(videoID)
	if err != nil {
		return nil, fmt.Errorf("error fetching video info: %w", err)
	}

	// Primeiro botão: Melhor qualidade de vídeo
	buttons := [][]telego.InlineKeyboardButton{
		{
			{
				Text:         "🥇 BEST - 🎥 MP4",
				CallbackData: fmt.Sprintf("yt_dl|%s|mp4+140|%d|v", videoID, userID),
			},
		},
		{
			{
				Text:         "🥇 BEST - 📀 320Kbps - MP3",
				CallbackData: fmt.Sprintf("yt_dl|%s|mp3|%d|a", videoID, userID),
			},
		},
	}

	qualList := []string{"1440p", "1080p", "720p", "480p", "360p", "240p", "144p"}
	qualDict := make(map[string]map[string]int64)
	audioDict := make(map[int]string)

	for _, format := range video.Formats {
		if format.MimeType == "video/mp4" {
			for _, qual := range qualList {
				if format.QualityLabel == qual {
					if qualDict[qual] == nil {
						qualDict[qual] = make(map[string]int64)
					}
					qualDict[qual][strconv.Itoa(format.ItagNo)] = format.ContentLength
				}
			}
		}

		if format.MimeType == "audio/mp4" && format.AudioChannels > 0 {
			bitrate := int(math.Round(float64(format.Bitrate) / 1000))
			audioDict[bitrate] = fmt.Sprintf("📀 %dKbps (%s)",
				bitrate,
				HumanBytes(format.ContentLength),
			)
		}
	}

	var videoButtons []telego.InlineKeyboardButton
	for _, qual := range qualList {
		if formats, ok := qualDict[qual]; ok && len(formats) > 0 {
			var maxItag string
			var maxSize int64
			for itag, size := range formats {
				if size > maxSize {
					maxSize = size
					maxItag = itag
				}
			}

			videoButtons = append(videoButtons, telego.InlineKeyboardButton{
				Text:         fmt.Sprintf("🎥 %s (%s)", qual, HumanBytes(maxSize)),
				CallbackData: fmt.Sprintf("yt_dl|%s|%s+140|%d|v", videoID, maxItag, userID),
			})
		}
	}
	if len(videoButtons) > 0 {
		buttons = append(buttons, SplitIntoRows(videoButtons, 2)...)
	}

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
		buttons = append(buttons, SplitIntoRows(audioBtns, 2)...)
	}

	return buttons, nil
}
