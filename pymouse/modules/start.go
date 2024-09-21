package modules

import (
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

func Start(bot *telego.Bot, update telego.Update) {
	bot.SendMessage(
		&telego.SendMessageParams{
			ChatID: telegoutil.ID(update.Message.Chat.ID),
			Text:   "Hello. I'am using TeleGO.",
		},
	)
}
