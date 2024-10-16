#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.

from typing import Union

from hydrogram import filters
from hydrogram.types import Message, CallbackQuery

from pymouse import PyMouse, Decorators, router
from ..utilities.gsmarena import gsm_arena_utils

@router.message(filters.command(["d", "specs"]))
@router.callback(filters.regex(r"device\|(.*)$"))
@Decorators().CatchError()
@Decorators().Locale()
async def GSMarenaHandle(c: PyMouse, union: Union[Message, CallbackQuery], i18n): # type: ignore
     # Run GSMarena with Handle
    await gsm_arena_utils.HandleGSMarena(
        c=c,
        union=union,
        i18n=i18n,
    )