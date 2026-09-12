package magisk

import (
	"encoding/json"
	"io"
	"pymouse/pymouse/helpers/rapidhttp"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type MagiskRelease struct {
	Variant      string
	DownloadURL  string
	ChangelogURL string
	Version      string
	VersionCode  string
}

type magiskResponse struct {
	Magisk struct {
		Link        string `json:"link"`
		Note        string `json:"note"`
		Version     string `json:"version"`
		VersionCode string `json:"versionCode"`
	} `json:"magisk"`
}

type magiskResultsCache struct {
	OwnerID int64
	Results []MagiskRelease
}

var magiskCache = cache.New(15*time.Minute, 30*time.Minute)

func setMagiskResults(key string, ownerID int64, results []MagiskRelease) {
	magiskCache.Set(key, magiskResultsCache{
		OwnerID: ownerID,
		Results: results,
	}, cache.DefaultExpiration)
}

func getMagiskResults(key string) (magiskResultsCache, bool) {
	value, found := magiskCache.Get(key)
	if !found {
		return magiskResultsCache{}, false
	}

	results, ok := value.(magiskResultsCache)
	return results, ok
}

func getMagisk() ([]MagiskRelease, bool) {
	const baseURL = "https://raw.githubusercontent.com/topjohnwu/magisk-files/master/"
	variants := []struct {
		name string
		path string
	}{
		{name: "stable", path: "stable.json"},
		{name: "beta", path: "beta.json"},
		{name: "canary", path: "canary.json"},
	}

	client := rapidhttp.GetHTTPClient()
	results := make([]MagiskRelease, 0, len(variants))
	for _, variant := range variants {
		resp, err := rapidhttp.Request(client, rapidhttp.HTTPStruct{
			Method: "GET",
			URL:    baseURL + variant.path,
		})
		if err != nil {
			logrus.Errorf("failed to fetch Magisk %s release: %v", variant.name, err)
			return nil, false
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
			logrus.Errorf("failed to fetch Magisk %s release: status=%s error=%v", variant.name, resp.Status, readErr)
			return nil, false
		}

		var data magiskResponse
		if err := json.Unmarshal(body, &data); err != nil {
			logrus.Errorf("failed to parse Magisk %s release: %v", variant.name, err)
			return nil, false
		}
		if data.Magisk.Link == "" || data.Magisk.Note == "" || data.Magisk.Version == "" {
			logrus.Errorf("Magisk %s release has incomplete data", variant.name)
			return nil, false
		}

		results = append(results, MagiskRelease{
			Variant:      variant.name,
			DownloadURL:  data.Magisk.Link,
			ChangelogURL: data.Magisk.Note,
			Version:      data.Magisk.Version,
			VersionCode:  data.Magisk.VersionCode,
		})
	}

	return results, true
}
