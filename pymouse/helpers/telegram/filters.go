package telegram

import (
	"context"
	"pymouse/pymouse/config"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func Command(cmd string) th.Predicate {
	prefixes := []string{"/", "!"}

	return func(ctx context.Context, update telego.Update) bool {
		if update.Message == nil {
			return false
		}

		text := strings.TrimSpace(update.Message.Text)
		if text == "" {
			return false
		}

		for _, p := range prefixes {
			full := p + cmd

			if text == full || strings.HasPrefix(text, full+" ") {
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