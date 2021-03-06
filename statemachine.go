package gogram

import (
	"github.com/thewug/gogram/data"

	"strings"
)

type State interface {
	Handle(*MessageCtx)
	HandleCallback(*CallbackCtx)

	StateExited(data.Sender)
	StateEntered(data.Sender)
}

type StateBase struct {
	StateMachine *MessageStateMachine
}

func MakeBase(machine *MessageStateMachine) StateBase {
	var sb StateBase
	sb.StateMachine = machine
	return sb
}

func (this *StateBase) HandleCallback(ctx *CallbackCtx) {
	return // default implementation does nothing
}

func (this *StateBase) Handle(ctx *CallbackCtx) {
	return // default implementation does nothing
}

func (this *StateBase) StateExited(sender data.Sender) {
	return // default implementation does nothing
}

func (this *StateBase) StateEntered(sender data.Sender) {
	return // default implementation does nothing
}

type MessageStateMachine struct {
	StateBase
	UserStates     map[data.Sender]State
	Handlers       map[string]State
	Default        State
}

func NewMessageStateMachine() (*MessageStateMachine) {
	csm := MessageStateMachine{
		UserStates: make(map[data.Sender]State),
		Handlers: make(map[string]State),
	}
	csm.Default = &csm

	return &csm
}

func (this *MessageStateMachine) AddCommand(cmd string, state State) {
	this.Handlers[strings.ToLower(cmd)] = state
}

func (this *MessageStateMachine) DelCommand(cmd string) {
	delete(this.Handlers, strings.ToLower(cmd))
}

func (this *MessageStateMachine) GetCommand(cmd string) State {
	return this.Handlers[strings.ToLower(cmd)]
}

func (this *MessageStateMachine) SetState(sender data.Sender, newstate State) {
	oldstate := this.UserStates[sender]
	if oldstate != nil {
		oldstate.StateExited(sender)
	}

	this.SetStateDirect(sender, newstate)

	if newstate != nil {
		newstate.StateEntered(sender)
	}
}

func (this *MessageStateMachine) SetStateDirect(sender data.Sender, newstate State) {

	if newstate != nil {
		this.UserStates[sender] = newstate
	} else {
		delete(this.UserStates, sender)
	}
}

func (this *MessageStateMachine) ProcessMessage(ctx *MessageCtx) {
	state, _ := this.UserStates[ctx.Msg.Sender()]
	if state == nil { state = this.Default }

	state.Handle(ctx)
}

func (this *MessageStateMachine) ProcessCallback(ctx *CallbackCtx) {
	state, _ := this.UserStates[ctx.Cb.Sender()]
	if state == nil { state = this.Default }

	state.HandleCallback(ctx)
}

func (this *MessageStateMachine) Handle(ctx *MessageCtx) {
	if !ctx.Bot.IsMyCommand(&ctx.Cmd) || len(ctx.Cmd.Command) == 0 {
		return
	}

	callback := this.Handlers[strings.ToLower(ctx.Cmd.Command)]
	if callback != nil { callback.Handle(ctx) }
}

func (this *MessageStateMachine) HandleCallback(ctx *CallbackCtx) {
	if !ctx.Bot.IsMyCommand(&ctx.Cmd) || len(ctx.Cmd.Command) == 0 {
		return
	}

	callback := this.Handlers[strings.ToLower(ctx.Cmd.Command)]
	if callback != nil { callback.HandleCallback(ctx) }
}
