#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.


from dataclasses import dataclass
from datetime import datetime
from PIL import Image, ImageDraw, ImageFont
from pytz import timezone
from typing import Union
from uuid import uuid4

from hydrogram.types import Message, InlineQuery

from pymouse import Config
from pymouse.assents import AllIcons, AllFonts
from pymouse.utils import HandleText, http
from pymouse.utils.tools.weather.exceptions import WeatherLocationNotProvidedError, WeatherLocationNotFound

Icons = {
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
    location_name: str | None
    timezone: str | None
    latitude: int | None
    longitude: int | None

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

@dataclass(frozen=True, slots=True)
class WeatherDailyForecastOverview:
    iconCode: int | None
    shortcast: str | None

@dataclass(frozen=True, slots=True)
class WeatherDailyForecastInfo:
    day: str | None
    max_temperature: int | None
    min_temperature: int | None
    overview: WeatherDailyForecastOverview

@dataclass(frozen=True, slots=True)
class WeatherDailyForecast:
    forecast: list[WeatherDailyForecastInfo]

@dataclass(frozen=True, slots=True)
class WeatherResponse:
    location: WeatherLocationInfo
    current_weather: WeatherInfo
    daily_forecast: WeatherDailyForecast

class Weather:
    def __init__(self) -> None:
        self.GetCoords: str = "https://api.weather.com/v3/location/search"
        self.GetWeatherUrl: str = "https://api.weather.com/v3/aggcommon/v3-wx-observations-current"
        self.Headers: dict = {"User-Agent": "Dalvik/2.1.0 (Linux; U; Android 12; M2012K11AG Build/SQ1D.211205.017)"}
        self.weatherAPIKey = "8de2d8b3a93542c9a2d8b3a935a2c909"
        self.ForecastDayURL = "https://api.weather.com/v1/geocode/{lat}/{lon}/forecast/daily/7day.json"

    @staticmethod
    def GetIconImage(IconCode: int) -> Image.open:
        """
        Get weather icon in image format by IconCode.

        Arguments:
            - `IconCode(int)`: Icon encoding number per Key.

        Returns:
            - `IconImage(Image.open)`: Represents the icon image.
        """
        FilePath: str = Icons.get(IconCode)
        pil_img = Image.open(
            fp=FilePath,
        )
        return pil_img.resize((150, 150), Image.Resampling.LANCZOS)

    @staticmethod
    def GetKeyDayorNight(get_coords: WeatherLocationInfo) -> str:
        """
        Get if is day or night with DateTime Object.

        Returns:
            - `str`: Return Key ("day" or "night").
        """
        zonetime = timezone(get_coords.timezone)
        now = datetime.now(zonetime)
        hour = now.hour
        if 6 <= hour < 18:
            return "day"
        else:
            return "night"

    @staticmethod
    def GetWeatherLocation(union: Union[Message, InlineQuery]) -> str:
        """
        Get the Location the user entered to check the Weather.

        Arguments:
            - `union(Message, InlineQuery)`: Instance of Telegram Message/InlineQuery.

        Returns:
            - `Localization(str)`: Return the Location entered to check the Weather.
        """
        Localization: str = HandleText().input_str(union)
        if not Localization:
            raise WeatherLocationNotProvidedError("The location was not entered. In this condition, I cannot perform the climate check.")
        else:
            return Localization

    @staticmethod
    def _parse_WeatherInfo(weather_dict: dict) -> WeatherInfo:
        """
        Formats the information in dict in an Information Class (WeatherInfo).

        Arguments:
            - `weather_dict(dict)`: Information Dictionary in the Json Language returned by the API.

        Returns:
            - `WeatherInfo`: Information Class that is based on the Dictionary.
        """
        return WeatherInfo(
            temperature=weather_dict.get("v3-wx-observations-current", {}).get("temperature"),
            temperature_FeelsLike=weather_dict.get("v3-wx-observations-current", {}).get("temperatureFeelsLike"),
            relativeHumidity=weather_dict.get("v3-wx-observations-current", {}).get("relativeHumidity"),
            windSpeed=weather_dict.get("v3-wx-observations-current", {}).get("windSpeed"),
            overview=WeatherOverview(
                iconCode=weather_dict.get("v3-wx-observations-current", {}).get("iconCode"),
                wxPhraseLong=weather_dict.get("v3-wx-observations-current", {}).get("wxPhraseLong"),
            ),
        )

    @staticmethod
    def _parse_WeatherLocationInfo(location_dict: dict) -> WeatherLocationInfo:
        """
        Formats the information in dict in an Information Class (WeatherLocationInfo).

        Arguments:
            - `location_dict(dict)`: Information Dictionary in the Json Language returned by the API.

        Returns:
            - WeatherLocationInfo: Information Class that is based on the Dictionary.
        """
        return WeatherLocationInfo(
            location_name=location_dict.get("address")[0],
            timezone=location_dict.get("ianaTimeZone")[0],
            latitude=location_dict.get("latitude")[0],
            longitude=location_dict.get("longitude")[0],
        )

    @staticmethod
    def _parse_WeatherDailyForecastInfo(
        day_or_night: str,
        forecast_dict: dict,
    ) -> WeatherDailyForecastInfo:
        return WeatherDailyForecastInfo(
                day=forecast_dict.get("dow"),
                max_temperature=forecast_dict.get("max_temp"),
                min_temperature=forecast_dict.get("min_temp"),
                overview=WeatherDailyForecastOverview(
                    iconCode=forecast_dict.get(day_or_night, {}).get("icon_code"),
                    shortcast=forecast_dict.get(day_or_night, {}).get("phrase_32char")
                ),
            )

    @staticmethod
    def _parse_WeatherDailyForecast(
        forecast_list: list[WeatherDailyForecastInfo],
    ) -> WeatherDailyForecast:
        return WeatherDailyForecast(
            forecast=forecast_list,
        )

    async def GetCoordsLocalization(self, union: Union[Message, InlineQuery], i18n: dict) -> WeatherLocationInfo:
        """
        Gets the coordinates of the specified location.

        Arguments:
            - `union(Message, InlineQuery)`: Instance of Telegram Message/InlineQuery.
            - `i18n(dict)`: Dictionary of international words.

        Returns:
            - WeatherLocationInfo: Information Class that is based on the Dictionary.
        """
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
            return loc_tuple
        else:
            raise WeatherLocationNotFound("Location Requested Not Found!")

    async def GetWeatherForecast(
        self,
        get_coords: WeatherLocationInfo,
        i18n: dict,
    ) -> list:
        """
        Gets the weather for every day of the week (7 periodic days).

        Arguments:
            - `get_coords(WeatherLocationInfo)`: Coordinates of the specified Location.
            - `i18n(dict)`: Dictionary of international words.

        Returns:
            - `list(dataclass)`: List of Information Class that is based on the Dictionary.
        """
        forecast_list = []
        weather_language = i18n["weather"]["language"]
        weather_units = i18n["weather"]["units"]
        # Generate URL with params #
        url = self.ForecastDayURL.format(
            lat=get_coords.latitude,
            lon=get_coords.longitude,
        )
        WeatherParams = dict(
            apiKey=self.weatherAPIKey,
            units=weather_units,
            language=weather_language,
        )
        # Make a Response List with HTTP Method #
        r = await http.get(
            url=url,
            headers=self.Headers,
            params=WeatherParams,
        )
        forecast_json = r.json()
        day_or_night = self.GetKeyDayorNight(get_coords=get_coords)
        for forecast in forecast_json["forecasts"]:
            current_forecast = self._parse_WeatherDailyForecastInfo(
                day_or_night=day_or_night,
                forecast_dict=forecast,
            )
            forecast_list.append(current_forecast)
        return self._parse_WeatherDailyForecast(forecast_list=forecast_list)

    async def GetCurrentWeather(
        self,
        get_coords: WeatherLocationInfo,
        i18n: dict
    ) -> WeatherInfo:
        """
        Gets the current weather based on the request made.

        Arguments:
            - `get_coords(WeatherLocationInfo)`: Coordinates of the specified Location.
            - `i18n(dict)`: Dictionary of international words.

        Returns:
            - `WeatherInfo`: Information Class that is based on the Dictionary.
        """
        weather_language = i18n["weather"]["language"]
        weather_units = i18n["weather"]["units"]
        WeatherParams = dict(
            apiKey=self.weatherAPIKey,
            format="json",
            language=weather_language,
            geocode="{lat},{lon}".format(
                lat=get_coords.latitude,
                lon=get_coords.longitude
            ),
            units=weather_units,
        )
        r = await http.get(
            url=self.GetWeatherUrl,
            headers=self.Headers,
            params=WeatherParams,
        )
        weather_json = r.json()
        tupled_infos = self._parse_WeatherInfo(
            weather_dict=weather_json,
        )
        return tupled_infos

    @staticmethod
    def AllinOneClass(
            location_class: WeatherLocationInfo,
            weather_class: WeatherInfo,
            forecast_class: WeatherDailyForecast,
    ) -> WeatherResponse:
        return WeatherResponse(
            location=location_class,
            current_weather=weather_class,
            daily_forecast=forecast_class,
        )

    def MakeWeatherInterface(
        self,
        weather: WeatherResponse,
        i18n: dict,
    ):
        hash_code = uuid4()
        width, height = 3840, 2160
        gradient = Image.new('RGB', (width, height), "#02455B")
        draw = ImageDraw.Draw(gradient)
        filename = str(Config.DOWNLOAD_PATH + "weather_{hashcode}.png").format(
            hashcode=hash_code,
        )

        # Make Background Color
        for y in range(height):
            color = (60, 8, 95 + int(80 * (y / height)))
            draw.line([(0, y), (width, y)], fill=color)

        # Get the Text Font
        ArialFont = AllFonts.ARIAL
        # === #
        BigFont = ImageFont.truetype(
            font=ArialFont,
            size=80,
        )
        UltraFont = ImageFont.truetype(
            font=ArialFont,
            size=65
        )
        MediumFont = ImageFont.truetype(
            font=ArialFont,
            size=60
        )
        SmallFont = ImageFont.truetype(
            font=ArialFont,
            size=35
        )
        # === #
        location = weather.location.location_name
        location_bbox = draw.textbbox(
            xy=(0, 0),
            text=location,
            font=BigFont,
        )
        XLocation = (width - (location_bbox[2] - location_bbox[0])) / 2
        draw.text(
            xy=(XLocation, 50),
            text=i18n["weather"]["location-name"].format(
                location=location,
            ),
            fill="white",
            font=BigFont,
        )
        # Draw Current Weather
        YCurrentWeather = 400
        CurrentWeatherIcon = self.GetIconImage(weather.current_weather.overview.iconCode)
        gradient.paste(
            im=CurrentWeatherIcon,
            box=(180, YCurrentWeather),
            mask=CurrentWeatherIcon
        )
        draw.text(
            xy=(380, YCurrentWeather),
            text=i18n["weather"]["temperature"].format(
                temperature=weather.current_weather.temperature,
            ),
            font=BigFont,
            fill="white"
        )
        draw.text(
            xy=(380, YCurrentWeather + 100),
            text=weather.current_weather.overview.wxPhraseLong,
            font=BigFont,
            fill="white",
        )
        draw.text(
            xy=(180, YCurrentWeather + 250),
            text=i18n["weather"]["temperature-FeelsLike"].format(
                feelslike=weather.current_weather.temperature_FeelsLike,
            ),
            font=UltraFont,
            fill="white",
        )
        draw.text(
            xy=(180, YCurrentWeather + 350),
            text=i18n["weather"]["humidity"].format(
                humidity=weather.current_weather.relativeHumidity,
            ),
            font=UltraFont,
            fill="white",
        )
        draw.text(
            xy=(180, YCurrentWeather + 450),
            text=i18n["weather"]["wind"].format(
                windspeed=weather.current_weather.windSpeed,
            ),
            font=UltraFont,
            fill="white",
        )
        # Draw Weekly Forecast
        YForecast = YCurrentWeather + 700
        XStart = 180
        XSpacing = 450
        for ForecastDaily in weather.daily_forecast.forecast:
            XIcon = XStart
            YIcon = YForecast
            XText = XStart
            YText = YForecast + 240
            ForecastDailyIcon = self.GetIconImage(ForecastDaily.overview.iconCode)
            # === #
            gradient.paste(
                im=ForecastDailyIcon,
                box=(XIcon, YIcon),
                mask=ForecastDailyIcon,
            )
            draw.text(
                xy=(XText, YText),
                text=ForecastDaily.day[:3],
                font=MediumFont,
                fill="white",
            )
            draw.text(
                xy=(XText, YText + 80),
                text="Max: {i18n_type}".format(
                    i18n_type=i18n["weather"]["temperature"].format(
                        temperature=ForecastDaily.max_temperature,
                    ),
                ),
                font=MediumFont,
                fill="white"
            )
            draw.text(
                xy=(XText, YText + 160),
                text="Min: {i18n_type}".format(
                    i18n_type=i18n["weather"]["temperature"].format(
                        temperature=ForecastDaily.min_temperature,
                    ),
                ),
                font=MediumFont,
                fill="white"
            )
            draw.text(
                xy=(XText, YText + 240),
                text=ForecastDaily.overview.shortcast[:18] + "..." if ForecastDaily.overview.shortcast >= 18 else ForecastDaily.overview.shortcast,
                font=SmallFont,
                fill="yellow"
            )
            XStart += XSpacing

        gradient.save(filename)
        return filename