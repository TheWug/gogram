package gogram

import (
	"github.com/thewug/gogram/data"

	"fmt"
)

// common state for the entire telegram package.

var userAgent string = "@KnottyBot (v1.0, operator: @wuuug id:68060168)"
var apiEndpoint string = "https://api.telegram.org/bot"
var apiFileEndpoint string = "https://api.telegram.org/file/bot"

func GetStringId(chat_id interface{}) (string) {
	switch t := chat_id.(type) {
	case data.UserID:
		return t.String()
	case data.ChatID:
		return t.String()
	case string:
		return t
	default:
		panic(fmt.Sprintf("Bad type for chat_id: %T (must be ChatID, UserID, or string)", chat_id))
	}

	return ""
}
