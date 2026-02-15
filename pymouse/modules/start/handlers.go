package start

import (
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
)

func Start(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	bot.SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID: telegoutil.ID(update.Message.Chat.ID),
			Text:   "Hello. I'am using TeleGO.",
		},
	)
	return nil
}
