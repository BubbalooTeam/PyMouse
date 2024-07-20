from shutil import rmtree

from hydrogram.types import Message, InlineQuery

from pymouse import Weather
from pymouse.utils.tools.weather.exceptions import (
    WeatherLocationNotFound,
    WeatherLocationNotProvidedError
)

async def HandleWeather(m: Message, i18n: dict):
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
    graphic = Weather().MakeWeatherInterface(
        weather=AllClassinOne,
        i18n=i18n,
    )
    await m.reply_photo(
        photo=graphic,
    )
    rmtree(
        path=graphic,
    )