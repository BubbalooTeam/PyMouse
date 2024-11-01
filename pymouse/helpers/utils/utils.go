package utils

import (
	"fmt"
	"strings"

	"github.com/mymmrac/telego"
)

func GetArgs(update telego.Update) string {
	MsgText := update.Message.Text
	SplitedText := strings.Split(MsgText, " ")
	Args := strings.Join(SplitedText[1:], " ")
	return Args
}

func TimeFormatter(seconds float64) string {
	totalSeconds := int(seconds)
	mins := totalSeconds / 60
	secs := totalSeconds % 60
	hrs := mins / 60
	mins = mins % 60
	days := hrs / 24
	hrs = hrs % 24

	var parts []string

	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hrs > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hrs))
	}
	if mins > 0 {
		parts = append(parts, fmt.Sprintf("%dm", mins))
	}
	if secs > 0 {
		parts = append(parts, fmt.Sprintf("%ds", secs))
	}

	return strings.Join(parts, ", ")
}
