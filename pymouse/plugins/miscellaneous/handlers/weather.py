#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.

from hydrogram import filters
from hydrogram.types import Message

from pymouse import PyMouse, Decorators, router
from ..utilities.weather import HandleWeather

@router.message(filters.command("weather"))
@Decorators().CatchError()
@Decorators().Locale()
async def weatherHandle(_, m: Message, i18n):
    # Run Weather Information's with Handle
    await HandleWeather(
        c=PyMouse,
        m=m,
        i18n=i18n
    )