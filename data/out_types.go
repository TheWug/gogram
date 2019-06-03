package data

import (
)

const HTML string = "HTML"
const Markdown string = "Markdown"

type OMessage struct {
	ChatID              interface{} // types: int, int64, or string
	TargetChatID        interface{} // "
	MessageID           int
	InlineID            string
	ReplyTo            *int
	Text                string
	Sticker             interface{} // types: string, io.Reader, []byte, or reqtify.FormFile
	ParseMode           string
	EnableWebPreview    bool
	ReplyMarkup         interface{} // types: TInlineKeyboard, TReplyMarkup, TReplyKeyboardRemove, or TForceReply
	DisableNotification bool
}

type OChatMember struct {
	ChatID interface{} // types: int, int64, or string
	UserID int
}

type OStickerSet struct {
	Name string
}

type OInlineQueryAnswer struct {
	QueryID    string
	Results  []interface{} // types: array of TInlineQueryResult*
	NextOffset string
	CacheTime  int
}

type OCallback struct {
	QueryID      string
	Notification string
	ShowAlert    bool
	CacheTime    int
	URL          string
}

type ORestrict struct {
	ChatID             interface{} // types: int, int64, or string
	UserID             int
	Until              int64
	CanSendMessages    bool
	CanSendMedia       bool
	CanSendInline      bool
	CanSendWebPreviews bool
}
