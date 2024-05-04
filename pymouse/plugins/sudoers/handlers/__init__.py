import os
import sys
import signal

from hydrogram.types import Message

from datetime import datetime

from pymouse import PyMouse, db, Decorators

class Sudoers_Plugins:
    @staticmethod
    @Decorators().require_dev()
    async def rr(c: PyMouse, m: Message):
        sent = await m.reply("<i>Restarting...</i>") 
        args = [sys.executable, "-m", "pymouse"]
        await sent.edit("<b>PyMouse is now Restarted!</b>")
        os.execl(sys.executable, *args)

    @staticmethod
    @Decorators().require_dev(only_owner=True)
    async def shutdown(c: PyMouse, m: Message):
        await m.reply("<b>PyMouse is now Offline!</b>")
        os.kill(os.getpid(), signal.SIGINT)

    @staticmethod
    @Decorators().require_dev()
    async def ping(c: PyMouse, m: Message):
        first = datetime.now()
        sent = await m.reply_text("<b>Pong!</b>")
        second = datetime.now()
        await sent.edit_text(
            f"<b>Pong!</b> <code>{(second - first).microseconds / 1000}</code>ms"
        )