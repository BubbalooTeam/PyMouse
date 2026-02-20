package pm_menu

import (
	"pymouse/pymouse/client"
	"pymouse/pymouse/helpers/telegram"
	"pymouse/pymouse/modules/pm_menu/about"
	"pymouse/pymouse/modules/pm_menu/help"
	"pymouse/pymouse/modules/pm_menu/localization"
	"pymouse/pymouse/modules/pm_menu/privacy"
	"pymouse/pymouse/modules/pm_menu/start"
	"regexp"

	th "github.com/mymmrac/telego/telegohandler"
)

func LoadModule(bS *client.BotStruct) {
	// About Handlers
	bS.Handler.Handle(about.AboutMessage, th.TextContains("/start about"))
	bS.Handler.Handle(about.AboutMessage, telegram.Command("about"))
	bS.Handler.Handle(about.AboutCallback, th.CallbackDataMatches(regexp.MustCompile(`^AboutMenu$`)))
	// Language Handlers
	bS.Handler.Handle(localization.ChangeLanguageMessage, th.TextContains("/start lang"))
	bS.Handler.Handle(localization.ChangeLanguageMessage, telegram.Command("lang"))
	bS.Handler.Handle(localization.ChangeLanguageCallback, th.CallbackDataMatches(regexp.MustCompile(`^LangMenu\|(.*)$`)))
	bS.Handler.Handle(localization.SelectLanguageCallback, th.CallbackDataMatches(regexp.MustCompile(`^ChangeLang\|(.*)$`)))
	bS.Handler.Handle(localization.SwitchLanguageCallback, th.CallbackDataMatches(regexp.MustCompile(`^SwitchLang\|(.*)$`)))
	// Help Handlers
	bS.Handler.Handle(help.HelpMenu, th.TextContains("/start help"))
	bS.Handler.Handle(help.HelpMenu, telegram.Command("help"))
	bS.Handler.Handle(help.HelpMenuCallback, th.CallbackDataMatches(regexp.MustCompile(`^HelpMenu$`)))
	bS.Handler.Handle(help.HelpModule, th.CallbackDataPrefix("help:"))
	// Privacy Handlers
	bS.Handler.Handle(privacy.PrivacyPolicyMessage, th.TextContains("/start privacy"))
	bS.Handler.Handle(privacy.PrivacyPolicyMessage, telegram.Command("privacy"))
	bS.Handler.Handle(privacy.PrivacyPolicyCallback, th.CallbackDataMatches(regexp.MustCompile(`^PrivacyPolicy$`)))
	// Start Handlers
	bS.Handler.Handle(start.StartMessage, telegram.Command("start"))
	bS.Handler.Handle(start.StartBackCallback, th.CallbackDataMatches(regexp.MustCompile(`^StartBack$`)))
}
