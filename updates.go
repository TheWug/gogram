package telegram

import (
	"fmt"
	"strconv"
	"errors"
	"net/http"
	"io/ioutil"
	"encoding/json"
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
