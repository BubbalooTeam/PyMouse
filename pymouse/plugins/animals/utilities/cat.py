#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.

from dataclasses import dataclass

from hydrogram.types import Message
from hydrogram.enums import ChatAction

from pymouse import PyMouse
from pymouse.utils import http

@dataclass(frozen=True, slots=True)
class TheCatImgResponse:
    image_url: str | None

class CatUtils:
    @staticmethod
    async def getCatImage():
        r = await http.get("https://api.thecatapi.com/v1/images/search")
        r = r.json()
        return TheCatImgResponse(
            image_url=r[0].get("url")
        )

    async def HandleCat(self, c: PyMouse, m: Message): # type: ignore
        r = await self.getCatImage()
        img = r.image_url
        if img is not None:
            return await m.reply_photo(
                photo=img,
                caption="<b>Meow!</b>"
            )
        else:
            return await m.reply(
                "<b>Meoow!</b>"
            )

cat_utils = CatUtils()