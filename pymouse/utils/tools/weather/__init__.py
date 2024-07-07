#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.


from dataclasses import dataclass
from typing import Union

from hydrogram.types import Message, InlineQuery

from pymouse.assents import AllIcons
from pymouse.utils import HandleText, http
from pymouse.utils.tools.weather.exceptions import WeatherLocationNotProvidedError, WeatherLocationNotFound

status_emojis = {
    0: AllIcons.Weather.RAIN_THUMDERSTORM,
    1: AllIcons.Weather.RAIN_THUMDERSTORM,
    2: AllIcons.Weather.RAIN_THUMDERSTORM,
    3: AllIcons.Weather.RAIN_THUMDERSTORM,
    4: AllIcons.Weather.RAIN_THUMDERSTORM,
    5: AllIcons.Weather.SNOW,
    6: AllIcons.Weather.SNOW,
    7: AllIcons.Weather.SNOW,
    8: AllIcons.Weather.SNOW,
    9: AllIcons.Weather.SNOW,
    10: AllIcons.Weather.SNOW,
    11: AllIcons.Weather.RAIN,
    12: AllIcons.Weather.RAIN,
    13: AllIcons.Weather.SNOW,
    14: AllIcons.Weather.SNOW,
    15: AllIcons.Weather.SNOW,
    16: AllIcons.Weather.SNOW,
    17: AllIcons.Weather.RAIN_THUMDERSTORM,
    18: AllIcons.Weather.RAIN,
    19: AllIcons.Weather.FOG,
    20: AllIcons.Weather.FOG,
    21: AllIcons.Weather.FOG,
    22: AllIcons.Weather.FOG,
    23: AllIcons.Weather.WIND,
    24: AllIcons.Weather.WIND,
    25: AllIcons.Weather.SNOW,
    26: AllIcons.Weather.CLOUDY,
    27: AllIcons.Weather.MOSTLY_CLOUDY,
    28: AllIcons.Weather.MOSTLY_CLOUDY,
    29: AllIcons.Weather.PARTLY_CLOUDY,
    30: AllIcons.Weather.PARTLY_CLOUDY,
    31: AllIcons.Weather.MOON,
    32: AllIcons.Weather.SUNNY,
    33: AllIcons.Weather.PARTLY_CLOUDY,
    34: AllIcons.Weather.PARTLY_CLOUDY,
    35: AllIcons.Weather.RAIN_THUMDERSTORM,
    36: AllIcons.Weather.HOT,
    37: AllIcons.Weather.THUMDERSTORM,
    38: AllIcons.Weather.THUMDERSTORM,
    39: AllIcons.Weather.RAIN,
    40: AllIcons.Weather.RAIN,
    41: AllIcons.Weather.SNOW,
    42: AllIcons.Weather.SNOW,
    43: AllIcons.Weather.SNOW,
    44: AllIcons.Weather.RESYNC,
    45: AllIcons.Weather.RAIN,
    46: AllIcons.Weather.SNOW,
    47: AllIcons.Weather.THUMDERSTORM,
}

@dataclass(frozen=True, slots=True)
class WeatherLocationInfo:
    latitude: int | None
    longitude: int | None

@dataclass(frozen=True, slots=True)
class WeatherCoords:
    geocode: str | None

@dataclass(frozen=True, slots=True)
class WeatherOverview:
    iconCode: int | None
    wxPhraseLong: str | None


@dataclass(frozen=True, slots=True)
class WeatherInfo:
    temperature: int | None
    temperature_FeelsLike: int | None
    relativeHumidity: int | None
    windSpeed: int | None
    overview: WeatherOverview


class Weather:
    def __init__(self) -> None:
        self.GetCoords: str = "https://api.weather.com/v3/location/search"
        self.GetWeatherUrl: str = "https://api.weather.com/v3/aggcommon/v3-wx-observations-current"
        self.Headers: dict = {"User-Agent": "Dalvik/2.1.0 (Linux; U; Android 12; M2012K11AG Build/SQ1D.211205.017)"}
        self.weatherAPIKey = "8de2d8b3a93542c9a2d8b3a935a2c909"

    @staticmethod
    def GetWeatherLocation(union: Union[Message, InlineQuery]) -> str:
        """
        Get the Location the user entered to check the Weather.

        Arguments:
            - `m(Message)`: Instance of Telegram Message.

        Returns:
            - Location(str): Return the Location entered to check the Weather.
        """
        Localization: str = HandleText().input_str(union)
        if not Localization:
            raise WeatherLocationNotProvidedError("The location was not entered. In this condition, I cannot perform the climate check.")
        else:
            return Localization

    @staticmethod
    def _parse_WeatherInfo(weather_dict: dict) -> WeatherInfo:
        return WeatherInfo(
            temperature=weather_dict.get("temperature"),
            temperature_FeelsLike=weather_dict.get("temperatureFeelsLike"),
            relativeHumidity=weather_dict.get("relativeHumidity"),
            windSpeed=weather_dict.get("windSpeed"),
            overview=WeatherOverview(
                iconCode=weather_dict.get("iconCode"),
                wxPhraseLong=weather_dict.get("wxPhraseLong"),
            ),
        )
    
    @staticmethod
    def _parse_WeatherLocationInfo(location_dict: dict) -> WeatherLocationInfo:
        return WeatherLocationInfo(
            latitude=location_dict.get("latitude")[0],
            longitude=location_dict.get("longitude")[0],
        )
    
    async def GetCoordsLocalization(self, union: Union[Message, InlineQuery], i18n: dict) -> WeatherCoords:
        getLocation = self.GetWeatherLocation(union=union)
        getCoordsParams = dict(
            apiKey=self.weatherAPIKey,
            format="json",
            language=i18n["weather"]["language"],
            query=getLocation
        )
        # Make a Request using HTTP Method
        r = await http.get(
            url=self.GetCoords,
            headers=self.Headers,
            params=getCoordsParams,
        )
        getted_json = r.json()
        loc_json = getted_json.get("location")
        if loc_json:
            loc_tuple = self._parse_WeatherLocationInfo(
                location_dict=loc_json,
            )
            return WeatherCoords(
                geocode="{latitude},{longitude}".format(
                    latitude=loc_tuple.latitude,
                    longitude=loc_tuple.longitude,
                )
            )
        else:
            raise WeatherLocationNotFound("Location Requested Not Found!")