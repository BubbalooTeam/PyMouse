#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.

from hydrogram import filters
from hydrogram.types import CallbackQuery

from pymouse import PyMouse, Decorators, router
from pymouse.plugins.pm_menu.utilities.localization import LocalizationInfo

@router.callback(filters.regex(r"^LangMenu\|(.*)$"))
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

@router.callback(filters.regex(r"ChangeLanguage\|(.*)$"))
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

@router.callback(filters.regex(r"SwitchLang\|(.*)$"))
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