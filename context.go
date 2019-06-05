package gogram

import (
	"github.com/thewug/gogram/data"

	"github.com/kballard/go-shellquote"

	"strings"
	"errors"
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

type ChatMemberCtx struct {
	Member *data.TChatMember
	Bot    *TelegramBot
}

type InlineCtx struct {
	Query *data.TInlineQuery
	Bot   *TelegramBot
}

type InlineResultCtx struct {
	Result *data.TChosenInlineResult
	Bot    *TelegramBot
}

type CallbackCtx struct {
	Cb  *data.TCallbackQuery
	Bot *TelegramBot
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
	m.ChatID = this.Msg.Chat.Id
	msg, err := this.Bot.Remote.SendMessage(m)
	return &MessageCtx {
		Msg: msg,
		Bot: this.Bot,
	}, err
}

func (this *MessageCtx) Reply(m data.OMessage) (*MessageCtx, error) {
	m.ChatID = this.Msg.Chat.Id
	m.ReplyTo = &this.Msg.Message_id
	msg, err := this.Bot.Remote.SendMessage(m)
	return &MessageCtx {
		Msg: msg,
		Bot: this.Bot,
	}, err
}

func (this *MessageCtx) ReplyOrPM(m data.OMessage) (*MessageCtx, error) {
	if this.Msg.Chat.Type == data.Channel { return nil, errors.New("Can't privately reply to a channel message!") }

	m.ChatID = this.Msg.From.Id
	if this.Msg.Chat.Type == data.Private {
		m.ReplyTo = &this.Msg.Message_id
	}
	msg, err := this.Bot.Remote.SendMessage(m)
	return &MessageCtx {
		Msg: msg,
		Bot: this.Bot,
	}, err
}

func (this *MessageCtx) EditText(m data.OMessage) (*MessageCtx, error) {
	m.ChatID = this.Msg.Chat.Id
	m.MessageID = this.Msg.Message_id
	msg, err := this.Bot.Remote.EditMessageText(m)
	return &MessageCtx {
		Msg: msg,
		Bot: this.Bot,
	}, err
}

func (this *MessageCtx) Delete() (error) {
	return this.Bot.Remote.DeleteMessage(data.OMessage{ChatID: this.Msg.Chat.Id, MessageID: this.Msg.Message_id})
}

func (this *MessageCtx) KickSender() (error) {
	if this.Msg.Chat.Type == data.Group || this.Msg.Chat.Type == data.Supergroup {
		return this.Bot.Remote.KickMember(data.OChatMember{ChatID: this.Msg.Chat.Id, UserID: this.Msg.From.Id})
	} else {
		return errors.New("Tried to kick message sender from channel or PM")
	}
}

func (this *MessageCtx) Member() (*ChatMemberCtx, error) {
	if this.Msg.Chat.Type == data.Group || this.Msg.Chat.Type == data.Supergroup {
		member, err := this.Bot.Remote.GetChatMember(data.OChatMember{ChatID: this.Msg.Chat.Id, UserID: this.Msg.From.Id})
		return &ChatMemberCtx{
			Member: member,
			Bot: this.Bot,
		}, err
	} else {
		return nil, errors.New("Tried to fetch chat info for sender from channel or PM")
	}
}

func (this *MessageCtx) RespondAsync(m data.OMessage, handler data.ResponseHandler) {
	m.ChatID = this.Msg.Chat.Id
	this.Bot.Remote.SendMessageAsync(m, handler)
}

func (this *MessageCtx) ReplyAsync(m data.OMessage, handler data.ResponseHandler) {
	m.ChatID = this.Msg.Chat.Id
	m.ReplyTo = &this.Msg.Message_id
	this.Bot.Remote.SendMessageAsync(m, handler)
}

func (this *MessageCtx) ReplyOrPMAsync(m data.OMessage, handler data.ResponseHandler) {
	if this.Msg.Chat.Type != data.Channel {
		m.ChatID = this.Msg.From.Id
		if this.Msg.Chat.Type == data.Private {
			m.ReplyTo = &this.Msg.Message_id
		}
		this.Bot.Remote.SendMessageAsync(m, handler)
	} else if handler != nil {
		handler.Callback(nil, false, errors.New("Can't privately reply to a channel message!"), 0)
	}
}

func (this *MessageCtx) EditTextAsync(m data.OMessage, handler data.ResponseHandler) {
	m.ChatID = this.Msg.Chat.Id
	m.MessageID = this.Msg.Message_id
	this.Bot.Remote.EditMessageTextAsync(m, handler)
}

func (this *MessageCtx) DeleteAsync(handler data.ResponseHandler) {
	this.Bot.Remote.DeleteMessageAsync(data.OMessage{ChatID: this.Msg.Chat.Id, MessageID: this.Msg.Message_id}, handler)
}

func (this *MessageCtx) KickSenderAsync(handler data.ResponseHandler) {
	if this.Msg.Chat.Type == data.Group || this.Msg.Chat.Type == data.Supergroup {
		this.Bot.Remote.KickMemberAsync(data.OChatMember{ChatID: this.Msg.Chat.Id, UserID: this.Msg.From.Id}, handler)
	} else if handler != nil {
		handler.Callback(nil, false, errors.New("Tried to kick message sender from channel or PM"), 0)
	}
}

func (this *MessageCtx) MemberAsync(handler data.ResponseHandler) {
	if this.Msg.Chat.Type == data.Group || this.Msg.Chat.Type == data.Supergroup {
		this.Bot.Remote.GetChatMemberAsync(data.OChatMember{ChatID: this.Msg.Chat.Id, UserID: this.Msg.From.Id}, handler)
	} else if handler != nil {
		handler.Callback(nil, false, errors.New("Tried to fetch chat info for sender from channel or PM"), 0)
	}
}

func (this *InlineCtx) Answer(o data.OInlineQueryAnswer) (error) {
	o.QueryID = this.Query.Id
	return this.Bot.Remote.AnswerInlineQuery(o)
}

func (this *InlineCtx) AnswerAsync(o data.OInlineQueryAnswer, handler data.ResponseHandler) {
	o.QueryID = this.Query.Id
	this.Bot.Remote.AnswerInlineQueryAsync(o, handler)
}
