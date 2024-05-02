from pymouse.utils import log
from pymouse.database.modules import DataBase
from pymouse.config import Config
from pymouse.client import PyMouse, START_TIME

log.info("Compiling DataBase...")
try:
    db = DataBase()
except Exception:
    log.info(f"Error in compiling DataBase...")
    exit()
log.info("DataBase compiled with sucessfully...")