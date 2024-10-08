#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.

import json
import asyncio
import sentry_sdk
import time

from apscheduler.schedulers.asyncio import AsyncIOScheduler
from datetime import datetime
from functools import partial
from io import BytesIO

from hydrogram import idle
from hydrogram.enums import ChatAction
from hydrogram.errors import FloodWait, Unauthorized
from HydroPatch import PatchManager

from pymouse import PyMouse, Config, db, log, localization, DownloadPaths, router
from pymouse.utils import http

sheduler = AsyncIOScheduler()
patcher = PatchManager(client=PyMouse)

async def BackupDB(c: PyMouse): # type: ignore
    log.info("Backing Up PyMouse DataBase...")
    dbData = db.db
    if dbData:
        dumpData = json.dumps(obj=dbData, indent=4)
        file = BytesIO(dumpData.encode())
        file.name = "dbBackup.json"

        await c.send_chat_action(
            chat_id=Config.LOG_CHANNEL,
            action=ChatAction.UPLOAD_DOCUMENT
        )
        await c.send_document(
            chat_id=Config.LOG_CHANNEL,
            document=file,
            caption=(
                "<b>#DATABASE #PyMOUSE #BACKUP</b>\n\n"
                + "<b>Hour:</b> <code>{hour}</code>\n".format(
                    hour=datetime.now().strftime("%H:%M:%S")
                )
                + "<b>Date:</b> <code>{date}</code>".format(
                    date=datetime.now().strftime("%d/%m/%Y")
                )
            )
        )
        log.info("The DB database has been backed up.")
    else:
        log.info("Not data in Local DataBase...")

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
    patcher.include_router(router=router)
    await PyMouse.start()
    await BackupDB(
        c=PyMouse
    )
    sheduler.add_job(
        func=partial(
            BackupDB,
            c=PyMouse
        ),
        trigger="cron",
        hour=0,
        minute=0,
        id="Backup Local DataBase."
    )
    sheduler.start()
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