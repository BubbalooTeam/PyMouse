package weather

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"pymouse/pymouse/config"
	"pymouse/pymouse/helpers/rapidhttp"
	"pymouse/pymouse/helpers/utils"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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

type ForecastPeriod struct {
	IconCode  int    `json:"icon_code"`
	Shortcast string `json:"phrase_32char"`
}

type WeatherDailyForecastInfo struct {
	Date           string          `json:"dow"`
	TemperatureMax *int            `json:"max_temp"`
	TemperatureMin *int            `json:"min_temp"`
	Day            *ForecastPeriod `json:"day,omitempty"`
	Night          *ForecastPeriod `json:"night,omitempty"`
}

type WeatherDailyForecast struct {
	Forecast []WeatherDailyForecastInfo `json:"forecasts"`
}

type WeatherResponse struct {
	LocationInfo   *WeatherLocationInfo
	CurrentWeather *WeatherInfo
	DailyForecast  *WeatherDailyForecast
}

const (
	getCoords     = "https://api.weather.com/v3/location/search"
	getWeather    = "https://api.weather.com/v3/aggcommon/v3-wx-observations-current"
	weatherAPIKey = "8de2d8b3a93542c9a2d8b3a935a2c909"
	forecastDay   = "https://api.weather.com/v1/geocode/%f/%f/forecast/daily/7day.json"
	iconBase      = "pymouse/assets/icons/weather/"
)

var (
	headers = map[string]string{
		"User-Agent": "Dalvik/2.1.0 (Linux; U; Android 12; M2012K11AG Build/SQ1D.211205.017)",
	}
	weatherIconPaths = map[int]string{
		0:  iconBase + "rain_with_thumderstorm.png",
		1:  iconBase + "rain_with_thumderstorm.png",
		2:  iconBase + "rain_with_thumderstorm.png",
		3:  iconBase + "rain_with_thumderstorm.png",
		4:  iconBase + "rain_with_thumderstorm.png",
		5:  iconBase + "snow.png",
		6:  iconBase + "snow.png",
		7:  iconBase + "snow.png",
		8:  iconBase + "snow.png",
		9:  iconBase + "snow.png",
		10: iconBase + "snow.png",
		11: iconBase + "rain.png",
		12: iconBase + "rain.png",
		13: iconBase + "snow.png",
		14: iconBase + "snow.png",
		15: iconBase + "snow.png",
		16: iconBase + "snow.png",
		17: iconBase + "rain_with_thumderstorm.png",
		18: iconBase + "rain.png",
		19: iconBase + "fog.png",
		20: iconBase + "fog.png",
		21: iconBase + "fog.png",
		22: iconBase + "fog.png",
		23: iconBase + "wind.png",
		24: iconBase + "wind.png",
		25: iconBase + "snow.png",
		26: iconBase + "cloudy.png",
		27: iconBase + "mostly_cloudy.png",
		28: iconBase + "mostly_cloudy.png",
		29: iconBase + "partly_cloudy.png",
		30: iconBase + "partly_cloudy.png",
		31: iconBase + "moon.png",
		32: iconBase + "sunny.png",
		33: iconBase + "partly_cloudy.png",
		34: iconBase + "partly_cloudy.png",
		35: iconBase + "rain_with_thumderstorm.png",
		36: iconBase + "hot.png",
		37: iconBase + "thumderstorm.png",
		38: iconBase + "thumderstorm.png",
		39: iconBase + "rain.png",
		40: iconBase + "rain.png",
		41: iconBase + "snow.png",
		42: iconBase + "snow.png",
		43: iconBase + "snow.png",
		44: iconBase + "resync.png",
		45: iconBase + "rain.png",
		46: iconBase + "snow.png",
		47: iconBase + "thumderstorm.png",
	}
)

var weatherIconCache = map[int]image.Image{}

func init() {
	for code, path := range weatherIconPaths {

		img, err := gg.LoadImage(path)
		if err != nil {
			logrus.Error("failed loading icon:", path, err)
			continue
		}

		w := float64(img.Bounds().Dx())
		h := float64(img.Bounds().Dy())

		im := gg.NewContext(150, 150)

		scale := math.Min(150/w, 150/h)

		im.Push()
		im.Scale(scale, scale)
		im.DrawImageAnchored(img, int(w/2), int(h/2), 0.5, 0.5)
		im.Pop()

		weatherIconCache[code] = im.Image()
	}
}

func GetIconImage(iconCode int) image.Image {
	img, ok := weatherIconCache[iconCode]
	if !ok {
		return weatherIconCache[44]
	}
	return img
}

func getWeatherLocationInfo(locationName string, language string) (*WeatherLocationInfo, error) {
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
	if err := json.NewDecoder(r.Body).Decode(&APILocationInfo); err != nil {
		return nil, err
	}

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

func getWeatherInfo(latitude, longitude float64, language string, units string) (*WeatherInfo, error) {
	var APIWeatherInfo WeatherInfo
	httpClient := rapidhttp.GetHTTPClient()
	getWeatherParams := map[string]string{
		"apiKey":   weatherAPIKey,
		"format":   "json",
		"language": language,
		"geocode":  fmt.Sprintf("%f,%f", latitude, longitude),
		"units":    units,
	}
	r, err := rapidhttp.Request(
		httpClient,
		rapidhttp.HTTPStruct{
			Method:  "GET",
			URL:     getWeather,
			Headers: headers,
			GETParams: &rapidhttp.HTTPGetStruct{
				Params: getWeatherParams,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&APIWeatherInfo); err != nil {
		return nil, err
	}

	return &APIWeatherInfo, nil
}

func getKeyDayOrNightbyTime(now time.Time) string {
	hour := now.Hour()
	if hour >= 6 && hour <= 18 {
		return "day"
	}
	return "night"
}

func getKeyDayOrNight(timezone string) string {
	zoneTime, err := time.LoadLocation(timezone)
	if err != nil {
		now := time.Now().UTC()
		r := getKeyDayOrNightbyTime(now)
		return r
	}
	now := time.Now().In(zoneTime)
	r := getKeyDayOrNightbyTime(now)
	return r
}

func getForecastInfo(latitude, longitude float64, language string, units string) (*WeatherDailyForecast, error) {
	var APIForecastDailyInfo WeatherDailyForecast

	httpClient := rapidhttp.GetHTTPClient()
	getForecastsParams := map[string]string{
		"apiKey":   weatherAPIKey,
		"language": language,
		"units":    units,
	}
	getForecastsURL := fmt.Sprintf(forecastDay, latitude, longitude)

	r, err := rapidhttp.Request(
		httpClient,
		rapidhttp.HTTPStruct{
			Method:  "GET",
			URL:     getForecastsURL,
			Headers: headers,
			GETParams: &rapidhttp.HTTPGetStruct{
				Params: getForecastsParams,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&APIForecastDailyInfo); err != nil {
		return nil, err
	}

	return &APIForecastDailyInfo, nil
}

func (w *WeatherDailyForecastInfo) GetForecastPeriod(period string) *ForecastPeriod {
	if strings.Contains(period, "night") && w.Night != nil {
		return w.Night
	}
	if strings.Contains(period, "day") && w.Day != nil {
		return w.Day
	}

	if w.Day != nil {
		return w.Day
	}
	return w.Night
}

func GetWeatherResponse(location string, l func(string) string) (*WeatherResponse, error) {
	localizationInfo, err := getWeatherLocationInfo(location, l("weather.language"))
	if err != nil {
		return nil, fmt.Errorf("failed to extract location informations: %v", err)
	}
	currentWeather, err := getWeatherInfo(localizationInfo.Latitude, localizationInfo.Longitude, l("weather.language"), l("weather.units"))
	if err != nil {
		return nil, fmt.Errorf("failed to extract weather informations: %v", err)
	}
	forecastWeather, err := getForecastInfo(localizationInfo.Latitude, localizationInfo.Longitude, l("weather.language"), l("weather.units"))
	if err != nil {
		return nil, fmt.Errorf("failed to extract weather forecast informations: %v", err)
	}
	return &WeatherResponse{
		LocationInfo:   localizationInfo,
		CurrentWeather: currentWeather,
		DailyForecast:  forecastWeather,
	}, nil
}

func MakeWeatherInterface(
	weather *WeatherResponse,
	l func(string) string,
) (string, error) {

	hashCode := uuid.New().String()
	width := 3840.0
	height := 2160.0

	im := gg.NewContext(int(width), int(height))

	isNight := getKeyDayOrNight(weather.LocationInfo.Timezone) == "night"

	grad := gg.NewLinearGradient(0, 0, 0, height)

	if isNight {
		grad.AddColorStop(0, color.RGBA{60, 8, 95, 255})
		grad.AddColorStop(1, color.RGBA{60, 8, 175, 255})
	} else {
		grad.AddColorStop(0, color.RGBA{30, 144, 255, 255})
		grad.AddColorStop(1, color.RGBA{0, 0, 139, 255})
	}

	im.SetFillStyle(grad)
	im.DrawRectangle(0, 0, width, height)
	im.Fill()

	if err := im.LoadFontFace("pymouse/assets/fonts/economica-italic.ttf", 80); err != nil {
		return "", err
	}

	locationText := weather.LocationInfo.LocationName

	im.SetRGB(1, 1, 1)
	im.DrawStringAnchored(locationText, width/2, 120, 0.5, 0.5)

	if err := im.LoadFontFace("pymouse/assets/fonts/notosans-bold.ttf", 80); err != nil {
		return "", err
	}
	yCurrent := 400.0

	periodNow := weather.CurrentWeather

	icon := GetIconImage(periodNow.WxObservations.IconCode)

	im.DrawImage(icon, 180, int(yCurrent))

	im.DrawString(
		fmt.Sprintf(l("weather.temperature"),
			periodNow.WxObservations.Temperature),
		380,
		yCurrent+80,
	)

	if err := im.LoadFontFace("pymouse/assets/fonts/economica-italic.ttf", 80); err != nil {
		return "", err
	}

	im.DrawString(
		periodNow.WxObservations.WXPhraseLong,
		380,
		yCurrent+180,
	)

	if err := im.LoadFontFace("pymouse/assets/fonts/arial.ttf", 65); err != nil {
		return "", err
	}

	im.DrawString(
		fmt.Sprintf(
			l("weather.temperature-feelsLike"),
			periodNow.WxObservations.TemperatureFeelsLike),
		180,
		yCurrent+300,
	)

	im.DrawString(
		fmt.Sprintf(l("weather.humidity"),
			periodNow.WxObservations.Humidity),
		180,
		yCurrent+400,
	)

	im.DrawString(
		fmt.Sprintf(l("weather.wind"),
			periodNow.WxObservations.WindSpeed),
		180,
		yCurrent+500,
	)

	yForecast := yCurrent + 700
	xStart := 180.0
	xSpacing := 450.0

	for _, daily := range weather.DailyForecast.Forecast {
		var period *ForecastPeriod

		if isNight && daily.Night != nil {
			period = daily.Night
		} else if !isNight && daily.Day != nil {
			period = daily.Day
		} else if daily.Day != nil {
			period = daily.Day
		} else {
			period = daily.Night
		}

		if period == nil {
			continue
		}

		icon := GetIconImage(period.IconCode)

		im.DrawImage(icon, int(xStart), int(yForecast))

		dayLabel := utils.FirstRunes(daily.Date, 3)

		if err := im.LoadFontFace("pymouse/assets/fonts/notosans-bold.ttf", 60); err != nil {
			return "", err
		}
		im.SetRGB(1, 1, 1)
		im.DrawString(dayLabel, xStart, yForecast+260)

		if err := im.LoadFontFace("pymouse/assets/fonts/arial.ttf", 60); err != nil {
			return "", err
		}

		if tempMax := daily.TemperatureMax; tempMax != nil {
			im.DrawString(
				fmt.Sprintf("Max: %d°", *tempMax),
				xStart,
				yForecast+340,
			)
		}

		if tempMin := daily.TemperatureMin; tempMin != nil {
			im.DrawString(
				fmt.Sprintf("Min: %d°", *tempMin),
				xStart,
				yForecast+420,
			)
		}

		shortcast := period.Shortcast
		if len(shortcast) > 18 {
			shortcast = shortcast[:18] + "..."
		}

		if err := im.LoadFontFace("pymouse/assets/fonts/arial.ttf", 35); err != nil {
			return "", err
		}

		im.SetRGB(1, 1, 0)
		im.DrawString(shortcast, xStart, yForecast+500)
		im.SetRGB(1, 1, 1)

		xStart += xSpacing
	}

	dir := fmt.Sprintf("%s/%s", config.DownloadPath, "weather")

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	filename := filepath.Join(
		dir,
		fmt.Sprintf("weather_%s.png", hashCode),
	)

	if err := im.SavePNG(filename); err != nil {
		return "", err
	}

	return filepath.Abs(filename)
}
