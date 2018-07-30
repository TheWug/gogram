package telegram

import (
	"errors"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"net/url"
)

var CallResponseChannel chan HandlerBox = make(chan HandlerBox, 10)

type HandlerBox struct {
	Success   bool
	Error     error
	Http_code int
	Handler   ResponseHandler
	Output   *json.RawMessage
	Bytes   []byte
}

// call this in a goroutine.
func DoAsyncFetch(client *http.Client, handler ResponseHandler, output *chan HandlerBox, apiurl string) {
	var hbox HandlerBox
	hbox.Handler = handler

	r, e := client.Get(apiurl)
	if r != nil {
		defer r.Body.Close()
	}
	if e != nil {
		hbox.Error = e
		if (output != nil) { *output <- hbox }
		return
	}

	hbox.Http_code = r.StatusCode
	hbox.Bytes, e = ioutil.ReadAll(r.Body)
	if e != nil {
		hbox.Error = e
		if (output != nil) { *output <- hbox }
		return
	}

	hbox.Success = true
	if (output != nil) { *output <- hbox }
	return
}

func DoAsyncCall(client *http.Client, handler ResponseHandler, output *chan HandlerBox, apiurl string) {
	temp := make(chan HandlerBox, 1)
	DoAsyncFetch(client, handler, &temp, apiurl)
	close(temp)
	hbox := <- temp
	hbox.Success = false

	var out TGenericResponse
	e = json.Unmarshal(hbox.Bytes, &out)

	if e != nil {
		hbox.Error = e
		if (output != nil) { *output <- hbox }
		return
	}

	e = HandleSoftError(&out)
	if e != nil {
		hbox.Error = e
		if (output != nil) { *output <- hbox }
		return
	}

	hbox.Output = out.Result
	hbox.Bytes = nil
	hbox.Success = true
	if (output != nil) { *output <- hbox }
	return
}

func DoFetch(client *http.Client, apiurl string) ([]byte, error) {
	ch := make(chan HandlerBox, 1)

	DoAsyncFetch(client, nil, &ch, apiurl)
	close(ch)
	output := <- ch

	return output.Bytes, output.Error
}

func DoCall(client *http.Client, apiurl string) (*json.RawMessage, error) {
	ch := make(chan HandlerBox, 1)

	DoAsyncCall(client, nil, &ch, apiurl)
	close(ch)
	output := <- ch

	return output.Output, output.Error
}

// Type Helpers

func OutputToMessage(raw *json.RawMessage, err error) (*TMessage, error) {
	if err != nil { return nil, err }

	var msg TMessage
	err = json.Unmarshal(*raw, &msg)

	if err != nil { return nil, err }
	if msg.Message_id == 0 { return nil, errors.New("Missing message") }
	return &msg, nil
}

func OutputToStickerSet(raw *json.RawMessage, err error) (*TStickerSet, error) {
	if err != nil { return nil, err }

	var set TStickerSet
	err = json.Unmarshal(*raw, &set)

	if err != nil { return nil, err }
	if set.Stickers == nil { return nil, errors.New("Missing sticker set") }
	return &set, nil
}

func OutputToChatMember(raw *json.RawMessage, err error) (*TChatMember, error) {
	if err != nil { return nil, err }

	var cm TChatMember
	err = json.Unmarshal(*raw, &cm)

	if err != nil { return nil, err }
	if cm.Status == "" { return nil, errors.New("Missing chat member") }
	return &cm, nil
}

func OutputToFile(raw *json.RawMessage, err error) (*TFile, error) {
	if err != nil { return nil, err }

	var out TFile
	err = json.Unmarshal(*raw, &out)

	if err != nil { return nil, err }
	if out.File_id == "" { return nil, errors.New("Missing file metadata") }
	return &out, nil
}
