package utils

import (
	"fmt"
	"strconv"
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

func FormatInteger(number int) string {
	switch {
	case number >= 1000000000000000000:
		return fmt.Sprintf("%.1fQn", float64(number)/1000000000000000000)
	case number >= 1000000000000000:
		return fmt.Sprintf("%.1fQd", float64(number)/1000000000000000)
	case number >= 1000000000000:
		return fmt.Sprintf("%.1fT", float64(number)/1000000000000)
	case number >= 1000000000:
		return fmt.Sprintf("%.1fB", float64(number)/1000000000)
	case number >= 1000000:
		return fmt.Sprintf("%.1fM", float64(number)/1000000)
	case number >= 1000:
		return fmt.Sprintf("%.1fK", float64(number)/1000)
	default:
		return strconv.Itoa(number)
	}
}

func HumanBytes(size int64) string {
	if size == 0 {
		return ""
	}

	power := float64(1024)
	units := []string{" ", "Ki", "Mi", "Gi", "Ti"}
	i := 0
	s := float64(size)

	for s >= power && i < len(units)-1 {
		s /= power
		i++
	}

	return fmt.Sprintf("%.2f %sB", s, units[i])
}

func SplitIntoRows(items []telego.InlineKeyboardButton, width int) (IKB [][]telego.InlineKeyboardButton) {
	for i := 0; i < len(items); i += width {
		end := i + width
		if end > len(items) {
			end = len(items)
		}
		IKB = append(IKB, items[i:end])
	}
	return IKB
}
