package checkers

import (
	"fmt"
	"strings"

	"pymouse/pymouse/database/repositories"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func SaveUsers(ctx *th.Context, update telego.Update) error {
	var userName string

	message := update.Message

	if message == nil {
		if update.CallbackQuery == nil {
			return ctx.Next(update)
		}

		if msg, ok := update.CallbackQuery.Message.(*telego.Message); ok {
			message = msg
		} else {
			return ctx.Next(update)
		}
	}

	if message.SenderChat != nil {
		return ctx.Next(update)
	}

	userID := message.From.ID
	firstName := message.From.FirstName

	if message.From.Username != "" {
		userName = fmt.Sprintf("%s", message.From.Username)
	}

	repositories.UpdateUser(userID, userName, firstName)

	return ctx.Next(update)
}

func SaveChats(ctx *th.Context, update telego.Update) error {
	var userName string

	message := update.Message

	if message == nil {
		if update.CallbackQuery == nil {
			return ctx.Next(update)
		}

		if msg, ok := update.CallbackQuery.Message.(*telego.Message); ok {
			message = msg
		} else {
			return ctx.Next(update)
		}
	}

	if strings.Contains(message.Chat.Type, "private") || message.SenderChat != nil {
		return ctx.Next(update)
	}

	chatID := message.Chat.ID
	chatTitle := message.Chat.Title

	if message.Chat.Username != "" {
		userName = fmt.Sprintf("%s", message.Chat.Username)
	}

	repositories.UpdateChat(chatID, userName, chatTitle)

	return ctx.Next(update)
}
