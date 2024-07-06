#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.


import asyncio
import time

from hydrogram import idle
from hydrogram.errors import FloodWait, Unauthorized

from pymouse import PyMouse, log, localization
from .services.load_handler.run import RunModules

async def run_mouse():
    localization.compile_locales()
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
        e = err.__traceback__
        log.error(err.with_traceback(e))
    finally:
        loop.stop()