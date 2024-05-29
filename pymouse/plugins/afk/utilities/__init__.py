from datetime import datetime

from hydrogram.types import Message, MessageEntity, User
from hydrogram.enums import MessageEntityType
from hydrogram.errors import UsernameInvalid, PeerIdInvalid

from pymouse import PyMouse, afkmodel_db
from pymouse.utils import UtilsTimer

class AFKUtils:
    @staticmethod
    async def get_user(entities: MessageEntity, c: PyMouse, m: Message) -> User | None: # type: ignore
        entoffset: int = entities.offset
        entlength: int = entities.length
        user = m.text[entoffset : entoffset + entlength]
        try:
            ent = await c.get_users(user)
        except (UsernameInvalid, PeerIdInvalid):
            return None
        return ent
    

    async def getMentioned(self, c: PyMouse, m: Message) -> User | None: # type: ignore
        u: User | None = None
        if m.entities:
            for y in m.entities:
                if y.type == MessageEntityType.MENTION:
                    u = await self.get_user(y, c, m)
                elif y.type == MessageEntityType.TEXT_MENTION:
                    u = y.user
        return u
    
    @staticmethod
    def getReplied(m: Message) -> User | None:
        replied = m.reply_to_message
        if replied and replied.from_user:
            return replied.from_user
        else:
            return None
    
    @staticmethod
    async def sender_afk(user: User, m: Message, i18n: dict) -> None:
        afktxt = ""
        gr: str | None = None
        gt: str | None = None
        afk = afkmodel_db.afk_db.getAFK(user.id)
        if afk.get("is_afk", False) != False:
            afktxt += i18n["afk"]["is-unavalaible"].format(user=user.mention)
            gr = afk.get("reason", None)
            gt = afk.get("time", None)
            if gr != None:
                afktxt += i18n["generic-strings"]["reason"].format(reason=gr)
            if afk.get("time", None) != None:
                afktxt += i18n["generic-strings"]["last-seen"].format(formattedtime=UtilsTimer().time_formatter(datetime.now().timestamp() - gt))
            await m.reply(afktxt)
        else:
            return None

    async def check_afk(self, c: PyMouse, m: Message, i18n: dict): # type: ignore
        user = await self.getMentioned(c, m)
        if user is not None:
            await self.sender_afk(user, m, i18n)
        user = self.getReplied(m)
        if user is not None:
            await self.sender_afk(user, m, i18n)
        return

    @staticmethod
    async def stop_afk(m: Message, i18n: dict):
        user = m.from_user
        if not user:
            return
        afkmodel_db.afk_db.unsetAFK(user.id)
        return await m.reply(i18n["afk"]["is-back"].format(user=user.mention))

afk_utils = AFKUtils()