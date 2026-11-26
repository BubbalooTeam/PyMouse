package set

import (
	"net/http"
	"net/url"
	"pymouse/pymouse/helpers/rapidhttp"
	"strings"
)

// StatsFMAPIURL is stats.fm's API, reverse-engineered from their own web app
// (no official public docs/SLA). It's unauthenticated for read-only profile
// data — no API key needed, same as ListenBrainz.
const StatsFMAPIURL = "https://api-v2.stats.fm/api/v1/"

// UserAgent is sent on every stats.fm request to identify the bot, since
// there's no official client library/etiquette doc to follow.
const UserAgent = "PyMouse-Bot (+https://github.com/BubbalooTeam/PyMouse)"

func normalizeUsername(username string) string {
	if strings.Contains(username, "@") {
		return strings.Replace(username, "@", "", 1)
	}
	return username
}

// CheckUsername reports whether username is a registered stats.fm custom ID
// (what a user sets under Settings -> Profile -> Custom ID; falls back to
// their opaque Spotify ID otherwise, which we don't support entering by hand).
func CheckUsername(httpClient *http.Client, username string) bool {
	r, err := rapidhttp.Request(
		httpClient,
		rapidhttp.HTTPStruct{
			Method:  "GET",
			URL:     StatsFMAPIURL + "users/" + url.PathEscape(username),
			Headers: map[string]string{"User-Agent": UserAgent},
		},
	)
	if err != nil || r.StatusCode != 200 {
		return false
	}
	return true
}
