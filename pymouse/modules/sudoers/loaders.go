package sudoers

import (
	"pymouse/pymouse/client"
	"pymouse/pymouse/helpers/telegram"
	"pymouse/pymouse/modules/sudoers/speedtst"

	"github.com/mymmrac/telego/telegohandler"
)

func LoadModules(bS *client.BotStruct) {
	bS.Handler.Handle(
		speedtst.SpeedTestMessage,
		telegohandler.And(
			telegram.Command("speedtest", bS.Client.Username()),
			telegram.IsDeveloper(),
		),
	)
}
