package data

import (
	"encoding/json"
)

type UserID int
type ChatID int64

type Sender struct {
	User UserID
	Channel ChatID
}

type ResponseHandler interface {
	Callback(result *json.RawMessage, success bool, err error, http_code int)
}
