package speedtst

import (
	"fmt"
	"pymouse/pymouse/helpers/utils"
	"time"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/showwin/speedtest-go/speedtest"
)

func SpeedTestMessage(ctx *telegohandler.Context, update telego.Update) error {
	bot := ctx.Bot()
	chatID := telegoutil.ID(update.Message.Chat.ID)

	sent, err := bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:    chatID,
		Text:      "<code>Running speed test...</code>",
		ParseMode: telego.ModeHTML,
		ReplyParameters: &telego.ReplyParameters{
			MessageID: update.Message.MessageID,
		},
	})
	if err != nil {
		return err
	}

	resultChan := make(chan string, 1)

	go func() {
		run := speedtest.New()

		servers, err := run.FetchServers()
		if err != nil {
			resultChan <- "<b>SpeedTest failed.</b>"
			return
		}

		available, err := servers.FindServer([]int{})
		if err != nil || len(available) == 0 {
			resultChan <- "<b>No servers available.</b>"
			return
		}

		for _, s := range available {
			if err := s.PingTest(nil); err != nil {
				continue
			}
			if err := s.DownloadTest(); err != nil {
				continue
			}
			if err := s.UploadTest(); err != nil {
				continue
			}

			resultChan <- fmt.Sprintf(
				"<b>SpeedTest completed!</b>\n"+
					"<b>🌀 Name:</b> <code>%s</code>\n"+
					"<b>🌐 Host:</b> <code>%s</code>\n"+
					"<b>🏁 Country:</b> <code>%s</code>\n\n"+
					"<b>⬇️ Download:</b> <code>%s/s</code>\n"+
					"<b>⬆️ Upload:</b> <code>%s/s</code>\n"+
					"<b>🏓 Ping:</b> <code>%d ms</code>",
				s.Name,
				s.Sponsor,
				s.Country,
				utils.HumanBytes(int64(s.DLSpeed)),
				utils.HumanBytes(int64(s.ULSpeed)),
				s.Latency.Milliseconds(),
			)
			return
		}

		resultChan <- "<b>SpeedTest failed.</b>"
	}()

	var result string

	select {
	case result = <-resultChan:
	case <-time.After(2 * time.Minute):
		result = "<b>SpeedTest timed out.</b>"
	}

	_, err = bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:    chatID,
		MessageID: sent.MessageID,
		Text:      result,
		ParseMode: telego.ModeHTML,
	})

	return err
}
