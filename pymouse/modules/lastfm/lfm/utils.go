package lfm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"net/http"
	"net/url"
	"os"
	"pymouse/pymouse/assets"
	"pymouse/pymouse/config"
	"pymouse/pymouse/helpers/rapidhttp"
	"pymouse/pymouse/modules/lastfm/set"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/google/uuid"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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
	Artist     string
	Track      string
	Loved      bool
	Playcount  int64
	Image      string
	Now        bool
	YouTubeURL string // direct watch?v=... link when resolvable; falls back to a search URL
}

type LastFMRecentTracksInfo struct {
	RecentTracks struct {
		Track []struct {
			Artist struct {
				Name string `json:"name,omitempty" default:"unknown"`
			} `json:"artist"`

			Name  string `json:"name,omitempty" default:"unknown"`
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

type Fonts struct {
	OpenSans string
	Poppins  string
	Arial    string
	// Unicode covers Latin/Cyrillic/Greek (used when text is non-ASCII
	// but contains no CJK characters).
	Unicode string
	// CJK covers Chinese/Japanese/Korean (and also Latin/Cyrillic,
	// so it is used alone when any CJK rune is present).
	CJK string
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
		Artist:     recentTracksInfo.RecentTracks.Track[0].Artist.Name,
		Track:      recentTracksInfo.RecentTracks.Track[0].Name,
		Loved:      recentTracksInfo.RecentTracks.Track[0].Loved == "1",
		Playcount:  userPlayCount,
		Image:      recentTracksInfo.RecentTracks.Track[0].Image[len(recentTracksInfo.RecentTracks.Track[0].Image)-1].URL,
		Now:        recentTracksInfo.RecentTracks.Track[0].Attr.NowPlaying == "true",
		YouTubeURL: resolveYouTubeURL(httpClient, recentTracksInfo.RecentTracks.Track[0].Artist.Name, recentTracksInfo.RecentTracks.Track[0].Name),
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

func loadFont(path string, size float64) font.Face {
	b, err := assets.ReadFile(path)
	if err != nil {
		panic(err)
	}

	ft, err := opentype.Parse(b)
	if err != nil {
		panic(err)
	}

	face, err := opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	if err != nil {
		panic(err)
	}

	return face
}

// pickTextColor returns white or black, whichever contrasts best with the
// mean luminance of the given region on the already-rendered background.
// Used so the track/artist copy stays readable on any cover — including the
// white Last.fm "no artwork" placeholder. Sampling is strided (every 2px);
// the region is small (~340x210).
func pickTextColor(img image.Image, region image.Rectangle) color.Color {
	bounds := img.Bounds()
	if !bounds.Intersect(region).Eq(region) {
		return color.White // defensive default on unexpected geometry
	}

	var sum uint64
	var n uint64
	for y := region.Min.Y; y < region.Max.Y; y += 2 {
		for x := region.Min.X; x < region.Max.X; x += 2 {
			r, g, b, _ := img.At(x, y).RGBA()
			// Rec. 601 luma, 8-bit range
			l := (299*uint64(r>>8) + 587*uint64(g>>8) + 114*uint64(b>>8)) / 1000
			sum += l
			n++
		}
	}
	if n == 0 {
		return color.White
	}

	luma := float64(sum) / float64(n) // 0..255
	if luma >= 128 {
		return color.Black // bright background -> dark text
	}
	return color.White // dark background -> light text
}

func darken(img image.Image, factor float64) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			dst.Set(x, y, color.RGBA{
				R: uint8(float64(r>>8) * factor),
				G: uint8(float64(g>>8) * factor),
				B: uint8(float64(b>>8) * factor),
				A: uint8(a >> 8),
			})
		}
	}

	return dst
}

func truncate(dc *gg.Context, text string, max float64) string {
	w, _ := dc.MeasureString(text)

	if w <= max {
		return text
	}

	runes := []rune(text)

	for len(runes) > 0 {
		runes = runes[:len(runes)-1]

		s := string(runes) + "..."

		w, _ := dc.MeasureString(s)

		if w <= max {
			return s
		}
	}

	return text
}

func checkUnicode(s string) bool {
	for _, r := range s {
		if r > 127 {
			return true
		}
	}

	return false
}

// containsCJK reports whether s contains any CJK ideograph, hiragana,
// katakana, hangul or CJK punctuation. These blocks are only covered by a
// dedicated CJK font (e.g. Noto Sans CJK), not by Latin/Cyrillic fonts.
func containsCJK(s string) bool {
	for _, r := range s {
		switch {
		case r >= 0x2E80 && r <= 0x9FFF: // CJK radicals + Kangxi + CJK ideographs
			return true
		case r >= 0xAC00 && r <= 0xD7AF: // Hangul Syllables
			return true
		case r >= 0xF900 && r <= 0xFAFF: // CJK Compatibility Ideographs
			return true
		case r >= 0xFF00 && r <= 0xFFEF: // Halfwidth/Fullwidth + Katakana
			return true
		case r >= 0x3000 && r <= 0x30FF: // CJK symbols + Hiragana + Katakana
			return true
		case r >= 0x3041 && r <= 0x309F: // Hiragana (redundant safety)
			return true
		case r >= 0xFF66 && r <= 0xFF9F: // Halfwidth Katakana
			return true
		}
	}

	return false
}

func DrawScrobble(
	glabClient *grab.Client,
	imgURL string,
	songName string,
	artistName string,
	userLastFM string,
	listening string,
	loved bool,
	fonts Fonts,
	l func(string) string,
) (string, error) {
	dir := fmt.Sprintf("%s/%s", config.DownloadPath, "lastfm")
	dc := gg.NewContext(600, 250)

	// background
	dc.SetRGB255(18, 18, 18)
	dc.Clear()

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	req, err := grab.NewRequest(dir, imgURL)
	if err != nil {
		fmt.Println("Failed to create download request for lastfm image.")
		return "", err
	}
	resp := glabClient.Do(req)
	if resp == nil || resp.Err() != nil {
		fmt.Printf("Failed to download lastfm image: %v\n", resp.Err())
		return "", err
	}
	imgPath := resp.Filename

	// album art
	if img, err := imaging.Open(imgPath); err == nil {

		// blurred bg
		bg := darken(img, 0.7)
		bg = imaging.Blur(bg, 20)
		bg = imaging.Resize(bg, 600, 600, imaging.Lanczos)

		dc.DrawImage(bg, 0, -250)

		// original cover
		cover := imaging.Resize(img, 200, 200, imaging.Lanczos)
		dc.DrawImage(cover, 25, 25)
	}

	// Text sits in the right-hand half of the card. On bright covers (or the
	// white Last.fm "no artwork" placeholder) the blurred background there can
	// be too light for the default white text. Pick the text colour based on
	// the mean luminance of the exact region the text will occupy: dark
	// background -> white text, bright background -> black text.
	// Geometry: every text element starts at x=248 and is capped at ~315px
	// wide; vertically they span y=25 (username top) to y=225 (loved bottom).
	textColor := pickTextColor(dc.Image(), image.Rect(238, 18, 573, 230))

	// fonts
	openSans := loadFont(fonts.OpenSans, 19)
	poppins := loadFont(fonts.Poppins, 18)
	arial := loadFont(fonts.Arial, 21)
	arial23 := loadFont(fonts.Arial, 17)

	// Pan-Unicode fallbacks for non-Latin scripts.
	// The Latin-only fonts (Poppins/OpenSans/Arial bundled here) do not
	// carry Cyrillic, Greek or CJK glyphs, so those runes render as
	// nothing. We swap in a font that actually covers the script:
	//   - any CJK rune  -> CJK font (also covers Latin/Cyrillic, so a
	//     mixed line like "愛 Love" still renders fully);
	//   - other non-ASCII (Cyrillic/Greek/etc.) -> Unicode font;
	//   - pure ASCII -> keep the original Latin font.
	var unicodeFace, unicodeFaceSm, cjkFace, cjkFaceSm font.Face
	if fonts.Unicode != "" {
		unicodeFace = loadFont(fonts.Unicode, 21)
		unicodeFaceSm = loadFont(fonts.Unicode, 17)
	}
	if fonts.CJK != "" {
		cjkFace = loadFont(fonts.CJK, 21)
		cjkFaceSm = loadFont(fonts.CJK, 17)
	}

	songFont := arial
	switch {
	case containsCJK(songName) && cjkFace != nil:
		songFont = cjkFace
	case checkUnicode(songName) && unicodeFace != nil:
		songFont = unicodeFace
	}

	artistFont := arial23
	switch {
	case containsCJK(artistName) && cjkFaceSm != nil:
		artistFont = cjkFaceSm
	case checkUnicode(artistName) && unicodeFaceSm != nil:
		artistFont = unicodeFaceSm
	}

	dc.SetColor(textColor)

	// username
	dc.SetFontFace(poppins)
	dc.DrawString(
		truncate(dc, userLastFM, 250),
		248,
		42,
	)

	// listening
	dc.SetFontFace(openSans)
	dc.DrawString(
		listening,
		248,
		77,
	)

	// song
	dc.SetFontFace(songFont)
	dc.DrawString(
		truncate(dc, songName, 315),
		248,
		140,
	)

	// artist
	dc.SetFontFace(artistFont)
	dc.DrawString(
		truncate(dc, artistName, 315),
		248,
		175,
	)

	// loved
	if loved {
		if heartBytes, err := assets.ReadFile("icons/lastfm/loved.png"); err == nil {
			if heart, derr := imaging.Decode(bytes.NewReader(heartBytes)); derr == nil {
				heart = imaging.Resize(heart, 25, 25, imaging.Lanczos)
				dc.DrawImage(heart, 248, 190)
			}
		}

		dc.SetFontFace(arial23)
		dc.DrawString(
			l("lastfm.nowplaying.loved"),
			278,
			210,
		)
	}

	// save
	filename := uuid.NewString() + ".jpg"

	out, err := os.Create(filename)
	if err != nil {
		return "", err
	}

	defer out.Close()

	err = jpeg.Encode(out, dc.Image(), &jpeg.Options{
		Quality: 95,
	})

	if err != nil {
		return "", err
	}

	defer os.Remove(imgPath)

	return filename, nil
}

// ytWatchRe matches the first youtube.com/watch?v=XXXX or youtu.be/XXXX link
// embedded in the Last.fm track page HTML ("Play track" button).
var ytWatchRe = regexp.MustCompile(`(?:https?://(?:www\.)?youtube\.com/watch\?v=|https?://youtu\.be/)([A-Za-z0-9_-]{11})`)

// resolveYouTubeURL tries to find a DIRECT YouTube link for a track, in order:
//  1. Scraping the Last.fm track page (last.fm/music/<artist>/_/<track>),
//     which exposes a curated "Play on YouTube" button for most tracks.
//  2. (skipped) yt-dlp search — too slow (seconds + needs node) for a /lfm
//     request; revisit only if scraping coverage proves insufficient.
//  3. Fallback: a YouTube search URL (artist + track). Never fails.
//
// The whole resolution runs with a short timeout so /lfm never blocks on it.
func resolveYouTubeURL(httpClient *http.Client, artist, track string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	// Layer 1: scrape the Last.fm track page.
	pageURL := fmt.Sprintf(
		"https://www.last.fm/music/%s/_/%s",
		url.PathEscape(artist),
		url.PathEscape(track),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, nil)
	if err == nil {
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; PyMouseBot/1.0)")
		if resp, rerr := httpClient.Do(req); rerr == nil {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20)) // 2 MiB cap
			resp.Body.Close()
			if m := ytWatchRe.FindStringSubmatch(string(body)); len(m) == 2 {
				return "https://www.youtube.com/watch?v=" + m[1]
			}
		}
	}

	// Layer 3 (fallback): a search URL. Always works.
	q := url.QueryEscape(strings.TrimSpace(artist + " " + track))
	return "https://www.youtube.com/results?search_query=" + q
}

