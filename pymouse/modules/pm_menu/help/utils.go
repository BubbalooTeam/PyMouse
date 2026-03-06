package help

import (
	"pymouse/pymouse/middlewares"
	"sort"
	"strings"

	"github.com/mymmrac/telego"
)

func FindModule(path []string, helpable []*middlewares.HelpEntry) *middlewares.HelpEntry {

	if len(path) == 0 || len(helpable) == 0 {
		return nil
	}

	rootSlug := path[0]

	var current *middlewares.HelpEntry

	for _, entry := range helpable {
		if middlewares.Slug(entry.Module) == rootSlug {
			current = entry
			break
		}
	}

	if current == nil {
		return nil
	}

	if len(path) == 1 {
		return current
	}

	if len(path) == 2 && path[0] == path[1] {
		return current
	}

	for i := 1; i < len(path); i++ {

		subSlug := path[i]
		found := false

		for _, sub := range current.Plugins {
			if middlewares.Slug(sub.Module) == subSlug {
				current = sub
				found = true
				break
			}
		}

		if !found {
			return nil
		}
	}

	return current
}

func GenerateHelpKeyboard(helpable []*middlewares.HelpEntry, l func(string) string) *telego.InlineKeyboardMarkup {

	if len(helpable) == 0 {
		helpable = middlewares.Help.GetHelpable()
	}

	sort.Slice(helpable, func(i, j int) bool {
		var a, b string

		if helpable[i].TitleI18n != "" {
			a = l(helpable[i].TitleI18n)
		} else {
			a = helpable[i].Module
		}

		if helpable[j].TitleI18n != "" {
			b = l(helpable[j].TitleI18n)
		} else {
			b = helpable[j].Module
		}

		return strings.ToLower(a) < strings.ToLower(b)
	})

	var rows [][]telego.InlineKeyboardButton
	var currentRow []telego.InlineKeyboardButton

	for _, node := range helpable {

		name := node.Module
		titleKey := node.TitleI18n

		var buttonTitle string

		if titleKey != "" {
			buttonTitle = l(titleKey)
		} else {
			buttonTitle = name
		}

		if buttonTitle == "" {
			continue
		}

		nameSlug := middlewares.Slug(name)
		callback := "help:" + nameSlug + "." + nameSlug

		currentRow = append(currentRow, telego.InlineKeyboardButton{
			Text:         buttonTitle,
			CallbackData: callback,
		})

		if len(currentRow) == 3 {
			rows = append(rows, currentRow)
			currentRow = nil
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}
