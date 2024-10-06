#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.

import sentry_sdk

from asyncio import get_event_loop
from typing import Union, Optional, Callable
from functools import partial, wraps
from re import sub as subs

from hydrogram import StopPropagation, StopTransmission
from hydrogram.types import Message, CallbackQuery, InlineQuery, ChatPrivileges
from hydrogram.enums import ChatType, ChatMemberStatus
from hydrogram.errors import ChatWriteForbidden

from pymouse import PyMouse, Config, db, usersmodel_db, chatsmodel_db, localization, FMT

class Decorators:
    def __init__(self):
        self.owner = Config.OWNER_ID

    def aiowrap(self):
        """
        Runs a synchronous function in the background, so that it does not create bugs in an asynchronous function
        """
        def decorator(func: Callable) -> Callable:
            @wraps(func)
            async def wrapper(*args, loop=None, executor=None, **kwargs):
                if loop is None:
                    loop = get_event_loop()
                pfunc = partial(func, *args, **kwargs)
                return await loop.run_in_executor(executor, pfunc)
            return wrapper
        return decorator

    def require_dev(self, only_owner: Optional[bool] = False):
        "Check if user, is a developer or owner of PyMouse."
        def decorator(func) -> Callable:
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
        "Save and Update a user in DataBase."
        def decorator(func) -> Callable:
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
        "Save and Update a chat in DataBase."
        def decorator(func) -> Callable:
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

    def CheckAdminRight(
        self,
        permissions: ChatPrivileges | None = None,
        accept_in_private: bool = False,
        complain_missing_permissions: bool = True
    ):
        "Check if the user is Administrator of the chat."
        def decorator(func: Callable) -> Callable:
            @wraps(func)
            async def wrapper(c: PyMouse, union: Union[Message, CallbackQuery], *args, **kwargs): # type: ignore
                patternSender = r"<(b|code|i|a.*?>)(.*?)</\1>"
                language = localization.get_localization_of_chat(union)
                i18n = localization.strings.get(language, {})
                m = union if isinstance(union, Message) else union.message
                MessageSender = union.reply_text if isinstance(union, Message) else partial(union.answer, show_alert=True)

                if m.chat.type == ChatType.PRIVATE:
                    if accept_in_private:
                        return await func(c, union, *args, **kwargs)
                    SenderText = i18n["enums"]["chatType"]["OnlyGroups"]
                    if isinstance(union, CallbackQuery):
                        SenderText = subs(patternSender, r"\2", SenderText)
                    return await MessageSender(SenderText)
                elif m.chat.type == ChatType.CHANNEL:
                    return await func(c, union, *args, **kwargs)

                user = await m.chat.get_member(union.from_user.id)
                if (
                    user.status == ChatMemberStatus.OWNER
                    or not permissions
                    and user.status == ChatMemberStatus.ADMINISTRATOR
                ):
                    return await func(c, union, *args, **kwargs)
                if user.status != ChatMemberStatus.ADMINISTRATOR:
                    if complain_missing_permissions:
                        SenderText = i18n["generic-strings"]["not-admin"]
                        if isinstance(union, CallbackQuery):
                            SenderText = subs(patternSender, r"\2", SenderText)
                        return await MessageSender(SenderText)
                    return None

                missing_permissions = [
                    permissions
                    for permissions, value in permissions.__dict__.items()
                    if value and not getattr(user.privileges, permissions)
                ]

                if not missing_permissions:
                    return await func(c, union, *args, **kwargs)

                if complain_missing_permissions:
                    SenderText = i18n["generic-strings"]["admin-need-permissions"].format(
                        permissions=FMT().FormatListtoString(
                            RelativeList=list(missing_permissions),
                            i18n=i18n
                        )
                    )
                    if isinstance(union, CallbackQuery):
                        SenderText = subs(patternSender, r"\2", SenderText)
                    return await MessageSender(SenderText)
            return wrapper
        return decorator

    def CatchError(self):
        """Captures the Error, Creates a Report and Sends it to the User."""
        def decorator(func) -> Callable:
            @wraps(func)
            async def wrapper(c: PyMouse, union: Union[Message, CallbackQuery], *args, **kwargs): # type: ignore
                language = localization.get_localization_of_chat(union)
                i18n = localization.strings.get(language, {})
                m = union if isinstance(union, Message) else union.message
                patternSender = r"<(b|code|i|a.*?>)(.*?)</\1>"
                MessageSender = union.reply_text if isinstance(union, Message) else partial(union.answer, show_alert=True)

                try:
                    return await func(c, union, *args, **kwargs)
                except ChatWriteForbidden:
                    return await PyMouse.leave_chat(
                        chat_id=m.chat.id
                    )
                except (StopTransmission, StopPropagation):
                    pass
                except Exception as e:
                    sentry_sdk.capture_exception(e)
                    SenderText = i18n["generic-strings"]["error-occurred"].format(
                        error=f"{e.__class__.__name__}: {e}\n-> Line: {e.__traceback__.tb_lineno}\n-> File: {e.__traceback__.tb_frame.f_code.co_filename}"
                    )
                    if isinstance(union, CallbackQuery):
                        SenderText = subs(patternSender, r"\2", SenderText)
                    return await MessageSender(SenderText)
            return wrapper
        return decorator

    def Locale(self):
        """Get strings from chat localization"""
        def decorator(func) -> Callable:
            @wraps(func)
            async def wrapper(c: PyMouse, union: Union[Message, CallbackQuery, InlineQuery], *args, **kwargs): # type: ignore
                language = localization.get_localization_of_chat(union)
                i18n = localization.strings.get(language, {})
                return await func(c, union, *args, i18n, **kwargs)
            return wrapper
        return decorator