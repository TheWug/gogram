package gogram

import (
	"strconv"
)

// common state for the entire telegram package.

var userAgent string = "@KnottyBot (v1.0, operator: @wuuug id:68060168)"
var apiEndpoint string = "https://api.telegram.org/bot"
var apiFileEndpoint string = "https://api.telegram.org/file/bot"

func GetStringId(chat_id interface{}) (string) {
	switch t := chat_id.(type) {
	case int:
		return strconv.FormatInt(int64(t), 10)
	case int64:
		return strconv.FormatInt(t, 10)
	case string:
		return t
	}

	return ""
}
