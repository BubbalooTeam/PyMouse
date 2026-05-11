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
}

type LastFMRecentTracksInfo struct {
	RecentTracks struct {
		Track []struct {
			Artist struct {
				Name string `json:"name"`
			} `json:"artist"`

			Name  string `json:"name"`
			Loved string `json:"loved"`

			Image []struct {
				URL string `json:"#text"`
			} `json:"image"`

			Attr struct {
				NowPlaying string `json:"nowplaying"`
			} `json:"@attr,omitempty"`
		} `json:"track"`
	} `json:"recenttracks"`
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
			URL:    "http://ws.audioscrobbler.com/2.0/",
			GETParams: &rapidhttp.HTTPGetStruct{
				Params: trackPlaysParams,
			},
		},
	)
	if err != nil || r.StatusCode != 200 {
		return 0, fmt.Errorf("failed to fetch track plays.")
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&trackPlaysInfo); err != nil {
		return 0, fmt.Errorf("failed to decode track plays information.")
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
			URL:    "http://ws.audioscrobbler.com/2.0/",
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
	}, nil
}
