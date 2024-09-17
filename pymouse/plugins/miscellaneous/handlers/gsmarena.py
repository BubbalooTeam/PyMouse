from typing import Union

from hydrogram import filters
from hydrogram.types import Message, CallbackQuery

from pymouse import Decorators, router
from ..utilities.gsmarena import HandleGSMarena

@router.message(filters.command(["d", "specs"]))
@router.callback(filters.regex(r"device\|(.*)$"))
@Decorators().Locale()
async def GSMarenaHandle(_, union: Union[Message, CallbackQuery], i18n):
     # Run GSMarena with Handle
    await HandleGSMarena(
        union=union,
        i18n=i18n,
    )