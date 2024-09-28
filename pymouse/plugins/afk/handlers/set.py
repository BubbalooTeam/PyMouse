#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.

from datetime import datetime

from hydrogram import filters
from hydrogram.types import Message
from hydrogram.enums import ChatAction

from pymouse import PyMouse, Decorators, afkmodel_db, router
from pymouse.utils import HandleText
from ..utilities.afk import afk_utils

@router.message(filters.command("afk") | filters.regex(r"^(?i:brb)(\s(.+))?"))
@Decorators().CatchError()
@Decorators().Locale()
async def setupAFK(c: PyMouse, m: Message, i18n): # type: ignore
    user = m.from_user
    avreason = True
    if not user:
        return

    is_afk = afkmodel_db.afk_db.getAFK(user.id).get("is_afk", False)
    if is_afk:
        await afk_utils.stop_afk(c, m, i18n)
        return
    # // Get informations for Toggling AFK to ON
    reason = HandleText().input_str(m)
    if not reason:
        reason = None
        avreason = False

    time = datetime.now().timestamp()
    # // Setting AFK in DataBase
    afkmodel_db.afk_db.setAFK(user.id, time, reason)

    await PyMouse.send_chat_action(
        chat_id=m.chat.id,
        action=ChatAction.TYPING,
    )
    afktext = i18n["afk"]["is-now-unavalaible"].format(user=user.mention)
    if avreason:
        afktext += i18n["generic-strings"]["reason"].format(reason=reason)
    await m.reply(afktext)