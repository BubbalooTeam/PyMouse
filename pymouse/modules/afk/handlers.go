package afk

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

func SetAFK(bot *telego.Bot, update telego.Update) {
	reason := extractReason(update.Message.Text)
	SetupAFK(
		update.Message.From.ID,
		time.Now(),
		reason,
	)
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

func AFKtest(bot *telego.Bot, update telego.Update) {
	checkAFK := GetAFK(update.Message.From.ID)
	if checkAFK != nil {
		if afkStr, ok := checkAFK["AFKbool"].(string); ok {
			if afk, err := strconv.ParseBool(afkStr); err != nil {
				fmt.Println("Error in convert AFKdata in bollean:", err)
			} else if afk {
				bot.SendMessage(
					&telego.SendMessageParams{
						ChatID: telegoutil.ID(update.Message.Chat.ID),
						Text:   fmt.Sprintf("%s is AFK!", update.Message.From.FirstName),
					},
				)
			} else {
				bot.SendMessage(
					&telego.SendMessageParams{
						ChatID: telegoutil.ID(update.Message.Chat.ID),
						Text:   fmt.Sprintf("%s is not AFK!", update.Message.From.FirstName),
					},
				)
			}
		} else {
			fmt.Println("boll of AfkData not returns in string format.")
		}
	} else {
		bot.SendMessage(
			&telego.SendMessageParams{
				ChatID: telegoutil.ID(update.Message.Chat.ID),
				Text:   fmt.Sprintf("%s is not AFK!", update.Message.From.FirstName),
			},
		)
	}
}
