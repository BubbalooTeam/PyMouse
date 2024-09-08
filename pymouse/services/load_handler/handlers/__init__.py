#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.


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
        self.regex_count = 0
        self.cmds_count = 0
        self.callbacks_count = 0
    # Message

    def add_cmd(
        self,
        function: Callable,
        commands: Union[str, list[str]],
    ):
        self.cmds.append(MessageHandler(callback=function, filters=Command(commands, Config.TRIGGER)))

    def add_cmd_no_trg(
            self,
            function: Callable,
            commands: Union[str, list[str]],
        ):
        self.cmds.append(MessageHandler(callback=function, filters=Command(commands, "")))

    def cmds_loader(self):
        dispatcher = self.cmds
        for handler in dispatcher:
            PyMouse.add_handler(handler=handler)
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
            PyMouse.add_handler(handler=handler)
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
            PyMouse.add_handler(handler=handler)
            self.callbacks_count += 1
        logging.info("%s Loaded CallBack Buttons Handler(s)!", self.callbacks_count)


load_modules = ModulesLoader()