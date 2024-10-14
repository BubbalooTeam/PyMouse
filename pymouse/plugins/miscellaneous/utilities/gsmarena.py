#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.

from re import sub as subs
from typing import Union

from hydrogram.types import Message, CallbackQuery
from hydrogram.enums import ChatAction
from hydrokeyboard import InlineKeyboard, InlineButton

from pymouse import PyMouse, GSMarena, GSMarenaDeviceBaseResult, GSMarenaSearchResults, log
from pymouse.utils import HandleText
from pymouse.utils.tools.gsm_arena.exceptions import GSMarenaDeviceNotFound

gsm_arena = GSMarena()

class GSMarenaUtils:
    @staticmethod
    def formatGSMarenaMessage(gsmarenaBaseResult: GSMarenaDeviceBaseResult, i18n: dict) -> str:
        phone = gsm_arena.parse_specifications(specifications=gsmarenaBaseResult.phone_details)
        deviceI18N = i18n["gsmarena"]["phoneFormatter"]

        attrs = [
            "<b>{spec_name}:</b> <i>{spec_value}</i>".format(
                spec_name=deviceI18N[key],
                spec_value=value
            )
            for key, value in phone.items()
            if value and value.strip() != "-" and deviceI18N.get(key) is not None
        ]


        formatted_message = (
            "<a href='{image}'>{space}</a>"
            "<a href='{url}'><b>{name}</b></a>\n\n{attrs}"
        ).format(
            image=gsmarenaBaseResult.image,
            space="\u2000",
            url=gsmarenaBaseResult.url,
            name=gsmarenaBaseResult.name,
            attrs="\n\n".join(attrs)
        )

        return formatted_message

    @staticmethod
    def GSMarenaCreateKeyboard(search: GSMarenaSearchResults, user_id: int):
        def formatDeviceName(name: str) -> str:
            # Insert spaces between lowercase letters and numbers
            name_with_spaces = subs(r'([a-z])([0-9])', r'\1 \2', name)
            # Insert spaces between lowercase and uppercase letters
            formatted_name = subs(r'([a-z])([A-Z])', r'\1 \2', name_with_spaces)

            # Capitalize each word
            final_name = formatted_name.title()

            return final_name
        keyboard = InlineKeyboard(row_width=2)
        buttons = []
        for device in search.results:
            buttons.append(
                InlineButton(
                    text=formatDeviceName(
                        name=device.name
                    ),
                    callback_data="device|{device_id}|{user_id}".format(
                        device_id=device.id,
                        user_id=user_id,
                    ),
                ),
            )
            if len(buttons) >= 30:
                break
        keyboard.add(*buttons)
        return keyboard


    async def HandleGSMarena(self, c: PyMouse, union: Union[Message, CallbackQuery], i18n: dict): # type: ignore
        chat = union.chat if isinstance(union,Message) else union.message.chat
        sender = union.reply if isinstance(union, Message) else union.edit_message_text

        async def GSMarenaBuildMessage(c: PyMouse, deviceID: str, i18n: dict): # type: ignore
            await c.send_chat_action(
                chat_id=chat.id,
                action=ChatAction.TYPING,
            )
            fetchDevice = await gsm_arena.fetch_device(device=deviceID)

            formatedMessage = self.formatGSMarenaMessage(gsmarenaBaseResult=fetchDevice, i18n=i18n)
            return await sender(
                text=formatedMessage,
                disable_web_page_preview=False,
            )

        if isinstance(union, Message):
            query = HandleText().input_str(union=union)
            if query:
                try:
                    getDevice = await gsm_arena.search_device(query=query)
                    if len(getDevice.results) >= 2:
                        keyboard = self.GSMarenaCreateKeyboard(search=getDevice, user_id=union.from_user.id)
                        await c.send_chat_action(
                            chat_id=chat.id,
                            action=ChatAction.TYPING,
                        )
                        return await sender(
                            text=i18n["gsmarena"]["reason"]["deviceLister"].format(
                                query=query,
                            ),
                            reply_markup=keyboard,
                        )
                    await GSMarenaBuildMessage(c=c, deviceID=getDevice.results[0].id, i18n=i18n)
                except GSMarenaDeviceNotFound as err:
                    await c.send_chat_action(
                        chat_id=chat.id,
                        action=ChatAction.TYPING,
                    )
                    await sender(i18n["gsmarena"]["reason"]["deviceNotFound"])
                    log.warning("%s: %s", err, query)
            else:
                await c.send_chat_action(
                    chat_id=chat.id,
                    action=ChatAction.TYPING,
                )
                await union.reply(i18n["gsmarena"]["reason"]["deviceNotProvided"])
        elif isinstance(union, CallbackQuery):
            inf = union.data.split("|")
            deviceID = inf[1]
            userID = inf[2]

            if union.from_user.id != int(userID):
                return await union.answer(i18n["gsmarena"]["checkers"]["notforYou"], show_alert=True)
            await GSMarenaBuildMessage(c=c, deviceID=deviceID, i18n=i18n)

gsm_arena_utils = GSMarenaUtils()