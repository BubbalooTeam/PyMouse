import json
import os
from pathlib import Path
from typing import Union, List

from hydrogram.types import Message, CallbackQuery, InlineQuery

from ..logger import log
from ...database.plugins.utilities.localization import localizationmodel_db

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
        locale = localizationmodel_db.localization_db.get_chat_language(union)
        if locale is None:
            return self.default_language
        else:
            return locale if locale in self.current_locales else self.default_language

localization = Localization()