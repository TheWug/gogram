package gogram

import (
	"github.com/thewug/gogram/data"

	"github.com/thewug/reqtify"

	"errors"
	"strconv"
	"log"
	"fmt"
	"time"
	"bytes"

	"encoding/json"

	"io/ioutil"
	"io"

	"net/http"
)


type Protocol struct {
	client reqtify.Reqtifier
	file_client reqtify.Reqtifier
	apiKey string
	me data.TUser
	current_async_id int
	mostRecentlyReceived int
}

func NewProtocol() (Protocol) {
	p := Protocol{
		client: reqtify.New("", nil, &http.Client{
			Transport: http.DefaultTransport,
			Timeout: 90 * time.Second,
		}, nil, userAgent),
		file_client: reqtify.New("", nil, &http.Client{
			Transport: http.DefaultTransport,
			Timeout: 90 * time.Second,
		}, nil, userAgent),
	}
	return p
}

func (this *Protocol) SetAPIKey(newKey string) () {
	this.apiKey = newKey
	this.client.Root = apiEndpoint + this.apiKey + "/"
	this.file_client.Root = apiFileEndpoint + this.apiKey + "/"
}

func (this *Protocol) GetMe() (data.TUser) {
	return this.me
}

func (this *Protocol) GetNextId() (int) {
	this.current_async_id = this.current_async_id + 1
	return this.current_async_id
}

func (this *Protocol) Test() (error) {
	req := this.client.New("getMe")
	r, e := req.Do()

	if r != nil {
    	log.Printf("[telegram] API call: %s (%s)\n", req.Path, r.Status)
		defer r.Body.Close()
	} else {
		log.Printf("[telegram] API call: %s (failed: %s)\n", req.Path, e.Error())
	}
	if e != nil { return e }

	b, e := ioutil.ReadAll(r.Body)
	if e != nil { return e }

	var resp data.TGenericResponse
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

func (this *Protocol) BuildGetChatMemberReq(chat_id interface{}, user_id int) (*reqtify.Request) {
	return this.client.New("getChatMember").
			   Arg("chat_id", GetStringId(chat_id)).
			   Arg("user_id", strconv.Itoa(user_id))
}

func (this *Protocol) BuildAnswerInlineQueryReq(q data.TInlineQuery, next_offset string, results []interface{}) (*reqtify.Request) {
	// next_offset should get stuck at -1 forever if pagination breaks somehow, to prevent infinite loops.

	b, e := json.Marshal(results)
	if e != nil { return nil }

	return this.client.New("answerInlineQuery").Method(reqtify.POST).
			   Arg("inline_query_id", q.Id).
			   Arg("cache_time", "30").
			   Arg("next_offset", next_offset).
			   Arg("results", string(b))
}

func (this *Protocol) BuildSendMessageReq(o data.OMessage) (*reqtify.Request) {
	req := this.client.New("sendMessage").Method(reqtify.POST).
			   Arg("chat_id", GetStringId(o.ChatID)).
			   Arg("text", o.Text).
			   ArgDefault("parse_mode", o.ParseMode, "").
			   ArgDefault("disable_web_page_preview", strconv.FormatBool(!o.EnableWebPreview), "false").
			   ArgDefault("disable_notification", strconv.FormatBool(o.DisableNotification), "false")
	if o.ReplyTo != nil {
		req.Arg("reply_to_message_id", strconv.Itoa(*o.ReplyTo))
	}
	if o.ReplyMarkup != nil {
		b, e := json.Marshal(o.ReplyMarkup)
		if e != nil { return nil }
		req.Arg("reply_markup", string(b))
	}

	return req
}

func (this *Protocol) BuildEditMessageReq(o data.OMessage) (*reqtify.Request) {
	req := this.client.New("editMessageText").Method(reqtify.POST).
			   ArgDefault("chat_id", GetStringId(o.ChatID), "").
			   ArgDefault("message_id", strconv.Itoa(o.MessageID), "0").
			   ArgDefault("inline_id", o.InlineID, "").
			   Arg("text", o.Text).
			   ArgDefault("parse_mode", o.ParseMode, "").
			   ArgDefault("disable_web_page_preview", strconv.FormatBool(!o.EnableWebPreview), "false")
	if o.ReplyMarkup != nil {
		b, e := json.Marshal(o.ReplyMarkup)
		if e != nil { return nil }
		req.Arg("reply_markup", string(b))
	}
	return req
}

func (this *Protocol) BuildDeleteMessageReq(o data.OMessage) (*reqtify.Request) {
	return this.client.New("deleteMessage").
			   Arg("chat_id", GetStringId(o.ChatID)).
			   Arg("message_id", strconv.Itoa(o.MessageID))
}

func (this *Protocol) BuildSendStickerReq(o data.OMessage) (*reqtify.Request) {
	req := this.client.New("sendSticker").
			   Arg("chat_id", GetStringId(o.ChatID)).
			   ArgDefault("disable_notification", strconv.FormatBool(o.DisableNotification), "false")
	if o.ReplyTo != nil {
		req.Arg("reply_to_message_id", strconv.Itoa(*o.ReplyTo))
	}
	if o.ReplyMarkup != nil {
		b, e := json.Marshal(o.ReplyMarkup)
		if e != nil { return nil }
		req.Arg("reply_markup", string(b))
	}
	switch t := o.Sticker.(type) {
	case string:
		req.Arg("sticker", t)
	case io.Reader:
		req.Method(reqtify.POST).FileArg("sticker", "sticker.webp", t)
	case []byte:
		req.Method(reqtify.POST).FileArg("sticker", "sticker.webp", bytes.NewReader(t))
	case reqtify.FormFile:
		req.Method(reqtify.POST).FileArg("sticker", t.Name, t.Data)
	}
	return req
}

func (this *Protocol) BuildForwardMessageReq(o data.OMessage) (*reqtify.Request) {
	return this.client.New("forwardMessage").
			   Arg("chat_id", GetStringId(o.TargetChatID)).
			   Arg("from_chat_id", GetStringId(o.ChatID)).
			   Arg("message_id", strconv.Itoa(o.MessageID)).
			   ArgDefault("disable_notification", strconv.FormatBool(o.DisableNotification), "false")
}

func (this *Protocol) BuildKickMemberReq(chat_id interface{}, member int) (*reqtify.Request) {
	return this.client.New("kickChatMember").
			   Arg("chat_id", GetStringId(chat_id)).
			   Arg("user_id", strconv.Itoa(member))
}

func (this *Protocol) BuildGetStickerSetReq(name string) (*reqtify.Request) {
	return this.client.New("getStickerSet").Arg("name", name)
}

func (this *Protocol) BuildRestrictChatMemberReq(chat_id interface{}, user_id int, until int64, messages, media, basic_media, web_previews bool) (*reqtify.Request) {
	return this.client.New("restrictChatMember").
		Arg("chat_id", GetStringId(chat_id)).
		Arg("user_id", strconv.Itoa(user_id)).
		Arg("until_date", strconv.FormatInt(until, 10)).
		Arg("can_send_messages", strconv.FormatBool(messages)).
		Arg("can_send_media_messages", strconv.FormatBool(media)).
		Arg("can_send_other_messages", strconv.FormatBool(basic_media)).
		Arg("can_send_web_page_previews", strconv.FormatBool(web_previews))
}

func (this *Protocol) BuildGetFileReq(file_id string) (*reqtify.Request) {
	return this.client.New("getFile").Arg("file_id", file_id)
}

func (this *Protocol) BuildDownloadFileReq(file_path string) (*reqtify.Request) {
	return this.file_client.New(file_path)
}

func (this *Protocol) BuildAnswerCallbackQueryReq(query_id, notification string, show_alert bool) (*reqtify.Request) {
	return this.client.New("answerCallbackQuery").Method(reqtify.POST).
			   Arg("callback_query_id", query_id).
			   Arg("text", notification).
			   ArgDefault("show_alert", strconv.FormatBool(show_alert), "false")
}

// Async calls

func (this *Protocol) AnswerInlineQueryAsync(q data.TInlineQuery, results []interface{}, offset string, rm data.ResponseHandler) {
	go DoAsyncCall(this.BuildAnswerInlineQueryReq(q, offset, results), rm)
}

func (this *Protocol) SendMessageAsync(o data.OMessage, sm data.ResponseHandler) () {
	go DoAsyncCall(this.BuildSendMessageReq(o), sm)
}

func (this *Protocol) EditMessageTextAsync(o data.OMessage, sm data.ResponseHandler) () {
	go DoAsyncCall(this.BuildEditMessageReq(o), sm)
}

func (this *Protocol) DeleteMessageAsync(o data.OMessage, sm data.ResponseHandler) () {
	go DoAsyncCall(this.BuildDeleteMessageReq(o), sm)
}

func (this *Protocol) SendStickerAsync(o data.OMessage, sm data.ResponseHandler) () {
	go DoAsyncCall(this.BuildSendStickerReq(o), sm)
}

func (this *Protocol) ForwardMessageAsync(o data.OMessage, sm data.ResponseHandler) () {
	go DoAsyncCall(this.BuildForwardMessageReq(o), sm)
}

func (this *Protocol) KickMemberAsync(chat_id interface{}, member int, si data.ResponseHandler) () {
	go DoAsyncCall(this.BuildKickMemberReq(chat_id, member), si)
}

func (this *Protocol) GetStickerSetAsync(name string, rm data.ResponseHandler) () {
	go DoAsyncCall(this.BuildGetStickerSetReq(name), rm)
}

func (this *Protocol) GetChatMemberAsync(chat_id interface{}, user_id int, rm data.ResponseHandler) () {
	go DoAsyncCall(this.BuildGetChatMemberReq(chat_id, user_id), rm, )
}

func (this *Protocol) RestrictChatMemberAsync(chat_id interface{}, user_id int, until int64, messages, media, basic_media, web_previews bool, rm data.ResponseHandler) () {
	go DoAsyncCall(this.BuildRestrictChatMemberReq(chat_id, user_id, until, messages, media, basic_media, web_previews), rm)
}

func (this *Protocol) GetFileAsync(file_id string, rm data.ResponseHandler) () {
	go DoAsyncCall(this.BuildGetFileReq(file_id), rm)
}

func (this *Protocol) DownloadFileAsync(file_path string, rm data.ResponseHandler) () {
	go DoAsyncFetch(this.BuildDownloadFileReq(file_path), rm)
}

func (this *Protocol) AnswerCallbackQueryAsync(query_id, notification string, show_alert bool, rm data.ResponseHandler) () {
	go DoAsyncCall(this.BuildAnswerCallbackQueryReq(query_id, notification, show_alert), rm)
}

// Synchronous calls

func (this *Protocol) AnswerInlineQuery(q data.TInlineQuery, results []interface{}, offset string) (error) {
	_, err := DoCall(this.BuildAnswerInlineQueryReq(q, offset, results))
	return err
}

func (this *Protocol) SendMessage(o data.OMessage) (*data.TMessage, error) {
	var m data.TMessage
	j, e := DoCall(this.BuildSendMessageReq(o))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) EditMessageText(o data.OMessage) (*data.TMessage, error) {
	var m data.TMessage
	j, e := DoCall(this.BuildEditMessageReq(o))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) DeleteMessage(o data.OMessage) (error) {
	_, err := DoCall(this.BuildDeleteMessageReq(o))
	return err
}

func (this *Protocol) SendSticker(o data.OMessage) (*data.TMessage, error) {
	var m data.TMessage
	j, e := DoCall(this.BuildSendStickerReq(o))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) ForwardMessage(o data.OMessage) (*data.TMessage, error) {
	var m data.TMessage
	j, e := DoCall(this.BuildForwardMessageReq(o))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) KickMember(chat_id interface{}, member int) (error) {
	_, err := DoCall(this.BuildKickMemberReq(chat_id, member))
	return err
}

func (this *Protocol) GetStickerSet(name string) (*data.TStickerSet, error) {
	var m data.TStickerSet
	j, e := DoCall(this.BuildGetStickerSetReq(name))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) GetChatMember(chat_id interface{}, user_id int) (*data.TChatMember, error) {
	var m data.TChatMember
	j, e := DoCall(this.BuildGetChatMemberReq(chat_id, user_id))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) RestrictChatMember(chat_id interface{}, user_id int, until int64, messages, media, basic_media, web_previews bool, rm data.ResponseHandler) (error) {
	_, err := DoCall(this.BuildRestrictChatMemberReq(chat_id, user_id, until, messages, media, basic_media, web_previews))
	return err
}

func (this *Protocol) GetFile(file_id string) (*data.TFile, error) {
	var m data.TFile
	j, e := DoCall(this.BuildGetFileReq(file_id))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) DownloadFile(file_path string) (io.ReadCloser, error) {
	return DoGetReader(this.BuildDownloadFileReq(file_path))
}

func (this *Protocol) AnswerCallbackQuery(query_id, notification string, show_alert bool) (error) {
	_, err := DoCall(this.BuildAnswerCallbackQueryReq(query_id, notification, show_alert))
	return err
}

// Updates

func (this *Protocol) GetUpdates() ([]data.TUpdate, error) {
	r, e := this.client.New("getUpdates").
			    Arg("offset", strconv.Itoa(this.mostRecentlyReceived)).
			    Arg("timeout", "3600").
			    Do()

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

	var out data.TGenericResponse
	e = json.Unmarshal(b, &out)

	if e != nil {
		return nil, e
	}

	e = HandleSoftError(&out)
	if e != nil {
		return nil, e
	}

	var updates []data.TUpdate
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
