from typing import Union
from hydrogram.types import Message, CallbackQuery
from hydrogram.enums import ChatType
from hydrokeyboard import InlineKeyboard, InlineButton

from pymouse import PyMouse, Decorators

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
                callback_data="HelpMenu",
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