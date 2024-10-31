package utils

import (
	"strings"

	"github.com/mymmrac/telego"
)

func GetArgs(update telego.Update) string {
	MsgText := update.Message.Text
	SplitedText := strings.Split(MsgText, " ")
	Args := strings.Join(SplitedText[1:], " ")
	return Args
}
