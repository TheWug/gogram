package gogram

import (
	"strconv"
)

// common state for the entire telegram package.

var userAgent string = "@KnottyBot (v1.0, operator: @wuuug id:68060168)"
var apiEndpoint string = "https://api.telegram.org/bot"
var apiFileEndpoint string = "https://api.telegram.org/file/bot"

func GetStringId(chat_id interface{}) (string) {
	var str_chat_id string
	switch t := chat_id.(type) {
	case int:
		str_chat_id = strconv.FormatInt(int64(t), 10)
	case int64:
		str_chat_id = strconv.FormatInt(t, 10)
	case string:
		str_chat_id = t
	}

	return str_chat_id
}
