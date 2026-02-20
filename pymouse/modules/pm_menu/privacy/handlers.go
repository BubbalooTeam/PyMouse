package privacy

import (
	"fmt"
	"pymouse/pymouse/helpers/i18n"
	"strings"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
)

func PrivacyPolicyMessage(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	msg := update.Message
	chat := msg.Chat
	l := i18n.Locale(update.Message.Chat)

	var keyboard [][]telego.InlineKeyboardButton
	botUser, err := bot.GetMe(ctx)
	if err != nil {
		return fmt.Errorf("Failed to retrieve bot information: %v", err)
	}

	if !strings.Contains(chat.Type, telego.ChatTypePrivate) {

		keyboard := &telego.InlineKeyboardMarkup{
			InlineKeyboard: [][]telego.InlineKeyboardButton{
				{
					{
						Text: l("buttons.go-to-pm"),
						URL:  fmt.Sprintf("https://t.me/%s?start=privacy", botUser.Username),
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

	keyboard = append(keyboard, []telego.InlineKeyboardButton{
		{
			Text: l("buttons.privacy-policy"),
			URL:  "https://telegram.org/privacy-tpa",
		},
	})
	bot.SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID:    telegoutil.ID(chat.ID),
			Text:      fmt.Sprintf(l("pm-menu.privacy-text"), botUser.FirstName),
			ParseMode: "HTML",
			ReplyMarkup: &telego.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		},
	)
	return nil
}

func PrivacyPolicyCallback(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	cb := update.CallbackQuery
	chat := cb.Message.GetChat()
	l := i18n.Locale(chat)

	botUser, err := bot.GetMe(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve bot information: %v", err)
	}

	keyboard := [][]telego.InlineKeyboardButton{
		{
			{
				Text: l("buttons.privacy-policy"),
				URL:  "https://telegram.org/privacy-tpa",
			},
		},
		{
			{
				Text:         l("buttons.back"),
				CallbackData: "StartBack",
			},
		},
	}

	_, err = bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:    telegoutil.ID(chat.ID),
		MessageID: cb.Message.GetMessageID(),
		Text:      fmt.Sprintf(l("pm-menu.privacy-text"), botUser.FirstName),
		ParseMode: "HTML",
		ReplyMarkup: &telego.InlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		},
	})

	return err
}
