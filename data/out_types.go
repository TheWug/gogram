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
