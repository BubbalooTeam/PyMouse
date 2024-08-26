#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.

from asyncio import sleep
from dataclasses import dataclass
from typing import Union, Optional, Pattern

from hydrogram.types import Message, CallbackQuery
from hydrokeyboard import InlineKeyboard, InlineButton

from pymouse import PyMouse, Decorators, localization

@dataclass
class LocalizationArgs:
    text: str
    buttons: InlineKeyboard

class LocalizationInfo:
    def get_changelang_buttons(
        self,
        i18n: dict,
        lang_callback: Optional[Union[str, Pattern]] = None,
        back_callback: Optional[Union[str, Pattern]] = None,
    ) -> InlineKeyboard:
        keyboard = InlineKeyboard(row_width=1)
        # === #
        keyboard.add(
            InlineButton(
                text=i18n["buttons"]["change-language"],
                callback_data=lang_callback,
            ),
            InlineButton(
                text=i18n["buttons"]["back"],
                callback_data=back_callback,
            ),
        )
        return keyboard


    async def get_changelang_text_and_buttons(
        self,
        c: PyMouse,  # type: ignore
        union: Union[Message, CallbackQuery],
        i18n: dict,
        lang_callback: Optional[Union[str, Pattern]] = None,
        back_callback: Optional[Union[str, Pattern]] = None,
    ):
        i18n_language = i18n["language"]
        chat_language = localization.get_localization_of_chat(union)
        msg_text = i18n_language["language-info"]["chat-language"].format(
            language_flag=i18n_language["flag"],
            language=i18n_language["name"]
        )
        keyboard = InlineKeyboard(row_width=1)
        # //
        statistics = localization.get_statistics(chat_language)
        # //
        msg_text += i18n_language["language-info"]["language-info"]
        msg_text += i18n_language["language-info"]["language-total-strings"].format(total_strings=statistics.total_strings)
        msg_text += i18n_language["language-info"]["language-translated"].format(translated_strings=statistics.strings_translated)
        msg_text += i18n_language["language-info"]["language-untranslated"].format(untranslated_strings=statistics.strings_untranslated)
        msg_text += "\n\n"
        percentage = statistics.percentage_translated
        if percentage == 100:
            msg_text += i18n_language["language-info"]["native-language"]
        else:
            msg_text += i18n_language["language-info"]["missing-translations"].format(strings_percentage=percentage)

        keyboard = self.get_changelang_buttons(
            i18n=i18n,
            lang_callback=lang_callback,
            back_callback=back_callback
        )

        return LocalizationArgs(
            text=msg_text,
            buttons=keyboard
        )

    def get_switchlang_buttons(
        self,
        i18n: dict,
        changemenu_back: Optional[Union[str, Pattern]] = None,
        back_callback: Optional[Union[str, Pattern]] = None,
    ):
        keyboard = InlineKeyboard()
        btns = list([])
        avl_languages = localization.current_locales
        for lng_temp in avl_languages:
            i18n_temp = localization.strings.get(lng_temp)
            btns.append(
                InlineButton(
                    text=str(i18n_temp["language"]["flag"] + i18n_temp["language"]["name"]),
                    callback_data="SwitchLang|{language}|{changemenu_back}".format(
                        language=lng_temp,
                        changemenu_back=changemenu_back
                    ),
                )
            )
        keyboard.add(*btns)
        keyboard.row(
            InlineButton(
                text=i18n["buttons"]["back"],
                callback_data=back_callback
            )
        )
        return keyboard

    async def get_switchlang_text_and_buttons(
        self,
        i18n: dict,
        changemenu_back: Optional[Union[str, Pattern]] = None,
        back_callback: Optional[Union[str, Pattern]] = None,
    ):
        i18n_language = i18n["language"]

        msg_text = i18n_language["switch-lang"]
        keyboard = self.get_switchlang_buttons(
            i18n=i18n,
            changemenu_back=changemenu_back,
            back_callback=back_callback,
        )
        return LocalizationArgs(
            text=msg_text,
            buttons=keyboard,
        )

    @staticmethod
    async def switchLanguage(
        union: Union[Message, CallbackQuery],
        language: str,
        sleeper: int
    ):
        localization.switchLanguage(
            union=union,
            language=language,
        )
        await sleep(sleeper)

    @staticmethod
    @Decorators().Locale()
    async def send_switchedlang_text_and_buttons(
        _,
        cb: CallbackQuery,
        changemenu_back: Optional[Union[str, Pattern]] = None,
        i18n: Optional[Union[dict, str]] = None,
    ):
        keyboard = InlineKeyboard()
        keyboard.add(
            InlineButton(
                text=i18n["buttons"]["back"],
                callback_data=changemenu_back,
            )
        )
        msg_text = i18n["language"]["switched-lang"].format(
            language=str(i18n["language"]["flag"] + i18n["language"]["name"])
        )
        return await cb.edit_message_text(
            text=msg_text,
            reply_markup=keyboard,
        )