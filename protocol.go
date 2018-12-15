package telegram

import (
	"errors"
	"strconv"
	"log"
	"fmt"
	"time"

	"encoding/json"

	"io/ioutil"

	"net/http"
	"net/url"
)


type Protocol struct {
	client http.Client
	apiKey string
	me TUser
	current_async_id int
	mostRecentlyReceived int
}

func NewProtocol() (Protocol) {
	p := Protocol{
		client: http.Client{
			Transport: http.DefaultTransport,
			Timeout: 90 * time.Second,
		},
	}
	return p
}

func (this *Protocol) SetAPIKey(newKey string) () {
	this.apiKey = newKey
}

func (this *Protocol) GetMe() (TUser) {
	return this.me
}

func (this *Protocol) GetNextId() (int) {
	this.current_async_id = this.current_async_id + 1
	return this.current_async_id
}

func (this *Protocol) Test() (error) {
	url := apiEndpoint + this.apiKey + "/getMe"
	r, e := this.client.Get(url)

	if r != nil {
    	log.Printf("[telegram] API call: %s (%s)\n", url, r.Status)
		defer r.Body.Close()
	} else {
    	log.Printf("[telegram] API call: %s (failed: %s)\n", url, e.Error())
    }
	if e != nil { return e }

	b, e := ioutil.ReadAll(r.Body)
	if e != nil { return e }

	var resp TGenericResponse
	e = json.Unmarshal(b, &resp)
	if e != nil { return e }

	e = HandleSoftError(&resp)
	if e != nil { return e }

	if resp.Result == nil {
		return errors.New("Missing required field (result)!")
	}

	e = json.Unmarshal(*resp.Result, &this.me)
	if e != nil { return e }

	log.Printf("Validated API key (%s)\n", *this.me.Username)
	return nil
}

// URL building functions

func (this *Protocol) BuildGetChatMemberURL(chat_id interface{}, user_id int) (string) {
	apiurl := apiEndpoint + this.apiKey + "/getChatMember?" + 
	       "chat_id=" + url.QueryEscape(GetStringId(chat_id)) +
	       "&user_id=" + strconv.Itoa(user_id)
	return apiurl
}

func (this *Protocol) BuildAnswerInlineQueryURL(q TInlineQuery, next_offset string) (string) {
	// next_offset should get stuck at -1 forever if pagination breaks somehow, to prevent infinite loops.

	return apiEndpoint + this.apiKey + "/answerInlineQuery?" +
		"inline_query_id=" + url.QueryEscape(q.Id) +
		"&next_offset=" + next_offset +
		"&cache_time=30" +
		"&results="
}

func (this *Protocol) BuildSendMessageURL(chat_id interface{}, text string, reply_to *int, mtype string, reply_markup interface{}, disable_preview bool) (string) {
	apiurl := apiEndpoint + this.apiKey + "/sendMessage?" + 
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
		if e != nil { return "" }
		apiurl = apiurl + "&reply_markup=" + url.QueryEscape(string(b))
	}
	if disable_preview {
		apiurl = apiurl + "&disable_web_page_preview=true"
	}

	return apiurl
}

func (this *Protocol) BuildEditMessageURL(chat_id interface{}, message_id int, inline_id string, text string, parse_mode string, reply_markup interface{}, disable_preview bool) (string) {
	apiurl := apiEndpoint + this.apiKey + "/editMessageText?" + 
	       "chat_id=" + url.QueryEscape(GetStringId(chat_id)) +
	       "&message_id=" + strconv.FormatInt(int64(message_id), 10) +
	       "&text=" + url.QueryEscape(text)
	if parse_mode != "" {
		apiurl = apiurl + "&parse_mode=" + url.QueryEscape(parse_mode)
	}
	if reply_markup != nil {
		b, e := json.Marshal(reply_markup)
		if e != nil { return "" }
		apiurl = apiurl + "&reply_markup=" + url.QueryEscape(string(b))
	}
	if disable_preview {
		apiurl = apiurl + "&disable_web_page_preview=true"
	}
	return apiurl
}

func (this *Protocol) BuildDeleteMessageURL(chat_id interface{}, message_id int) (string) {
	apiurl := apiEndpoint + this.apiKey + "/deleteMessage?" + 
	       "chat_id=" + url.QueryEscape(GetStringId(chat_id)) +
	       "&message_id=" + strconv.FormatInt(int64(message_id), 10)
	return apiurl
}

func (this *Protocol) BuildSendStickerURL(chat_id interface{}, sticker_id string, reply_to *int, disable_notification bool) (string) {
	url := apiEndpoint + this.apiKey + "/sendSticker?" + 
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

func (this *Protocol) BuildForwardMessageURL(chat_id interface{}, from_chat_id interface{}, message_id int, disable_notification bool) (string) {
	url := apiEndpoint + this.apiKey + "/forwardMessage?" + 
	       "chat_id=" + url.QueryEscape(GetStringId(chat_id)) +
	       "&from_chat_id=" + url.QueryEscape(GetStringId(from_chat_id)) + 
	       "&message_id=" + strconv.FormatInt(int64(message_id), 10)
	if disable_notification {
		url = url + "&disable_notification=true"
	}
	return url
}

func (this *Protocol) BuildKickMemberURL(chat_id interface{}, member int) (string) {
	return apiEndpoint + this.apiKey + "/kickChatMember?" + 
	       "chat_id=" + url.QueryEscape(GetStringId(chat_id)) +
	       "&user_id=" + strconv.FormatInt(int64(member), 10)
}

func (this *Protocol) BuildGetStickerSetURL(name string) (string) {
	return apiEndpoint + this.apiKey + "/getStickerSet?" + 
	       "name=" + url.QueryEscape(name)
}

func (this *Protocol) BuildRestrictChatMemberURL(chat_id interface{}, user_id int, until int64, messages, media, basic_media, web_previews bool) (string) {
	return apiEndpoint + this.apiKey + "/restrictChatMember?" + 
	       "chat_id=" + url.QueryEscape(GetStringId(chat_id)) +
	       "&user_id=" + strconv.FormatInt(int64(user_id), 10) +
	       "&until_date=" + strconv.FormatInt(int64(until), 10) +
	       "&can_send_messages=" + strconv.FormatBool(messages) +
	       "&can_send_media_messages=" + strconv.FormatBool(media) +
	       "&can_send_other_messages=" + strconv.FormatBool(basic_media) +
	       "&can_send_web_page_previews=" + strconv.FormatBool(web_previews)
}

func (this *Protocol) BuildGetFileURL(file_id string) (string) {
	return apiEndpoint + this.apiKey + "/getFile?" + 
	       "file_id=" + url.QueryEscape(file_id)
}

func (this *Protocol) BuildDownloadFileURL(file_path string) (string) {
	return apiFileEndpoint + this.apiKey + "/" + file_path
}

func (this *Protocol) BuildAnswerCallbackQueryURL(query_id, notification string, show_alert bool) (string) {
	apiurl := apiEndpoint + this.apiKey + "/answerCallbackQuery?" + 
	          "callback_query_id=" + url.QueryEscape(query_id) +
	          "&text=" + url.QueryEscape(notification)

	if show_alert {
	       apiurl = apiurl + "&show_alert=true"
	}

	return apiurl
}

// Async calls

func (this *Protocol) AnswerInlineQueryAsync(q TInlineQuery, out []interface{}, offset string, rm ResponseHandler) (error) {
	b, e := json.Marshal(out)
	if e != nil { return e }
	surl := this.BuildAnswerInlineQueryURL(q, offset)
	go DoAsyncCall(&this.client, rm, &CallResponseChannel, surl + url.QueryEscape(string(b)))
	return nil
}


func (this *Protocol) SendMessageAsync(chat_id interface{}, text string, reply_to *int, parse_mode string, reply_markup interface{}, disable_preview bool, sm ResponseHandler) () {
	go DoAsyncCall(&this.client, sm, &CallResponseChannel, this.BuildSendMessageURL(chat_id, text, reply_to, parse_mode, reply_markup, disable_preview))
}

func (this *Protocol) EditMessageTextAsync(chat_id interface{}, message_id int, inline_id string, text string, parse_mode string, reply_markup interface{}, disable_preview bool, sm ResponseHandler) () {
	go DoAsyncCall(&this.client, sm, &CallResponseChannel, this.BuildEditMessageURL(chat_id, message_id, inline_id, text, parse_mode, reply_markup, disable_preview))
}

func (this *Protocol) DeleteMessageAsync(chat_id interface{}, message_id int, sm ResponseHandler) () {
	go DoAsyncCall(&this.client, sm, &CallResponseChannel, this.BuildDeleteMessageURL(chat_id, message_id))
}

func (this *Protocol) SendStickerAsync(chat_id interface{}, sticker_id string, reply_to *int, disable_notification bool, sm ResponseHandler) () {
	go DoAsyncCall(&this.client, sm, &CallResponseChannel, this.BuildSendStickerURL(chat_id, sticker_id, reply_to, disable_notification))
}

func (this *Protocol) ForwardMessageAsync(chat_id interface{}, from_chat_id interface{}, message_id int, disable_notification bool, sm ResponseHandler) () {
	go DoAsyncCall(&this.client, sm, &CallResponseChannel, this.BuildForwardMessageURL(chat_id, from_chat_id, message_id, disable_notification))
}

func (this *Protocol) KickMemberAsync(chat_id interface{}, member int, si ResponseHandler) () {
	go DoAsyncCall(&this.client, si, &CallResponseChannel, this.BuildKickMemberURL(chat_id, member))
}

func (this *Protocol) GetStickerSetAsync(name string, rm ResponseHandler) () {
	go DoAsyncCall(&this.client, rm, &CallResponseChannel, this.BuildGetStickerSetURL(name))
}

func (this *Protocol) GetChatMemberAsync(chat_id interface{}, user_id int, rm ResponseHandler) () {
	go DoAsyncCall(&this.client, rm, &CallResponseChannel, this.BuildGetChatMemberURL(chat_id, user_id))
}

func (this *Protocol) RestrictChatMemberAsync(chat_id interface{}, user_id int, until int64, messages, media, basic_media, web_previews bool, rm ResponseHandler) () {
	go DoAsyncCall(&this.client, rm, &CallResponseChannel, this.BuildRestrictChatMemberURL(chat_id, user_id, until, messages, media, basic_media, web_previews))
}

func (this *Protocol) GetFileAsync(file_id string, rm ResponseHandler) () {
	go DoAsyncCall(&this.client, rm, &CallResponseChannel, this.BuildGetFileURL(file_id))
}

func (this *Protocol) DownloadFileAsync(file_path string, rm ResponseHandler) () {
	go DoAsyncFetch(&this.client, rm, &CallResponseChannel, this.BuildDownloadFileURL(file_path))
}

func (this *Protocol) AnswerCallbackQueryAsync(query_id, notification string, show_alert bool, rm ResponseHandler) () {
	go DoAsyncCall(&this.client, rm, &CallResponseChannel, this.BuildAnswerCallbackQueryURL(query_id, notification, show_alert))
}

// Synchronous calls

func (this *Protocol) AnswerInlineQuery(q TInlineQuery, out []interface{}, offset string) (error) {
	b, err := json.Marshal(out)
	if err != nil { return err }

	_, err = DoCall(&this.client, this.BuildAnswerInlineQueryURL(q, offset) + url.QueryEscape(string(b)))
	return err
}

func (this *Protocol) SendMessage(chat_id interface{}, text string, reply_to *int, mtype string, reply_markup interface{}, disable_preview bool) (*TMessage, error) {
	return OutputToMessage(DoCall(&this.client, this.BuildSendMessageURL(chat_id, text, reply_to, mtype, reply_markup, disable_preview)))
}

func (this *Protocol) EditMessageText(chat_id interface{}, message_id int, inline_id string, text string, parse_mode string, reply_markup interface{}, disable_preview bool) (*TMessage, error) {
	return OutputToMessage(DoCall(&this.client, this.BuildEditMessageURL(chat_id, message_id, inline_id, text, parse_mode, reply_markup, disable_preview)))
}

func (this *Protocol) DeleteMessage(chat_id interface{}, message_id int) (error) {
	_, err := DoCall(&this.client, this.BuildDeleteMessageURL(chat_id, message_id))
	return err
}

func (this *Protocol) SendSticker(chat_id interface{}, sticker_id string, reply_to *int, disable_notification bool) (*TMessage, error) {
	return OutputToMessage(DoCall(&this.client, this.BuildSendStickerURL(chat_id, sticker_id, reply_to, disable_notification)))
}

func (this *Protocol) ForwardMessage(chat_id interface{}, from_chat_id interface{}, message_id int, disable_notification bool) (*TMessage, error) {
	return OutputToMessage(DoCall(&this.client, this.BuildForwardMessageURL(chat_id, from_chat_id, message_id, disable_notification)))
}

func (this *Protocol) KickMember(chat_id interface{}, member int) (error) {
	_, err := DoCall(&this.client, this.BuildKickMemberURL(chat_id, member))
	return err
}

func (this *Protocol) GetStickerSet(name string) (*TStickerSet, error) {
	return OutputToStickerSet(DoCall(&this.client, this.BuildGetStickerSetURL(name)))
}

func (this *Protocol) GetChatMember(chat_id interface{}, user_id int) (*TChatMember, error) {
	return OutputToChatMember(DoCall(&this.client, this.BuildGetChatMemberURL(chat_id, user_id)))
}

func (this *Protocol) RestrictChatMember(chat_id interface{}, user_id int, until int64, messages, media, basic_media, web_previews bool, rm ResponseHandler) (error) {
	_, err := DoCall(&this.client, this.BuildRestrictChatMemberURL(chat_id, user_id, until, messages, media, basic_media, web_previews))
	return err
}

func (this *Protocol) GetFile(file_id string) (*TFile, error) {
	return OutputToFile(DoCall(&this.client, this.BuildGetFileURL(file_id)))
}

func (this *Protocol) DownloadFile(file_path string) ([]byte, error) {
	return DoFetch(&this.client, this.BuildDownloadFileURL(file_path))
}

func (this *Protocol) AnswerCallbackQuery(query_id, notification string, show_alert bool) (error) {
	_, err := DoCall(&this.client, this.BuildAnswerCallbackQueryURL(query_id, notification, show_alert))
	return err
}

// Updates

func (this *Protocol) GetUpdates() ([]TUpdate, error) {
	url := apiEndpoint + this.apiKey + "/getUpdates?" + 
	       "offset=" + strconv.Itoa(this.mostRecentlyReceived) +
	       "&timeout=3600"
	r, e := this.client.Get(url)

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
		this.mostRecentlyReceived = updates[len(updates) - 1].Update_id + 1
	}

	return updates, nil
}
