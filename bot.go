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

	update_channel chan *data.TUpdate
	update_confirm_channel chan bool
	maintenance_ticker *time.Ticker
	settings InitSettings

	Log *log.Logger
	ErrorLog *log.Logger

	Remote Protocol
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

func (this *TelegramBot) asyncUpdateLoop() () {
	for {
		updates, e := this.Remote.GetUpdates()
		if e != nil {
			this.ErrorLog.Printf("Error (async update loop): %s\n", e.Error())
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}

		for _, update := range updates {
			this.update_channel <- &update
			<- this.update_confirm_channel
			this.Remote.markUpdateProcessed(&update)
		}

		if len(updates) == 0 {
			this.Remote.unmarkProcessedUpdate()
		}
	}
}

func (this *TelegramBot) Init(filename string, s InitSettings) (error) {
	this.Log = log.New(os.Stdout, "", log.LstdFlags)
	this.ErrorLog = this.Log
	this.Remote = NewProtocol(this)

	this.update_channel = make(chan *data.TUpdate)
	this.update_confirm_channel = make(chan bool)

	bytes, e := ioutil.ReadFile(filename)
	if e != nil { return e }

	e = json.Unmarshal(bytes, s)
	if e != nil { return e }

	this.settings = s
	e = this.settings.InitializeAll(this)
	return e
}

func (this *TelegramBot) MainLoop() {
	go this.asyncUpdateLoop()

	this.maintenance_ticker = time.NewTicker(time.Second)

	var seconds int64 = 0
	for {
		select {
		case u := <- this.update_channel:
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
				this.inline_callback.ProcessInlineQuery(&InlineCtx{Bot: this, Query: u.Inline_query})
			}
			if u.Chosen_inline_result != nil && this.inline_callback != nil {
				this.inline_callback.ProcessInlineQueryResult(&InlineResultCtx{Bot: this, Result: u.Chosen_inline_result})
			}
			if u.Callback_query != nil && this.callback_callback != nil {
				this.callback_callback.ProcessCallback(&CallbackCtx{Bot: this, Cb: u.Callback_query})
			}
			this.update_confirm_channel <- true
		case <- this.maintenance_ticker.C:
			for _, m := range this.maintenance_callbacks {
				if (seconds % m.GetInterval() == 0) { m.DoMaintenance(this) }
			}
			seconds++
		case hbox := <- call_response_channel:
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
	return strings.ToLower(*this.Remote.GetMe().Username) == strings.ToLower(cmd.Target) || len(cmd.Target) == 0
}
