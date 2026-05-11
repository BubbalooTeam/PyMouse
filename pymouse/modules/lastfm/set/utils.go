package set

import (
	"net/http"
	"pymouse/pymouse/config"
	"pymouse/pymouse/helpers/rapidhttp"
	"strings"
)

func normalizeUsername(username string) string {
	if strings.Contains(username, "@") {
		return strings.Replace(username, "@", "", 1)
	}
	return username
}

func CheckUsername(httpClient *http.Client, username string) bool {
	userLfmParams := map[string]string{
		"method":  "user.getinfo",
		"user":    username,
		"api_key": config.LastFMAPIKey,
		"format":  "json",
	}
	r, err := rapidhttp.Request(
		httpClient,
		rapidhttp.HTTPStruct{
			Method: "GET",
			URL:    "http://ws.audioscrobbler.com/2.0/",
			GETParams: &rapidhttp.HTTPGetStruct{
				Params: userLfmParams,
			},
		},
	)
	if err != nil || r.StatusCode != 200 {
		return false
	}
	return true
}
