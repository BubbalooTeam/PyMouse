package i18n

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"pymouse/locales"
	"pymouse/pymouse/database/repositories"
	"strings"

	"github.com/mymmrac/telego"
)

type LocalizationStats struct {
	TotalStrings         int     // Total number of strings in the localization file
	TranslatedStrings    int     // Number of strings that have been translated
	UntranslatedStrings  int     // Number of strings that are still untranslated
	PercentageTranslated float64 // Percentage of strings that have been translated
}

var AvalaibleLanguages []string
var defaultLanguage = "en_us"
var StringsCache = make(map[string]map[string]interface{})

func CompileLocales() error {
	// Load locale files from the embedded locales.Files first; fall back to
	// the on-disk "locales/" directory so `go run` without the embed still
	// works during development.
	load := func(name string, data []byte) error {
		langCode := strings.TrimSuffix(name, filepath.Ext(name))

		langMap := make(map[string]interface{})
		if err := json.Unmarshal(data, &langMap); err != nil {
			return fmt.Errorf("Failed to unmarshal data from localization file %q.", name)
		}

		StringsCache[langCode] = langMap
		AvalaibleLanguages = append(AvalaibleLanguages, langCode)
		return nil
	}

	walkErr := fs.WalkDir(locales.Files, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("Failed to retrieve localization file: %v", err)
		}
		if d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		data, rerr := locales.Files.ReadFile(path)
		if rerr != nil {
			return fmt.Errorf("Failed to read localization file: %v", rerr)
		}
		return load(filepath.Base(path), data)
	})
	if walkErr != nil {
		// Fallback: try the on-disk "locales/" directory (e.g. unbundled go run).
		return filepath.Walk("locales", func(path string, info os.FileInfo, ferr error) error {
			if ferr != nil || info.IsDir() || filepath.Ext(path) != ".json" {
				return ferr
			}
			data, rerr := os.ReadFile(path)
			if rerr != nil {
				return fmt.Errorf("Failed to read localization file: %v", rerr)
			}
			return load(filepath.Base(path), data)
		})
	}
	return nil
}

func GetStringFromNestedMap(langMap map[string]interface{}, key string) string {
	keys := strings.Split(key, ".")
	currentMap := langMap

	for _, k := range keys {
		value, ok := currentMap[k]
		if !ok {
			return "STRING_UNAVALAIBLE"
		}

		if nestedMap, isMap := value.(map[string]interface{}); isMap {
			currentMap = nestedMap
		} else if strValue, isString := value.(string); isString {
			return strValue
		} else {
			return "STRING_UNAVALAIBLE"
		}
	}

	return "STRING_UNAVALAIBLE"
}

func Locale(chat telego.Chat) func(string) string {
	language := repositories.GetChatLanguage(chat)
	if language == "" {
		language = defaultLanguage
	}

	langMap, ok := StringsCache[language]
	if !ok {
		langMap = StringsCache[defaultLanguage]
	}

	if langMap == nil {
		return func(string) string {
			return "STRING_UNAVAILABLE"
		}
	}

	return func(key string) string {
		if val := GetStringFromNestedMap(langMap, key); val != "" {
			return val
		}
		return "STRING_UNAVAILABLE"
	}
}

func GetLocalizationStats(langCode string) (LocalizationStats, error) {
	defaultLangMap, ok := StringsCache[defaultLanguage]
	if !ok {
		return LocalizationStats{}, fmt.Errorf("default localization not found")
	}

	langMap, ok := StringsCache[langCode]
	if !ok {
		return LocalizationStats{}, fmt.Errorf("localization not found")
	}

	var stats LocalizationStats

	var recursive func(current map[string]interface{}, path string)

	recursive = func(current map[string]interface{}, path string) {
		for key, value := range current {

			fullPath := key
			if path != "" {
				fullPath = path + "." + key
			}

			switch v := value.(type) {

			case string:
				stats.TotalStrings++

				defaultStr := GetStringFromNestedMap(defaultLangMap, fullPath)
				langStr := GetStringFromNestedMap(langMap, fullPath)

				if langCode != defaultLanguage {
					if langStr != defaultStr &&
						langStr != "STRING_UNAVALAIBLE" {
						stats.TranslatedStrings++
					}
				} else {
					if langStr == defaultStr {
						stats.TranslatedStrings++
					}
				}

			case map[string]interface{}:
				if fullPath != "help.titles" {
					recursive(v, fullPath)
				}
			}
		}
	}

	recursive(defaultLangMap, "")

	stats.UntranslatedStrings = stats.TotalStrings - stats.TranslatedStrings

	if stats.TotalStrings > 0 {
		stats.PercentageTranslated =
			(float64(stats.TranslatedStrings) / float64(stats.TotalStrings)) * 100
	}

	return stats, nil
}
