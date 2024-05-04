import time

from hydrogram import Client, __version__
from hydrogram.raw.all import layer
from hydrogram.enums import ParseMode

from pymouse import Config, log

START_TIME = time.time()

class PyMouseBOT(Client):
    def __init__(self):
        super().__init__(
            name="PyMouseBOT",
            app_version=Config.VERSION,
            api_id=Config.API_ID,
            api_hash=Config.API_HASH,
            bot_token=Config.BOT_TOKEN,
            parse_mode=ParseMode.HTML,
            workers=Config.WORKERS,
            ipv6=Config.IPV6,
            in_memory=True,
        )

    async def start(self):
        await super().start()
        self.me = await self.get_me()
        # //
        text = "<b>PyMouse v{version} started!</b>\n\n"
        text += "<b>System:</b> <code>{sys}</code>\n"
        text += "<b>Hydrogram:</b> <code>{hydrogram}</code> (Layer {layer})"
        # //
        await self.send_message(
            chat_id=Config.LOG_CHANNEL,
            text=text.format(
                version=Config.VERSION,
                sys=self.system_version,
                hydrogram=__version__,
                layer=layer
            )
        )
        log.info(
            "PyMouse running with Hydrogram v%s (Layer %s) started on @%s. Hello!",
            __version__,
            layer,
            self.me.username
        )
    # //
    async def stop(self):
        text_ = "<b>PyMouse stopped with v{version}!</b>\n\n"
        text_ += "<b>System</b>: <code>{sys}</code>"
        await self.send_message(
            chat_id=Config.LOG_CHANNEL,
            text=text_.format(
                version=Config.VERSION,
                sys=self.system_version
            )
        )
        await super().stop()
        log.info("PyMouse is stopped!")


PyMouse = PyMouseBOT()