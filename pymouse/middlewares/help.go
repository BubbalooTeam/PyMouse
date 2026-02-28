package middlewares

import (
	"strings"
	"sync"
)

type HelpEntry struct {
	Module          string       `json:"module"`
	DescriptionI18n string       `json:"description_i18n"`
	TitleI18n       string       `json:"title_i18n"`
	Plugins         []*HelpEntry `json:"plugins,omitempty"`
}

type HelpMiddleware struct {
	mu       sync.Mutex
	helpable []*HelpEntry
}

var Help = NewHelp()

func NewHelp() *HelpMiddleware {
	return &HelpMiddleware{
		helpable: []*HelpEntry{},
	}
}

func Slug(name string) string {
	replacer := strings.NewReplacer(
		"ã", "a",
		"õ", "o",
		"ç", "c",
		"é", "e",
		"á", "a",
		"ó", "o",
		"í", "i",
		"ú", "u",
	)
	name = strings.ToLower(name)
	name = replacer.Replace(name)
	return strings.ReplaceAll(name, " ", "")
}

func buildI18nKey(path []string, name string) string {
	s := Slug(name)
	return strings.Join(append([]string{"help", "descriptions"}, append(path, s)...), ".")
}

func buildTitleKey(path []string, name string) string {
	s := Slug(name)
	return strings.Join(append([]string{"help", "titles"}, append(path, s)...), ".")
}

type submodulesOption struct {
	items []string
}

func (h *HelpMiddleware) WithSubmodules(items ...string) submodulesOption {
	return submodulesOption{items: items}
}

func (h *HelpMiddleware) RegisterHelp(main string, args ...any) {

	h.mu.Lock()
	defer h.mu.Unlock()

	mainSlug := Slug(main)

	var mainEntry *HelpEntry
	for _, entry := range h.helpable {
		if strings.EqualFold(entry.Module, main) {
			mainEntry = entry
			break
		}
	}

	if mainEntry == nil {
		mainEntry = &HelpEntry{
			Module:          strings.Title(main),
			DescriptionI18n: buildI18nKey([]string{mainSlug}, mainSlug),
			TitleI18n:       buildTitleKey([]string{mainSlug}, mainSlug),
		}
		h.helpable = append(h.helpable, mainEntry)
	}

	if len(args) == 0 {
		return
	}

	if sub, ok := args[0].(string); ok {

		subSlug := Slug(sub)

		var subEntry *HelpEntry
		for _, p := range mainEntry.Plugins {
			if strings.EqualFold(p.Module, sub) {
				subEntry = p
				break
			}
		}

		if subEntry == nil {
			subEntry = &HelpEntry{
				Module:          strings.Title(sub),
				DescriptionI18n: buildI18nKey([]string{mainSlug}, subSlug),
				TitleI18n:       buildTitleKey([]string{mainSlug}, subSlug),
			}
			mainEntry.Plugins = append(mainEntry.Plugins, subEntry)
		}

		if len(args) > 1 {
			if opt, ok := args[1].(submodulesOption); ok {

				for _, child := range opt.items {

					childSlug := Slug(child)

					exists := false
					for _, c := range subEntry.Plugins {
						if strings.EqualFold(c.Module, child) {
							exists = true
							break
						}
					}

					if !exists {
						subEntry.Plugins = append(subEntry.Plugins, &HelpEntry{
							Module:          strings.Title(child),
							DescriptionI18n: buildI18nKey([]string{mainSlug, subSlug}, childSlug),
							TitleI18n:       buildTitleKey([]string{mainSlug, subSlug}, childSlug),
						})
					}
				}
			}
		}

		return
	}

	if opt, ok := args[0].(submodulesOption); ok {

		for _, child := range opt.items {

			childSlug := Slug(child)

			exists := false
			for _, c := range mainEntry.Plugins {
				if strings.EqualFold(c.Module, child) {
					exists = true
					break
				}
			}

			if !exists {
				mainEntry.Plugins = append(mainEntry.Plugins, &HelpEntry{
					Module:          strings.Title(child),
					DescriptionI18n: buildI18nKey([]string{mainSlug}, childSlug),
					TitleI18n:       buildTitleKey([]string{mainSlug}, childSlug),
				})
			}
		}
	}
}

func (h *HelpMiddleware) GetHelpable() []*HelpEntry {
	return h.helpable
}
