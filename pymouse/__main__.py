#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.


import asyncio
import sentry_sdk
import time

from traceback import format_exc

from hydrogram import idle
from hydrogram.errors import FloodWait, Unauthorized

from pymouse import PyMouse, Config, log, localization, DownloadPaths
from pymouse.utils import http
from .services.load_handler.run import RunModules

def RestartClean():
    Ddown_path = DownloadPaths().Delete_DownloadPath()
    if Ddown_path in ["Path doesn't exists..", "Deleted directory and all content there.."]:
        DownloadPaths().Make_DownloadPath()
    else:
        log.error("Failed to delete download path.")

def IntegrateSentry(sentry_dsn: str, version: str):
    log.info("Initializing Sentry integration...")
    sentry_sdk.init(
        dsn=sentry_dsn,
        release=version,
    )
    log.info("Sentry started!")

async def run_mouse():
    localization.compile_locales()
    RestartClean()
    IntegrateSentry(
        sentry_dsn=Config.SENTRY_DSN,
        version=Config.VERSION,
    )
    RunModules()
    await PyMouse.start()
    await idle()
    await PyMouse.stop()
    await http.aclose()

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
        log.error(err)
    finally:
        loop.stop()