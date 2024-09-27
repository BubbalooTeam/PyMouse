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
from hydrogram.types import Message, CallbackQuery, ChatPrivileges
from hydrogram.enums import ChatType
from hydrokeyboard import InlineKeyboard, InlineButton

from pymouse import PyMouse, Decorators, usersmodel_db, chatsmodel_db, router

@router.message(filters.command("privacy"))
@router.callback(filters.regex(r"^PrivacyPolicy$"))
@Decorators().CatchError()
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
    )
    if isinstance(union, Message):
        await union.reply(
            text=privacypolicyText,
            reply_markup=keyboard,
        )
    elif isinstance(union, CallbackQuery):
        keyboard.row(
            InlineButton(
                text=i18n["buttons"]["back"],
                callback_data="StartBack"
            )
        )
        await union.edit_message_text(
            text=privacypolicyText,
            reply_markup=keyboard,
        )

@router.callback(filters.regex("^PrivacyData$"))
@Decorators().CatchError()
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
@Decorators().CatchError()
@Decorators().CheckAdminRight(
    permissions=ChatPrivileges(can_change_info=True),
    accept_in_private=True
)
async def ReadYourData(c: PyMouse, cb: CallbackQuery): # type: ignore
    data = {}
    TGid=cb.from_user.id if cb.message.chat.type == ChatType.PRIVATE else cb.message.chat.id
    # GET data in the DataBase.
    if cb.message.chat.type == ChatType.PRIVATE:
        data = usersmodel_db.users_db.getuser_dict(
            user_id=TGid,
        )
    else:
        data = chatsmodel_db.chats_db.get_chat_dict(
            chat_id=TGid,
        )
    # Make document with informations.
    file = BytesIO(str.encode(str(data)))
    file.name = "{TGid}_dataInfos.txt".format(
        TGid=TGid,
    )
    await cb.edit_message_text("Sending chat data...")
    await c.send_document(
        chat_id=cb.from_user.id,
        document=file,
    )
    await cb.edit_message_text("Your details have been sent to you, check your messages in your chat with me.")
