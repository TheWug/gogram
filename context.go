package gogram

import (
	"github.com/thewug/gogram/data"

	"github.com/kballard/go-shellquote"

	"strings"
)

type CommandData struct {
	Command    string
	Target     string
	Line	   string
	Argstr     string
	Args     []string
	ParseError error
}

func ParseCommand(m *data.TMessage) (CommandData) {
	var line string
	if m.Text != nil && *m.Text != "" {
		line = *m.Text
	} else if m.Caption != nil {
		line = *m.Caption
	}

	return ParseCommandFromString(line)
}

func ParseCommandFromString(line string) (CommandData) {
	var c CommandData

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
	c.Args, c.ParseError = shellquote.Split(line)
	return c
}

type MessageCtx struct {
	Msg          *data.TMessage
	Edited        bool
	Cmd           CommandData
	Bot          *TelegramBot
	Machine      *MessageStateMachine
}

func NewMessageCtx(msg *data.TMessage, edited bool, bot *TelegramBot) (*MessageCtx) {
	return &MessageCtx{
		Msg: msg,
		Edited: edited,
		Cmd: ParseCommand(msg),
		Bot: bot,
		Machine: bot.state_machine,
	}
}

func (this *MessageCtx) SetState(newstate State) {
	if this.Machine == nil {
		panic("Tried to set state, but there was no state machine!")
	}
	this.Machine.SetState(this.Msg.Sender(), newstate)
}

func (this *MessageCtx) GetState() (State) {
	if this.Machine == nil {
		panic("Tried to get state, but there was no state machine!")
	}
	state, _ := this.Machine.UserStates[this.Msg.Sender()]
	return state
}

func (this *MessageCtx) Respond(m data.OMessage) (*MessageCtx, error) {
	msg, err := this.Bot.Remote().SendMessage(this.Msg.Chat.Id, m.Text, m.ReplyTo, m.ParseMode, m.ReplyMarkup, !m.EnableWebPreview)
	return &MessageCtx {
		Msg: msg,
		Bot: this.Bot,
	}, err
}

func (this *MessageCtx) Reply(m data.OMessage) (*MessageCtx, error) {
	msg, err := this.Bot.Remote().SendMessage(this.Msg.Chat.Id, m.Text, &this.Msg.Message_id, m.ParseMode, m.ReplyMarkup, !m.EnableWebPreview)
	return &MessageCtx {
		Msg: msg,
		Bot: this.Bot,
	}, err
}

func (this *MessageCtx) Delete() (error) {
	return this.Bot.Remote().DeleteMessage(this.Msg.Chat.Id, this.Msg.Message_id)
}

func (this *MessageCtx) RespondAsync(m data.OMessage, handler data.ResponseHandler) {
	this.Bot.Remote().SendMessageAsync(this.Msg.Chat.Id, m.Text, m.ReplyTo, m.ParseMode, m.ReplyMarkup, !m.EnableWebPreview, handler)
}

func (this *MessageCtx) ReplyAsync(m data.OMessage, handler data.ResponseHandler) {
	this.Bot.Remote().SendMessageAsync(this.Msg.Chat.Id, m.Text, &this.Msg.Message_id, m.ParseMode, m.ReplyMarkup, !m.EnableWebPreview, handler)
}

func (this *MessageCtx) DeleteAsync(handler data.ResponseHandler) {
	this.Bot.Remote().DeleteMessageAsync(this.Msg.Chat.Id, this.Msg.Message_id, handler)
}
