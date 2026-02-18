package help

import (
	"fmt"
	"pymouse/pymouse/helpers/i18n"
	"pymouse/pymouse/middlewares"
	"strings"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/sirupsen/logrus"
)

func HelpMenu(ctx *telegohandler.Context, update telego.Update) error {
	if update.Message == nil {
		return nil
	}

	bot := ctx.Bot()

	botUser, err := bot.GetMe(ctx)
	if err != nil {
		logrus.Error("Failed to retrieve basic bot information.")
		return nil
	}

	l := i18n.Locale(update.Message.Chat)
	chat := update.Message.Chat

	if !strings.Contains(chat.Type, telego.ChatTypePrivate) {

		keyboard := &telego.InlineKeyboardMarkup{
			InlineKeyboard: [][]telego.InlineKeyboardButton{
				{
					{
						Text: l("buttons.go-to-pm"),
						URL:  fmt.Sprintf("https://t.me/%s?start=help", botUser.Username),
					},
				},
			},
		}

		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID:      telegoutil.ID(chat.ID),
				Text:        fmt.Sprintf(l("pm-menu.start-group"), botUser.FirstName),
				ParseMode:   "HTML",
				ReplyMarkup: keyboard,
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			},
		)

		return nil
	}

	// typing
	bot.SendChatAction(ctx, &telego.SendChatActionParams{
		ChatID: telegoutil.ID(chat.ID),
		Action: telego.ChatActionTyping,
	})

	helpText := fmt.Sprintf(l("help.help-intro"), botUser.FirstName)

	keyboard := GenerateHelpKeyboard(nil, l)

	bot.SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID:      telegoutil.ID(chat.ID),
			Text:        helpText,
			ParseMode:   "HTML",
			ReplyMarkup: keyboard,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		},
	)

	return nil
}

func HelpMenuCallback(ctx *telegohandler.Context, update telego.Update) error {
	if update.CallbackQuery == nil {
		return nil
	}

	cb := update.CallbackQuery.Message.(*telego.Message)
	bot := ctx.Bot()

	l := i18n.Locale(cb.Chat)

	botUser, err := bot.GetMe(ctx)
	if err != nil {
		logrus.Error("Failed to retrieve basic bot information.")
		return nil
	}

	helpText := fmt.Sprintf(l("help.help-intro"), botUser.FirstName)

	keyboard := GenerateHelpKeyboard(nil, l)

	rows := keyboard.InlineKeyboard
	rows = append(rows, []telego.InlineKeyboardButton{
		{
			Text:         l("buttons.back"),
			CallbackData: "StartBack",
		},
	})
	keyboard.InlineKeyboard = rows

	bot.EditMessageText(
		ctx,
		&telego.EditMessageTextParams{
			ChatID:      telegoutil.ID(cb.Chat.ID),
			MessageID:   cb.MessageID,
			Text:        helpText,
			ParseMode:   "HTML",
			ReplyMarkup: keyboard,
		},
	)

	return nil
}

func HelpModule(ctx *telegohandler.Context, update telego.Update) error {

	if update.CallbackQuery == nil {
		return nil
	}

	cb := update.CallbackQuery

	if cb.Data == "" || !strings.HasPrefix(cb.Data, "help:") {
		return nil
	}

	bot := ctx.Bot()
	l := i18n.Locale(cb.Message.GetChat())

	callbackData := strings.TrimPrefix(cb.Data, "help:")
	modulePath := strings.Split(callbackData, ".")

	if len(modulePath) >= 2 && modulePath[0] == modulePath[1] {
		modulePath = modulePath[:1]
	}

	moduleInfo := FindModule(modulePath, middlewares.Help.GetHelpable())
	if moduleInfo == nil {

		bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
			CallbackQueryID: cb.ID,
			Text:            l("help.module-info-error"),
			ShowAlert:       true,
		})

		return nil
	}

	moduleNode := moduleInfo
	moduleName := moduleNode.Module

	titleKey := moduleNode.TitleI18n
	if titleKey == "" {
		if len(modulePath) == 1 {
			titleKey = fmt.Sprintf("help.titles.%s.%s", modulePath[0], modulePath[0])
		} else {
			titleKey = "help.titles." + strings.Join(modulePath, ".")
		}
	}

	moduleTitle := l(titleKey)
	if strings.TrimSpace(moduleTitle) == "" {
		moduleTitle = moduleName
	}

	// -------- TEXTO --------

	helpText := fmt.Sprintf(
		l("help.here-is-help"),
		moduleTitle,
	)

	description := l(moduleNode.DescriptionI18n)
	if strings.TrimSpace(description) == "" {
		description = fmt.Sprintf(
			l("help.module-description-not-found"),
			moduleName,
		)
	}

	helpText += description

	// -------- KEYBOARD --------

	var rows [][]telego.InlineKeyboardButton
	var currentRow []telego.InlineKeyboardButton

	for _, sub := range moduleNode.Plugins {

		subName := sub.Module
		subTitleKey := sub.TitleI18n

		if subName == "" || subTitleKey == "" {
			continue
		}

		subTitle := l(subTitleKey)
		if strings.TrimSpace(subTitle) == "" {
			continue
		}

		subSlug := middlewares.Slug(subName)

		callback := fmt.Sprintf(
			"help:%s",
			strings.Join(append(modulePath, subSlug), "."),
		)

		currentRow = append(currentRow, telego.InlineKeyboardButton{
			Text:         subTitle,
			CallbackData: callback,
		})

		if len(currentRow) == 3 {
			rows = append(rows, currentRow)
			currentRow = nil
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	// -------- BOTÃO VOLTAR --------

	var backCallback string

	if len(modulePath) <= 1 {
		backCallback = "HelpMenu"
	} else {

		parent := modulePath[:len(modulePath)-1]

		if len(parent) == 1 {
			backCallback = fmt.Sprintf("help:%s.%s", parent[0], parent[0])
		} else {
			backCallback = "help:" + strings.Join(parent, ".")
		}
	}

	rows = append(rows, []telego.InlineKeyboardButton{
		{
			Text:         l("buttons.back"),
			CallbackData: backCallback,
		},
	})

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}

	bot.EditMessageText(
		ctx,
		&telego.EditMessageTextParams{
			ChatID:      telegoutil.ID(cb.Message.GetChat().ID),
			MessageID:   cb.Message.GetMessageID(),
			Text:        helpText,
			ParseMode:   "HTML",
			ReplyMarkup: keyboard,
		},
	)

	return nil
}
