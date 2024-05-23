from ..exceptions import LoadModulesError
from ..handlers import load_modules

class RunModules:
    def __init__(self):
        try:
            from pymouse.plugins import (
                afk,
                sudoers,
            )
            # Sudoers
            load_modules.add_cmd(afk.AFK_Plugins().setupAFK, "afk", "AFK", "Mark yourself away (away from the keyboard).")
            load_modules.add_regex(afk.AFK_Plugins().setupAFK, r"^(?i:brb)(\s(.+))?")
            load_modules.add_cmd(sudoers.Sudoers_Plugins().ping, "ping")
            load_modules.add_cmd(sudoers.Sudoers_Plugins().rr, ["rr", "restart"])
            load_modules.add_cmd(sudoers.Sudoers_Plugins().shutdown, "shutdown")
            load_modules.add_cmd(sudoers.Sudoers_Plugins().speed_test, "speedtest")
            # Final loader => Leave "load_modules.cmds_loader()" last, if you put it before the commands will not be added correctly, causing failures!
            load_modules.cmds_loader()
            # Final loader => Leave "load_modules.regex_loader()" last, if you put it before the regex will not be added correctly, causing failures!
            load_modules.regex_loader()
        except LoadModulesError as e:
            raise LoadModulesError("[services/load_handler/run/RunModules]: Modules are not loaded correctly for reason: {error}".format(error=e))