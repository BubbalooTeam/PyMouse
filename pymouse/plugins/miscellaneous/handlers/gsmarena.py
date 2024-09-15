from hydrogram import filters
from hydrogram.types import Message

from pymouse import Decorators, router
from ..utilities.gsmarena import HandleGSMarena

@router.message(filters.command(["d", "specs"]))
@Decorators().Locale()
async def GSMarenaHandle(_, m: Message, i18n):
     # Run GSMarena with Handle
    await HandleGSMarena(
        m=m,
        i18n=i18n,
    )