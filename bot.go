package gogram

import (
	"github.com/thewug/gogram/data"

	"log"
)


type InlineQueryable interface {
	ProcessInlineQuery(Bot, *data.TInlineQuery)
	ProcessInlineQueryResult(Bot, *data.TChosenInlineResult)
}


type Callbackable interface {
	ProcessCallback(Bot, *data.TCallbackQuery)
}


type Messagable interface {
	ProcessMessage(Bot, *data.TMessage, bool)
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
	M      *data.TMessage
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
	Remote() (*Protocol)

	// DEPRECATED! instead, use MessageStateMachine.
	AddCommand(cmd string, cb Command)
	HandleCommand(m *data.TMessage)
}

