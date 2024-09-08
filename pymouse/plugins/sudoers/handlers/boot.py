#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.

import os
import sys
import signal

from hydrogram import filters
from hydrogram.types import Message

from pymouse import PyMouse, Decorators, router

@router.message(filters.command(["rr", "restart"]))
@Decorators().require_dev()
async def rr(c: PyMouse, m: Message): #type: ignore
    sent = await m.reply("<i>Restarting...</i>")
    args = [sys.executable, "-m", "pymouse"]
    await sent.edit("<b>PyMouse is now Restarted!</b>")
    os.execl(sys.executable, *args)

@router.message(filters.command("shutdown"))
@Decorators().require_dev(only_owner=True)
async def shutdown(c: PyMouse, m: Message): # type: ignore
    await m.reply("<b>PyMouse is now Offline!</b>")
    os.kill(os.getpid(), signal.SIGINT)