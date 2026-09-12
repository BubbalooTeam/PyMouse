package android

import (
	"pymouse/pymouse/client"
	"pymouse/pymouse/helpers/telegram"
	"pymouse/pymouse/middlewares"
	"pymouse/pymouse/modules/android/devices"
	"pymouse/pymouse/modules/android/lineageos"
	"pymouse/pymouse/modules/android/magisk"
	"pymouse/pymouse/modules/android/twrp"
	"regexp"

	"github.com/mymmrac/telego/telegohandler"
)

func LoadModule(bS *client.BotStruct) {
	middlewares.Help.RegisterHelp("Android", middlewares.Help.WithSubmodules("LineageOS", "Magisk", "TWRP", "Variants"))

	// WhatIs Handlers
	bS.Handler.Handle(devices.Variants, telegram.Command("variants", bS.Client.Username()))
	bS.Handler.Handle(devices.WhatIs, telegram.Command("whatis", bS.Client.Username()))
	bS.Handler.Handle(magisk.MagiskHandler, telegram.Command("magisk", bS.Client.Username()))
	bS.Handler.Handle(lineageos.LineageOSHandler,
		telegohandler.Or(
			telegram.Command("los", bS.Client.Username()),
			telegram.Command("lineageos", bS.Client.Username()),
		))
	bS.Handler.Handle(
		magisk.MagiskVariantCallbackHandler,
		telegohandler.CallbackDataMatches(regexp.MustCompile(`^magisk_variant\|.*$`)),
	)
	bS.Handler.Handle(
		lineageos.LineageOSPageCallbackHandler,
		telegohandler.CallbackDataMatches(regexp.MustCompile(`^lineageos_page\|.*$`)),
	)
	bS.Handler.Handle(twrp.TWRPHandler, telegram.Command("twrp", bS.Client.Username()))
	bS.Handler.Handle(
		twrp.TWRPPageCallbackHandler,
		telegohandler.CallbackDataMatches(regexp.MustCompile(`^twrp_page\|.*$`)),
	)
}
