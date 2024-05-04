from .utils import log
from pymouse.config import Config
from pymouse.client import PyMouse, START_TIME
from pymouse.database.modules import DataBase


log.info("Compiling DataBase...")
try:
    db = DataBase()
except Exception:
    log.info(f"Error in compiling DataBase...")
    exit()
log.info("DataBase compiled with sucessfully...")

# if use database input here
from pymouse.decorators import Decorators