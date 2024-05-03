import os

from pymouse.utils.logger import log
from dotenv import load_dotenv

# === #
class MakeConfig:
    def make_config(self):
        with open("config.env", "a") as file:
            for var in [
                "API_ID", 
                "API_HASH", 
                "BOT_TOKEN", 
                "WORKERS", 
                "LOG_CHANNEL", 
                "TRIGGER",
                "IPV6"
            ]:
                inp = input("Enter {variable}\n: ".format(variable=var))
                file.write("{variable}='{inputed}'\n".format(variable=var, inputed=inp))
        log.info("Sucessfully created config.env file...")
# === #
if os.path.exists("config.env"):
    load_dotenv("config.env")
else:
    MakeConfig().make_config()

class Config:
    API_ID = int(os.getenv("API_ID"))
    API_HASH = os.getenv("API_HASH")
    BOT_TOKEN = os.getenv("BOT_TOKEN")
    WORKERS = int(os.getenv("WORKERS", 24))
    LOG_CHANNEL = int(os.getenv("LOG_CHANNEL"))
    TRIGGER = os.getenv("TRIGGER", "/ !".split())
    IPV6 = bool(os.getenv("IPV6", False))
    VERSION = "1.0.0"