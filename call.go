package telegram

import (
	"errors"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var CallResponseChannel chan HandlerBox = make(chan HandlerBox, 10)

type HandlerBox struct {
	Success   bool
	Error     error
	Http_code int
	Handler  ResponseHandler
	Output   *json.RawMessage
}

// call this in a goroutine.
func DoAsyncCall(client *http.Client, handler ResponseHandler, output *chan HandlerBox, apiurl string) {
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
	b, e := ioutil.ReadAll(r.Body)
	if e != nil {
		hbox.Error = e
		if (output != nil) { *output <- hbox }
		return
	}

	var out TGenericResponse
	e = json.Unmarshal(b, &out)

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
	hbox.Success = true
	if (output != nil) { *output <- hbox }
	return
}

func DoGenericAPICall(apiurl string, out_obj interface{}) (error) {
	r, e := http.Get(apiurl)
	if r != nil {
		defer r.Body.Close()
	}
	if e != nil {
		return e
	}

	b, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return e
	}

	var out TGenericResponse
	e = json.Unmarshal(b, &out)

	if e != nil {
		return e
	}

	e = HandleSoftError(&out)
	if e != nil {
		return e
	}

	e = json.Unmarshal(*out.Result, out_obj)

	if e != nil {
		return e
	}

	return nil
}

func DoCall(client *http.Client, apiurl string) (*json.RawMessage, error) {
	ch := make(chan HandlerBox, 1)

	DoAsyncCall(client, nil, &ch, apiurl)
	close(ch)
	output := <- ch

	return output.Output, output.Error
}

func DeleteMessage(chat_id interface{}, message_id int) (error) {
	var str_chat_id string
	switch t := chat_id.(type) {
	case int:
		str_chat_id = strconv.FormatInt(int64(t), 10)
	case int64:
		str_chat_id = strconv.FormatInt(t, 10)
	}

	apiurl := apiEndpoint + apiKey + "/deleteMessage?" + 
	       "chat_id=" + url.QueryEscape(str_chat_id) +
	       "&message_id=" + strconv.FormatInt(int64(message_id), 10)

	return DoSendAPICall(apiurl)
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
	if cm.Status == "" { return nil, errors.New("Missing sticker set") }
	return &cm, nil
}

func GetFile(file_id string) (*TFile, error) {
	url := apiEndpoint + apiKey + "/getFile?" + 
	       "file_id=" + url.QueryEscape(file_id)
	
	var f TFile
	err := DoGenericAPICall(url, &f)
	return &f, err
}

func DownloadFile(file_path string) ([]byte, error) {
	apiurl := apiFileEndpoint + apiKey + "/" + file_path
	r, e := http.Get(apiurl)
	if r != nil {
		defer r.Body.Close()
	}
	if e != nil {
		return nil, e
	}

	return ioutil.ReadAll(r.Body)
}

func AnswerCallbackQuery(query_id, notification string, show_alert bool) (error) {
	url := apiEndpoint + apiKey + "/answerCallbackQuery?" + 
	       "callback_query_id=" + url.QueryEscape(query_id) +
	       "&text=" + url.QueryEscape(notification)

	if show_alert {
	       url = url + "&show_alert=true"
	}

	return DoSendAPICall(url)
}
