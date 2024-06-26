from .utils import log
from pymouse.config import Config
from pymouse.client import PyMouse, START_TIME
from pymouse.database.modules import DataBase


log.info("Compiling DataBase...")
try:
    db = DataBase()
except Exception:
    log.info("Error in compiling DataBase...")
    exit()
log.info("DataBase compiled with successfully...")


# Bot utils
from pymouse.database.plugins.utilities.users import usersmodel_db
from pymouse.database.plugins.utilities.afk import afkmodel_db
from pymouse.database.plugins.utilities.chats import chatsmodel_db
from pymouse.database.plugins.utilities.localization import localizationmodel_db

from .utils.localization import localization
from pymouse.decorators import Decorators