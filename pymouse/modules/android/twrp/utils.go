package twrp

import (
	"fmt"
	"net/http"
	"pymouse/pymouse/helpers/rapidhttp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type TWRPBaseResult struct {
	Updated     string
	FileName    string
	FileSize    string
	DownloadURL string
}

type twrpResultsCache struct {
	Device  string
	OwnerID int64
	Results []TWRPBaseResult
}

var twrpCache = cache.New(15*time.Minute, 30*time.Minute)

func setTWRPResults(key string, device string, ownerID int64, results []TWRPBaseResult) {
	twrpCache.Set(key, twrpResultsCache{
		Device:  device,
		OwnerID: ownerID,
		Results: results,
	}, cache.DefaultExpiration)
}

func getTWRPResults(key string) (twrpResultsCache, bool) {
	value, found := twrpCache.Get(key)
	if !found {
		return twrpResultsCache{}, false
	}

	results, ok := value.(twrpResultsCache)
	return results, ok
}

func getTWRP(codename string) ([]TWRPBaseResult, bool) {
	httpClient := rapidhttp.GetHTTPClient()

	url := fmt.Sprintf("https://eu.dl.twrp.me/%s/", codename)
	resp, err := rapidhttp.Request(
		httpClient,
		rapidhttp.HTTPStruct{
			Method: "GET",
			URL:    url,
		},
	)
	if err != nil {
		logrus.Errorf("failed to fetch twrp for device: %s", codename)
		return nil, false
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		logrus.Errorf("failed to fetch twrp for device: %s (status: %s)", codename, resp.Status)
		return nil, false
	}

	page, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logrus.Errorf("failed to parse twrp page for device: %s: %v", codename, err)
		return nil, false
	}

	updated := strings.TrimSpace(page.Find("em").First().Text())
	rows := page.Find("table").First().Find("tr")
	if updated == "" || rows.Length() == 0 {
		return nil, false
	}

	results := make([]TWRPBaseResult, 0, rows.Length())
	rows.Each(func(_ int, selected *goquery.Selection) {
		download := selected.Find("a").First()
		href, exists := download.Attr("href")
		fileSize := strings.TrimSpace(selected.Find("span.filesize").First().Text())
		if download.Length() == 0 || !exists || href == "" || fileSize == "" {
			return
		}

		results = append(results, TWRPBaseResult{
			Updated:     updated,
			FileName:    strings.ToLower(strings.TrimSpace(download.Text())),
			FileSize:    fileSize,
			DownloadURL: fmt.Sprintf("https://dl.twrp.me%s", href),
		})
	})

	if len(results) == 0 {
		return nil, false
	}
	return results, true
}
