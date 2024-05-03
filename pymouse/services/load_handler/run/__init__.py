from ..handlers import load_modules

class RunModules:
    def __init__(self):
        from pymouse.plugins import (
            sudoers
        )
        # Sudoers
        load_modules.add_cmd(sudoers.Sudoers_Plugins().ping, "ping")
        load_modules.add_cmd(sudoers.Sudoers_Plugins().rr, ["rr", "restart"])
        load_modules.add_cmd(sudoers.Sudoers_Plugins().shutdown, "shutdown")
        # Final loader => Leave "load_modules.cmds_loader()" last, if you put it before the commands will not be added correctly, causing failures!
        load_modules.cmds_loader()