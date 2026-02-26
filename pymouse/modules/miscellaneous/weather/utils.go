package weather

import (
	"pymouse/pymouse/helpers/rapidhttp"
)

type WeatherLocationInfo struct {
	LocationName string
	Timezone     string
	Latitude     float64
	Longitude    float64
}

type WeatherOverview struct {
	IconCode      string
	WXPhreaseLong string
}

type WeatherInfo struct {
	Temperature          int
	TemperatureFeelsLike int
	Humidity             int
	WindSpeed            int
	Overview             WeatherOverview
}

type WeatherDailyForecastOverview struct {
	IconCode  string
	Shortcast string
}

type WeatherDailyForecastInfo struct {
	Date           string
	TemperatureMax int
	TemperatureMin int
	Overview       WeatherDailyForecastOverview
}

type WeatherDailyForecast struct {
	Forecast []WeatherDailyForecastInfo
}

type WeatherResponse struct {
	LocationInfo   WeatherLocationInfo
	CurrentWeather WeatherInfo
	DailyForecast  WeatherDailyForecast
}

const (
	getCoords     = "https://api.weather.com/v3/location/search"
	getWeather    = "https://api.weather.com/v3/aggcommon/v3-wx-observations-current"
	weatherAPIKey = "8de2d8b3a93542c9a2d8b3a935a2c909"
	forecastDay   = "https://api.weather.com/v1/geocode/%s/%s/forecast/daily/7day.json"
)

var headers = map[string]string{
	"User-Agent": "Dalvik/2.1.0 (Linux; U; Android 12; M2012K11AG Build/SQ1D.211205.017)",
}

func getWeatherLocationInfo(locationName string, language string) (WeatherLocationInfo, error) {
	httpClient := rapidhttp.GetHTTPClient()
	getCoordsParams := map[string]string{
		"apiKey":   weatherAPIKey,
		"format":   "json",
		"language": language,
		"query":    locationName,
	}

	rapidhttp.Request(
		httpClient,
		rapidhttp.HTTPStruct{
			Method:  "GET",
			URL:     getCoords,
			Headers: headers,
			GETParams: &rapidhttp.HTTPGetStruct{
				Params: getCoordsParams,
			},
		},
	)
	return WeatherLocationInfo{}, nil
}
