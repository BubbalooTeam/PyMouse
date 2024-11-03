package youtube

import (
	"fmt"

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
	Vidviews int,
	publishedTime string,
	CreatorofContentUrl string,
	CreatorofContentName string,
) string {
	headInfo := fmt.Sprintf("<b><a href=\"%s\">%s</a></b>\n\n", Vidurl, Vidtitle)

	out := fmt.Sprintf("<b>Published:</b> %s\n", publishedTime)
	out += fmt.Sprintf("<b>Duration:</b> %s\n", Vidduration)
	out += fmt.Sprintf("<b>Views:</b> %d\n", Vidviews)

	CreatorOfContent := fmt.Sprintf("<a href=\"%s\">%s</a>", CreatorofContentUrl, CreatorofContentName)
	out += fmt.Sprintf("<b>Creator:</b> %s\n", CreatorOfContent)

	return headInfo + out
}
