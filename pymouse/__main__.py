import asyncio
import time

from hydrogram import idle
from hydrogram.errors import FloodWait, Unauthorized

from pymouse import PyMouse, log
from .services.load_handler.run import RunModules

async def run_mouse():
    RunModules()
    await PyMouse.start()
    await idle()
    await PyMouse.stop()

if __name__ == "__main__" :
    loop = asyncio.get_event_loop()
    try:
        loop.run_until_complete(run_mouse())
    except Unauthorized:
        log.error("Without authorization! Hydrogram made initialization attempts and it was not allowed, invalid BOT TOKEN or LOG CHANNEL.")
    except FloodWait as flood:
        log.warning("Hydrogram is now waiting %s seconds to restart the processes!", flood.value)
        time.sleep(flood.value)
        loop.run_until_complete(run_mouse())
    except KeyboardInterrupt:
        pass
    except Exception as err:
        log.error(err.with_traceback(None))
    finally:
        loop.stop()