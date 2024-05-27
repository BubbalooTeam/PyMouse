from hydrogram.types import Message

from pymouse import PyMouse, Decorators

class PMMenu_Plugins:
    @staticmethod
    @Decorators().Locale()
    async def start(c: PyMouse, m: Message, i18n): # type: ignore
        await m.reply(i18n["pm-menu"]["start-private"].format(user=m.from_user.mention, bot=c.me.mention))