package lfm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pymouse/pymouse/config"
	"pymouse/pymouse/helpers/rapidhttp"
	"pymouse/pymouse/modules/lastfm/set"
	"strconv"
)

const (
	LastFMAPIURL = "http://ws.audioscrobbler.com/2.0/"
)

type trackPlaysInformations struct {
	Track struct {
		UserPlayCount string `json:"userplaycount"`
	} `json:"track"`
}

type LastFMTrackInformations struct {
	Artist    string
	Track     string
	Loved     bool
	Playcount int64
	Image     string
	Now       bool
	LastFMURL string // last.fm track page URL, used by the "Open on Last.fm" button
}

type LastFMRecentTracksInfo struct {
	RecentTracks struct {
		Track []struct {
			Artist struct {
				Name string `json:"name,omitempty" default:"unknown"`
			} `json:"artist"`

			Name  string `json:"name,omitempty" default:"unknown"`
			Loved string `json:"loved"`
			URL   string `json:"url,omitempty"`

			Image []struct {
				URL string `json:"#text"`
			} `json:"image"`

			Attr struct {
				NowPlaying string `json:"nowplaying"`
			} `json:"@attr,omitempty"`
		} `json:"track"`
	} `json:"recenttracks"`
}

// getRecentTracks fetches the user's last `limit` scrobbled tracks. Returns a
// slice of LastFMTrackInformations (Image/Playcount/YouTubeURL may be zero —
// this call is lighter than getTrack, which also hits track.getinfo per item).
func getRecentTracks(httpClient *http.Client, username string, limit int) ([]LastFMTrackInformations, error) {
	var info LastFMRecentTracksInfo
	params := map[string]string{
		"method":   "user.getrecenttracks",
		"user":     username,
		"api_key":  config.LastFMAPIKey,
		"format":   "json",
		"limit":    strconv.Itoa(limit),
		"extended": "1",
	}
	r, err := rapidhttp.Request(httpClient, rapidhttp.HTTPStruct{
		Method:    "GET",
		URL:       LastFMAPIURL,
		GETParams: &rapidhttp.HTTPGetStruct{Params: params},
	})
	if err != nil || r.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch recent tracks.")
	}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode recent tracks information.")
	}
	out := make([]LastFMTrackInformations, 0, limit)
	for _, t := range info.RecentTracks.Track {
		if len(out) >= limit {
			break // last.fm sometimes returns limit+1 when a track is now-playing
		}
		img := ""
		if n := len(t.Image); n > 0 {
			img = t.Image[n-1].URL
		}
		out = append(out, LastFMTrackInformations{
			Artist: t.Artist.Name,
			Track:  t.Name,
			Loved:  t.Loved == "1",
			Image:  img,
			Now:    t.Attr.NowPlaying == "true",
		})
	}
	return out, nil
}

func trackPlays(httpClient *http.Client, username string, artist string, track string) (int64, error) {
	var trackPlaysInfo trackPlaysInformations
	trackPlaysParams := map[string]string{
		"method":  "track.getinfo",
		"artist":  artist,
		"track":   track,
		"user":    username,
		"api_key": config.LastFMAPIKey,
		"format":  "json",
	}
	r, err := rapidhttp.Request(
		httpClient,
		rapidhttp.HTTPStruct{
			Method: "GET",
			URL:    LastFMAPIURL,
			GETParams: &rapidhttp.HTTPGetStruct{
				Params: trackPlaysParams,
			},
		},
	)
	if err != nil || r.StatusCode != 200 {
		fmt.Printf("Failed to fetch track plays: %v\n", err)
		return 0, fmt.Errorf("failed to fetch track plays.")
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&trackPlaysInfo); err != nil {
		return 0, fmt.Errorf("failed to decode track plays information.")
	}
	if trackPlaysInfo.Track.UserPlayCount == "" {
		return 0, nil
	}
	playCount, err := strconv.ParseInt(trackPlaysInfo.Track.UserPlayCount, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse track plays information.")
	}

	return playCount, nil
}

func getTrack(httpClient *http.Client, username string) (LastFMTrackInformations, error) {
	var recentTracksInfo LastFMRecentTracksInfo
	if !set.CheckUsername(httpClient, username) {
		return LastFMTrackInformations{}, fmt.Errorf("invalid username.")
	}
	recentTracksParams := map[string]string{
		"method":   "user.getrecenttracks",
		"user":     username,
		"api_key":  config.LastFMAPIKey,
		"format":   "json",
		"limit":    "1",
		"extended": "1",
	}
	r, err := rapidhttp.Request(
		httpClient,
		rapidhttp.HTTPStruct{
			Method: "GET",
			URL:    LastFMAPIURL,
			GETParams: &rapidhttp.HTTPGetStruct{
				Params: recentTracksParams,
			},
		},
	)
	if err != nil || r.StatusCode != 200 {
		return LastFMTrackInformations{}, fmt.Errorf("failed to fetch recent tracks.")
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&recentTracksInfo); err != nil {
		return LastFMTrackInformations{}, fmt.Errorf("failed to decode recent tracks information.")
	}

	if len(recentTracksInfo.RecentTracks.Track) == 0 {
		return LastFMTrackInformations{}, fmt.Errorf("no scrobbles found.")
	}
	userPlayCount, err := trackPlays(
		httpClient,
		username,
		recentTracksInfo.RecentTracks.Track[0].Artist.Name,
		recentTracksInfo.RecentTracks.Track[0].Name,
	)
	if err != nil {
		return LastFMTrackInformations{}, fmt.Errorf("failed to fetch track plays.")
	}
	defer r.Body.Close()

	return LastFMTrackInformations{
		Artist:    recentTracksInfo.RecentTracks.Track[0].Artist.Name,
		Track:     recentTracksInfo.RecentTracks.Track[0].Name,
		Loved:     recentTracksInfo.RecentTracks.Track[0].Loved == "1",
		Playcount: userPlayCount,
		Image:     recentTracksInfo.RecentTracks.Track[0].Image[len(recentTracksInfo.RecentTracks.Track[0].Image)-1].URL,
		Now:       recentTracksInfo.RecentTracks.Track[0].Attr.NowPlaying == "true",
		LastFMURL: recentTracksInfo.RecentTracks.Track[0].URL,
	}, nil
}

func getListeningText(trackInfo LastFMTrackInformations, l func(string) string) string {
	textKey := "lastfm.nowplaying.was-listening"
	if trackInfo.Now {
		textKey = "lastfm.nowplaying.is-listening"
	}
	listeningText := l(textKey)
	if trackInfo.Playcount > 0 {
		listeningText += fmt.Sprintf(l("lastfm.nowplaying.userplaycount"), trackInfo.Playcount)
	}
	return fmt.Sprintf("%s.", listeningText)
}
