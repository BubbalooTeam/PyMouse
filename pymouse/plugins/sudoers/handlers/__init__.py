import os
import sys
import signal

from hydrogram.types import Message, InputMediaPhoto

from datetime import datetime

from pymouse import PyMouse, db, Decorators, log
from pymouse.utils import NetworkUtils, NetworkEvents

class Sudoers_Plugins:
    @staticmethod
    @Decorators().require_dev()
    async def rr(c: PyMouse, m: Message): #type: ignore
        sent = await m.reply("<i>Restarting...</i>") 
        args = [sys.executable, "-m", "pymouse"]
        await sent.edit("<b>PyMouse is now Restarted!</b>")
        os.execl(sys.executable, *args)

    @staticmethod
    @Decorators().require_dev(only_owner=True)
    async def shutdown(c: PyMouse, m: Message): # type: ignore
        await m.reply("<b>PyMouse is now Offline!</b>")
        os.kill(os.getpid(), signal.SIGINT)

    @staticmethod
    async def ping(c: PyMouse, m: Message): # type: ignore
        first = datetime.now()
        sent = await m.reply_text("<b>Pong!</b>")
        second = datetime.now()
        await sent.edit_text(
            f"<b>Pong!</b> <code>{(second - first).microseconds / 1000}</code>ms"
        )

    @staticmethod
    @Decorators().require_dev()
    async def speed_test(c: PyMouse, m: Message): # type: ignore
        msg = await m.reply_photo(photo=NetworkEvents.RUNNING_SPEEDTEST, caption="<b>Running SpeedTest...</b>")
        try:
            # Get infos with speedtest and unpacking infos
            dl, ul, name, host, ping, isp, country, cc, path = await NetworkUtils().speedtest_performer()
            await msg.edit_media(
                media=InputMediaPhoto(
                    media=path,
                    caption=f"🌀 <b>Name:</b> <code>{name}</code>\n🌐 <b>Host:</b> <code>{host}</code>\n🏁 <b>Country:</b> <code>{country}, {cc}</code>\n\n🏓 <b>Ping:</b> <code>{ping} ms</code>\n🔽 <b>Download:</b> <code>{dl} Mbps</code>\n🔼 <b>Upload:</b> <code>{ul} Mbps</code>\n🖥  <b>ISP:</b> <code>{isp}</code>"
                )
            )
        except Exception:
            log.error("[speedtest/handler]: Error in performing speedtest...")