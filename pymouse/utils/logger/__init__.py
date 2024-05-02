import logging
from logging.handlers import RotatingFileHandler

logging.basicConfig(
    level=logging.INFO,
    format="%(levelname)s | \u001B[35m%(name)s |\u001B[33m %(asctime)s | \u001B[37m%(message)s",
    datefmt="%m/%d %H:%M:%S",
    handlers=[
        RotatingFileHandler(
        "PyMouse.log", maxBytes=2801100, backupCount=1),
        logging.StreamHandler()
    ]
)

# To avoid some annoying log
logging.getLogger("hydrogram").setLevel(logging.WARNING)
logging.getLogger("hydrogram.parser.html").setLevel(logging.ERROR)
logging.getLogger("hydrogram.session.session").setLevel(logging.ERROR)
logging.getLogger("httpx").setLevel(logging.ERROR)

log: logging.Logger = logging.getLogger(__name__)

log.info("\033[1m\033[35mPyMouse\033[0m")
log.info("\033[96mProject maintained by:\033[0m BubbalooTeam")