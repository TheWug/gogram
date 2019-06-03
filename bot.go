package gogram

import (
	"github.com/thewug/gogram/data"
)


type InlineQueryable interface {
	ProcessInlineQuery(*TelegramBot, *data.TInlineQuery)
	ProcessInlineQueryResult(*TelegramBot, *data.TChosenInlineResult)
}


type Callbackable interface {
	ProcessCallback(*TelegramBot, *data.TCallbackQuery)
}


type Messagable interface {
	ProcessMessage(*TelegramBot, *data.TMessage, bool)
}


type Maintainer interface {
	DoMaintenance(*TelegramBot)
	GetInterval() int64
}


type InitSettings interface {
	InitializeAll() error
}


// DEPRECATED
type Command interface {
	Callback(*CommandData)
}

type CommandData struct {
	M      *data.TMessage
	Command string
	Target  string
	Line	string
	Argstr  string
	Args  []string
}