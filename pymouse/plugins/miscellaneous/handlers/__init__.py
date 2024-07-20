from hydrogram.types import Message

from pymouse import Decorators
from pymouse.plugins.miscellaneous.utilities import HandleWeather

class Misccellaneous:
    @staticmethod
    @Decorators().Locale()
    async def weatherHandle(_, m: Message, i18n):
        # Run Weather Information's with Handle
        await HandleWeather(
            m=m,
            i18n=i18n
        )