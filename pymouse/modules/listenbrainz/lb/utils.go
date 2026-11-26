package lb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"pymouse/pymouse/helpers/rapidhttp"
	"pymouse/pymouse/modules/listenbrainz/set"
	"regexp"
	"strings"
)

// ListenBrainzTrackInformations mirrors just the fields the now-playing card
// and the recent-tracks list need. Unlike Last.fm, ListenBrainz exposes no
// per-track playcount or loved/feedback state through a plain, keyless GET,
// so the card simply doesn't render those rows for this service.
type ListenBrainzTrackInformations struct {
	Artist string
	Track  string
	Image  string
	Now    bool
}

type mbidMapping struct {
	CAAID          int64  `json:"caa_id"`
	CAAReleaseMBID string `json:"caa_release_mbid"`
}

type additionalInfo struct {
	OriginURL string `json:"origin_url"`
}

type trackMetadata struct {
	ArtistName     string         `json:"artist_name"`
	TrackName      string         `json:"track_name"`
	MBIDMapping    *mbidMapping   `json:"mbid_mapping"`
	AdditionalInfo additionalInfo `json:"additional_info"`
}

type listen struct {
	TrackMetadata trackMetadata `json:"track_metadata"`
}

type listensResponse struct {
	Payload struct {
		Listens []listen `json:"listens"`
	} `json:"payload"`
}

// coverArtURL builds a Cover Art Archive URL from a resolved mbid_mapping, or
// "" when ListenBrainz couldn't resolve this listen against MusicBrainz (very
// common for listens imported from streaming services without MBIDs). This is
// the cheapest option — the data already came back on the listen itself, no
// extra request — so it's tried first.
func coverArtURL(m *mbidMapping) string {
	if m == nil || m.CAAID == 0 || m.CAAReleaseMBID == "" {
		return ""
	}
	return fmt.Sprintf("https://coverartarchive.org/release/%s/%d-250.jpg", m.CAAReleaseMBID, m.CAAID)
}

// deezerSearchResponse is the (trimmed) shape of https://api.deezer.com/search.
type deezerSearchResponse struct {
	Data []struct {
		Album struct {
			CoverBig string `json:"cover_big"` // 500x500 square cover, no letterboxing
		} `json:"album"`
	} `json:"data"`
}

// deezerCoverURL looks up artist+track on Deezer's public, keyless search API
// and returns the matched release's square cover art, or "" on no match/error.
// This is the fallback for listens MusicBrainz has no mapping for — in
// practice most of the Brazilian funk/trap catalog YTM scrobbles, which
// Deezer's catalog covers far better than MusicBrainz does.
func deezerCoverURL(httpClient *http.Client, artist, track string) string {
	if artist == "" || track == "" {
		return ""
	}
	// Quoted artist:/track: filters do an exact(-ish) match instead of a loose
	// keyword search, which is what keeps this from returning a cover for the
	// wrong song entirely.
	q := fmt.Sprintf(`artist:"%s" track:"%s"`, deezerEscape(artist), deezerEscape(track))
	r, err := rapidhttp.Request(httpClient, rapidhttp.HTTPStruct{
		Method: "GET",
		URL:    "https://api.deezer.com/search",
		GETParams: &rapidhttp.HTTPGetStruct{
			Params: map[string]string{"q": q, "limit": "1"},
		},
	})
	if err != nil || r.StatusCode != 200 {
		return ""
	}
	defer r.Body.Close()

	var out deezerSearchResponse
	if err := json.NewDecoder(r.Body).Decode(&out); err != nil || len(out.Data) == 0 {
		return ""
	}
	return out.Data[0].Album.CoverBig
}

// deezerEscape drops double quotes from a search term so it can't break out
// of the artist:"..."/track:"..." filter it's embedded in.
func deezerEscape(s string) string {
	return strings.ReplaceAll(s, `"`, "'")
}

// youtubeVideoIDRe matches an 11-char YouTube video ID out of any watch URL
// shape ListenBrainz's YTM/YouTube scrobblers submit as origin_url:
// https://music.youtube.com/watch?v=ID, https://www.youtube.com/watch?v=ID,
// https://youtu.be/ID.
var youtubeVideoIDRe = regexp.MustCompile(`(?:[?&]v=|youtu\.be/)([\w-]{11})`)

// youtubeThumbnail returns the video's thumbnail URL when originURL looks
// like a YouTube/YouTube Music watch link, "" otherwise. Last-resort fallback
// only: YTM thumbnails are 16:9 with ugly letterboxing on the 600x250 card,
// so this only runs when neither MusicBrainz nor Deezer have the release.
func youtubeThumbnail(originURL string) string {
	m := youtubeVideoIDRe.FindStringSubmatch(originURL)
	if len(m) < 2 {
		return ""
	}
	return "https://img.youtube.com/vi/" + m[1] + "/hqdefault.jpg"
}

// resolveImage tries, in order of cost: the mbid_mapping ListenBrainz already
// resolved, a Deezer search, and finally the source YouTube video's thumbnail.
func resolveImage(httpClient *http.Client, t trackMetadata) string {
	if img := coverArtURL(t.MBIDMapping); img != "" {
		return img
	}
	if img := deezerCoverURL(httpClient, t.ArtistName, t.TrackName); img != "" {
		return img
	}
	return youtubeThumbnail(t.AdditionalInfo.OriginURL)
}

func toTrackInfo(httpClient *http.Client, l listen, now bool) ListenBrainzTrackInformations {
	return ListenBrainzTrackInformations{
		Artist: l.TrackMetadata.ArtistName,
		Track:  l.TrackMetadata.TrackName,
		Image:  resolveImage(httpClient, l.TrackMetadata),
		Now:    now,
	}
}

// fetchListens is the shared GET for both /playing-now and /listens: both
// return the same listensResponse shape (JSON Documentation — Fetching listen
// JSON). statusCode lets callers tell "user not found" (404) apart from "no
// listens yet" (200, empty payload).
func fetchListens(httpClient *http.Client, path string) (listensResponse, int, error) {
	var out listensResponse
	r, err := rapidhttp.Request(httpClient, rapidhttp.HTTPStruct{
		Method:  "GET",
		URL:     set.ListenBrainzAPIURL + path,
		Headers: map[string]string{"User-Agent": set.UserAgent},
	})
	if err != nil {
		return out, 0, fmt.Errorf("failed to reach ListenBrainz.")
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return out, r.StatusCode, nil
	}
	if err := json.NewDecoder(r.Body).Decode(&out); err != nil {
		return out, r.StatusCode, fmt.Errorf("failed to decode ListenBrainz response.")
	}
	return out, r.StatusCode, nil
}

// getRecentEntries returns up to `limit` tracks for username, most recent
// first. When the user is currently listening to something, that track is
// entry 0 with Now=true (ListenBrainz keeps "playing now" separate from
// historical listens, so this stitches the two together to mirror Last.fm's
// combined recenttracks list).
func getRecentEntries(httpClient *http.Client, username string, limit int) ([]ListenBrainzTrackInformations, error) {
	escaped := url.PathEscape(username)

	nowResp, status, err := fetchListens(httpClient, "user/"+escaped+"/playing-now")
	if err != nil {
		return nil, err
	}
	if status == 404 {
		return nil, fmt.Errorf("invalid username.")
	}

	out := make([]ListenBrainzTrackInformations, 0, limit)
	if status == 200 && len(nowResp.Payload.Listens) > 0 {
		out = append(out, toTrackInfo(httpClient, nowResp.Payload.Listens[0], true))
	}

	if remaining := limit - len(out); remaining > 0 {
		recentResp, status, err := fetchListens(httpClient, fmt.Sprintf("user/%s/listens?count=%d", escaped, remaining))
		if err != nil {
			return nil, err
		}
		if status == 404 {
			return nil, fmt.Errorf("invalid username.")
		}
		for _, l := range recentResp.Payload.Listens {
			out = append(out, toTrackInfo(httpClient, l, false))
		}
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("no listens found.")
	}
	if len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

// getTrack returns the current (or, failing that, the most recent) listen —
// the single entry /lb and /listenbrainz show.
func getTrack(httpClient *http.Client, username string) (ListenBrainzTrackInformations, error) {
	entries, err := getRecentEntries(httpClient, username, 1)
	if err != nil {
		return ListenBrainzTrackInformations{}, err
	}
	return entries[0], nil
}

func getListeningText(trackInfo ListenBrainzTrackInformations, l func(string) string) string {
	textKey := "listenbrainz.nowplaying.was-listening"
	if trackInfo.Now {
		textKey = "listenbrainz.nowplaying.is-listening"
	}
	return fmt.Sprintf("%s.", l(textKey))
}
