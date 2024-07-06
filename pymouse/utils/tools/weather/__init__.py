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
from pymouse.utils import HandleText
from pymouse.utils.tools.weather.exceptions import WeatherLocationNotProvidedError

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
    19: "🌫",
    20: "🌫",
    21: "🌫",
    22: "🌫",
    23: "🌬",
    24: "🌬",
    25: "🌨",
    26: "☁️",
    27: "🌥",
    28: "🌥",
    29: "⛅️",
    30: "⛅️",
    31: "🌙",
    32: AllIcons.Weather.SUNNY,
    33: "🌤",
    34: "🌤",
    35: "⛈",
    36: "🔥",
    37: "🌩",
    38: "🌩",
    39: "🌧",
    40: "🌧",
    41: "❄️",
    42: "❄️",
    43: "❄️",
    44: "n/a",
    45: "🌧",
    46: "🌨",
    47: "🌩",
}

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
        self.get_coords: str = "https://api.weather.com/v3/location/search"
        self.get_weather_url: str = "https://api.weather.com/v3/aggcommon/v3-wx-observations-current"
        self.headers: dict = {"User-Agent": "Dalvik/2.1.0 (Linux; U; Android 12; M2012K11AG Build/SQ1D.211205.017)"}

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