package sf

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"pymouse/pymouse/helpers/rapidhttp"
	"pymouse/pymouse/modules/statsfm/set"
	"strings"
)

// StatsFMTrackInformations mirrors just the fields the now-playing card and
// the recent-tracks list need. Unlike ListenBrainz, stats.fm always has a
// proper album cover (synced straight from Spotify/Apple Music) and an exact
// lifetime playcount per track — no cover-art fallback chain needed here.
type StatsFMTrackInformations struct {
	Artist    string
	Track     string
	Image     string
	Now       bool
	TrackID   int64 // stats.fm's own numeric track id, used to look up Playcount
	Playcount int64
}

type sfmArtist struct {
	Name string `json:"name"`
}

type sfmAlbum struct {
	Image string `json:"image"`
}

type sfmTrack struct {
	ID      int64       `json:"id"`
	Name    string      `json:"name"`
	Artists []sfmArtist `json:"artists"`
	Albums  []sfmAlbum  `json:"albums"`
}

type sfmStream struct {
	Track sfmTrack `json:"track"`
}

type sfmCurrentResponse struct {
	Item *sfmStream `json:"item"`
}

type sfmRecentResponse struct {
	Items []sfmStream `json:"items"`
}

type sfmTopTracksResponse struct {
	Items []struct {
		Streams int64    `json:"streams"`
		Track   sfmTrack `json:"track"`
	} `json:"items"`
}

func artistNames(t sfmTrack) string {
	names := make([]string, 0, len(t.Artists))
	for _, a := range t.Artists {
		names = append(names, a.Name)
	}
	return strings.Join(names, ", ")
}

func albumImage(t sfmTrack) string {
	if len(t.Albums) == 0 {
		return ""
	}
	return t.Albums[0].Image
}

func toTrackInfo(s sfmStream, now bool) StatsFMTrackInformations {
	return StatsFMTrackInformations{
		Artist:  artistNames(s.Track),
		Track:   s.Track.Name,
		Image:   albumImage(s.Track),
		Now:     now,
		TrackID: s.Track.ID,
	}
}

// fetchJSON is the shared GET + JSON-decode for every stats.fm endpoint used
// here. statusCode lets callers tell "user not found" (404) apart from a
// user with privacy/sync settings hiding that particular section (also a
// non-200, handled the same way: treat as empty).
func fetchJSON(httpClient *http.Client, path string, params map[string]string, out any) (int, error) {
	r, err := rapidhttp.Request(httpClient, rapidhttp.HTTPStruct{
		Method:    "GET",
		URL:       set.StatsFMAPIURL + path,
		GETParams: &rapidhttp.HTTPGetStruct{Params: params},
		Headers:   map[string]string{"User-Agent": set.UserAgent},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to reach stats.fm.")
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return r.StatusCode, nil
	}
	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		return r.StatusCode, fmt.Errorf("failed to decode stats.fm response.")
	}
	return r.StatusCode, nil
}

// getRecentEntries returns up to `limit` tracks for username, most recent
// first. When the user is currently playing something, that track is entry 0
// with Now=true (stats.fm keeps "currently playing" separate from the
// recent-streams history, same split ListenBrainz has).
func getRecentEntries(httpClient *http.Client, username string, limit int) ([]StatsFMTrackInformations, error) {
	escaped := url.PathEscape(username)

	var cur sfmCurrentResponse
	status, err := fetchJSON(httpClient, "users/"+escaped+"/streams/current", nil, &cur)
	if err != nil {
		return nil, err
	}
	if status == 404 {
		return nil, fmt.Errorf("invalid username.")
	}

	out := make([]StatsFMTrackInformations, 0, limit)
	if status == 200 && cur.Item != nil {
		out = append(out, toTrackInfo(*cur.Item, true))
	}

	if remaining := limit - len(out); remaining > 0 {
		var recent sfmRecentResponse
		status, err := fetchJSON(httpClient, "users/"+escaped+"/streams/recent", nil, &recent)
		if err != nil {
			return nil, err
		}
		if status == 404 {
			return nil, fmt.Errorf("invalid username.")
		}
		for _, s := range recent.Items {
			if len(out) >= limit {
				break
			}
			out = append(out, toTrackInfo(s, false))
		}
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("no listens found.")
	}
	return out, nil
}

// getTrackPlaycount looks up trackID's lifetime stream count among the
// user's 250 most-played tracks. Matched by stats.fm's own numeric track id
// (exact, unlike the fuzzy artist/track-name matching ListenBrainz needs) —
// "" / not found just means it's outside their top 250, handled like a
// Last.fm playcount of 0 (the row is omitted, not an error).
func getTrackPlaycount(httpClient *http.Client, username string, trackID int64) (int64, bool) {
	if trackID == 0 {
		return 0, false
	}
	var top sfmTopTracksResponse
	status, err := fetchJSON(
		httpClient,
		"users/"+url.PathEscape(username)+"/top/tracks",
		map[string]string{"range": "lifetime", "limit": "250"},
		&top,
	)
	if err != nil || status != 200 {
		return 0, false
	}
	for _, item := range top.Items {
		if item.Track.ID == trackID {
			return item.Streams, true
		}
	}
	return 0, false
}

// getTrack returns the current (or, failing that, the most recent) stream —
// the single entry /sf and /statsfm show — enriched with its lifetime
// playcount when the track is popular enough to be in the user's top 250.
func getTrack(httpClient *http.Client, username string) (StatsFMTrackInformations, error) {
	entries, err := getRecentEntries(httpClient, username, 1)
	if err != nil {
		return StatsFMTrackInformations{}, err
	}
	track := entries[0]
	if count, ok := getTrackPlaycount(httpClient, username, track.TrackID); ok {
		track.Playcount = count
	}
	return track, nil
}

func getListeningText(trackInfo StatsFMTrackInformations, l func(string) string) string {
	textKey := "statsfm.nowplaying.was-listening"
	if trackInfo.Now {
		textKey = "statsfm.nowplaying.is-listening"
	}
	listeningText := l(textKey)
	if trackInfo.Playcount > 0 {
		listeningText += fmt.Sprintf(l("statsfm.nowplaying.userplaycount"), trackInfo.Playcount)
	}
	return fmt.Sprintf("%s.", listeningText)
}
