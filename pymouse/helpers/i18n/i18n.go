package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"pymouse/pymouse/database/utilitiesdb"
	"strings"

	"github.com/mymmrac/telego"
)

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

func getChatLanguage(chat telego.Chat) string {
	var chatLanguage string
	if strings.Contains(chat.Type, telego.ChatTypePrivate) {
		chatLanguage = utilitiesdb.FindUser(chat.ID, "").Language
	} else {
		chatLanguage = utilitiesdb.FindChat(chat.ID).Language
	}
	return chatLanguage
}

func getStringFromNestedMap(langMap map[string]interface{}, key string) string {
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
	language := getChatLanguage(chat)

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
		if val := getStringFromNestedMap(langMap, key); val != "" {
			return val
		}
		return "STRING_UNAVAILABLE"
	}
}
