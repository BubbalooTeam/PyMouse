package localization

import (
	"fmt"
	"pymouse/pymouse/database/repositories"
	"pymouse/pymouse/helpers/i18n"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/sirupsen/logrus"
)

type LocalizationArgs struct {
	Text    string
	Buttons *telego.InlineKeyboardMarkup
}

func getChangeLangButtons(
	l func(string) string,
	langCallback string,
	backCallback string,
) *telego.InlineKeyboardMarkup {

	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{
					Text:         l("buttons.change-language"),
					CallbackData: langCallback,
				},
			},
			{
				{
					Text:         l("buttons.back"),
					CallbackData: backCallback,
				},
			},
		},
	}
}

func GetChangeLangTextAndButtons(
	chat telego.Chat,
	l func(string) string,
	langCallback string,
	backCallback string,
) LocalizationArgs {

	chatLanguage := repositories.GetChatLanguage(chat)
	stats, _ := i18n.GetLocalizationStats(chatLanguage)

	msgText := fmt.Sprintf(
		l("language.language-info.chat-language"),
		l("language.flag"),
		l("language.name"),
	)

	msgText += l("language.language-info.language-info")
	msgText += fmt.Sprintf(
		l("language.language-info.language-total-strings"),
		stats.TotalStrings,
	)
	msgText += fmt.Sprintf(
		l("language.language-info.language-translated"),
		stats.TranslatedStrings,
	)
	msgText += fmt.Sprintf(
		l("language.language-info.language-untranslated"),
		stats.UntranslatedStrings,
	)

	msgText += "\n\n"

	if int(stats.PercentageTranslated) == 100 {
		msgText += l("language.language-info.native-language")
	} else {
		msgText += fmt.Sprintf(
			l("language.language-info.missing-translations"),
			int(stats.PercentageTranslated),
		)
	}

	keyboard := getChangeLangButtons(
		l,
		langCallback,
		backCallback,
	)

	return LocalizationArgs{
		Text:    msgText,
		Buttons: keyboard,
	}
}

func getSwitchLangButtons(
	update telego.Update,
	l func(string) string,
	changeMenuBack string,
	backCallback string,
) *telego.InlineKeyboardMarkup {
	var rows [][]telego.InlineKeyboardButton
	var currentRow []telego.InlineKeyboardButton

	chatLanguage := repositories.GetChatLanguage(update.CallbackQuery.Message.GetChat())
	for i, lang := range i18n.AvalaibleLanguages {
		langMap, ok := i18n.StringsCache[lang]
		if !ok {
			logrus.Error("Failed to retrieve the language translation map.")
			continue
		}

		flag := i18n.GetStringFromNestedMap(langMap, "language.flag")
		name := i18n.GetStringFromNestedMap(langMap, "language.name")
		if lang == chatLanguage {
			flag = "[✓] - " + flag
		}

		btn := telego.InlineKeyboardButton{
			Text: flag + " " + name,
			CallbackData: fmt.Sprintf(
				"SwitchLang|%s|%s",
				lang,
				changeMenuBack,
			),
		}

		currentRow = append(currentRow, btn)
		if (i+1)%3 == 0 {
			rows = append(rows, currentRow)
			currentRow = []telego.InlineKeyboardButton{}
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	rows = append(rows, []telego.InlineKeyboardButton{
		{
			Text:         l("buttons.back"),
			CallbackData: backCallback,
		},
	})

	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func GetSwitchLangTextAndButtons(
	update telego.Update,
	l func(string) string,
	changeMenuBack string,
	backCallback string,
) LocalizationArgs {

	msgText := l("language.switch-lang")

	keyboard := getSwitchLangButtons(
		update,
		l,
		changeMenuBack,
		backCallback,
	)

	return LocalizationArgs{
		Text:    msgText,
		Buttons: keyboard,
	}
}

func SendLanguageSwitched(
	ctx *telegohandler.Context,
	update telego.Update,
	changeMenuBack string,
) error {

	if update.CallbackQuery == nil {
		return nil
	}

	bot := ctx.Bot()
	cb := update.CallbackQuery
	chat := cb.Message.GetChat()
	l := i18n.Locale(chat)

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

	_, err := bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      telegoutil.ID(chat.ID),
		MessageID:   cb.Message.GetMessageID(),
		Text:        msgText,
		ReplyMarkup: keyboard,
		ParseMode:   "HTML",
	})

	return err
}
