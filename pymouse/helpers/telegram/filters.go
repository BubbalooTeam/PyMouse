package telegram

import (
	"context"
	"pymouse/pymouse/config"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func Command(cmd string, botUsername string) th.Predicate {
	prefixes := []string{"/", "!"}
	cmd = strings.ToLower(cmd)

	return func(ctx context.Context, update telego.Update) bool {
		var mention string

		if update.Message == nil {
			return false
		}

		text := strings.TrimSpace(update.Message.Text)
		if text == "" {
			return false
		}

		text = strings.ToLower(text)

		parts := strings.Fields(text)
		if len(parts) == 0 {
			return false
		}

		command := parts[0]

		if i := strings.Index(command, "@"); i != -1 {
			mention = command[i+1:]
			command = command[:i]
		}

		botUsername = strings.ToLower(botUsername)
		if mention != "" && mention != botUsername {
			return false
		}

		for _, p := range prefixes {
			if command == p+cmd {
				return true
			}
		}

		return false
	}
}

func IsDeveloper() th.Predicate {
	return func(ctx context.Context, update telego.Update) bool {
		if update.Message != nil {
			return update.Message.From.ID == config.OwnerID
		}

		if update.CallbackQuery != nil {
			return update.CallbackQuery.From.ID == config.OwnerID
		}

		return false
	}
}
