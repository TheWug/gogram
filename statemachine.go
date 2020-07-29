package gogram

import (
	"github.com/thewug/gogram/data"

	"strings"
)

type MessageStateMachine struct {
	UserStates     map[data.Sender]State
	Handlers       map[string]State
	Default        State
}

type State interface {
	Handle(*MessageCtx)
	HandleCallback(*CallbackCtx)
}

type StateBase struct {
	StateMachine *MessageStateMachine
}

type StateIgnoreCallbacks struct {
}

func (this *StateIgnoreCallbacks) HandleCallback(ctx *CallbackCtx) {
	return // do nothing
}

func (this *StateBase) HandleCallback(ctx *CallbackCtx) {
	return // default implementation does nothing
}

type StateIgnoreMessages struct {
}

func (this *StateBase) Handle(ctx *CallbackCtx) {
	return // default implementation does nothing
}

func (this *StateIgnoreMessages) HandleCallback(ctx *CallbackCtx) {
	return // do nothing
}

func (this *StateBase) SetState(sender data.Sender, state State) {
	this.StateExited(sender)
	this.StateMachine.SetState(sender, state)
	if state != nil {
		state.StateEntered(sender)
	}
}

func (this *StateBase) StateExited(sender data.Sender) {
	return // default implementation does nothing
}

func (this *StateBase) StateEntered(sender data.Sender) {
	return // default implementation does nothing
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

func (this *MessageStateMachine) SetState(sender data.Sender, state State) {
	if state != nil {
		this.UserStates[sender] = state
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
