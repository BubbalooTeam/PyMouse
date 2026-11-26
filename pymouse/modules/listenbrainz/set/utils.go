package set

import (
	"net/http"
	"net/url"
	"pymouse/pymouse/helpers/rapidhttp"
	"strings"
)

// ListenBrainzAPIURL is the public, keyless ListenBrainz API root. Unlike
// Last.fm, reading listens/user data requires no API key — only submitting
// listens does (see https://listenbrainz.readthedocs.io/en/latest/users/api/core.html).
const ListenBrainzAPIURL = "https://api.listenbrainz.org/1/"

// UserAgent is sent on every ListenBrainz API request, as asked by their API
// etiquette (https://listenbrainz.readthedocs.io/en/latest/).
const UserAgent = "PyMouse-Bot (+https://github.com/BubbalooTeam/PyMouse)"

func normalizeUsername(username string) string {
	if strings.Contains(username, "@") {
		return strings.Replace(username, "@", "", 1)
	}
	return username
}

// CheckUsername reports whether username is a registered ListenBrainz (i.e.
// MusicBrainz) account, using the lightweight listen-count endpoint: 200 means
// the user exists, 404 means it doesn't.
func CheckUsername(httpClient *http.Client, username string) bool {
	r, err := rapidhttp.Request(
		httpClient,
		rapidhttp.HTTPStruct{
			Method:  "GET",
			URL:     ListenBrainzAPIURL + "user/" + url.PathEscape(username) + "/listen-count",
			Headers: map[string]string{"User-Agent": UserAgent},
		},
	)
	if err != nil || r.StatusCode != 200 {
		return false
	}
	return true
}
