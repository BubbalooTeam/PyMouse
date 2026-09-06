// Package locales embeds the JSON localization files shipped in this
// directory so the compiled binary carries them without needing the
// locales/ folder on disk at runtime.
//
// Files are listed explicitly (rather than via a glob like *.json) because
// go:embed patterns cannot use globs when the directory also contains Go
// files — every locale file must be named here. Add new locales to this
// list when they are introduced.
package locales

import "embed"

//go:embed en_us.json pt_br.json
var Files embed.FS
