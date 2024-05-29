from datetime import datetime

from hydrogram import filters
from hydrogram.types import Message

from pymouse import PyMouse, afkmodel_db, Decorators
from pymouse.utils import HandleText
from ..utilities import afk_utils

class AFK_Plugins:
    @staticmethod
    @Decorators.Locale
    async def setupAFK(c: PyMouse, m: Message, i18n): # type: ignore
        user = m.from_user
        avreason = True
        if not user:
            return

        is_afk = afkmodel_db.afk_db.getAFK(user.id).get("is_afk", False)
        if is_afk:
            await afk_utils.stop_afk(m)
            return
        # // Get informations for Toggling AFK to ON
        reason = HandleText().input_str(m)
        if not reason:
            reason = None
            avreason = False

        time = datetime.now().timestamp()
        # // Setting AFK in DataBase
        afkmodel_db.afk_db.setAFK(user.id, time, reason)
        afktext = i18n["afk"]["is-now-unavalaible"].format(user=user.mention)
        if avreason:
            afktext += i18n["generic-strings"]["reason"].format(reason=reason)
        await m.reply(afktext)

    @staticmethod
    @PyMouse.on_message(~filters.private & ~filters.bot & filters.all, group=2)
    @Decorators().Locale()
    async def handleAFK(c: PyMouse, m: Message, i18n): # type: ignore
        user = m.from_user
        if not user:
            return
        # check if text is AFK command
        try:
            if m.text:
                if m.text.startswith(("brb", "/afk", "!afk")):
                    return
        except AttributeError:
            return
        
        is_afk = afkmodel_db.afk_db.getAFK(user.id).get("is_afk", False)
        if is_afk:
            await afk_utils.stop_afk(m, i18n)
            return
        
        await afk_utils.check_afk(c, m, i18n)