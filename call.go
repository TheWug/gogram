package telegram

import (
	"strconv"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

var CallResponseChannel chan HandlerBox = make(chan HandlerBox, 10)

type HandlerBox struct {
	Success   bool
	Error     error
	Http_code int
	Handler  ResponseHandler
	Output   *json.RawMessage
}

func DoAsyncCall(handler ResponseHandler, output *chan HandlerBox, apiurl string) {
	var hbox HandlerBox
	hbox.Handler = handler

	r, e := http.Get(apiurl)
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

func DoCall(apiurl string) (*json.RawMessage, error) {
	ch := make(chan HandlerBox, 1)

	DoAsyncCall(nil, &ch, apiurl)
	close(ch)
	output := <- ch

	return output.Output, output.Error
}

// URL building functions

func BuildGetChatMemberURL(chat_id interface{}, user_id int) (string) {
	apiurl := apiEndpoint + apiKey + "/getChatMember?" + 
	       "chat_id=" + url.QueryEscape(GetStringId(chat_id)) +
	       "&user_id=" + strconv.Itoa(user_id)
	return apiurl
}

func BuildAnswerInlineQueryURL(q TInlineQuery, next_offset string) (string) {
	// next_offset should get stuck at -1 forever if pagination breaks somehow, to prevent infinite loops.

	return apiEndpoint + apiKey + "/answerInlineQuery?" +
		"inline_query_id=" + url.QueryEscape(q.Id) +
		"&next_offset=" + next_offset +
		"&cache_time=30" +
		"&results="
}

func BuildSendMessageURL(chat_id interface{}, text string, reply_to *int, mtype string, reply_markup interface{}) (string) {
	apiurl := apiEndpoint + apiKey + "/sendMessage?" + 
	       "chat_id=" + url.QueryEscape(GetStringId(chat_id)) +
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

	return apiurl
}

func BuildEditMessageURL(chat_id interface{}, message_id int, inline_message_id string, text string, parse_mode string, reply_markup interface{}) (string) {
	apiurl := apiEndpoint + apiKey + "/editMessageText?" + 
	       "chat_id=" + url.QueryEscape(GetStringId(chat_id)) +
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

	return apiurl
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

func BuildSendStickerURL(chat_id interface{}, sticker_id string, reply_to *int, disable_notification bool) (string) {
	url := apiEndpoint + apiKey + "/sendSticker?" + 
	       "chat_id=" + url.QueryEscape(GetStringId(chat_id)) +
	       "&sticker=" + url.QueryEscape(sticker_id)
	if reply_to != nil {
		url = url + "&reply_to_message_id=" + strconv.Itoa(*reply_to)
	}
	if disable_notification {
		url = url + "&disable_notification=true"
	}
	return url
}

func BuildForwardMessageURL(chat_id interface{}, from_chat_id interface{}, message_id int, disable_notification bool) (string) {
	url := apiEndpoint + apiKey + "/forwardMessage?" + 
	       "chat_id=" + url.QueryEscape(GetStringId(chat_id)) +
	       "&from_chat_id=" + url.QueryEscape(GetStringId(from_chat_id)) + 
	       "&message_id=" + strconv.FormatInt(int64(message_id), 10)
	if disable_notification {
		url = url + "&disable_notification=true"
	}
	return url
}

func BuildKickMemberURL(chat_id interface{}, member int) (string) {
	return apiEndpoint + apiKey + "/kickChatMember?" + 
	       "chat_id=" + url.QueryEscape(GetStringId(chat_id)) +
	       "&user_id=" + strconv.FormatInt(int64(member), 10)
}

// Type Helpers

func OutputToMessage(raw *json.RawMessage, err error) (*TMessage, error) {
	if err != nil { return nil, err }

	var msg TMessage
	err = json.Unmarshal(*raw, &msg)

	if err != nil { return nil, err }
	return &msg, nil
}

// Async calls

func GetChatMemberAsync(chat_id interface{}, user_id int, rm ResponseHandler) () {
	go DoAsyncCall(rm, &CallResponseChannel, BuildGetChatMemberURL(chat_id, user_id))
}

func AnswerInlineQueryAsync(q TInlineQuery, out []interface{}, offset string, rm ResponseHandler) (error) {
	b, e := json.Marshal(out)
	if e != nil { return e }
	surl := BuildAnswerInlineQueryURL(q, offset)
	go DoAsyncCall(rm, &CallResponseChannel, surl + url.QueryEscape(string(b)))
	return nil
}


func SendMessageAsync(chat_id interface{}, text string, reply_to *int, mtype string, sm ResponseHandler) () {
	go DoAsyncCall(sm, &CallResponseChannel, BuildSendMessageURL(chat_id, text, reply_to, mtype))
}

func EditMessageTextAsync(chat_id interface{}, message_id int, _ string, text string, parse_mode string, sm ResponseHandler) () {
	go DoAsyncCall(sm, &CallResponseChannel, BuildEditMessageURL(chat_id, message_id, "", text, parse_mode))
}

func SendStickerAsync(chat_id interface{}, sticker_id string, reply_to *int, disable_notification bool, sm ResponseHandler) () {
	go DoAsyncCall(sm, &CallResponseChannel, BuildSendStickerURL(chat_id, sticker_id, reply_to, disable_notification))
}

func ForwardMessageAsync(chat_id interface{}, from_chat_id interface{}, message_id int, disable_notification bool, sm ResponseHandler) () {
	go DoAsyncCall(sm, &CallResponseChannel, BuildForwardMessageURL(chat_id, from_chat_id, message_id, disable_notification))
}

func KickMemberAsync(chat_id interface{}, member int, si ResponseHandler) () {
	go DoAsyncCall(si, &CallResponseChannel, BuildKickMemberURL(chat_id, member))
}

// Synchronous calls

func GetChatMember(chat_id interface{}, user_id int) (*TChatMember, error) {
	raw, err := DoCall(BuildGetChatMemberURL(chat_id, user_id))
	if err != nil { return nil, err }

	var chatmember TChatMember
	err = json.Unmarshal(*raw, &chatmember)

	if err != nil { return nil, err }
	return &chatmember, nil
}

func AnswerInlineQuery(q TInlineQuery, out []interface{}, offset string) (error) {
	b, err := json.Marshal(out)
	if err != nil { return err }

	_, err = DoCall(BuildAnswerInlineQueryURL(q, offset) + url.QueryEscape(string(b)))
	return err
}

func SendMessage(chat_id interface{}, text string, reply_to *int, mtype string, reply_markup interface{}) (*TMessage, error) {
	return OutputToMessage(DoCall(BuildSendMessageURL(chat_id, text, reply_to, mtype, reply_markup)))
}

func EditMessageText(chat_id interface{}, message_id int, inline_message_id string, text string, parse_mode string, reply_markup interface{}) (*TMessage, error) {
	return OutputToMessage(DoCall(BuildEditMessageURL(chat_id, message_id, "", text, parse_mode, reply_markup)))
}

func SendSticker(chat_id interface{}, sticker_id string, reply_to *int, disable_notification bool) (*TMessage, error) {
	return OutputToMessage(DoCall(BuildSendStickerURL(chat_id, sticker_id, reply_to, disable_notification)))
}

func ForwardMessage(chat_id interface{}, from_chat_id interface{}, message_id int, disable_notification bool) (*TMessage, error) {
	return OutputToMessage(DoCall(BuildForwardMessageURL(chat_id, from_chat_id, message_id, disable_notification)))
}

func KickMember(chat_id interface{}, member int) (error) {
	_, err := DoCall(BuildKickMemberURL(chat_id, member))
	return err
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

