package bot

import (
	"github.com/thewug/gogram"
	"io"
)


type Protocol interface {
	SetAPIKey(newKey string)
	GetMe() (gogram.TUser)
	GetNextId() (int)
	Test() (error)

	AnswerInlineQueryAsync(q gogram.TInlineQuery, results []interface{}, offset string, rm gogram.ResponseHandler)
	SendMessageAsync(chat_id interface{}, text string, reply_to *int, parse_mode string, reply_markup interface{}, disable_preview bool, sm gogram.ResponseHandler)
	EditMessageTextAsync(chat_id interface{}, message_id int, inline_id string, text string, parse_mode string, reply_markup interface{}, disable_preview bool, sm gogram.ResponseHandler)
	DeleteMessageAsync(chat_id interface{}, message_id int, sm gogram.ResponseHandler)
	SendStickerAsync(chat_id interface{}, sticker_id string, reply_to *int, disable_notification bool, sm gogram.ResponseHandler)
	ForwardMessageAsync(chat_id interface{}, from_chat_id interface{}, message_id int, disable_notification bool, sm gogram.ResponseHandler)
	KickMemberAsync(chat_id interface{}, member int, si gogram.ResponseHandler)
	GetStickerSetAsync(name string, rm gogram.ResponseHandler)
	GetChatMemberAsync(chat_id interface{}, user_id int, rm gogram.ResponseHandler)
	RestrictChatMemberAsync(chat_id interface{}, user_id int, until int64, messages, media, basic_media, web_previews bool, rm gogram.ResponseHandler)
	GetFileAsync(file_id string, rm gogram.ResponseHandler)
	DownloadFileAsync(file_path string, rm gogram.ResponseHandler)
	AnswerCallbackQueryAsync(query_id, notification string, show_alert bool, rm gogram.ResponseHandler)

	AnswerInlineQuery(q gogram.TInlineQuery, results []interface{}, offset string) (error)
	SendMessage(chat_id interface{}, text string, reply_to *int, mtype string, reply_markup interface{}, disable_preview bool) (*gogram.TMessage, error)
	EditMessageText(chat_id interface{}, message_id int, inline_id string, text string, parse_mode string, reply_markup interface{}, disable_preview bool) (*gogram.TMessage, error)
	DeleteMessage(chat_id interface{}, message_id int) (error)
	SendSticker(chat_id interface{}, sticker_id string, reply_to *int, disable_notification bool) (*gogram.TMessage, error)
	ForwardMessage(chat_id interface{}, from_chat_id interface{}, message_id int, disable_notification bool) (*gogram.TMessage, error)
	KickMember(chat_id interface{}, member int) (error)
	GetStickerSet(name string) (*gogram.TStickerSet, error)
	GetChatMember(chat_id interface{}, user_id int) (*gogram.TChatMember, error)
	RestrictChatMember(chat_id interface{}, user_id int, until int64, messages, media, basic_media, web_previews bool, rm gogram.ResponseHandler) (error)
	GetFile(file_id string) (*gogram.TFile, error)
	DownloadFile(file_path string) (io.ReadCloser, error)
	AnswerCallbackQuery(query_id, notification string, show_alert bool) (error)

	GetUpdates() ([]gogram.TUpdate, error)
}

