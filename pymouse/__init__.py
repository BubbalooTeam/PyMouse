#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.

from HydroPatch.router import Router
from pymouse.utils import log
from pymouse.config import Config
from pymouse.client import PyMouse, START_TIME
from pymouse.database.modules import DataBase

router = Router()

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

from pymouse.utils.tools.DownloadPath import DownloadPaths
from pymouse.utils.localization import localization
from pymouse.decorators import Decorators
from pymouse.utils.tools.weather import Weather
from pymouse.utils.tools.youtube_downloader import YouTubeDownloader, YT_DLP