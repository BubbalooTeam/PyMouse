from typing import Union

from hydrogram import filters
from hydrogram.filters import Message, CallbackQuery
from hydrokeyboard import InlineKeyboard, InlineButton

from pymouse import PyMouse, Decorators, router

@router.message(filters.command("help"))
@router.callback(filters.regex(r"HelpMenu$"))
@Decorators().Locale()
async def HelpMenu(c: PyMouse, union: Union[Message, CallbackQuery], i18n): # type: ignore
    HelpText = i18n["pm-menu"]["help-below"].format(
        bot=c.me.first_name,
    )
    keyboard = InlineKeyboard(row_width=1)
    keyboard.add(
        InlineButton(
            text=i18n["buttons"]["help"],
            url="book.bubbalooteam.me/pymouse/help"
        )
    )
    if isinstance(union, Message):
        await union.reply(
            text=HelpText,
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
            text=HelpText,
            reply_markup=keyboard,
        )
