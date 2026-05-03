package devices

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
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			},
		)
		return nil
	}

	results, found := getDevices(query)
	if !found {
		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				Text:      l("android.whatis.checkers.device-not-found"),
				ParseMode: "HTML",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			},
		)
		return nil
	}

	result := results[0]
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
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		},
	)
	return nil
}

func Variants(ctx *telegohandler.Context, update telego.Update) error {
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
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			},
		)
		return nil
	}
	if strings.HasPrefix(strings.ToLower(query), "sm-") {
		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				Text:      fmt.Sprintf(l("android.whatis.checkers.variants-not-supported"), strings.ToUpper(query)),
				ParseMode: "HTML",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			},
		)
		return nil
	}

	results, found := getDevices(query)
	if !found {
		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				Text:      l("android.whatis.checkers.device-not-found"),
				ParseMode: "HTML",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			},
		)
		return nil
	}

	var variantsText string
	for _, result := range results {
		variant := fmt.Sprintf("%s %s", result.Brand, result.Name)
		model := result.Device
		if model != "" {
			variant += fmt.Sprintf(" (%s)", model)
		}
		if variant != "" {
			variantsText += fmt.Sprintf("- %s\n", variant)
		}
	}

	bot.SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID: telegoutil.ID(update.Message.Chat.ID),
			Text: fmt.Sprintf(
				l("android.whatis.device-variants"),
				strings.ToLower(query),
				variantsText,
			),
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		},
	)
	return nil
}
