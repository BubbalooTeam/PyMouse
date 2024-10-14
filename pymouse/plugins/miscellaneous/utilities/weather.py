#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.

from hydrogram.types import Message
from hydrogram.enums import ChatAction

from pymouse import PyMouse, Weather, DownloadPaths, log
from pymouse.utils.tools.weather.exceptions import (
    WeatherLocationNotFound,
    WeatherLocationNotProvidedError
)

class WeatherUtils:
    @staticmethod
    async def HandleWeather(c: PyMouse, m: Message, i18n: dict): # type: ignore
        await c.send_chat_action(
            chat_id=m.chat.id,
            action=ChatAction.TYPING,
        )
        try:
            location = await Weather().GetCoordsLocalization(
                union=m,
                i18n=i18n,
            )
        except WeatherLocationNotProvidedError:
            return await m.reply(i18n["weather"]["reason"]["weatherNotProvided"])
        except WeatherLocationNotFound:
            return await m.reply(i18n["weather"]["reason"]["weatherLocationNotFound"])
        # === #
        current_weather = await Weather().GetCurrentWeather(
            get_coords=location,
            i18n=i18n,
        )
        forecast_weather = await Weather().GetWeatherForecast(
            get_coords=location,
            i18n=i18n,
        )
        AllClassinOne = Weather().AllinOneClass(
            location_class=location,
            weather_class=current_weather,
            forecast_class=forecast_weather,
        )
        ### Make graphic interface
        await c.send_chat_action(
            chat_id=m.chat.id,
            action=ChatAction.UPLOAD_PHOTO,
        )
        graphic = Weather().MakeWeatherInterface(
            weather=AllClassinOne,
            i18n=i18n,
        )
        await m.reply_photo(
            photo=graphic,
        )

        rmpath = DownloadPaths().Delete_FileName(graphic)
        if rmpath in ["Path doesn't exists..", "Deleted filename and all content there.."]:
            pass
        else:
            log.error("Fail in Delete Filename..")

weather_utils = WeatherUtils()