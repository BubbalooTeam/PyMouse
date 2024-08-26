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

from hydrogram.types import Message, CallbackQuery
from hydrogram.enums import ChatType
from hydrokeyboard import InlineKeyboard, InlineButton

from pymouse import PyMouse, Decorators, usersmodel_db
from pymouse.plugins.pm_menu.utilities.localization import LocalizationInfo

class PMMenu_Plugins:
    @staticmethod
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
            )
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
            await union.edit_message_text(
                text=start_text,
                reply_markup=keyboard,
            )

    @staticmethod
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

    @staticmethod
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

    @staticmethod
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

    @staticmethod
    @Decorators().Locale()
    async def ChangeLanguageMenu(c: PyMouse, cb: CallbackQuery, i18n): # type: ignore
        inf = cb.data.split("|")
        changemenu_back = inf[1]
        changelang_back = inf[2]

        text_and_buttons = await LocalizationInfo().get_changelang_text_and_buttons(
            c=c,
            union=cb,
            i18n=i18n,
            lang_callback="ChangeLanguage|{changemenu_back}|{changelang_back}".format(
                changemenu_back=changemenu_back,
                changelang_back=changelang_back
            ),
            back_callback=changemenu_back
        )
        await cb.edit_message_text(
            text=text_and_buttons.text,
            reply_markup=text_and_buttons.buttons,
        )

    @staticmethod
    @Decorators().Locale()
    async def SelectLanguageMenu(_, cb: CallbackQuery, i18n): # type: ignore
        inf = cb.data.split("|")
        changemenu_back = inf[1]
        changelang_back = inf[2]

        text_and_buttons = await LocalizationInfo().get_switchlang_text_and_buttons(
            i18n=i18n,
            changemenu_back=changemenu_back,
            back_callback="{changelang_back}|{changemenu_back}|LangMenu".format(
                changelang_back=changelang_back,
                changemenu_back=changemenu_back,
            ),
        )
        await cb.edit_message_text(
            text=text_and_buttons.text,
            reply_markup=text_and_buttons.buttons,
        )

    @staticmethod
    @Decorators().Locale()
    async def SwitchLanguage(_, cb: CallbackQuery, i18n):
        inf = cb.data.split("|")
        language = inf[1]
        changemenu_back = inf[2]

        # == #
        await LocalizationInfo().switchLanguage(
            union=cb,
            language=language,
            sleeper=0.5
        )

        # === #
        await LocalizationInfo().send_switchedlang_text_and_buttons(_, cb, changemenu_back)