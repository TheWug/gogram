package gogram

import (
	"github.com/thewug/gogram/data"

	"log"
	"time"
	"os"
	"strings"

	"io/ioutil"
	"encoding/json"
)


type TelegramBot struct {
	callback_callback Callbackable
	message_callback Messagable
	inline_callback InlineQueryable
	maintenance_callbacks []Maintainer

	state_machine *MessageStateMachine

	update_channel chan []data.TUpdate
	maintenance_ticker *time.Ticker
	settings InitSettings

	log *log.Logger
	errorlog *log.Logger

	remote Protocol
}

func (this *TelegramBot) Log() (*log.Logger) {
	return this.log
}

func (this *TelegramBot) SetLog(log *log.Logger) {
	this.log = log
}

func (this *TelegramBot) SetErrorLog(log *log.Logger) {
	this.errorlog = log
}

func (this *TelegramBot) ErrorLog() (*log.Logger) {
	return this.errorlog
}

func (this *TelegramBot) Remote() (*Protocol) {
	return &this.remote
}

func (this *TelegramBot) SetStateMachine(m *MessageStateMachine) {
	this.state_machine = m
}

func (this *TelegramBot) SetMessageCallback(cb Messagable) {
	this.message_callback = cb
}

func (this *TelegramBot) SetInlineCallback(cb InlineQueryable) {
	this.inline_callback = cb
}

func (this *TelegramBot) SetCallbackCallback(cb Callbackable) {
	this.callback_callback = cb
}

func (this *TelegramBot) AddMaintenanceCallback(cb Maintainer) {
	this.maintenance_callbacks = append(this.maintenance_callbacks, cb)
}

func (this *TelegramBot) asyncUpdateLoop(output chan []data.TUpdate) () {
	for {
		updates, e := this.remote.GetUpdates()
		if e != nil {
			this.errorlog.Printf("Error (async update loop): %s\n", e.Error())
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}

		output <- updates
	}
}

func (this *TelegramBot) Init(filename string, s InitSettings) (error) {
	this.log = log.New(os.Stdout, "", log.LstdFlags)
	this.errorlog = this.log
	this.remote = NewProtocol()

	bytes, e := ioutil.ReadFile(filename)
	if e != nil { return e }

	e = json.Unmarshal(bytes, s)
	if e != nil { return e }

	this.settings = s
	e = this.settings.InitializeAll(this)
	return e
}

func (this *TelegramBot) MainLoop() {
	this.update_channel = make(chan []data.TUpdate, 3)
	go this.asyncUpdateLoop(this.update_channel)

	this.maintenance_ticker = time.NewTicker(time.Second)

	var seconds int64 = 0
	for {
		select {
		case updates := <- this.update_channel:
			for _, u := range updates {
				if u.Message != nil && this.message_callback != nil {
					this.message_callback.ProcessMessage(NewMessageCtx(u.Message, false, this))
				}
				if u.Edited_message != nil && this.message_callback != nil {
					this.message_callback.ProcessMessage(NewMessageCtx(u.Edited_message, true, this))
				}
				if u.Channel_post != nil && this.message_callback != nil {
					this.message_callback.ProcessMessage(NewMessageCtx(u.Channel_post, false, this))
				}
				if u.Edited_channel_post != nil && this.message_callback != nil {
					this.message_callback.ProcessMessage(NewMessageCtx(u.Edited_channel_post, true, this))
				}
				if u.Inline_query != nil && this.inline_callback != nil {
					this.inline_callback.ProcessInlineQuery(this, u.Inline_query)
				}
				if u.Chosen_inline_result != nil && this.inline_callback != nil {
					this.inline_callback.ProcessInlineQueryResult(this, u.Chosen_inline_result)
				}
				if u.Callback_query != nil && this.callback_callback != nil {
					this.callback_callback.ProcessCallback(this, u.Callback_query)
				}
			}
		case <- this.maintenance_ticker.C:
			for _, m := range this.maintenance_callbacks {
				if (seconds % m.GetInterval() == 0) { m.DoMaintenance(this) }
			}
			seconds++
		case hbox := <- CallResponseChannel:
			if hbox.Error != nil {
				log.Println(hbox.Error.Error())
			}
			if hbox.Handler != nil {
				hbox.Handler.Callback(hbox.Output, hbox.Success, hbox.Error, hbox.Http_code)
			}
		}
	}
}

func (this *TelegramBot) IsMyCommand(cmd *CommandData) (bool) {
	return strings.ToLower(*this.remote.GetMe().Username) == strings.ToLower(cmd.Target) || len(cmd.Target) == 0
}

type MessageStateMachine struct {
	UserStates     map[data.Sender]State
	Handlers       map[string]State
	Default        State
}

type State interface {
	Handle(*MessageCtx)
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

func (this *MessageStateMachine) Handle(ctx *MessageCtx) {
	if !ctx.Bot.IsMyCommand(&ctx.Cmd) || len(ctx.Cmd.Command) == 0 {
		return
	}

	callback := this.Handlers[strings.ToLower(ctx.Cmd.Command)]
	if callback != nil { callback.Handle(ctx) }
}
