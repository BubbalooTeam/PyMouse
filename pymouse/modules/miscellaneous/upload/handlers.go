package upload

import (
	"fmt"
	"os"
	"pymouse/pymouse/helpers/i18n"
	"pymouse/pymouse/helpers/utils"
	"time"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
)

func Upload(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	l := i18n.Locale(update.Message.Chat)

	query := utils.GetArgs(update)
	if query == "" {
		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				Text:      l("upload.checkers.url-not-provided"),
				ParseMode: "HTML",
			},
		)
		return nil
	}
	msg, _ := bot.SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      l("upload.downloading"),
			ParseMode: "HTML",
		},
	)
	now := time.Now()
	filename, err := downloadByURL(query)
	if err != nil {
		bot.EditMessageText(
			ctx,
			&telego.EditMessageTextParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				MessageID: msg.MessageID,
				Text:      fmt.Sprintf(l("upload.checkers.download-failed"), err),
				ParseMode: "HTML",
			},
		)
		return nil
	}
	bot.EditMessageText(
		ctx,
		&telego.EditMessageTextParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			MessageID: msg.MessageID,
			Text:      l("upload.uploading"),
			ParseMode: "HTML",
		},
	)
	file, err := os.Open(filename)
	if err != nil {
		bot.EditMessageText(
			ctx,
			&telego.EditMessageTextParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				MessageID: msg.MessageID,
				Text:      fmt.Sprintf(l("upload.checkers.upload-failed"), err),
				ParseMode: "HTML",
			},
		)
		return nil
	}
	defer func() {
		bot.DeleteMessage(
			ctx,
			&telego.DeleteMessageParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				MessageID: msg.MessageID,
			},
		)
		file.Close()
		os.Remove(filename)
	}()
	_, err = bot.SendDocument(
		ctx,
		&telego.SendDocumentParams{
			ChatID: telegoutil.ID(update.Message.Chat.ID),
			Document: telego.InputFile{
				File: file,
			},
			Caption:   fmt.Sprintf(l("upload.uploaded"), utils.TimeFormatter(time.Until(now).Abs().Seconds())),
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		},
	)
	return nil
}
