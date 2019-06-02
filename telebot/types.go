package telebot

import (
	"github.com/thewug/gogram"
	"github.com/thewug/gogram/bot"

	"log"
	"time"
	"os"
	"strings"

	"io/ioutil"
	"encoding/json"

	"github.com/kballard/go-shellquote"
)


type TelegramBot struct {
	callback_callback bot.Callbackable
	message_callback bot.Messagable
	inline_callback bot.InlineQueryable
	maintenance_callbacks []bot.Maintainer

	update_channel chan []gogram.TUpdate
	maintenance_ticker *time.Ticker
	settings bot.InitSettings

	log *log.Logger
	errorlog *log.Logger

	remote gogram.Protocol
	commands map[string]bot.Command
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

func (this *TelegramBot) Remote() (bot.Protocol) {
	return &this.remote
}

func (this *TelegramBot) SetMessageCallback(cb bot.Messagable) {
	this.message_callback = cb
}

func (this *TelegramBot) SetInlineCallback(cb bot.InlineQueryable) {
	this.inline_callback = cb
}

func (this *TelegramBot) SetCallbackCallback(cb bot.Callbackable) {
	this.callback_callback = cb
}

func (this *TelegramBot) AddMaintenanceCallback(cb bot.Maintainer) {
	this.maintenance_callbacks = append(this.maintenance_callbacks, cb)
}

func (this *TelegramBot) AddCommand(cmd string, cb bot.Command) {
	if this.commands == nil {
		this.commands = make(map[string]bot.Command)
	}
	this.commands[strings.ToLower(cmd)] = cb
}

func (this *TelegramBot) asyncUpdateLoop(output chan []gogram.TUpdate) () {
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

func (this *TelegramBot) Init(filename string, s bot.InitSettings) (error) {
	this.log = log.New(os.Stdout, "", log.LstdFlags)
	this.errorlog = this.log
	this.remote = gogram.NewProtocol()

	bytes, e := ioutil.ReadFile(filename)
	if e != nil { return e }

	e = json.Unmarshal(bytes, s)
	if e != nil { return e }

	this.settings = s
	e = this.settings.InitializeAll()
	return e
}

func (this *TelegramBot) MainLoop() {
	this.update_channel = make(chan []gogram.TUpdate, 3)
	go this.asyncUpdateLoop(this.update_channel)

	this.maintenance_ticker = time.NewTicker(time.Second)

	var seconds int64 = 0
	for {
		select {
		case updates := <- this.update_channel:
			for _, u := range updates {
				if u.Inline_query != nil && this.inline_callback != nil {
					this.inline_callback.ProcessInlineQuery(this, u.Inline_query)
				}
				if u.Chosen_inline_result != nil && this.inline_callback != nil {
					this.inline_callback.ProcessInlineQueryResult(this, u.Chosen_inline_result)
				}
				if u.Message != nil && this.message_callback != nil {
					this.message_callback.ProcessMessage(this, u.Message, false)
				}
				if u.Edited_message != nil && this.message_callback != nil {
					this.message_callback.ProcessMessage(this, u.Edited_message, true)
				}
				if u.Callback_query != nil && this.callback_callback != nil {
					this.callback_callback.ProcessCallback(this, u.Callback_query)
				}
				if u.Channel_post != nil && this.message_callback != nil {
					this.message_callback.ProcessMessage(this, u.Channel_post, false)
				}
				if u.Edited_channel_post != nil && this.message_callback != nil {
					this.message_callback.ProcessMessage(this, u.Edited_channel_post, true)
				}
			}
		case <- this.maintenance_ticker.C:
			for _, m := range this.maintenance_callbacks {
				if (seconds % m.GetInterval() == 0) { m.DoMaintenance(this) }
			}
			seconds++
		case hbox := <- gogram.CallResponseChannel:
			if hbox.Error != nil {
				log.Println(hbox.Error.Error())
			}
			if hbox.Handler != nil {
				hbox.Handler.Callback(hbox.Output, hbox.Success, hbox.Error, hbox.Http_code)
			}
		}
	}
}

// bot command main handler.
func (this *TelegramBot) HandleCommand(m *gogram.TMessage) () {
	cmd, err := ParseCommand(m)

	if err != nil { return }

	// command directed at another user
	if !this.IsMyCommand(&cmd) { return }

	if this.commands == nil { return }
	callback, has := this.commands[cmd.Command]

	// tried to use nonexistent command
	if !has { return }

	// run the command
	callback.Callback(&cmd)
}

func (this *TelegramBot) IsMyCommand(cmd *bot.CommandData) (bool) {
	return strings.ToLower(*this.remote.GetMe().Username) == strings.ToLower(cmd.Target) || len(cmd.Target) == 0
}

type MsgContext struct {
	Cmd      bot.CommandData
	CmdError error
	Msg      gogram.TMessage
	MsgEdit  bool
	Bot     *TelegramBot
	Machine *MessageStateMachine
}

func (this *MsgContext) SetState(newstate State) {
	this.Machine.SetState(this.Msg.Sender(), newstate)
}

func (this *MsgContext) GetState() (State) {
	state, _ := this.Machine.UserStates[this.Msg.Sender()]
	return state
}

type MessageStateMachine struct {
	UserStates     map[gogram.Sender]State
	Handlers       map[string]State
	Default        State
}

type State interface {
	Handle(*MsgContext)
}

func NewMessageStateMachine() (*MessageStateMachine) {
	csm := MessageStateMachine{
		UserStates: make(map[gogram.Sender]State),
		Handlers: make(map[string]State),
	}
	csm.Default = &csm

	return &csm
}

func (this *MessageStateMachine) AddCommand(cmd string, state State) {
	this.Handlers[strings.ToLower(cmd)] = state
}

func (this *MessageStateMachine) SetState(sender gogram.Sender, state State) {
	if state != nil {
		this.UserStates[sender] = state
	} else {
		delete(this.UserStates, sender)
	}
}

func (this *MessageStateMachine) ProcessMessage(bot *TelegramBot, msg *gogram.TMessage, edited bool) {
	var ctx MsgContext
	ctx.Msg = *msg
	ctx.MsgEdit = edited
	ctx.Cmd, ctx.CmdError = ParseCommand(&ctx.Msg)
	ctx.Bot = bot
	ctx.Machine = this

	this.FeedContext(&ctx)
}

func (this *MessageStateMachine) FeedContext(ctx *MsgContext) {
	state, _ := this.UserStates[ctx.Msg.Sender()]
	if state == nil { state = this.Default }

	state.Handle(ctx)
}

func (this *MessageStateMachine) Handle(ctx *MsgContext) {
	if !ctx.Bot.IsMyCommand(&ctx.Cmd) || len(ctx.Cmd.Command) == 0 {
		return
	}

	callback := this.Handlers[strings.ToLower(ctx.Cmd.Command)]
	if callback != nil { callback.Handle(ctx) }
}


func ParseCommand(m *gogram.TMessage) (bot.CommandData, error) {
	var line string
	if m.Text != nil && *m.Text != "" {
		line = *m.Text
	} else if m.Caption != nil {
		line = *m.Caption
	}

	c, e := ParseCommandFromString(line)
	c.M = m
	return c, e
}

func ParseCommandFromString(line string) (bot.CommandData, error) {
	var c bot.CommandData
	var err error

	c.Line = line
	if strings.HasPrefix(line, "/") {
		tokens := strings.SplitN(line, " ", 2)
		if len(tokens) == 2 {
			line = tokens[1]
		} else {
			line = ""
		}
		command := tokens[0]
		tokens = strings.SplitN(command, "@", 2)
		if len(tokens) == 2 {
			c.Target = tokens[1]
		} else {
			c.Target = ""
		}
		c.Command = tokens[0]
	}

	c.Argstr = line
	c.Args, err = shellquote.Split(line)
	return c, err
}
