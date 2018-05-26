package telebot

import (
	"telegram"
)

type CommandData struct {
	M      *telegram.TMessage
	Command string
	Target  string
	Line	string
	Argstr  string
	Args  []string
}

type Command interface {
	Callback(*CommandData)
}
