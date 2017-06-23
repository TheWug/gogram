package telegram

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"encoding/json"

	"io/ioutil"

	"net/http"
)

func GetUpdates() ([]TUpdate, error) {
	url := apiEndpoint + apiKey + "/getUpdates?" + 
	       "offset=" + strconv.Itoa(mostRecentlyReceived) +
	       "&timeout=3600"
	r, e := http.Get(url)

	if r != nil {
		defer r.Body.Close()
		if r.StatusCode != 200 {
			return nil, errors.New(fmt.Sprintf("API %s", r.Status))
		}
	}
	if e != nil {
		return nil, e
	}

	b, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return nil, e
	}

	var out TGenericResponse
	e = json.Unmarshal(b, &out)

	if e != nil {
		return nil, e
	}

	e = HandleSoftError(&out)
	if e != nil {
		return nil, e
	}

	var updates []TUpdate
	e = json.Unmarshal(*out.Result, &updates)

	if e != nil {
		return nil, e
	}

	// track the next update to request
	if len(updates) != 0 {
		mostRecentlyReceived = updates[len(updates) - 1].Update_id + 1
	}

	return updates, nil
}

func AsyncUpdateLoop(output chan []TUpdate) () {
	for {
		updates, e := GetUpdates()
		if e != nil {
			// if an error occurred, sleep before trying again to avoid a tight error loop
			//common.ErrorLog("telegram", "telegram.GetUpdates()", e, "GetUpdates failed: ")
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}

		output <- updates
	}
}
