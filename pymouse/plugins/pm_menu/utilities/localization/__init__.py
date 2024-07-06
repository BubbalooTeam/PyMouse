#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.


from typing import Union, Optional, Pattern

from hydrogram.types import Message, CallbackQuery
from hydrokeyboard import InlineKeyboard, InlineButton

from pymouse import PyMouse, Decorators, localization

class LocalizationInfo:
    def get_changelang_buttons(
        i18n: dict,
        lang_callback: Optional[Union[str, Pattern]] = None,
        back_callback: Optional[Union[str, Pattern]] = None,
    ):
        keyboard = InlineKeyboard(row_width=1)
        # === #
        keyboard.add(
            InlineButton(
                text=i18n["buttons"]["change-language"],
                
            )
        )
        
        
    async def get_text_and_buttons(
        self, 
        c: PyMouse,  # type: ignore
        union: Union[Message, CallbackQuery], 
        i18n: dict,
        clang_menu: Optional[Union[str, Pattern]] = None
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
            
        keyboard.row(InlineButton(
            text=i18n["buttons"]["change-language"],
            callback_data=clang_menu
        ))
        