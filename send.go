package telegram

import (
	"strconv"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
)

func DoAsyncSendAPICall(sent SentItem, output *chan SentItem, apiurl string) {
	r, e := http.Get(apiurl)
	if r != nil {
		defer r.Body.Close()
	}
	if e != nil {
		sent.Error = e
		if (output != nil) { *output <- sent }
		return
	}

	sent.Http_code = r.StatusCode
	sent.Success = true
	if (output != nil) { *output <- sent }
	return
}

func DoAsyncSendMessageAPICall(sent SentMessage, output *chan SentMessage, apiurl string) {
	r, e := http.Get(apiurl)
	if r != nil {
		defer r.Body.Close()
	}
	if e != nil {
		sent.Error = e
		if (output != nil) { *output <- sent }
		return
	}

	sent.Http_code = r.StatusCode
	b, e := ioutil.ReadAll(r.Body)
	if e != nil {
		sent.Error = e
		if (output != nil) { *output <- sent }
		return
	}

	var out TGenericResponse
	e = json.Unmarshal(b, &out)

	if e != nil {
		sent.Error = e
		if (output != nil) { *output <- sent }
		return
	}

	e = HandleSoftError(&out)
	if e != nil {
		sent.Error = e
		if (output != nil) { *output <- sent }
		return
	}

	var message *TMessage = new(TMessage)
	e = json.Unmarshal(*out.Result, message)

	if e != nil {
		sent.Error = e
		if (output != nil) { *output <- sent }
		return
	}

	sent.Message = message
	sent.Success = true
	if (output != nil) { *output <- sent }
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

func DoSendAPICall(apiurl string) (error) {
	ch := make(chan SentItem, 1)
	var sent SentItem
	sent.Id = -1

	DoAsyncSendAPICall(sent, &ch, apiurl)
	close(ch)
	sent = <- ch

	return sent.Error
}

func DoSendMessageAPICall(apiurl string) (*TMessage, error) {
	ch := make(chan SentMessage, 1)
	var sent SentMessage
	sent.Id = -1

	DoAsyncSendMessageAPICall(sent, &ch, apiurl)
	close(ch)
	sent = <- ch

	return sent.Message, sent.Error
}

// URL building functions

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

// Async calls

func SendMessageAsync(chat_id interface{}, text string, reply_to *int, mtype string, output *chan SentMessage, sm *SentMessage) (int) {
	if sm == nil { sm = new(SentMessage) }
	sm.Id = GetNextId()
	go DoAsyncSendMessageAPICall(*sm, output, BuildSendMessageURL(chat_id, text, reply_to, mtype))
	return sm.Id
}

func EditMessageTextAsync(chat_id interface{}, message_id int, _ string, text string, parse_mode string, output *chan SentMessage, sm *SentMessage) (int) {
	if sm == nil { sm = new(SentMessage) }
	sm.Id = GetNextId()
	go DoAsyncSendMessageAPICall(*sm, output, BuildEditMessageURL(chat_id, message_id, "", text, parse_mode))
	return sm.Id
}

func SendStickerAsync(chat_id interface{}, sticker_id string, reply_to *int, disable_notification bool, output *chan SentMessage, sm *SentMessage) (int) {
	if sm == nil { sm = new(SentMessage) }
	sm.Id = GetNextId()
	go DoAsyncSendMessageAPICall(*sm, output, BuildSendStickerURL(chat_id, sticker_id, reply_to, disable_notification))
	return sm.Id
}

func ForwardMessageAsync(chat_id interface{}, from_chat_id interface{}, message_id int, disable_notification bool, output *chan SentMessage, sm *SentMessage) (int) {
	if sm == nil { sm = new(SentMessage) }
	sm.Id = GetNextId()
	go DoAsyncSendMessageAPICall(*sm, output, BuildForwardMessageURL(chat_id, from_chat_id, message_id, disable_notification))
	return sm.Id
}

func KickMemberAsync(chat_id interface{}, member int, output *chan SentItem, si *SentItem) (int) {
	if si == nil { si = new(SentItem) }
	si.Id = GetNextId()
	go DoAsyncSendAPICall(*si, output, BuildKickMemberURL(chat_id, member))
	return si.Id
}

// Synchronous calls

func SendMessage(chat_id interface{}, text string, reply_to *int, mtype string, reply_markup interface{}) (*TMessage, error) {
	return DoSendMessageAPICall(BuildSendMessageURL(chat_id, text, reply_to, mtype, reply_markup))
}

func EditMessageText(chat_id interface{}, message_id int, inline_message_id string, text string, parse_mode string, reply_markup interface{}) {
	DoSendMessageAPICall(BuildEditMessageURL(chat_id, message_id, inline_message_id, text, parse_mode, reply_markup))
}

func SendSticker(chat_id interface{}, sticker_id string, reply_to *int, disable_notification bool) (message *TMessage, e error) {
	return DoSendMessageAPICall(BuildSendStickerURL(chat_id, sticker_id, reply_to, disable_notification))
}

func ForwardMessage(chat_id interface{}, from_chat_id interface{}, message_id int, disable_notification bool) (message *TMessage, e error) {
	return DoSendMessageAPICall(BuildForwardMessageURL(chat_id, from_chat_id, message_id, disable_notification))
}

func KickMember(chat_id interface{}, member int) (error) {
	return DoSendAPICall(BuildKickMemberURL(chat_id, member))
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
