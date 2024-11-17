package checkers

import (
	"fmt"
	"pymouse/pymouse/database/utilitiesdb"

	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func SaveUsers(bot *telego.Bot, update telego.Update, next th.Handler) {
	var UserName string

	message := update.Message
	if message == nil {
		if update.CallbackQuery == nil {
			next(bot, update)
			return
		}
		message = update.CallbackQuery.Message.(*telego.Message)
	}
	if message.SenderChat != nil {
		next(bot, update)
		return
	}
	// Telegram user informations
	UserID := message.From.ID
	FirstName := message.From.FirstName
	if message.From.Username != "" {
		UserName = fmt.Sprintf("@%s", message.From.Username)
	}

	// Update or Insert User Informations
	utilitiesdb.UpdateUser(UserID, UserName, FirstName)

	// Pass to next handler
	next(bot, update)
}

func SaveChats(bot *telego.Bot, update telego.Update, next th.Handler) {
	var UserName string

	message := update.Message
	if message == nil {
		if update.CallbackQuery == nil {
			next(bot, update)
			return
		}
		message = update.CallbackQuery.Message.(*telego.Message)
	}
	if strings.Contains(message.Chat.Type, "private") || message.SenderChat != nil {
		next(bot, update)
		return
	}

	// Telegram Chat Informations
	ChatID := message.Chat.ID
	ChatTitle := message.Chat.Title
	if message.Chat.Username != "" {
		UserName = fmt.Sprintf("@%s", message.Chat.Username)
	}

	// Update or Insert Chat Informations
	utilitiesdb.UpdateChat(ChatID, UserName, ChatTitle)

	// Pass to next handler
	next(bot, update)
}
