package data

import (
	"encoding/json"
)

type Sender struct {
	User UserID
	Channel ChatID
}

type ResponseHandler interface {
	Callback(result *json.RawMessage, success bool, err error, http_code int)
}
