// Package scrobblecard renders the "now playing" / "recent tracks" card image
// shared by the scrobbling integrations (Last.fm, ListenBrainz, ...). It knows
// nothing about any particular service's API — callers fetch the track data
// themselves and pass already-localized strings in.
package scrobblecard

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"strconv"

	"pymouse/pymouse/assets"
	"pymouse/pymouse/config"

	"github.com/cavaliergopher/grab/v3"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/google/uuid"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// Fonts holds the bundled font paths (relative to pymouse/assets) used to
// render a card. Unicode/CJK are fallbacks for scripts the Latin fonts don't
// cover; either may be left empty to skip that fallback.
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

// RecentEntry is one row of a DrawRecentScrobble list.
type RecentEntry struct {
	Track  string
	Artist string
	Now    bool // currently playing -> labeled with nowLabel instead of a number
	Loved  bool // prefixed with a heart glyph
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
// Used so the track/artist copy stays readable on any cover — including a
// bright "no artwork" placeholder. Sampling is strided (every 2px); the
// region is small (~340x210).
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

// DrawScrobble renders the 600x250 now-playing card: cover art (or a blurred
// version of it as background), username, listening status, song/artist, and
// an optional "loved" row.
//
// downloadDir is the directory the cover art gets downloaded into (a
// service-specific subdirectory of config.DownloadPath, e.g. "lastfm" or
// "listenbrainz"). heartIconPath/lovedText are both optional — leave them
// empty to render the card without a loved row (used by services that don't
// expose per-track love/feedback state).
func DrawScrobble(
	glabClient *grab.Client,
	downloadDir string,
	imgURL string,
	songName string,
	artistName string,
	username string,
	listening string,
	loved bool,
	heartIconPath string,
	lovedText string,
	fonts Fonts,
) (string, error) {
	dc := gg.NewContext(600, 250)

	// background
	dc.SetRGB255(18, 18, 18)
	dc.Clear()

	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return "", err
	}

	req, err := grab.NewRequest(downloadDir, imgURL)
	if err != nil {
		fmt.Println("Failed to create download request for scrobble image.")
		return "", err
	}
	resp := glabClient.Do(req)
	if resp == nil || resp.Err() != nil {
		fmt.Printf("Failed to download scrobble image: %v\n", resp.Err())
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
	// default "no artwork" placeholder) the blurred background there can be
	// too light for the default white text. Pick the text colour based on
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
		truncate(dc, username, 250),
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
	if loved && heartIconPath != "" {
		if heartBytes, err := assets.ReadFile(heartIconPath); err == nil {
			if heart, derr := imaging.Decode(bytes.NewReader(heartBytes)); derr == nil {
				heart = imaging.Resize(heart, 25, 25, imaging.Lanczos)
				dc.DrawImage(heart, 248, 190)
			}
		}

		dc.SetFontFace(arial23)
		dc.DrawString(
			lovedText,
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

// DrawRecentScrobble renders a vertical list of recent tracks for the "+"
// expand view. Layout: same dark background as the now-playing card, with each
// entry on its own line ("1. track — artist"), the first one flagged as
// now-playing if applicable. Height grows with the number of entries.
func DrawRecentScrobble(
	downloadDir string,
	username string,
	entries []RecentEntry,
	fonts Fonts,
	header string,
	nowLabel string,
) (string, error) {
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return "", err
	}

	const (
		width     = 600
		rowH      = 38
		headerH   = 95 // leaves a gap below the header text (subtitle at y=62)
		marginX   = 30.0
		listMaxPx = 520.0 // truncate width for "N. track — artist"
	)
	height := headerH + rowH*len(entries) + 20

	dc := gg.NewContext(width, height)
	dc.SetRGB255(18, 18, 18)
	dc.Clear()

	// Header: username + caption (e.g. "recent tracks").
	poppins := loadFont(fonts.Poppins, 22)
	arial := loadFont(fonts.Arial, 18)

	dc.SetColor(color.White)
	dc.SetFontFace(poppins)
	dc.DrawString(username, marginX, 40)

	dc.SetFontFace(arial)
	dc.SetRGB255(170, 170, 170)
	dc.DrawString(header, marginX, 62)
	dc.SetColor(color.White)

	// Rows. The now-playing entry (if any) is labeled with nowLabel instead of
	// a number; the remaining entries are numbered sequentially, so a 5-row
	// list reads: current, 2, 3, 4, 5.
	for i, t := range entries {
		y := float64(headerH + i*rowH)
		var prefix string
		switch {
		case t.Now:
			prefix = nowLabel + " — "
		case t.Loved:
			prefix = "♥ " + strconv.Itoa(i+1) + ". "
		default:
			prefix = strconv.Itoa(i+1) + ". "
		}
		label := prefix + t.Track + " — " + t.Artist
		dc.SetFontFace(arial)
		dc.DrawString(truncate(dc, label, listMaxPx), marginX, y)
	}

	filename := downloadDir + "/" + uuid.NewString() + ".jpg"
	out, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer out.Close()
	if err := jpeg.Encode(out, dc.Image(), &jpeg.Options{Quality: 95}); err != nil {
		return "", err
	}
	return filename, nil
}

// ScrobbleDir returns the config.DownloadPath subdirectory a service should
// download its cover art into.
func ScrobbleDir(service string) string {
	return fmt.Sprintf("%s/%s", config.DownloadPath, service)
}
