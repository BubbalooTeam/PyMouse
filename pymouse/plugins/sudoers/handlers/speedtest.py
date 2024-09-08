#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.

from hydrogram import filters
from hydrogram.types import Message, InputMediaPhoto

from pymouse import PyMouse, Decorators, log, router
from pymouse.utils import NetworkUtils, NetworkEvents

@router.message(filters.command("speedtest"))
@Decorators().require_dev()
async def speed_test(c: PyMouse, m: Message): # type: ignore
    msg = await m.reply_photo(photo=NetworkEvents.RUNNING_SPEEDTEST, caption="<b>Running SpeedTest...</b>")
    try:
        # Get infos with speedtest and unpacking infos
        dl, ul, name, host, ping, isp, country, cc, path = await NetworkUtils().speedtest_performer()
        await msg.edit_media(
            media=InputMediaPhoto(
                media=path,
                caption=f"🌀 <b>Name:</b> <code>{name}</code>\n🌐 <b>Host:</b> <code>{host}</code>\n🏁 <b>Country:</b> <code>{country}, {cc}</code>\n\n🏓 <b>Ping:</b> <code>{ping} ms</code>\n🔽 <b>Download:</b> <code>{dl} Mbps</code>\n🔼 <b>Upload:</b> <code>{ul} Mbps</code>\n🖥  <b>ISP:</b> <code>{isp}</code>"
            )
        )
    except Exception:
        log.error("[speedtest/handler]: Error in performing speedtest...")