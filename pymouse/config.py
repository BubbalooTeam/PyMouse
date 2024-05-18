import os

from pymouse import log
from dotenv import load_dotenv

# === #
vars = {
    "API_ID": "Required",
    "API_HASH": "Required",
    "BOT_TOKEN": "Required",
    "WORKERS": "Optional",
    "LOG_CHANNEL": "Required",
    "TRIGGER": "Optional",
    "IPV6": "Optional",
    "OWNER_ID": "Required"
}

class MakeConfig:
    def make_config(self):
        with open("config.env", "a") as file:
            for var in vars.keys():
                inp = input("Enter {variable}\n: ".format(variable=var))
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
    WORKERS = int(os.getenv("WORKERS", 24))
    LOG_CHANNEL = int(os.getenv("LOG_CHANNEL"))
    TRIGGER = os.getenv("TRIGGER", "/ !".split())
    IPV6 = bool(os.getenv("IPV6", False))
    OWNER_ID = int(os.getenv("OWNER_ID"))
    VERSION = "1.0.0"