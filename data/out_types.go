package data

import (
)

type OMessage struct {
	ChatID              interface{}
	Text                string
	ParseMode           string
	EnableWebPreview    bool	// I decided to flip this one because in general, disabling it is a more sensible default.
	DisableNotification bool
	ReplyTo            *int
	ReplyMarkup         interface{}
}
