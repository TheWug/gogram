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

func (this *Protocol) BuildGetChatMemberReq(o data.OChatMember) (*reqtify.Request) {
	return this.client.New("getChatMember").
			   Arg("chat_id", GetStringId(o.ChatID)).
			   Arg("user_id", strconv.Itoa(o.UserID))
}

func (this *Protocol) BuildAnswerInlineQueryReq(o data.OInlineQueryAnswer) (*reqtify.Request) {
	// next_offset should get stuck at -1 forever if pagination breaks somehow, to prevent infinite loops.

	b, e := json.Marshal(o.Results)
	if e != nil { return nil }

	return this.client.New("answerInlineQuery").Method(reqtify.POST).
			   Arg("inline_query_id", o.QueryID).
			   ArgDefault("cache_time", strconv.Itoa(o.CacheTime), "0").
			   Arg("next_offset", o.NextOffset).
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

func (this *Protocol) BuildKickMemberReq(o data.OChatMember) (*reqtify.Request) {
	return this.client.New("kickChatMember").
			   Arg("chat_id", GetStringId(o.ChatID)).
			   Arg("user_id", strconv.Itoa(o.UserID))
}

func (this *Protocol) BuildGetStickerSetReq(o data.OStickerSet) (*reqtify.Request) {
	return this.client.New("getStickerSet").Arg("name", o.Name)
}

func (this *Protocol) BuildRestrictChatMemberReq(o data.ORestrict) (*reqtify.Request) {
	return this.client.New("restrictChatMember").
		Arg("chat_id", GetStringId(o.ChatID)).
		Arg("user_id", strconv.Itoa(o.UserID)).
		ArgDefault("until_date", strconv.FormatInt(o.Until, 10), "0").
		Arg("can_send_messages", strconv.FormatBool(o.CanSendMessages)).
		Arg("can_send_media_messages", strconv.FormatBool(o.CanSendMedia)).
		Arg("can_send_other_messages", strconv.FormatBool(o.CanSendInline)).
		Arg("can_send_web_page_previews", strconv.FormatBool(o.CanSendWebPreviews))
}

func (this *Protocol) BuildGetFileReq(o data.OGetFile) (*reqtify.Request) {
	return this.client.New("getFile").Arg("file_id", o.FileID)
}

func (this *Protocol) BuildDownloadFileReq(o data.OFile) (*reqtify.Request) {
	return this.file_client.New(o.FilePath)
}

func (this *Protocol) BuildAnswerCallbackQueryReq(o data.OCallback) (*reqtify.Request) {
	return this.client.New("answerCallbackQuery").Method(reqtify.POST).
			   Arg("callback_query_id", o.QueryID).
			   ArgDefault("text", o.Notification, "").
			   ArgDefault("show_alert", strconv.FormatBool(o.ShowAlert), "false").
			   ArgDefault("cache_time", strconv.Itoa(o.CacheTime), "0").
			   ArgDefault("url", o.URL, "")
}

// Async calls

func (this *Protocol) AnswerInlineQueryAsync(o data.OInlineQueryAnswer, rm data.ResponseHandler) {
	go DoAsyncCall(this.BuildAnswerInlineQueryReq(o), rm)
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

func (this *Protocol) KickMemberAsync(o data.OChatMember, si data.ResponseHandler) () {
	go DoAsyncCall(this.BuildKickMemberReq(o), si)
}

func (this *Protocol) GetStickerSetAsync(o data.OStickerSet, rm data.ResponseHandler) () {
	go DoAsyncCall(this.BuildGetStickerSetReq(o), rm)
}

func (this *Protocol) GetChatMemberAsync(o data.OChatMember, rm data.ResponseHandler) () {
	go DoAsyncCall(this.BuildGetChatMemberReq(o), rm)
}

func (this *Protocol) RestrictChatMemberAsync(o data.ORestrict, rm data.ResponseHandler) () {
	go DoAsyncCall(this.BuildRestrictChatMemberReq(o), rm)
}

func (this *Protocol) GetFileAsync(o data.OGetFile, rm data.ResponseHandler) () {
	go DoAsyncCall(this.BuildGetFileReq(o), rm)
}

func (this *Protocol) DownloadFileAsync(o data.OFile, rm data.ResponseHandler) () {
	go DoAsyncFetch(this.BuildDownloadFileReq(o), rm)
}

func (this *Protocol) AnswerCallbackQueryAsync(o data.OCallback, rm data.ResponseHandler) () {
	go DoAsyncCall(this.BuildAnswerCallbackQueryReq(o), rm)
}

// Synchronous calls

func (this *Protocol) AnswerInlineQuery(o data.OInlineQueryAnswer) (error) {
	_, err := DoCall(this.BuildAnswerInlineQueryReq(o))
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

func (this *Protocol) KickMember(o data.OChatMember) (error) {
	_, err := DoCall(this.BuildKickMemberReq(o))
	return err
}

func (this *Protocol) GetStickerSet(o data.OStickerSet) (*data.TStickerSet, error) {
	var m data.TStickerSet
	j, e := DoCall(this.BuildGetStickerSetReq(o))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) GetChatMember(o data.OChatMember) (*data.TChatMember, error) {
	var m data.TChatMember
	j, e := DoCall(this.BuildGetChatMemberReq(o))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) RestrictChatMember(o data.ORestrict, rm data.ResponseHandler) (error) {
	_, err := DoCall(this.BuildRestrictChatMemberReq(o))
	return err
}

func (this *Protocol) GetFile(o data.OGetFile) (*data.TFile, error) {
	var m data.TFile
	j, e := DoCall(this.BuildGetFileReq(o))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) DownloadFile(o data.OFile) (io.ReadCloser, error) {
	return DoGetReader(this.BuildDownloadFileReq(o))
}

func (this *Protocol) AnswerCallbackQuery(o data.OCallback) (error) {
	_, err := DoCall(this.BuildAnswerCallbackQueryReq(o))
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
