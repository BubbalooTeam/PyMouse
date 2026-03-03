package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"pymouse/pymouse/helpers/rapidhttp"
)

type weatherLocationResponse struct {
	Location struct {
		Address   []string  `json:"address"`
		Timezone  []string  `json:"ianaTimezone"`
		Latitude  []float64 `json:"latitude"`
		Longitude []float64 `json:"longitude"`
	} `json:"location"`
}

type WeatherLocationInfo struct {
	LocationName string
	Timezone     string
	Latitude     float64
	Longitude    float64
}

type WeatherInfo struct {
	WxObservations struct {
		Temperature          int    `json:"temperature"`
		TemperatureFeelsLike int    `json:"temperatureFeelsLike"`
		Humidity             int    `json:"relativeHumidity"`
		WindSpeed            int    `json:"windSpeed"`
		IconCode             int    `json:"iconCode"`
		WXPhraseLong         string `json:"wxPhraseLong"`
	} `json:"v3-wx-observations-current"`
}

type WeatherDailyForecastOverview struct {
	IconCode  int
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

func GetWeatherLocationInfo(locationName string, language string) (*WeatherLocationInfo, error) {
	var APILocationInfo weatherLocationResponse
	httpClient := rapidhttp.GetHTTPClient()
	getCoordsParams := map[string]string{
		"apiKey":   weatherAPIKey,
		"format":   "json",
		"language": language,
		"query":    locationName,
	}

	r, err := rapidhttp.Request(
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
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(rBody, &APILocationInfo)

	if len(APILocationInfo.Location.Address) == 0 {
		return nil, fmt.Errorf("location not found")
	}

	return &WeatherLocationInfo{
		LocationName: APILocationInfo.Location.Address[0],
		Timezone:     APILocationInfo.Location.Timezone[0],
		Latitude:     APILocationInfo.Location.Latitude[0],
		Longitude:    APILocationInfo.Location.Longitude[0],
	}, nil
}

func GetWeatherInfo(latitude, longitude float64, language string) (*WeatherInfo, error) {
	var APIWeatherInfo WeatherInfo
	httpClient := rapidhttp.GetHTTPClient()
	getWeatherParams := map[string]string{
		"apiKey":   weatherAPIKey,
		"format":   "json",
		"language": language,
		"geocode":  fmt.Sprintf("%f,%f", latitude, longitude),
		"units":    "m",
	}
	r, err := rapidhttp.Request(
		httpClient,
		rapidhttp.HTTPStruct{
			Method: "GET",
			URL:    getWeather,
			GETParams: &rapidhttp.HTTPGetStruct{
				Params: getWeatherParams,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(rBody, &APIWeatherInfo)
	return &APIWeatherInfo, nil
}
