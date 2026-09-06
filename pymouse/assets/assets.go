// Package assets embeds the bot's static resources (fonts and icons) into the
// binary so the compiled executable is self-contained and does not need the
// pymouse/assets/ tree on disk at runtime.
//
// Callers address files using paths relative to this directory, e.g.
// "fonts/arial.ttf" or "icons/lastfm/loved.png". The legacy
// "pymouse/assets/..." disk paths are also accepted (the prefix is stripped),
// which keeps `go run` development working when files are read from disk.
package assets

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
)

// Fonts holds every bundled font file under fonts/.
//
//go:embed fonts/*.ttf fonts/*.otf
var Fonts embed.FS

// Icons holds every bundled icon under icons/ (lastfm, weather, ...).
//
//go:embed icons
var Icons embed.FS

// diskRoot is the on-disk prefix the rest of the codebase historically used to
// address assets (relative to the repo root / the bot's working directory).
const diskRoot = "pymouse/assets"

// normalize converts either a legacy disk path ("pymouse/assets/fonts/arial.ttf")
// or an already-embed-relative path ("fonts/arial.ttf") into the embed-relative
// form ("fonts/arial.ttf").
func normalize(path string) string {
	p := strings.TrimPrefix(path, diskRoot+"/")
	return strings.TrimPrefix(p, "./")
}

// readFirstEmbed tries each embed.FS in turn, returning the first that contains
// the (embed-relative) name. On miss it falls back to disk so that unbundled
// `go run` development still works.
func readFirstEmbed(name string, fss ...embed.FS) ([]byte, error) {
	for _, fs := range fss {
		if b, err := fs.ReadFile(name); err == nil {
			return b, nil
		}
	}
	return os.ReadFile(filepath.Join(diskRoot, name))
}

// ReadFile returns the bytes of any bundled asset (font or icon), addressed by
// embed-relative path or legacy disk path. It tries the embedded copy first
// and falls back to disk.
func ReadFile(path string) ([]byte, error) {
	name := normalize(path)
	return readFirstEmbed(name, Fonts, Icons)
}
