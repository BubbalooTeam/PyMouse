package start

import (
	"fmt"
	"pymouse/pymouse/helpers/i18n"
	"strings"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/sirupsen/logrus"
)

func StartMessage(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	msg := update.Message
	chat := msg.Chat
	l := i18n.Locale(chat)

	botUser, err := bot.GetMe(ctx)
	if err != nil {
		logrus.Error("Failed to retrieve bot information.")
		return nil
	}

	bot.SendChatAction(ctx, &telego.SendChatActionParams{
		ChatID: telegoutil.ID(chat.ID),
		Action: telego.ChatActionTyping,
	})

	if !strings.Contains(chat.Type, telego.ChatTypePrivate) {

		startText := fmt.Sprintf(
			l("pm-menu.start-group"),
			botUser.FirstName,
		)

		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID:    telegoutil.ID(chat.ID),
				Text:      startText,
				ParseMode: "HTML",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: msg.MessageID,
				},
			},
		)

		return nil
	}

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{
					Text:         l("buttons.language"),
					CallbackData: "LangMenu|StartBack|LangMenu",
				},
				{
					Text:         l("buttons.help"),
					CallbackData: "HelpMenu",
				},
			},
			{
				{
					Text:         l("buttons.about"),
					CallbackData: "AboutMenu",
				},
				{
					Text:         l("buttons.privacy"),
					CallbackData: "PrivacyPolicy",
				},
			},
		},
	}

	startText := fmt.Sprintf(
		l("pm-menu.start-private"),
		msg.From.FirstName,
		botUser.FirstName,
	)

	bot.SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID:      telegoutil.ID(chat.ID),
			Text:        startText,
			ParseMode:   "HTML",
			ReplyMarkup: keyboard,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: msg.MessageID,
			},
		},
	)

	return nil
}

func StartBackCallback(ctx *telegohandler.Context, update telego.Update) error {
	cb := update.CallbackQuery
	bot := ctx.Bot()
	chat := cb.Message.GetChat()
	l := i18n.Locale(chat)

	botUser, err := bot.GetMe(ctx)
	if err != nil {
		logrus.Error("Failed to retrieve bot information.")
		return nil
	}

	if !strings.Contains(chat.Type, telego.ChatTypePrivate) {

		startText := fmt.Sprintf(
			l("pm-menu.start-group"),
			botUser.FirstName,
		)

		bot.EditMessageText(
			ctx,
			&telego.EditMessageTextParams{
				ChatID:    telegoutil.ID(chat.ID),
				MessageID: cb.Message.GetMessageID(),
				Text:      startText,
				ParseMode: "HTML",
			},
		)

		return nil
	}

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{
					Text:         l("buttons.language"),
					CallbackData: "LangMenu|StartBack|LangMenu",
				},
				{
					Text:         l("buttons.help"),
					CallbackData: "HelpMenu",
				},
			},
			{
				{
					Text:         l("buttons.about"),
					CallbackData: "AboutMenu",
				},
				{
					Text:         l("buttons.privacy"),
					CallbackData: "PrivacyPolicy",
				},
			},
		},
	}

	startText := fmt.Sprintf(
		l("pm-menu.start-private"),
		cb.From.FirstName,
		botUser.FirstName,
	)

	bot.EditMessageText(
		ctx,
		&telego.EditMessageTextParams{
			ChatID:      telegoutil.ID(chat.ID),
			MessageID:   cb.Message.GetMessageID(),
			Text:        startText,
			ParseMode:   "HTML",
			ReplyMarkup: keyboard,
		},
	)

	return nil
}
