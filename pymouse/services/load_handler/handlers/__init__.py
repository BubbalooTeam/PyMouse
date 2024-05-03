import logging

from typing import Union, Callable, Pattern, Optional

from hydrogram.filters import (
    command as Command,
    regex as Regex,
)

from hydrogram.handlers import MessageHandler, CallbackQueryHandler

from pymouse import PyMouse, Config

class ModulesLoader:
    def __init__(self):
        self.cmds = []
        self.callbacks = []
        self.regex = []
        self.helpable = {}
        self.regex_count = 0
        self.cmds_count = 0
        self.callbacks_count = 0
    # Message

    def add_cmd(
        self,
        function: Callable,
        commands: Union[str, list[str]],
        plugin_name: Optional[str] = "",
        about_commands: Optional[str] = "",
    ):
        """
        Handler-Mod system, which helps in handling commands in a simple and quick way to configure, automating information saving processes in a TableDict.

        Arguments:
            `function`: Callable,
            `commands`: Union[str, list[str]],
            `plugin_name`: Optional[str],
            `about_commands`: Optional[str],

        `Use something like: load_modules.add_cmd(func, ["cmd1", "cmd2"], "Plugin1", "Help for Plugin1")`
        """
        self.cmds.append(MessageHandler(callback=function, filters=Command(commands, Config.TRIGGER)))

        if not plugin_name:
            return

        table_cmd = {
            "functions": [
                {
                    "commands": commands,
                    "about_commands": about_commands,
                }
            ]
        }

        if plugin_name in self.helpable:
            self.helpable[plugin_name]["functions"].append(table_cmd["functions"])
        else:
            self.helpable[plugin_name] = table_cmd

    def add_cmd_no_trg(
            self,
            function: Callable,
            commands: Union[str, list[str]],
        ):
        self.cmds.append(MessageHandler(callback=function, filters=Command(commands, "")))

    def cmds_loader(self):
        dispatcher = self.cmds
        for handler in dispatcher:
            PyMouse().add_handler(handler=handler)
            self.cmds_count += 1
        logging.info("%s Loaded Command Handler(s)!", self.cmds_count)

    # Message Regex

    def add_regex(
            self,
            function: Callable,
            pattern: Union[str, Pattern],
    ):
        self.regex.append(MessageHandler(function, Regex(pattern)))

    def regex_loader(self):
        dispatcher = self.regex
        for handler in dispatcher:
            PyMouse().add_handler(handler=handler)
            self.regex_count += 1
        logging.info("%s Loaded Regex Filters Handler(s)!", self.regex_count)

    # Callback

    def add_callback_btn(
            self,
            function: Callable,
            pattern: Union[str, Pattern],
    ):
        self.callbacks.append(CallbackQueryHandler(function, Regex(pattern)))

    def callbacks_loader(self):
        dispatcher = self.callbacks
        for handler in dispatcher:
            PyMouse().add_handler(handler=handler)
            self.callbacks_count += 1
        logging.info("%s Loaded CallBack Buttons Handler(s)!", self.callbacks_count)
    

load_modules = ModulesLoader()