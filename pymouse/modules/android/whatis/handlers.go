package whatis

import (
	"fmt"
	"pymouse/pymouse/helpers/i18n"
	"pymouse/pymouse/helpers/utils"
	"strings"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
)

func WhatIs(ctx *telegohandler.Context, update telego.Update) error {
	var modelText string
	bot := ctx.Bot()
	l := i18n.Locale(update.Message.Chat)

	query := utils.GetArgs(update)
	if query == "" {
		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				Text:      l("android.whatis.checkers.device-not-provided"),
				ParseMode: "HTML",
			},
		)
		return nil
	}

	result, found := getDevice(query)
	if !found {
		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				Text:      l("android.whatis.checkers.device-not-found"),
				ParseMode: "HTML",
			},
		)
		return nil
	}

	model := result.Device
	if model != "" {
		modelText += fmt.Sprintf(" (%s)", model)
	}
	bot.SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID: telegoutil.ID(update.Message.Chat.ID),
			Text: fmt.Sprintf(
				l("android.whatis.device-info"),
				strings.ToLower(query),
				result.Brand,
				result.Name,
				modelText,
			),
			ParseMode: "HTML",
		},
	)
	return nil
}
