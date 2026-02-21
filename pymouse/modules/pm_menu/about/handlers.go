package about

import (
	"fmt"
	"pymouse/pymouse/helpers/i18n"
	"strings"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/sirupsen/logrus"
)

func AboutMessage(ctx *telegohandler.Context, update telego.Update) error {
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
		keyboard := &telego.InlineKeyboardMarkup{
			InlineKeyboard: [][]telego.InlineKeyboardButton{
				{
					{
						Text: l("buttons.go-to-pm"),
						URL:  fmt.Sprintf("https://t.me/%s?start=about", botUser.Username),
					},
				},
			},
		}

		groupText := fmt.Sprintf(
			l("pm-menu.start-group"),
			botUser.FirstName,
		)

		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID:      telegoutil.ID(chat.ID),
				Text:        groupText,
				ParseMode:   "HTML",
				ReplyMarkup: keyboard,
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
					Text: l("buttons.source-code"),
					URL:  "https://github.com/BubbalooTeam/PyMouse",
				},
			},
			{
				{
					Text: l("buttons.donate"),
					URL:  "https://livepix.gg/bubbalooteam",
				},
				{
					Text: l("buttons.channel-news"),
					URL:  "https://t.me/PyMouseNews",
				},
			},
		},
	}

	aboutText := fmt.Sprintf(
		l("pm-menu.about-text"),
		botUser.FirstName,
	)

	bot.SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID:    telegoutil.ID(chat.ID),
			Text:      aboutText,
			ParseMode: "HTML",
			LinkPreviewOptions: &telego.LinkPreviewOptions{
				IsDisabled: true,
			},
			ReplyMarkup: keyboard,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: msg.MessageID,
			},
		},
	)

	return nil
}

func AboutCallback(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	cb := update.CallbackQuery
	chat := cb.Message.GetChat()
	l := i18n.Locale(chat)

	botUser, err := bot.GetMe(ctx)
	if err != nil {
		logrus.Error("Failed to retrieve bot information.")
		return nil
	}

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{
					Text: l("buttons.source-code"),
					URL:  "https://github.com/BubbalooTeam/PyMouse",
				},
			},
			{
				{
					Text: l("buttons.donate"),
					URL:  "https://livepix.gg/bubbalooteam",
				},
				{
					Text: l("buttons.channel-news"),
					URL:  "https://t.me/PyMouseNews",
				},
			},
			{
				{
					Text:         l("buttons.back"),
					CallbackData: "StartBack",
				},
			},
		},
	}

	aboutText := fmt.Sprintf(
		l("pm-menu.about-text"),
		botUser.FirstName,
	)

	bot.EditMessageText(
		ctx,
		&telego.EditMessageTextParams{
			ChatID:    telegoutil.ID(chat.ID),
			MessageID: cb.Message.GetMessageID(),
			Text:      aboutText,
			ParseMode: "HTML",
			LinkPreviewOptions: &telego.LinkPreviewOptions{
				IsDisabled: true,
			},
			ReplyMarkup: keyboard,
		},
	)

	bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: cb.ID,
	})

	return nil
}
