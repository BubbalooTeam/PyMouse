#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.


import os

from pymouse import log
from dotenv import load_dotenv

# === #
vars = {
    "API_ID": "Required",
    "API_HASH": "Required",
    "BOT_TOKEN": "Required",
    "SENTRY_DSN": "Required",
    "CROWDIN_URL": "Optional",
    "IPV6": "Optional",
    "LOG_CHANNEL": "Required",
    "OWNER_ID": "Required",
    "TRIGGER": "Optional",
    "WORKERS": "Optional",
}

class MakeConfig:
    def make_config(self):
        with open("config.env", "a") as file:
            for var in vars.keys():
                inp = input("Enter {variable} -> ({boolean})\n: ".format(
                    variable=var,
                    boolean=vars.get(var, "Optional")
                    )
                )
                if inp == "" and vars.get(var, "Optional") != "Required":
                    continue
                else:
                    while inp == "":
                        inp = input("Please enter a value for {variable}\n: ".format(variable=var))
                    file.write('{variable}="{inputed}"\n'.format(variable=var, inputed=inp))
        log.info("Successfully created config.env file...")
# === #


if os.path.exists("config.env"):
    load_dotenv("config.env")
else:
    MakeConfig().make_config()
    load_dotenv("config.env")


class Config:
    API_ID = int(os.getenv("API_ID"))
    API_HASH = os.getenv("API_HASH")
    BOT_TOKEN = os.getenv("BOT_TOKEN")
    SENTRY_DSN = os.getenv("SENTRY_DSN")
    CROWDIN_URL = os.getenv("CROWDIN_URL", "https://crowdin.com/project/pymouse")
    DOWNLOAD_PATH = "pymouse/downloads/"
    IPV6 = bool(os.getenv("IPV6", False))
    LOG_CHANNEL = int(os.getenv("LOG_CHANNEL"))
    OWNER_ID = int(os.getenv("OWNER_ID"))
    TRIGGER = os.getenv("TRIGGER", "/ !".split())
    VERSION = "1.0.0"
    WORKERS = int(os.getenv("WORKERS", 24))