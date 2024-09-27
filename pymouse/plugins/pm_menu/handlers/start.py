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
from hydrogram.enums import ChatType
from hydrokeyboard import InlineKeyboard, InlineButton

from pymouse import PyMouse, Decorators, router

@router.message(filters.command("start"))
@router.callback(filters.regex(r"^StartBack$"))
@Decorators().CatchError()
@Decorators().Locale()
async def start_(c: PyMouse, union: Union[Message, CallbackQuery], i18n): # type: ignore
    keyboard = InlineKeyboard()

    # Row button's for start message
    keyboard.row(
        InlineButton(
            text=i18n["buttons"]["language"],
            callback_data="LangMenu|StartBack|LangMenu",
        ),
        InlineButton(
            text=i18n["buttons"]["help"],
            callback_data="HelpMenu",
        ),
    )
    # Second row button's for start message
    keyboard.row(
        InlineButton(
            text=i18n["buttons"]["about"],
            callback_data="AboutMenu",
        ),
        InlineButton(
            text=i18n["buttons"]["privacy"],
            callback_data="PrivacyPolicy"
        ),
    )

    start_text = i18n["pm-menu"]["start-private"].format(user=union.from_user.mention, bot=c.me.mention)
    if isinstance(union, Message):
        if union.chat.type != ChatType.PRIVATE:
            start_text = i18n["pm-menu"]["start-group"].format(bot=c.me.mention)
            keyboard = None
        await c.send_message(
            chat_id=union.chat.id,
            text=start_text,
            reply_markup=keyboard
        )
    elif isinstance(union, CallbackQuery):
        if union.message.chat.type != ChatType.PRIVATE:
            start_text = i18n["pm-menu"]["start-group"].format(bot=c.me.mention)
            keyboard = None
        await union.edit_message_text(
            text=start_text,
            reply_markup=keyboard,
        )