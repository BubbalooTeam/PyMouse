from typing import Union

from hydrogram import filters
from hydrogram.types import Message, CallbackQuery

from pymouse import PyMouse, Decorators, router
from ..utilities.gsmarena import HandleGSMarena

@router.message(filters.command(["d", "specs"]))
@router.callback(filters.regex(r"device\|(.*)$"))
@Decorators().CatchError()
@Decorators().Locale()
async def GSMarenaHandle(c: PyMouse, union: Union[Message, CallbackQuery], i18n): # type: ignore
     # Run GSMarena with Handle
    await HandleGSMarena(
        c=PyMouse,
        union=union,
        i18n=i18n,
    )