package gogram

import (
	"github.com/thewug/gogram/data"

	"github.com/thewug/reqtify"

	"errors"
	"strconv"
	"log"
	"time"
	"bytes"

	"encoding/json"

	"io/ioutil"
	"io"

	"net/http"
)


type Protocol struct {
	client      reqtify.Reqtifier
	file_client reqtify.Reqtifier

	apiKey           string
	me               data.TUser
	current_async_id int

	nextUpdateOffset     data.UpdateID
	nextUpdateOffsetless bool

	bot *TelegramBot
}

func NewProtocol(bound *TelegramBot) (Protocol) {
	p := Protocol{
		nextUpdateOffsetless: true,
		bot: bound,
	}
	return p
}

func (this *Protocol) SetAPIKey(newKey string) () {
	this.apiKey = newKey
	this.client = reqtify.New(apiEndpoint + this.apiKey + "/", nil, &http.Client{
		Transport: http.DefaultTransport,
		Timeout: 90 * time.Second,
	}, nil, userAgent)
	this.file_client = reqtify.New(apiFileEndpoint + this.apiKey + "/", nil, &http.Client{
		Transport: http.DefaultTransport,
		Timeout: 90 * time.Second,
	}, nil, userAgent)
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
	log.Printf("[telegram] API call: %s (%s)\n", req.GetPath(), r.Status)
		defer r.Body.Close()
	} else {
		log.Printf("[telegram] API call: %s (failed: %s)\n", req.GetPath(), e.Error())
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

func (this *Protocol) BuildGetUpdatesReq() (reqtify.Request) {
	return this.client.New("getUpdates").
			   ArgDefault("offset", this.getNextUpdateOffset(), "").
			   Arg("timeout", "3600")
}

func (this *Protocol) BuildGetChatMemberReq(o data.OChatMember) (reqtify.Request) {
	return this.client.New("getChatMember").
			   Arg("chat_id", GetStringId(o.ChatId)).
			   Arg("user_id", o.UserId.String())
}

func (this *Protocol) BuildGetChatReq(o data.OChatMember) (reqtify.Request) {
	return this.client.New("getChat").
			   Arg("chat_id", GetStringId(o.ChatId))
}

func (this *Protocol) BuildAnswerInlineQueryReq(o data.OInlineQueryAnswer) (reqtify.Request) {
	// next_offset should get stuck at -1 forever if pagination breaks somehow, to prevent infinite loops.

	// manually encode, instead of using Marshal. Size matters here, and without SetEscapeHTML(false),
	// every < and > (of which there are tons in the inline HTML tags) expands into "\uxxxx", 6 times longer.
	// experimentation has shown average space savings of 20% is possible by simply skipping this escaping,
	// which also saves time, and is a compatibility workaround for browsers anyway, so totally unnecessary here.
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	e := encoder.Encode(o.Results)
	b := buf.Bytes()

	if e != nil { return nil }

	return this.client.New("answerInlineQuery").Method(reqtify.POST).
			   Arg("inline_query_id", o.Id.String()).
			   ArgDefault("cache_time", strconv.Itoa(o.CacheTime), "0").
			   Arg("next_offset", o.NextOffset).
			   Arg("results", string(b)).
			   ArgDefault("is_personal", o.IsPersonal, false).
			   ArgDefault("switch_pm_text", o.SwitchPMText, "").
			   ArgDefault("switch_pm_parameter", o.SwitchPMParam, "").
			   Multipart() // the alternative of multipart is URL encoded, which escapes a lot of stuff and thusly takes much more space
}

func (this *Protocol) BuildSendMessageReq(o data.OMessage) (reqtify.Request) {
	req := this.client.New("sendMessage").Method(reqtify.POST).
			   Arg("chat_id", GetStringId(o.ChatId)).
			   Arg("text", o.Text).
			   ArgDefault("parse_mode", o.ParseMode.String(), "").
			   ArgDefault("disable_web_page_preview", strconv.FormatBool(o.DisableWebPagePreview), "false").
			   ArgDefault("disable_notification", o.DisableNotification.String(), "false")
	if o.ReplyToId != nil {
		req.Arg("reply_to_message_id", o.ReplyToId.String())
	}
	if o.ReplyMarkup != nil {
		b, e := json.Marshal(o.ReplyMarkup)
		if e != nil { return nil }
		req.Arg("reply_markup", string(b))
	}

	return req
}

func (this *Protocol) BuildEditMessageReq(o data.OMessageEdit) (reqtify.Request) {
	req := this.client.New("editMessageText").Method(reqtify.POST).
			   ArgDefault("chat_id", GetStringId(o.SourceChatId), "").
			   ArgDefault("message_id", o.SourceMessageId.String(), "0").
			   ArgDefault("inline_id", o.SourceInlineId.String(), "").
			   Arg("text", o.Text).
			   ArgDefault("parse_mode", o.ParseMode.String(), "").
			   ArgDefault("disable_web_page_preview", strconv.FormatBool(o.DisableWebPagePreview), "false")
	if o.ReplyMarkup != nil {
		b, e := json.Marshal(o.ReplyMarkup)
		if e != nil { return nil }
		req.Arg("reply_markup", string(b))
	}
	return req
}

func (this *Protocol) BuildEditCaptionReq(o data.OCaptionEdit) (reqtify.Request) {
	req := this.client.New("editMessageCaption").Method(reqtify.POST).
			   ArgDefault("chat_id", GetStringId(o.SourceChatId), "").
			   ArgDefault("message_id", o.SourceMessageId.String(), "0").
			   ArgDefault("inline_id", o.SourceInlineId.String(), "").
			   Arg("caption", o.Text).
			   ArgDefault("parse_mode", o.ParseMode.String(), "")
	if o.ReplyMarkup != nil {
		b, e := json.Marshal(o.ReplyMarkup)
		if e != nil { return nil }
		req.Arg("reply_markup", string(b))
	}
	return req
}

func (this *Protocol) BuildEditReplyMarkupReq(o data.OMessageEdit) (reqtify.Request) {
	req := this.client.New("editMessageReplyMarkup").Method(reqtify.POST).
			   ArgDefault("chat_id", GetStringId(o.ChatId), "").
			   ArgDefault("message_id", o.SourceMessageId.String(), "0").
			   ArgDefault("inline_id", o.SourceInlineId.String(), "")
	if o.ReplyMarkup != nil {
		b, e := json.Marshal(o.ReplyMarkup)
		if e != nil { return nil }
		req.Arg("reply_markup", string(b))
	}
	return req
}

func (this *Protocol) BuildDeleteMessageReq(o data.ODelete) (reqtify.Request) {
	return this.client.New("deleteMessage").
			   Arg("chat_id", GetStringId(o.SourceChatId)).
			   Arg("message_id", o.SourceMessageId.String())
}

// This resolves the allowed file attachment modes. 'object' can be any of the following:
// FileID: passed directly to API, represents an internal (already uploaded) file.
// string: passed directly to API, assumed to be a URL from which telegram will download the file.
// io.Reader: marshalled into a form file, where it will be read to completion by the HTTP request. filename backup is used.
// byte array: same as above. filename backup is used.
// reqtify.FormFile: same as above, but uses included filename
func applyFile(req reqtify.Request, tag, filename, backup_filename string, object interface{}) {
	if filename == "" { filename = backup_filename }
	switch file := object.(type) {
	case data.FileID:
		req.Arg(tag, string(file))
	case string:
		req.Arg(tag, file)
	case io.Reader:
		req.Method(reqtify.POST).FileArg(tag, filename, file)
	case []byte:
		req.Method(reqtify.POST).FileArg(tag, filename, bytes.NewReader(file))
	case reqtify.FormFile:
		req.Method(reqtify.POST).FileArg(tag, file.Name, file.Data)
	default:
		panic("unsupported file attachment mode")
	}
}

func (this *Protocol) BuildSendStickerReq(o data.OSticker) (reqtify.Request) {
	req := this.client.New("sendSticker").
			   Arg("chat_id", GetStringId(o.ChatId)).
			   ArgDefault("disable_notification", o.DisableNotification.String(), "false")
	if o.ReplyToId != nil {
		req.Arg("reply_to_message_id", o.ReplyToId.String())
	}
	if o.ReplyMarkup != nil {
		b, e := json.Marshal(o.ReplyMarkup)
		if e != nil { return nil }
		req.Arg("reply_markup", string(b))
	}
	applyFile(req, "sticker", o.FileName, "sticker.webp", o.File)
	return req
}

func (this *Protocol) BuildSendPhotoReq(o data.OPhoto) (reqtify.Request) {
	req := this.client.New("sendPhoto").
			   Arg("chat_id", GetStringId(o.ChatId)).
			   ArgDefault("caption", o.Text, "").
			   ArgDefault("parse_mode", o.ParseMode.String(), "").
			   ArgDefault("disable_notification", o.DisableNotification.String(), "false")
	if o.ReplyToId != nil {
		req.Arg("reply_to_message_id", o.ReplyToId.String())
	}
	if o.ReplyMarkup != nil {
		b, e := json.Marshal(o.ReplyMarkup)
		if e != nil { return nil }
		req.Arg("reply_markup", string(b))
	}
	applyFile(req, "photo", o.FileName, "photo.jpg", o.File)
	return req
}

func (this *Protocol) BuildSendAnimationReq(o data.OAnimation) (reqtify.Request) {
	req := this.client.New("sendAnimation").
			   Arg("chat_id", GetStringId(o.ChatId)).
			   ArgDefault("caption", o.Text, "").
			   ArgDefault("parse_mode", o.ParseMode.String(), "").
			   ArgDefault("disable_notification", o.DisableNotification.String(), "false").
			   ArgDefault("height", o.Height, 0).
			   ArgDefault("width", o.Width, 0).
			   ArgDefault("duration", o.Duration, 0)
	if o.ReplyToId != nil {
		req.Arg("reply_to_message_id", o.ReplyToId.String())
	}
	if o.ReplyMarkup != nil {
		b, e := json.Marshal(o.ReplyMarkup)
		if e != nil { return nil }
		req.Arg("reply_markup", string(b))
	}
	applyFile(req, "animation", o.FileName, "animation.mp4", o.File)
	return req
}

func (this *Protocol) BuildSendDocumentReq(o data.ODocument) (reqtify.Request) {
	req := this.client.New("sendDocument").
			   Arg("chat_id", GetStringId(o.ChatId)).
			   ArgDefault("caption", o.Text, "").
			   ArgDefault("parse_mode", o.ParseMode.String(), "").
			   ArgDefault("disable_notification", o.DisableNotification.String(), "false")
	if o.ReplyToId != nil {
		req.Arg("reply_to_message_id", o.ReplyToId.String())
	}
	if o.ReplyMarkup != nil {
		b, e := json.Marshal(o.ReplyMarkup)
		if e != nil { return nil }
		req.Arg("reply_markup", string(b))
	}
	applyFile(req, "document", o.FileName, "file", o.File)
	return req
}

func (this *Protocol) BuildSendAudioReq(o data.OAudio) (reqtify.Request) {
	req := this.client.New("sendAudio").
			   Arg("chat_id", GetStringId(o.ChatId)).
			   ArgDefault("caption", o.Text, "").
			   ArgDefault("parse_mode", o.ParseMode.String(), "").
			   ArgDefault("disable_notification", o.DisableNotification.String(), "false")
	if o.ReplyToId != nil {
		req.Arg("reply_to_message_id", o.ReplyToId.String())
	}
	if o.ReplyMarkup != nil {
		b, e := json.Marshal(o.ReplyMarkup)
		if e != nil { return nil }
		req.Arg("reply_markup", string(b))
	}
	applyFile(req, "photo", o.FileName, "photo.jpg", o.File)
	return req
}

func (this *Protocol) BuildForwardMessageReq(o data.OForward) (reqtify.Request) {
	return this.client.New("forwardMessage").
			   Arg("chat_id", GetStringId(o.ChatId)).
			   Arg("from_chat_id", GetStringId(o.SourceChatId)).
			   Arg("message_id", o.SourceMessageId.String()).
			   ArgDefault("disable_notification", o.DisableNotification.String(), "false")
}

func (this *Protocol) BuildKickMemberReq(o data.OChatMember) (reqtify.Request) {
	return this.client.New("kickChatMember").
			   Arg("chat_id", GetStringId(o.ChatId)).
			   Arg("user_id", o.UserId.String())
}

func (this *Protocol) BuildGetStickerSetReq(o data.OStickerSet) (reqtify.Request) {
	return this.client.New("getStickerSet").Arg("name", o.Name)
}

func (this *Protocol) BuildRestrictChatMemberReq(o data.ORestrict) (reqtify.Request) {
	req := this.client.New("restrictChatMember").Method(reqtify.POST).
			   Arg("chat_id", GetStringId(o.ChatID)).
			   Arg("user_id", strconv.Itoa(o.UserID)).
			   ArgDefault("until_date", strconv.FormatInt(o.Until, 10), "0")
	
	b, e := json.Marshal(o.ToTChatPermissions())
	if e != nil { panic(e.Error()) }
	req.Arg("permissions", string(b))

	return req
}

func (this *Protocol) BuildGetFileReq(o data.OGetFile) (reqtify.Request) {
	return this.client.New("getFile").Arg("file_id", o.Id.String())
}

func (this *Protocol) BuildDownloadFileReq(o data.OFile) (reqtify.Request) {
	return this.file_client.New(o.FilePath)
}

func (this *Protocol) BuildAnswerCallbackQueryReq(o data.OCallback) (reqtify.Request) {
	return this.client.New("answerCallbackQuery").Method(reqtify.POST).
			   Arg("callback_query_id", o.Id.String()).
			   ArgDefault("text", o.Notification, "").
			   ArgDefault("show_alert", strconv.FormatBool(o.ShowAlert), "false").
			   ArgDefault("cache_time", strconv.Itoa(o.CacheTime), "0").
			   ArgDefault("url", o.URL, "")
}

func (this *Protocol) BuildSetChatPermissionsReq(o data.ORestrict) (reqtify.Request) {
	req := this.client.New("setChatPermissions").Method(reqtify.POST).
			  Arg("chat_id", GetStringId(o.ChatID))

	b, e := json.Marshal(o.ToTChatPermissions())
	if e != nil { panic(e.Error()) }
	req.Arg("permissions", string(b))

	return req
}

// Async calls

func (this *Protocol) SendMessageAsync(o data.OMessage, sm data.ResponseHandler) () {
	go DoAsyncCall(this.bot.Log, this.BuildSendMessageReq(o), sm)
}

func (this *Protocol) ForwardMessageAsync(o data.OForward, sm data.ResponseHandler) () {
	go DoAsyncCall(this.bot.Log, this.BuildForwardMessageReq(o), sm)
}

func (this *Protocol) SendStickerAsync(o data.OSticker, sm data.ResponseHandler) () {
	go DoAsyncCall(this.bot.Log, this.BuildSendStickerReq(o), sm)
}

func (this *Protocol) SendDocumentAsync(o data.ODocument, sm data.ResponseHandler) () {
	go DoAsyncCall(this.bot.Log, this.BuildSendDocumentReq(o), sm)
}

func (this *Protocol) SendPhotoAsync(o data.OPhoto, sm data.ResponseHandler) () {
	go DoAsyncCall(this.bot.Log, this.BuildSendPhotoReq(o), sm)
}

func (this *Protocol) SendAnimationAsync(o data.OAnimation, sm data.ResponseHandler) () {
	go DoAsyncCall(this.bot.Log, this.BuildSendAnimationReq(o), sm)
}

func (this *Protocol) AnswerInlineQueryAsync(o data.OInlineQueryAnswer, rm data.ResponseHandler) {
	go DoAsyncCall(this.bot.Log, this.BuildAnswerInlineQueryReq(o), rm)
}

func (this *Protocol) EditMessageTextAsync(o data.OMessageEdit, sm data.ResponseHandler) () {
	go DoAsyncCall(this.bot.Log, this.BuildEditMessageReq(o), sm)
}

func (this *Protocol) EditMessageCaptionAsync(o data.OCaptionEdit, sm data.ResponseHandler) () {
	go DoAsyncCall(this.bot.Log, this.BuildEditCaptionReq(o), sm)
}

func (this *Protocol) EditReplyMarkupAsync(o data.OMessageEdit, sm data.ResponseHandler) () {
	go DoAsyncCall(this.bot.Log, this.BuildEditReplyMarkupReq(o), sm)
}

func (this *Protocol) DeleteMessageAsync(o data.ODelete, sm data.ResponseHandler) () {
	go DoAsyncCall(this.bot.Log, this.BuildDeleteMessageReq(o), sm)
}

func (this *Protocol) KickMemberAsync(o data.OChatMember, si data.ResponseHandler) () {
	go DoAsyncCall(this.bot.Log, this.BuildKickMemberReq(o), si)
}

func (this *Protocol) GetStickerSetAsync(o data.OStickerSet, rm data.ResponseHandler) () {
	go DoAsyncCall(this.bot.Log, this.BuildGetStickerSetReq(o), rm)
}

func (this *Protocol) GetChatAsync(o data.OChatMember, rm data.ResponseHandler) () {
	go DoAsyncCall(this.bot.Log, this.BuildGetChatReq(o), rm)
}

func (this *Protocol) GetChatMemberAsync(o data.OChatMember, rm data.ResponseHandler) () {
	go DoAsyncCall(this.bot.Log, this.BuildGetChatMemberReq(o), rm)
}

func (this *Protocol) RestrictChatMemberAsync(o data.ORestrict, rm data.ResponseHandler) () {
	go DoAsyncCall(this.bot.Log, this.BuildRestrictChatMemberReq(o), rm)
}

func (this *Protocol) GetFileAsync(o data.OGetFile, rm data.ResponseHandler) () {
	go DoAsyncCall(this.bot.Log, this.BuildGetFileReq(o), rm)
}

func (this *Protocol) DownloadFileAsync(o data.OFile, rm data.ResponseHandler) () {
	go DoAsyncFetch(this.bot.Log, this.BuildDownloadFileReq(o), rm)
}

func (this *Protocol) AnswerCallbackQueryAsync(o data.OCallback, rm data.ResponseHandler) () {
	go DoAsyncCall(this.bot.Log, this.BuildAnswerCallbackQueryReq(o), rm)
}

func (this *Protocol) SetChatPermissionsAsync(o data.ORestrict, rm data.ResponseHandler) () {
	go DoAsyncCall(this.bot.Log, this.BuildSetChatPermissionsReq(o), rm)
}

// Synchronous calls

func (this *Protocol) SendMessage(o data.OMessage) (*data.TMessage, error) {
	var m data.TMessage
	j, e := DoCall(this.bot.Log, this.BuildSendMessageReq(o))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) ForwardMessage(o data.OForward) (*data.TMessage, error) {
	var m data.TMessage
	j, e := DoCall(this.bot.Log, this.BuildForwardMessageReq(o))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) SendDocument(o data.ODocument) (*data.TMessage, error) {
	var m data.TMessage
	j, e := DoCall(this.bot.Log, this.BuildSendDocumentReq(o))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) SendSticker(o data.OSticker) (*data.TMessage, error) {
	var m data.TMessage
	j, e := DoCall(this.bot.Log, this.BuildSendStickerReq(o))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) SendPhoto(o data.OPhoto) (*data.TMessage, error) {
	var m data.TMessage
	j, e := DoCall(this.bot.Log, this.BuildSendPhotoReq(o))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) SendAnimation(o data.OAnimation) (*data.TMessage, error) {
	var m data.TMessage
	j, e := DoCall(this.bot.Log, this.BuildSendAnimationReq(o))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) AnswerInlineQuery(o data.OInlineQueryAnswer) (error) {
	_, err := DoCall(this.bot.Log, this.BuildAnswerInlineQueryReq(o))
	return err
}

func (this *Protocol) EditMessageText(o data.OMessageEdit) (*data.TMessage, error) {
	var m data.TMessage
	j, e := DoCall(this.bot.Log, this.BuildEditMessageReq(o))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) EditMessageCaption(o data.OCaptionEdit) (*data.TMessage, error) {
	var m data.TMessage
	j, e := DoCall(this.bot.Log, this.BuildEditCaptionReq(o))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) EditReplyMarkup(o data.OMessageEdit) (*data.TMessage, error) {
	var m data.TMessage
	j, e := DoCall(this.bot.Log, this.BuildEditReplyMarkupReq(o))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) DeleteMessage(o data.ODelete) (error) {
	_, err := DoCall(this.bot.Log, this.BuildDeleteMessageReq(o))
	return err
}

func (this *Protocol) KickMember(o data.OChatMember) (error) {
	_, err := DoCall(this.bot.Log, this.BuildKickMemberReq(o))
	return err
}

func (this *Protocol) GetStickerSet(o data.OStickerSet) (*data.TStickerSet, error) {
	var m data.TStickerSet
	j, e := DoCall(this.bot.Log, this.BuildGetStickerSetReq(o))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) GetChat(o data.OChatMember) (*data.TChat, error) {
	var m data.TChat
	j, e := DoCall(this.bot.Log, this.BuildGetChatReq(o))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) GetChatMember(o data.OChatMember) (*data.TChatMember, error) {
	var m data.TChatMember
	j, e := DoCall(this.bot.Log, this.BuildGetChatMemberReq(o))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) RestrictChatMember(o data.ORestrict, rm data.ResponseHandler) (error) {
	_, err := DoCall(this.bot.Log, this.BuildRestrictChatMemberReq(o))
	return err
}

func (this *Protocol) GetFile(o data.OGetFile) (*data.TFile, error) {
	var m data.TFile
	j, e := DoCall(this.bot.Log, this.BuildGetFileReq(o))
	return &m, OutputToObject(j, e, &m)
}

func (this *Protocol) DownloadFile(o data.OFile) (io.ReadCloser, error) {
	return DoGetReader(this.bot.Log, this.BuildDownloadFileReq(o))
}

func (this *Protocol) AnswerCallbackQuery(o data.OCallback) (error) {
	_, err := DoCall(this.bot.Log, this.BuildAnswerCallbackQueryReq(o))
	return err
}

func (this *Protocol) SetChatPermissions(o data.ORestrict) (error) {
	_, err := DoCall(this.bot.Log, this.BuildSetChatPermissionsReq(o))
	return err
}

// Updates

func (this *Protocol) GetUpdates() ([]data.TUpdate, error) {
	var updates []data.TUpdate
	r := this.BuildGetUpdatesReq()
	j, e := DoCall(nil, r)
	err := OutputToObject(j, e, &updates)
	return updates, err
}

func (this *Protocol) markUpdateProcessed(update *data.TUpdate) {
	if update == nil {
		panic("Tried to confirm a nil update!")
	}

	if this.nextUpdateOffset < update.Id + data.UpdateID(1) || this.nextUpdateOffsetless {
		this.nextUpdateOffset = update.Id + data.UpdateID(1)
		this.nextUpdateOffsetless = false
	}
}

func (this *Protocol) unmarkProcessedUpdate() {
	this.nextUpdateOffsetless = true
}

func (this *Protocol) getNextUpdateOffset() (string) {
	if this.nextUpdateOffsetless {
		return ""
	}

	return this.nextUpdateOffset.String()
}
