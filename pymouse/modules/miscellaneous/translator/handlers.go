package translator

import (
	"context"
	"fmt"
	"pymouse/pymouse/helpers/i18n"
	"pymouse/pymouse/helpers/utils"
	"strings"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/sirupsen/logrus"
	gt "gopkg.gilang.dev/translator/v2"
)

func Translate(ctx *telegohandler.Context, update telego.Update) error {
	var text string
	bot := ctx.Bot()
	l := i18n.Locale(update.Message.Chat)

	args := utils.GetArgs(update)

	if update.Message.ReplyToMessage != nil {
		replyText := update.Message.ReplyToMessage.Text

		if replyText == "" {
			replyText = update.Message.ReplyToMessage.Caption
		}

		if args != "" {
			text = args + " " + replyText
		} else {
			text = replyText
		}
	} else {
		text = args
	}

	targetLang := getTranslatorLanguage(text, update.Message.Chat)
	if strings.HasPrefix(text, targetLang) {
		text = strings.TrimSpace(strings.Replace(text, targetLang, "", 1))
	}
	if text == "" {
		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				Text:      l("translator.checkers.translator-text-not-provided"),
				ParseMode: "HTML",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			},
		)
		return nil
	}

	msg, _ := bot.SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      l("translator.translating"),
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		},
	)
	translated, err := gt.Translate(
		context.Background(),
		text,
		targetLang,
	)
	if err != nil {
		bot.EditMessageText(
			ctx,
			&telego.EditMessageTextParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				MessageID: msg.MessageID,
				Text:      l("translator.checkers.translator-failed"),
				ParseMode: "HTML",
			},
		)
		logrus.Errorf("failed to translate text: %v", err)
		return nil
	}
	if translated.From.Language.Iso == targetLang && translated.From.Language.Iso != "en" {
		targetLang = "en"
		translated, err = gt.Translate(
			context.Background(),
			text,
			targetLang,
		)
		if err != nil {
			bot.EditMessageText(
				ctx,
				&telego.EditMessageTextParams{
					ChatID:    telegoutil.ID(update.Message.Chat.ID),
					MessageID: msg.MessageID,
					Text:      l("translator.checkers.translator-failed"),
					ParseMode: "HTML",
				},
			)
			logrus.Errorf("failed to translate text: %v", err)
			return nil
		}
	}

	bot.EditMessageText(
		ctx,
		&telego.EditMessageTextParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			MessageID: msg.MessageID,
			Text:      fmt.Sprintf("<b>%s -> %s</b>\n<blockquote><code>%s</code></blockquote>", translated.From.Language.Iso, targetLang, translated.Text),
			ParseMode: "HTML",
		},
	)
	return nil
}
