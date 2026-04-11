package upload

import (
	"fmt"
	"path"
	"pymouse/pymouse/config"
	"pymouse/pymouse/helpers/utils"
	"strings"

	"github.com/cavaliergopher/grab/v3"
)

func isAllowedMime(mime string) bool {
	allowed := []string{
		"image/", "video/", "audio/",
		"text/",
		"application/pdf",
		"application/json",
		"application/zip",
		"application/x-rar-compressed",
		"application/octet-stream",
	}

	for _, a := range allowed {
		if strings.HasPrefix(mime, a) {
			return true
		}
	}
	return false
}

func downloadByURL(url string) (string, error) {
	client := grab.NewClient()

	allowedExtensions := []string{
		".jpg", ".jpeg", ".png", ".gif",
		".bmp", ".webp", ".svg", ".psd",
		".ai", ".eps", ".ico",
		".mp4", ".mov", ".avi", ".mkv",
		".flv", ".m3u8",
		".go", ".py", ".java", ".c",
		".cpp", ".h", ".cs", ".rb",
		".php", ".sql", ".js", ".ts",
		".html", ".css", ".xml", ".json",
		".yml", ".yaml", ".md",
		".pdf", ".docx", ".xlsx", ".pptx",
		".txt", ".csv", ".ini", ".conf",
		".log", ".vtt", ".srt",
		".zip", ".rar", ".7z", ".tar", ".gz",
		".mp3", ".wav", ".ogg", ".flac",
		".aac", ".opus", ".m4a",
		".ttf", ".otf", ".woff", ".woff2",
	}

	cleanURL := strings.Split(url, "?")[0]
	lowerURL := strings.ToLower(cleanURL)

	var ext string
	if strings.HasSuffix(lowerURL, ".tar.gz") {
		ext = ".tar.gz"
	} else {
		ext = strings.ToLower(path.Ext(lowerURL))
	}

	if !utils.StringInSlice(ext, allowedExtensions) && ext != ".tar.gz" {
		return "", fmt.Errorf("extension not allowed: %s", ext)
	}

	req, err := grab.NewRequest(config.DownloadPath, url)
	if err != nil {
		return "", err
	}

	resp := client.Do(req)

	contentType := resp.HTTPResponse.Header.Get("Content-Type")

	if !isAllowedMime(contentType) {
		return "", fmt.Errorf("mime type not allowed: %s", contentType)
	}

	if err := resp.Err(); err != nil {
		return "", err
	}

	return resp.Filename, nil
}
