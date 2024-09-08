#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.


import json
import os
from pathlib import Path
from typing import Union, List
from dataclasses import dataclass

from hydrogram.types import Message, CallbackQuery, InlineQuery

from ..logger import log
from ...database.plugins.utilities.localization import localizationmodel_db

@dataclass(frozen=True, slots=True)
class LocalizationStats:
    """
    Represents a Localization Statistics:
        * Get informations with `number/percentage of strings translated/untranslated`
    """
    total_strings: int
    """
    The total number of strings.

    `type:` int
    """

    strings_translated: int
    """
    The number of translated strings.

    `type:` int
    """

    strings_untranslated: int
    """
    The number of untranslated strings.

    `type:` int
    """

    percentage_translated: int
    """
    The percentage of translated strings.

    `type:` int
    """

class Localization:
    def __init__(self):
        self.strings_path = "localization/"
        self.current_locales: List[str] = [
            "pt_br", # Brazilian Portuguese (BRA)
            "en_us", # American English (USA)
        ]
        self.default_language: str = "en_us" # American English (USA) is a default language
        self.strings = {}

    def get_all_files(self):
        """Get all files from strings."""
        path = Path(self.strings_path)
        return [i.absolute() for i in path.glob("**/*")]

    def compile_locales(self):
        """Compile bot languages."""
        log.info("Compiling Localization...")
        all_files = self.get_all_files()
        for file in all_files:
            with open(file=file, mode="r", encoding="utf-8") as file:
                language = os.path.basename(file.name).replace(".json", "")
                jsonloader: dict = json.load(file)
                self.strings[language] = jsonloader
        log.info("Localization compiled successfully!")

    def get_localization_of_chat(self, union: Union[Message, CallbackQuery, InlineQuery]):
        """
        Get chat Localization.

        Arguments:
            `union (Union[Message, CallbackQuery, InlineQuery])`: The received message or query.

        Returns:
            `language(str)`: Representing Chat Language.
        """
        locale = localizationmodel_db.localization_db.get_chat_language(union)
        if locale is None:
            return self.default_language
        else:
            return locale if locale in self.current_locales else self.default_language

    @staticmethod
    def switchLanguage(
        union: Union[Message, CallbackQuery],
        language: str,
    ):
        """
        Switch chat localization/language.

        Arguments:
            `union (Union[Message, CallbackQuery, InlineQuery])`: The received message or query.
            `language(str)`: The switched chat language.
        """
        localizationmodel_db.localization_db.set_chat_language(
            union=union,
            language=language
        )


    def get_statistics(self, language: str):
        """
        Gets string translation statistics for the specified language.

        Arguments:
            `language(str)`: Specified Language.

        Returns:
            `LocalizationStats`: Statistics for the specified language.
        """
        default_language = self.strings.get(self.default_language, {})
        requested_language = self.strings.get(language, {})

        def recursive_count_strings(default_value, requested_value):
            if isinstance(requested_value, str):
                nonlocal total_strings, translated_strings
                total_strings += 1
                if language != self.default_language:
                    if requested_value != default_value:
                        translated_strings += 1
                else:
                    if requested_value == default_value:
                        translated_strings += 1
            elif isinstance(requested_value, dict):
                for subkey in requested_value:
                    if subkey in default_value:
                        recursive_count_strings(default_value[subkey], requested_value[subkey])

        total_strings = 0
        translated_strings = 0

        for key in default_language:
            if key in requested_language:
                default_value = default_language[key]
                requested_value = requested_language[key]
                recursive_count_strings(default_value, requested_value)

        if total_strings > 0:
            translation_percentage = (translated_strings / total_strings) * 100
        else:
            translation_percentage = 0

        untranslated_strings = total_strings - translated_strings

        return LocalizationStats(
            total_strings=total_strings,
            strings_translated=translated_strings,
            strings_untranslated=untranslated_strings,
            percentage_translated=round(translation_percentage),
        )


localization = Localization()