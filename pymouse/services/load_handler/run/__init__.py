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
            from pymouse.plugins.sudoers.handlers import Sudoers_Plugins
            # Sudoers
            load_modules.add_cmd(Sudoers_Plugins().ping, "ping")
            load_modules.add_cmd(Sudoers_Plugins().rr, ["rr", "restart"])
            load_modules.add_cmd(Sudoers_Plugins().shutdown, "shutdown")
            load_modules.add_cmd(Sudoers_Plugins().speed_test, "speedtest")
            # Final loader => Leave "load_modules.cmds_loader()" last, if you put it before the commands will not be added correctly, causing failures!
            load_modules.cmds_loader()
            # Final loader => Leave "load_modules.regex_loader()" last, if you put it before the regex will not be added correctly, causing failures!
            load_modules.regex_loader()
        except (ImportError, re.error) as e:
            raise LoadModulesError("Modules are not loaded correctly for reason: {error}".format(error=e))