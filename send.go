package telegram

import (
	"strconv"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
)

func DoSendAPICall(apiurl string) (error) {
	r, e := http.Get(apiurl)
	if r != nil {
		defer r.Body.Close()
	}
	if e != nil {
		return e
	}
	return nil
}

func DoSendMessageAPICall(apiurl string) (*TMessage, error) {
	r, e := http.Get(apiurl)
	if r != nil {
		defer r.Body.Close()
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

	var message *TMessage = new(TMessage)
	e = json.Unmarshal(*out.Result, message)

	if e != nil {
		return nil, e
	}

	return message, nil
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

func SendMessage(chat_id interface{}, text string, reply_to *int, mtype string, reply_markup interface{}) (*TMessage, error) {
	var str_chat_id string
	switch t := chat_id.(type) {
	case int:
		str_chat_id = strconv.FormatInt(int64(t), 10)
	case int64:
		str_chat_id = strconv.FormatInt(t, 10)
	}

	apiurl := apiEndpoint + apiKey + "/sendMessage?" + 
	       "chat_id=" + url.QueryEscape(str_chat_id) +
	       "&text=" + url.QueryEscape(text)

	if mtype != "" {
		apiurl = apiurl + "&parse_mode=" + url.QueryEscape(mtype)
	}
	       
	if reply_to != nil {
		apiurl = apiurl + "&reply_to_message_id=" + strconv.Itoa(*reply_to)
	}

	if reply_markup != nil {
		b, e := json.Marshal(reply_markup)
		if e != nil { return nil, e }
		apiurl = apiurl + "&reply_markup=" + url.QueryEscape(string(b))
	}

	return DoSendMessageAPICall(apiurl)
}

func EditMessageText(chat_id interface{}, message_id int, inline_message_id string, text string, parse_mode string, reply_markup interface{}) (*TMessage, error) {
	var str_chat_id string
	switch t := chat_id.(type) {
	case int:
		str_chat_id = strconv.FormatInt(int64(t), 10)
	case int64:
		str_chat_id = strconv.FormatInt(t, 10)
	}

	apiurl := apiEndpoint + apiKey + "/editMessageText?" + 
	       "chat_id=" + url.QueryEscape(str_chat_id) +
	       "&message_id=" + strconv.FormatInt(int64(message_id), 10) +
	       "&text=" + url.QueryEscape(text)

	if parse_mode != "" {
		apiurl = apiurl + "&parse_mode=" + url.QueryEscape(parse_mode)
	}

	if reply_markup != nil {
		b, e := json.Marshal(reply_markup)
		if e != nil { return nil, e }
		apiurl = apiurl + "&reply_markup=" + url.QueryEscape(string(b))
	}

	return DoSendMessageAPICall(apiurl)
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

func SendSticker(chat_id interface{}, sticker_id string, reply_to *int, disable_notification bool) (message *TMessage, e error) {
	var str_chat_id string
	switch t := chat_id.(type) {
	case string:
		str_chat_id = t
	case int:
		str_chat_id = strconv.FormatInt(int64(t), 10)
	case int64:
		str_chat_id = strconv.FormatInt(t, 10)
	}

	url := apiEndpoint + apiKey + "/sendSticker?" + 
	       "chat_id=" + url.QueryEscape(str_chat_id) +
	       "&sticker=" + url.QueryEscape(sticker_id)

	if reply_to != nil {
		url = url + "&reply_to_message_id=" + strconv.Itoa(*reply_to)
	}

	if disable_notification {
		url = url + "&disable_notification=true"
	}

	return DoSendMessageAPICall(url)
}

func ForwardMessage(chat_id interface{}, from_chat_id interface{}, message_id int, disable_notification bool) (message *TMessage, e error) {
	var str_chat_id string
	switch t := chat_id.(type) {
	case string:
		str_chat_id = t
	case int:
		str_chat_id = strconv.FormatInt(int64(t), 10)
	case int64:
		str_chat_id = strconv.FormatInt(t, 10)
	}

	var str_from_chat_id string
	switch t := from_chat_id.(type) {
	case string:
		str_from_chat_id = t
	case int:
		str_from_chat_id = strconv.FormatInt(int64(t), 10)
	case int64:
		str_from_chat_id = strconv.FormatInt(t, 10)
	}

	url := apiEndpoint + apiKey + "/forwardMessage?" + 
	       "chat_id=" + url.QueryEscape(str_chat_id) +
	       "&from_chat_id=" + url.QueryEscape(str_from_chat_id) + 
	       "&message_id=" + strconv.FormatInt(int64(message_id), 10)

	if disable_notification {
		url = url + "&disable_notification=true"
	}

	return DoSendMessageAPICall(url)
}

func KickMember(chat_id interface{}, member int) {
	var str_chat_id string
	switch t := chat_id.(type) {
	case string:
		str_chat_id = t
	case int:
		str_chat_id = strconv.FormatInt(int64(t), 10)
	case int64:
		str_chat_id = strconv.FormatInt(t, 10)
	}

	url := apiEndpoint + apiKey + "/kickChatMember?" + 
	       "chat_id=" + url.QueryEscape(str_chat_id) +
	       "&user_id=" + strconv.FormatInt(int64(member), 10)

	DoSendAPICall(url)
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
