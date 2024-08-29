#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.


import re
from ..exceptions import LoadModulesError
from ..handlers import load_modules

class RunModules:
    def __init__(self):
        try:
            from pymouse.plugins.afk.handlers import AFK_Plugins
            from pymouse.plugins.medias.handlers import Medias_Plugins
            from pymouse.plugins.miscellaneous.handlers import Misccellaneous_Plugins
            from pymouse.plugins.pm_menu.handlers import PMMenu_Plugins
            from pymouse.plugins.sudoers.handlers import Sudoers_Plugins
            # Sudoers
            load_modules.add_cmd(AFK_Plugins().setupAFK, "afk")
            load_modules.add_regex(AFK_Plugins().setupAFK, r"^(?i:brb)(\s(.+))?")
            load_modules.add_cmd(Medias_Plugins().ytdl_handler, "ytdl")
            load_modules.add_cmd(Misccellaneous_Plugins().weatherHandle, "weather")
            load_modules.add_cmd(PMMenu_Plugins().privacyPolicy, "privacy")
            load_modules.add_cmd(PMMenu_Plugins().start_, "start")
            load_modules.add_cmd(Sudoers_Plugins().ping, "ping")
            load_modules.add_cmd(Sudoers_Plugins().rr, ["rr", "restart"])
            load_modules.add_cmd(Sudoers_Plugins().shutdown, "shutdown")
            load_modules.add_cmd(Sudoers_Plugins().speed_test, "speedtest")
            # Final loader => Leave "load_modules.cmds_loader()" last, if you put it before the commands will not be added correctly, causing failures!
            load_modules.cmds_loader()
            # Final loader => Leave "load_modules.regex_loader()" last, if you put it before the regex will not be added correctly, causing failures!
            load_modules.regex_loader()
            # Callbacks
            load_modules.add_callback_btn(PMMenu_Plugins().start_, r"^StartBack$")
            load_modules.add_callback_btn(PMMenu_Plugins().privacyPolicy, r"^PrivacyPolicy$")
            load_modules.add_callback_btn(PMMenu_Plugins().privacyPolicyRead, r"^PrivacyData$")
            load_modules.add_callback_btn(PMMenu_Plugins().ReadYourData, r"^YourDataCollected$")
            load_modules.add_callback_btn(PMMenu_Plugins().ChangeLanguageMenu, r"^LangMenu\|(.*)$")
            load_modules.add_callback_btn(PMMenu_Plugins().SelectLanguageMenu, r"ChangeLanguage\|(.*)$")
            load_modules.add_callback_btn(PMMenu_Plugins().SwitchLanguage, r"SwitchLang\|(.*)$")
            load_modules.add_callback_btn(Medias_Plugins().download_handler, r"yt_(gen|dl)\|(.*)$")
            load_modules.add_callback_btn(Medias_Plugins().ytdl_scroll_callback, r"ytdl_scroll\|(.*)$")

            load_modules.callbacks_loader()
        except (ImportError, re.error) as e:
            raise LoadModulesError("Modules are not loaded correctly for reason: {error}".format(error=e))