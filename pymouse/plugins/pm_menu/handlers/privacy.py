#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.

from io import BytesIO
from typing import Union

from hydrogram import filters
from hydrogram.types import Message, CallbackQuery
from hydrokeyboard import InlineKeyboard, InlineButton

from pymouse import PyMouse, Decorators, usersmodel_db, router

@router.message(filters.command("privacy"))
@router.callback(filters.regex(r"^PrivacyPolicy$"))
@Decorators().Locale()
async def privacyPolicy(c: PyMouse, union: Union[Message, CallbackQuery], i18n): # type: ignore
    privacypolicyText = i18n["pm-menu"]["privacy-policy"].format(
        bot=c.me.first_name,
    )
    keyboard = InlineKeyboard(row_width=1)
    keyboard.add(
        InlineButton(
            text=i18n["buttons"]["privacy-data"],
            callback_data="PrivacyData",
        ),
        InlineButton(
            text=i18n["buttons"]["back"],
            callback_data="StartBack"
        )
    )
    if isinstance(union, Message):
        await union.reply(
            text=privacypolicyText,
            reply_markup=keyboard,
        )
    elif isinstance(union, CallbackQuery):
        await union.edit_message_text(
            text=privacypolicyText,
            reply_markup=keyboard,
        )

@router.callback(filters.regex("^PrivacyData$"))
@Decorators().Locale()
async def privacyPolicyRead(c: PyMouse, cb: CallbackQuery, i18n): # type: ignore
    privacypolicyReadText = i18n["pm-menu"]["privacy-policyRead"].format(
        bot=c.me.first_name,
    )
    keyboard = InlineKeyboard(row_width=2)
    keyboard.add(
        InlineButton(
            text=i18n["buttons"]["open-privacy-policy"],
            url="https://book.bubbalooteam.me/pymouse/privacy-policy",
        ),
        InlineButton(
            text=i18n["buttons"]["your-data"],
            callback_data="YourDataCollected",
        ),
    )
    await cb.edit_message_text(
        text=privacypolicyReadText,
        reply_markup=keyboard,
    )

@router.callback(filters.regex(r"^YourDataCollected$"))
@Decorators().Locale()
async def ReadYourData(c: PyMouse, cb: CallbackQuery, i18n): # type: ignore
    user_id = cb.from_user.id
    # GET data in the DataBase.
    data = usersmodel_db.users_db.getuser_dict(
        user_id=user_id,
    )
    # Make document with informations.
    file = BytesIO(str.encode(str(data)))
    file.name = "{user_id}_dataInfos.txt".format(
        user_id=user_id,
    )
    await cb.edit_message_text("Sending your data...")
    await c.send_document(
        chat_id=user_id,
        document=file,
    )
    await cb.edit_message_text("Your details have been sent to you, check your messages in PyMouse profile.")
