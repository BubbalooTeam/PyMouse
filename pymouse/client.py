import time

from hydrogram import Client, __version__
from hydrogram.raw.all import layer
from hydrogram.enums import ParseMode

from pymouse import Config, log

START_TIME = time.time()

class PyMouse(Client):
    def __init__(self):
        self.name: str = self.__class__.__name__.lower()
        self.version: str = Config.VERSION
        self.api_id: str = Config.API_ID
        self.api_hash: str = Config.API_HASH
        self.bot_token: str = Config.BOT_TOKEN
        self.logger: str = Config.LOG_CHANNEL
        self.workers: str = Config.WORKERS
        self.ipv6: bool = Config.IPV6
        super().__init__(
            name=self.name,
            app_version=self.version,
            api_id=self.api_id,
            api_hash=self.api_hash,
            bot_token=self.bot_token,
            parse_mode=ParseMode.HTML,
            workers=self.workers,
            ipv6=self.ipv6,
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
            chat_id=self.logger,
            text=text.format(
                version=self.version,
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
        text_ = "<b>PyMouse stopped with {version}!</b>\n\n"
        text_ += "<b>System</b>: <code>{sys}</code>"
        await self.send_message(
            chat_id=Config.GP_LOGS,
            text=text_.format(
                version=self.version,
                sys=self.system_version
            )
        )
        await super().stop()
        log.info("WhiterRobot is stopped!")