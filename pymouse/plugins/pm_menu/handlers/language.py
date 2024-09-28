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
from hydrogram.enums import ChatAction
from hydrogram.types import Message, CallbackQuery, ChatPrivileges

from pymouse import PyMouse, Decorators, router
from pymouse.plugins.pm_menu.utilities.localization import LocalizationInfo

@router.message(filters.command(["setlang", "lang"]))
@router.callback(filters.regex(r"^LangMenu\|(.*)$"))
@Decorators().CatchError()
@Decorators().CheckAdminRight(
    permissions=ChatPrivileges(can_change_info=True),
    accept_in_private=True
)
@Decorators().Locale()
async def ChangeLanguageMenu(c: PyMouse, union: Union[Message, CallbackQuery], i18n): # type: ignore
    chat = union.chat if isinstance(union, Message) else union.message.chat
    sender = union.edit_message_text if isinstance(union, CallbackQuery) else union.reply
    if isinstance(union, CallbackQuery):
        inf = union.data.split("|")
        changemenu_back = inf[1]
        changelang_back = inf[2]
    elif isinstance(union, Message):
        changemenu_back = "StartBack"
        changelang_back = "LangMenu"

    await c.send_chat_action(
        chat_id=chat.id,
        action=ChatAction.TYPING,
    )
    text_and_buttons = await LocalizationInfo().get_changelang_text_and_buttons(
        c=c,
        union=union,
        i18n=i18n,
        lang_callback="ChangeLanguage|{changemenu_back}|{changelang_back}".format(
            changemenu_back=changemenu_back,
            changelang_back=changelang_back
        ),
        back_callback=changemenu_back
    )
    await sender(
        text=text_and_buttons.text,
        reply_markup=text_and_buttons.buttons,
    )

@router.callback(filters.regex(r"ChangeLanguage\|(.*)$"))
@Decorators().CatchError()
@Decorators().CheckAdminRight(
    permissions=ChatPrivileges(can_change_info=True),
    accept_in_private=True
)
@Decorators().Locale()
async def SelectLanguageMenu(c: PyMouse, cb: CallbackQuery, i18n): # type: ignore
    inf = cb.data.split("|")
    changemenu_back = inf[1]
    changelang_back = inf[2]

    await c.send_chat_action(
        chat_id=cb.message.chat.id,
        action=ChatAction.TYPING,
    )
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

@router.callback(filters.regex(r"SwitchLang\|(.*)$"))
@Decorators().CatchError()
@Decorators().CheckAdminRight(
    permissions=ChatPrivileges(can_change_info=True),
    accept_in_private=True
)
@Decorators().Locale()
async def SwitchLanguage(c: PyMouse, cb: CallbackQuery, i18n): # type: ignore
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
    await LocalizationInfo().send_switchedlang_text_and_buttons(c, cb, changemenu_back)