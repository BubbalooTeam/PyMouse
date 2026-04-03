package translator

import (
	"pymouse/pymouse/database/repositories"
	"strings"

	"github.com/mymmrac/telego"
)

var translatorAvalaibleLanguages = []string{
	`af`, `sq`, `am`, `ar`,
	`hy`, `as`, `ay`, `az`,
	`bm`, `eu`, `be`, `bn`,
	`bho`, `bs`, `bg`, `ca`,
	`ceb`, `zh`, `co`, `hr`,
	`cs`, `da`, `dv`, `doi`,
	`nl`, `en`, `eo`, `et`,
	`ee`, `fil`, `fi`, `fr`,
	`fy`, `gl`, `ka`, `de`,
	`el`, `gn`, `gu`, `ht`,
	`ha`, `haw`, `he`, `iw`,
	`hi`, `hmn`, `hu`, `is`,
	`ig`, `ilo`, `id`, `ga`,
	`it`, `ja`, `jv`, `jw`,
	`kn`, `kk`, `km`, `rw`,
	`gom`, `ko`, `kri`, `ku`,
	`ckb`, `ky`, `lo`, `la`,
	`lv`, `ln`, `lt`, `lg`,
	`lb`, `mk`, `mai`, `mg`,
	`ms`, `ml`, `mt`, `mi`,
	`mr`, `mni`, `lus`, `mn`,
	`my`, `ne`, `no`, `ny`,
	`or`, `om`, `ps`, `fa`,
	`pl`, `pt`, `pa`, `qu`,
	`ro`, `ru`, `sm`, `sa`,
	`gd`, `nso`, `sr`, `st`,
	`sn`, `sd`, `si`, `sk`,
	`sl`, `so`, `es`, `su`,
	`sw`, `sv`, `tl`, `tg`,
	`ta`, `tt`, `te`, `th`,
	`ti`, `ts`, `tr`, `tk`,
	`ak`, `uk`, `ur`, `ug`,
	`uz`, `vi`, `cy`, `xh`,
	`yi`, `yo`, `zu`,
}

func getTranslatorLanguage(text string, chat telego.Chat) string {
	checkLang := func(lang string) bool {
		for _, l := range translatorAvalaibleLanguages {
			if l == lang {
				return true
			}
		}
		return false
	}

	chatLang := repositories.GetChatLanguage(chat)
	chatParts := strings.Split(chatLang, "_")

	defaultLang := "en"
	if len(chatParts) > 0 {
		defaultLang = chatParts[0]
	}

	fields := strings.Fields(text)
	if len(fields) == 0 {
		return defaultLang
	}

	lang := fields[0]
	langParts := strings.Split(lang, "_")

	if !checkLang(langParts[0]) {
		return defaultLang
	}

	if len(langParts) > 1 {
		if !checkLang(langParts[1]) {
			if len(chatParts) > 1 {
				return chatParts[0]
			}
			return "en"
		}
	}

	return lang
}
