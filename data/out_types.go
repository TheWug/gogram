package data

import (
)

const HTML string = "HTML"
const Markdown string = "Markdown"

type OMessage struct {
	// used for all outgoing messages
	ChatID              interface{} // not used for editing messages if InlineID is present
	Text                string
	ParseMode           string
	EnableWebPreview    bool	// I decided to flip this one because in general, disabling it is a more sensible default.
	ReplyMarkup         interface{}

	// only used for new messages
	DisableNotification bool
	ReplyTo            *int

	// only used for message edits
	MessageID           int
	InlineID            string
}
