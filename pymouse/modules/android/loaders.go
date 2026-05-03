package android

import (
	"pymouse/pymouse/client"
	"pymouse/pymouse/helpers/telegram"
	"pymouse/pymouse/middlewares"
	"pymouse/pymouse/modules/android/devices"
)

func LoadModule(bS *client.BotStruct) {
	middlewares.Help.RegisterHelp("Android")

	// WhatIs Handlers
	bS.Handler.Handle(devices.Variants, telegram.Command("variants", bS.Client.Username()))
	bS.Handler.Handle(devices.WhatIs, telegram.Command("whatis", bS.Client.Username()))
}
