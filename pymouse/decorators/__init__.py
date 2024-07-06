#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.


from typing import Union, Optional
from functools import wraps
from hydrogram.types import Message, CallbackQuery, InlineQuery
from hydrogram.enums import ChatType

from pymouse import PyMouse, Config, db, usersmodel_db, chatsmodel_db, localization

class Decorators:
    def __init__(self):
        self.owner = Config.OWNER_ID

    def require_dev(self, only_owner: Optional[bool] = False):
        def decorator(func):
            @wraps(func)
            async def wrapper(c: PyMouse, m: Message, *args, **kwargs): # type: ignore
                dev_db = db.GetCollection("sudoers")
                if not m.from_user:
                    return None
                user_id = m.from_user.id
                if user_id == self.owner:
                    return await func(c, m, *args, **kwargs)
                get_dev_by_id = dev_db.find_one({"dev_users": user_id})
                if get_dev_by_id and only_owner is False:
                    return await func(c, m, *args, **kwargs)
                else:
                    return None
            return wrapper
        return decorator
    
    def SaveUsers(self):
        def decorator(func):
            @wraps(func)
            async def wrapper(c: PyMouse, m: Message, *args, **kwargs): # type: ignore
                if not m.from_user:
                    return await func(c, m, *args, **kwargs)
                user_id = m.from_user.id
                actual_username = m.from_user.username
                actual_uname = m.from_user.first_name + " " + m.from_user.last_name if not m.from_user.last_name is None else m.from_user.first_name
                actual_languagebytg = m.from_user.language_code
                usersmodel_db.users_db.update_user(user_id, actual_uname, actual_username, actual_languagebytg)
                return await func(c, m, *args, **kwargs)
            return wrapper
        return decorator
    
    def SaveChats(self):
        def decorator(func):
            @wraps(func)
            async def wrapper(c: PyMouse, m: Message, *args, **kwargs): # type: ignore
                if m.chat.type == ChatType.PRIVATE:
                    return await func(c, m, *args, **kwargs)
                chat_id = m.chat.id
                actual_chatname = m.chat.username
                actual_chattitle = m.chat.title
                chatsmodel_db.chats_db.update_chat(chat_id, actual_chattitle, actual_chatname)
                return await func(c, m, *args, **kwargs)
            return wrapper
        return decorator
    
    def Locale(self):
        """Get strings from chat localization"""
        def decorator(func):
            @wraps(func)
            async def wrapper(c: PyMouse, union: Union[Message, CallbackQuery, InlineQuery], *args, **kwargs): # type: ignore
                language = localization.get_localization_of_chat(union)
                i18n = localization.strings.get(language, {})
                return await func(c, union, *args, i18n, **kwargs)
            return wrapper
        return decorator