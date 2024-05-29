from typing import Union
from hydrogram.types import Message, CallbackQuery
from hydrogram.enums import ChatType
from hydrokeyboard import InlineKeyboard, InlineButton

from pymouse import PyMouse, Decorators, localization

class PMMenu_Plugins:
    @staticmethod
    @Decorators().Locale()
    async def start_(c: PyMouse, union: Union[Message, CallbackQuery], i18n): # type: ignore
        msg = union.message if isinstance(union, CallbackQuery) else union
        keyboard = InlineKeyboard()

        # Row button's for start message
        keyboard.row(
            InlineButton(
                text=i18n["buttons"]["language"],
                callback_data="LangConfigOpts.start_back",
            ),
            InlineButton(
                text=i18n["buttons"]["help"],
                callback_data="pm_menu:HelpMenu",
            ),
        )
        # Second row button's for start message
        keyboard.row(
            InlineButton(
                text=i18n["buttons"]["about"],
                callback_data="AboutMenu",
            )
        )

        start_text = i18n["pm-menu"]["start-private"].format(user=union.from_user.mention, bot=c.me.mention)
        if isinstance(union, Message):
            if union.chat.type != ChatType.PRIVATE:
                start_text = i18n["pm-menu"]["start-group"].format(bot=c.me.mention)
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

    async def glang_info(c: PyMouse, m: Message): # type: ignore
        chat_language = localization.get_localization_of_chat(m)
        lang_info = localization.get_statistics(chat_language)
        text = "<b>Language Info:</b>\n\n<b>Total of Strings:</b> <code>{total_strings}</code>\n<b>Translated Strings:</b> <code>{translated_strings}</code>\n<b>Untranslated Strings:</b> <code>{untranslated_strings}</code>".format(total_strings=lang_info.total_strings, translated_strings=lang_info.strings_translated, untranslated_strings=lang_info.strings_untranslated)
        text += "\n\nOh! All strings already translated!" if lang_info.percentage_translated >= 100 else "\n\nThere are still translations missing! Currently this bot is <code>{strings_percentage} % </code> translated.".format(strings_percentage=lang_info.percentage_translated)
        await m.reply(text)