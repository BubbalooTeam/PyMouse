package android

import (
	"pymouse/pymouse/client"
	"pymouse/pymouse/helpers/telegram"
	"pymouse/pymouse/middlewares"
	"pymouse/pymouse/modules/android/whatis"
)

func LoadModule(bS *client.BotStruct) {
	middlewares.Help.RegisterHelp("Android")

	// WhatIs Handlers
	bS.Handler.Handle(whatis.WhatIs, telegram.Command("whatis", bS.Client.Username()))
}
