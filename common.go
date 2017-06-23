package telegram

import (
	"errors"
	"log"
	"strconv"

	"encoding/json"

	"io/ioutil"

	"net/http"
)

// common state for the entire telegram package.

var apiEndpoint string = "https://api.telegram.org/bot"
var apiFileEndpoint string = "https://api.telegram.org/file/bot"
var apiKey string = "CHANGEME"
var me TUser
var current_async_id int = 0

var mostRecentlyReceived int = 0

func SetAPIKey(newKey string) () {
	apiKey = newKey
}

func Test() (error) {
	url := apiEndpoint + apiKey + "/getMe"
	r, e := http.Get(url)

	if r != nil {
		defer r.Body.Close()
		log.Printf("[telegram] API call: %s (%s)\n", url, r.Status)
	}
	if e != nil {
		return e
	}

	b, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return e
	}

	var resp TGenericResponse
	e = json.Unmarshal(b, &resp)

	if e != nil {
		return e
	}

	e = HandleSoftError(&resp)
	if e != nil {
		return e
	}

	if resp.Result == nil {
		return errors.New("Missing required field (result)!")
	}

	e = json.Unmarshal(*resp.Result, &me)

	if e != nil {
		return e
	}

	log.Printf("Validated API key (%s)\n", *me.Username)
	return nil
}

func GetMe() (TUser) {
	return me
}

func GetNextId() (int) {
	current_async_id = current_async_id + 1
	return current_async_id
}

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
