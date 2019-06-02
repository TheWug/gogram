package bot

import (
	"github.com/thewug/gogram"

	"log"
)


type InlineQueryable interface {
	ProcessInlineQuery(Bot, *gogram.TInlineQuery)
	ProcessInlineQueryResult(Bot, *gogram.TChosenInlineResult)
}


type Callbackable interface {
	ProcessCallback(Bot, *gogram.TCallbackQuery)
}


type Messagable interface {
	ProcessMessage(Bot, *gogram.TMessage, bool)
}


type Maintainer interface {
	DoMaintenance(Bot)
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
	M      *gogram.TMessage
	Command string
	Target  string
	Line	string
	Argstr  string
	Args  []string
}


type Bot interface {
	SetMessageCallback(cb Messagable)
	SetInlineCallback(cb InlineQueryable)
	SetCallbackCallback(cb Callbackable)
	AddMaintenanceCallback(cb Maintainer)

	Log() (*log.Logger)
	ErrorLog() (*log.Logger)
	Remote() (Protocol)

	// DEPRECATED! instead, use MessageStateMachine.
	AddCommand(cmd string, cb Command)
	HandleCommand(m *gogram.TMessage)
}

