package upload

import (
	"fmt"
	"os"
	"pymouse/pymouse/helpers/utils"
	"time"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
)

func Upload(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()

	query := utils.GetArgs(update)
	if query == "" {
		bot.SendMessage(
			ctx,
			&telego.SendMessageParams{
				ChatID: telegoutil.ID(update.Message.Chat.ID),
				Text:   "Please provide a URL to upload.",
			},
		)
		return nil
	}
	msg, _ := bot.SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID: telegoutil.ID(update.Message.Chat.ID),
			Text:   "Downloading the file, please wait...",
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
				Text:      "Failed to download the file.",
			},
		)
		return nil
	}
	bot.EditMessageText(
		ctx,
		&telego.EditMessageTextParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			MessageID: msg.MessageID,
			Text:      "Uploading file...",
		},
	)
	file, err := os.Open(filename)
	if err != nil {
		bot.EditMessageText(
			ctx,
			&telego.EditMessageTextParams{
				ChatID:    telegoutil.ID(update.Message.Chat.ID),
				MessageID: msg.MessageID,
				Text:      "Failed to open the file.",
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
			Caption: fmt.Sprintf("<b>Time:</b> <code>%s</code>", utils.TimeFormatter(time.Until(now).Abs().Seconds())),
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		},
	)
	return nil
}
