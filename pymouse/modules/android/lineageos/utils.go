package lineageos

import (
	"encoding/json"
	"io"
	"pymouse/pymouse/helpers/rapidhttp"
	"strconv"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type LineageBuild struct {
	URL      string
	Filename string
	Size     string
	Version  string
	Datetime string
	ROMType  string
}

type lineageResponse struct {
	Response []struct {
		Filename string `json:"filename"`
		URL      string `json:"url"`
		Size     int64  `json:"size"`
		Version  string `json:"version"`
		Datetime int64  `json:"datetime"`
		ROMType  string `json:"romtype"`
	} `json:"response"`
}

type lineageResultsCache struct {
	OwnerID int64
	Device  string
	Builds  []LineageBuild
}

var lineageCache = cache.New(15*time.Minute, 30*time.Minute)

func setLineageResults(key string, ownerID int64, device string, builds []LineageBuild) {
	lineageCache.Set(key, lineageResultsCache{
		OwnerID: ownerID,
		Device:  device,
		Builds:  builds,
	}, cache.DefaultExpiration)
}

func getLineageResults(key string) (lineageResultsCache, bool) {
	value, found := lineageCache.Get(key)
	if !found {
		return lineageResultsCache{}, false
	}

	results, ok := value.(lineageResultsCache)
	return results, ok
}

func getLineageBuilds(device string) ([]LineageBuild, bool) {
	client := rapidhttp.GetHTTPClient()
	url := "https://download.lineageos.org/api/v1/" + device + "/nightly/*"
	resp, err := rapidhttp.Request(client, rapidhttp.HTTPStruct{
		Method: "GET",
		URL:    url,
	})
	if err != nil {
		logrus.Errorf("failed to fetch LineageOS builds for device %s: %v", device, err)
		return nil, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("failed to read LineageOS response for device %s: %v", device, err)
		return nil, false
	}

	var data lineageResponse
	if err := json.Unmarshal(body, &data); err != nil {
		logrus.Errorf("failed to parse LineageOS response for device %s: %v", device, err)
		return nil, false
	}
	if len(data.Response) == 0 {
		return nil, false
	}

	builds := make([]LineageBuild, 0, len(data.Response))
	for _, build := range data.Response {
		if build.URL == "" || build.Filename == "" {
			continue
		}
		builds = append(builds, LineageBuild{
			URL:      build.URL,
			Filename: build.Filename,
			Size:     formatSize(build.Size),
			Version:  build.Version,
			Datetime: time.Unix(build.Datetime, 0).Format("02/01/2006 15:04"),
			ROMType:  strings.ToUpper(build.ROMType),
		})
	}

	return builds, len(builds) > 0
}

func formatSize(size int64) string {
	if size <= 0 {
		return ""
	}
	units := []string{"B", "KB", "MB", "GB", "TB"}
	value := float64(size)
	unit := 0
	for value >= 1024 && unit < len(units)-1 {
		value /= 1024
		unit++
	}
	return strings.TrimRight(strings.TrimRight(formatFloat(value), "0"), ".") + " " + units[unit]
}

func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}
