package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	dir := "locales"

	err := filepath.Walk(
		dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("Failed to retrieve localization file: %v", err)
			}

			if !info.IsDir() && filepath.Ext(path) == ".json" {
				langCode := filepath.Base(path[:len(path)-len(filepath.Ext(path))])

				data, err := os.ReadFile(path)
				if err != nil {
					return fmt.Errorf("Failed to read localization file: %v", err)
				}

				langMap := make(map[string]interface{})
				err = json.Unmarshal(data, &langMap)
				if err != nil {
					return fmt.Errorf("Failed to unmarshal data from localization file.")
				}

				StringsCache[langCode] = langMap
				AvalaibleLanguages = append(AvalaibleLanguages, langCode)
			}
			return nil
		},
	)
	return err
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
