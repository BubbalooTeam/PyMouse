from typing import Any

from hydrogram.types import Message

from pymouse import GSMarena, GSMarenaDeviceBaseResult, log
from pymouse.utils import HandleText
from pymouse.utils.tools.gsm_arena.exceptions import GSMarenaDeviceNotFound

gsm_arena = GSMarena()

def formatGSMarenaMessage(gsmarenaBaseResult: GSMarenaDeviceBaseResult, i18n: dict) -> str:
    phone = gsm_arena.parse_specifications(specifications=gsmarenaBaseResult.phone_details)
    deviceI18N = i18n["gsmarena"]["phoneFormatter"]

    attrs = [
        "<b>{spec_name}:</b> <i>{spec_value}</i>".format(
            spec_name=deviceI18N[key],
            spec_value=value
        )
        for key, value in phone.items()
        if value.strip() and value.strip() != "-" and deviceI18N.get(key) is not None
    ]


    formatted_message = f"<a href='{gsmarenaBaseResult.image}'>{gsmarenaBaseResult.name}</a>\n\n{"\n\n".join(attrs)}"
    return formatted_message

async def HandleGSMarena(m: Message, i18n: dict):
    query = HandleText().input_str(union=m)
    if query:
        try:
            getDevice = await gsm_arena.search_device(query=query)
            fetchDevice = await gsm_arena.fetch_device(device=getDevice.results[0].id)
            formatedMessage = formatGSMarenaMessage(gsmarenaBaseResult=fetchDevice, i18n=i18n)
            await m.reply(formatedMessage)
        except GSMarenaDeviceNotFound as err:
            await m.reply(i18n["gsmarena"]["reason"]["deviceNotFound"])
            log.warning("%s: %s", err, query)
    else:
        await m.reply(i18n["gsmarena"]["reason"]["deviceNotProvided"])