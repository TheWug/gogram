package telebot

import (
	"github.com/thewug/gogram"
)

type CommandData struct {
	M      *gogram.TMessage
	Command string
	Target  string
	Line	string
	Argstr  string
	Args  []string
}

type Command interface {
	Callback(*CommandData)
}
