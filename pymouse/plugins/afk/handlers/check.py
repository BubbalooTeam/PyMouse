from hydrogram import filters
from hydrogram.types import Message

from pymouse import PyMouse, Decorators, afkmodel_db, router
from ..utilities.afk import afk_utils

@router.message(~filters.private & ~filters.bot & filters.all, group=2)
@Decorators().CatchError()
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