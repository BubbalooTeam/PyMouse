package youtube

import (
	"fmt"
	"log"
	"pymouse/pymouse/helpers/rapidhttp"

	"github.com/google/uuid"
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
