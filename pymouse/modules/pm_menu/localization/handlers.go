package localization

import (
	"fmt"
	"pymouse/pymouse/database/utilitiesdb"
	"pymouse/pymouse/helpers/admin"
	"pymouse/pymouse/helpers/i18n"
	"strings"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/sirupsen/logrus"
)

func ChangeLanguageMessage(ctx *telegohandler.Context, update telego.Update) error {
	if update.Message == nil {
		return nil
	}

	bot := ctx.Bot()
	msg := update.Message
	chat := msg.Chat
	l := i18n.Locale(chat)

	if !admin.CheckAdmin(ctx, update, bot, admin.PermChangeInfo, l) {
		return nil
	}

	changeMenuBack := "StartBack"
	changeLangBack := "LangMenu"

	args := GetChangeLangTextAndButtons(
		chat,
		l,
		fmt.Sprintf("ChangeLang|%s|%s", changeMenuBack, changeLangBack),
		changeMenuBack,
	)

	bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:      telegoutil.ID(chat.ID),
		Text:        args.Text,
		ReplyMarkup: args.Buttons,
		ParseMode:   "HTML",
		ReplyParameters: &telego.ReplyParameters{
			MessageID: msg.MessageID,
		},
	})

	return nil
}

func ChangeLanguageCallback(ctx *telegohandler.Context, update telego.Update) error {
	if update.CallbackQuery == nil {
		return nil
	}

	bot := ctx.Bot()
	cb := update.CallbackQuery
	chat := cb.Message.GetChat()
	l := i18n.Locale(chat)

	if !admin.CheckAdmin(ctx, update, bot, admin.PermChangeInfo, l) {
		return nil
	}

	parts := strings.Split(cb.Data, "|")

	changeMenuBack := parts[1]
	changeLangBack := parts[2]

	args := GetChangeLangTextAndButtons(
		chat,
		l,
		fmt.Sprintf("ChangeLang|%s|%s", changeMenuBack, changeLangBack),
		changeMenuBack,
	)

	_, err := bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      telegoutil.ID(chat.ID),
		MessageID:   cb.Message.GetMessageID(),
		Text:        args.Text,
		ReplyMarkup: args.Buttons,
		ParseMode:   "HTML",
	})
	if err != nil {
		logrus.Error(err)
		return nil
	}

	bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: cb.ID,
	})

	return nil
}

func SelectLanguageCallback(ctx *telegohandler.Context, update telego.Update) error {
	if update.CallbackQuery == nil {
		return nil
	}

	bot := ctx.Bot()
	cb := update.CallbackQuery
	chat := cb.Message.GetChat()
	l := i18n.Locale(chat)

	if !admin.CheckAdmin(ctx, update, bot, admin.PermChangeInfo, l) {
		return nil
	}

	parts := strings.Split(cb.Data, "|")

	changeMenuBack := parts[1]
	changeLangBack := parts[2]

	args := GetSwitchLangTextAndButtons(
		update,
		l,
		changeMenuBack,
		fmt.Sprintf("%s|%s|LangMenu", changeLangBack, changeMenuBack),
	)

	bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      telegoutil.ID(chat.ID),
		MessageID:   cb.Message.GetMessageID(),
		Text:        args.Text,
		ReplyMarkup: args.Buttons,
		ParseMode:   "HTML",
	})

	bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: cb.ID,
	})

	return nil
}

func SwitchLanguageCallback(ctx *telegohandler.Context, update telego.Update) error {
	if update.CallbackQuery == nil {
		return nil
	}

	bot := ctx.Bot()
	cb := update.CallbackQuery
	chat := cb.Message.GetChat()
	l := i18n.Locale(chat)

	if !admin.CheckAdmin(ctx, update, bot, admin.PermChangeInfo, l) {
		return nil
	}

	parts := strings.Split(cb.Data, "|")

	language := parts[1]
	changeMenuBack := parts[2]

	switchLanguage := utilitiesdb.SetChatLanguage(chat, language)
	if !switchLanguage {
		bot.AnswerCallbackQuery(
			ctx,
			&telego.AnswerCallbackQueryParams{
				CallbackQueryID: cb.ID,
				Text:            fmt.Sprintf(l("language.already-set"), l("language.name")),
				ShowAlert:       true,
			},
		)
		return nil
	}

	// Reload locale after changing language
	l = i18n.Locale(chat)

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{
					Text:         l("buttons.back"),
					CallbackData: changeMenuBack,
				},
			},
		},
	}

	msgText := fmt.Sprintf(
		l("language.switched-lang"),
		l("language.flag")+l("language.name"),
	)

	bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      telegoutil.ID(chat.ID),
		MessageID:   cb.Message.GetMessageID(),
		Text:        msgText,
		ReplyMarkup: keyboard,
		ParseMode:   "HTML",
	})

	bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: cb.ID,
	})

	return nil
}
