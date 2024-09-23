package afk

import (
	"regexp"
	"strconv"
	"time"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

func SetAFK(bot *telego.Bot, update telego.Update) {
	reason := extractReason(update.Message.Text)
	SetupAFK(strconv.FormatInt(update.Message.From.ID, 10), time.Now(), reason)
	bot.SendMessage(&telego.SendMessageParams{
		ChatID: telegoutil.ID(update.Message.Chat.ID),
		Text:   "is now AFK!",
	})
}

func extractReason(text string) string {
	matches := regexp.MustCompile(`^(?:brb|\/afk)\s(.+)$`).FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
