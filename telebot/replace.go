package telebot

import (
	"strings"
)

func NormalizeDashes(cmd string) (string) {
	return strings.Replace(cmd, "—", "--", -1)
}
