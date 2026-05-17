package lfm

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"net/http"
	"os"
	"pymouse/pymouse/config"
	"pymouse/pymouse/helpers/rapidhttp"
	"pymouse/pymouse/modules/lastfm/set"
	"strconv"

	"github.com/cavaliergopher/grab/v3"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/google/uuid"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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
	b, err := os.ReadFile(path)
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

	// fonts
	openSans := loadFont(fonts.OpenSans, 19)
	poppins := loadFont(fonts.Poppins, 18)
	arial := loadFont(fonts.Arial, 21)
	arial23 := loadFont(fonts.Arial, 17)

	songFont := poppins
	if !checkUnicode(songName) {
		songFont = arial
	}

	artistFont := openSans
	if !checkUnicode(artistName) {
		artistFont = arial23
	}

	dc.SetColor(color.White)

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
		if heart, err := imaging.Open("pymouse/assets/icons/lastfm/loved.png"); err == nil {
			heart = imaging.Resize(heart, 25, 25, imaging.Lanczos)
			dc.DrawImage(heart, 248, 190)
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
