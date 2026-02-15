package gsmarena

import (
	"fmt"
	"pymouse/pymouse/helpers/utils"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/patrickmn/go-cache"
)

var searchCache = cache.New(15*time.Minute, 30*time.Minute)

func DeviceSearch(ctx *th.Context, update telego.Update) error {
	bot := ctx.Bot()
	query := utils.GetArgs(update)

	if query == "" {
		bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      "<b>Device not provided!</b>\n<i>Uso:</i> <code>/d Moto G34</code>",
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return nil
	}

	results := searchDevice(query).Results

	switch len(results) {

	case 0:
		bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      "<b>Device not found!</b>",
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return nil

	case 1:
		deviceSpecs := fetchDevice(results[0].ID)
		formattedMsg := formatGSMarenaMessage(deviceSpecs)

		bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID:    telegoutil.ID(update.Message.Chat.ID),
			Text:      formattedMsg,
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return nil

	default:
		searchID := uuid.NewString()[:8]
		searchCache.Set(searchID, results, cache.DefaultExpiration)

		btns := GSMarenaCreateKeyboard(
			results,
			int(update.Message.From.ID),
			searchID,
			1,
		)

		bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: telegoutil.ID(update.Message.Chat.ID),
			Text: fmt.Sprintf(
				"I was returned different variations of different devices with the query <b>%s</b>, choose a model from these devices using the buttons below.",
				query,
			),
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
			ReplyMarkup: btns,
		})
		return nil
	}
}

func DeviceSearchPagination(ctx *th.Context, update telego.Update) error {
	bot := ctx.Bot()
	callback := update.CallbackQuery

	parts := strings.Split(callback.Data, "|")
	if len(parts) != 4 {
		return nil
	}

	page, _ := strconv.Atoi(parts[1])
	userID, _ := strconv.Atoi(parts[2])
	searchID := parts[3]

	// Segurança: apenas o dono pode usar
	if int(callback.From.ID) != userID {
		bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
			Text:            "You cannot use this pagination.",
			ShowAlert:       true,
		})
		return nil
	}

	cached, found := searchCache.Get(searchID)
	if !found {
		bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
			Text:            "Search expired. Please search again.",
			ShowAlert:       true,
		})
		return nil
	}

	results := cached.([]GSMArenaDeviceSearchResult)

	btns := GSMarenaCreateKeyboard(results, userID, searchID, page)

	bot.EditMessageReplyMarkup(ctx, &telego.EditMessageReplyMarkupParams{
		ChatID:      telegoutil.ID(callback.Message.GetChat().ID),
		MessageID:   callback.Message.GetMessageID(),
		ReplyMarkup: btns,
	})

	bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: callback.ID,
	})

	return nil
}

func DeviceSearchSelect(ctx *th.Context, update telego.Update) error {
	bot := ctx.Bot()
	callback := update.CallbackQuery

	parts := strings.Split(callback.Data, "|")
	if len(parts) != 4 {
		return nil
	}

	deviceID := parts[1]

	userID, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil
	}

	searchID := parts[3]

	if int(callback.From.ID) != userID {
		bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
			Text:            "You cannot select this device.",
			ShowAlert:       true,
		})
		return nil
	}

	cached, found := searchCache.Get(searchID)
	if !found {
		bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
			Text:            "Search expired. Please search again.",
			ShowAlert:       true,
		})
		return nil
	}

	results := cached.([]GSMArenaDeviceSearchResult)

	valid := false
	for _, d := range results {
		if d.ID == deviceID {
			valid = true
			break
		}
	}
	if !valid {
		return nil
	}

	deviceSpecs := fetchDevice(deviceID)
	formattedMsg := formatGSMarenaMessage(deviceSpecs)

	msg := callback.Message
	if msg == nil {
		return nil
	}

	bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:    telegoutil.ID(msg.GetChat().ID),
		MessageID: msg.GetMessageID(),
		Text:      formattedMsg,
		ParseMode: "HTML",
	})

	bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: callback.ID,
	})

	return nil
}
