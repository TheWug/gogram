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
	Cb      *data.TCallbackQuery
	Bot     *TelegramBot

	answered bool
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

func (this *MessageCtx) DispatchCommand() bool {
	if len(this.Cmd.Command) > 0 {
		if this.Bot.IsMyCommand(&this.Cmd) {
			callback := this.Machine.Handlers[strings.ToLower(this.Cmd.Command)]
			if callback != nil {
				callback.Handle(this)
				return true
			}
		}
	}
	return false
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

// Send a message to the same chat that a message came from but without directly replying to it.
func (this *MessageCtx) Respond(m data.OMessage) (*MessageCtx, error) {
	m.ChatId = this.Msg.Chat.Id
	msg, err := this.Bot.Remote.SendMessage(m)
	return &MessageCtx {
		Msg: msg,
		Bot: this.Bot,
		Machine: this.Machine,
	}, err
}

// Reply to the specified message.
func (this *MessageCtx) Reply(m data.OMessage) (*MessageCtx, error) {
	m.ChatId = this.Msg.Chat.Id
	m.ReplyToId = &this.Msg.Id
	msg, err := this.Bot.Remote.SendMessage(m)
	return &MessageCtx {
		Msg: msg,
		Bot: this.Bot,
		Machine: this.Machine,
	}, err
}

// Reply to the specified message if it's a PM, otherwise PM the sender.
func (this *MessageCtx) ReplyOrPM(m data.OMessage) (*MessageCtx, error) {
	if this.Msg.Chat.Type == data.Channel { return nil, errors.New("Can't privately reply to a channel message!") }

	m.ChatId = this.Msg.From.Id
	if this.Msg.Chat.Type == data.Private {
		m.ReplyToId = &this.Msg.Id
	}
	msg, err := this.Bot.Remote.SendMessage(m)
	return &MessageCtx {
		Msg: msg,
		Bot: this.Bot,
		Machine: this.Machine,
	}, err
}

// pm the message sender.
func (this *MessageCtx) PM(m data.OMessage) (*MessageCtx, error) {
	if this.Msg.Chat.Type == data.Channel { return nil, errors.New("Can't privately reply to a channel message!") }

	m.ChatId = this.Msg.From.Id
	msg, err := this.Bot.Remote.SendMessage(m)
	return &MessageCtx {
		Msg: msg,
		Bot: this.Bot,
		Machine: this.Machine,
	}, err
}

// forward this message to another chat
func (this *MessageCtx) Forward(f data.OForward) (*MessageCtx, error) {
	f.SourceChatId = this.Msg.Chat.Id
	f.SourceMessageId = this.Msg.Id
	msg, err := this.Bot.Remote.ForwardMessage(f)
	return &MessageCtx {
		Msg: msg,
		Bot: this.Bot,
		Machine: this.Machine,
	}, err
}

// edit this message.
func (this *MessageCtx) EditText(e data.OMessageEdit) (*MessageCtx, error) {
	e.SourceChatId = this.Msg.Chat.Id
	e.SourceMessageId = this.Msg.Id
	msg, err := this.Bot.Remote.EditMessageText(e)
	return &MessageCtx {
		Msg: msg,
		Bot: this.Bot,
		Machine: this.Machine,
	}, err
}

// delete this message
func (this *MessageCtx) Delete() (error) {
	return this.Bot.Remote.DeleteMessage(data.ODelete{SourceChatId: this.Msg.Chat.Id, SourceMessageId: this.Msg.Id})
}

// kick the sender of this message from the group they sent it to
func (this *MessageCtx) KickSender() (error) {
	if this.Msg.Chat.Type == data.Group || this.Msg.Chat.Type == data.Supergroup {
		return this.Bot.Remote.KickMember(data.OChatMember{TargetData: data.TargetData{ChatId: this.Msg.Chat.Id}, UserId: this.Msg.From.Id})
	} else {
		return errors.New("Tried to kick message sender from channel or PM")
	}
}

// fetch info about the sender of this message
func (this *MessageCtx) Member() (*ChatMemberCtx, error) {
	if this.Msg.Chat.Type == data.Group || this.Msg.Chat.Type == data.Supergroup {
		member, err := this.Bot.Remote.GetChatMember(data.OChatMember{TargetData: data.TargetData{ChatId: this.Msg.Chat.Id}, UserId: this.Msg.From.Id})
		return &ChatMemberCtx{
			Member: member,
			Bot: this.Bot,
		}, err
	} else {
		return nil, errors.New("Tried to fetch chat info for sender from channel or PM")
	}
}

func (this *MessageCtx) RespondAsync(m data.OMessage, handler data.ResponseHandler) {
	m.ChatId = this.Msg.Chat.Id
	this.Bot.Remote.SendMessageAsync(m, handler)
}

func (this *MessageCtx) ReplyAsync(m data.OMessage, handler data.ResponseHandler) {
	m.ChatId = this.Msg.Chat.Id
	m.ReplyToId = &this.Msg.Id
	this.Bot.Remote.SendMessageAsync(m, handler)
}

func (this *MessageCtx) ReplyOrPMAsync(m data.OMessage, handler data.ResponseHandler) {
	if this.Msg.Chat.Type != data.Channel {
		m.ChatId = this.Msg.From.Id
		if this.Msg.Chat.Type == data.Private {
			m.ReplyToId = &this.Msg.Id
		}
		this.Bot.Remote.SendMessageAsync(m, handler)
	} else if handler != nil {
		handler.Callback(nil, false, errors.New("Can't privately reply to a channel message!"), 0)
	}
}

func (this *MessageCtx) PMAsync(m data.OMessage, handler data.ResponseHandler) {
	if this.Msg.Chat.Type != data.Channel {
		m.ChatId = this.Msg.From.Id
		this.Bot.Remote.SendMessageAsync(m, handler)
	} else if handler != nil {
		handler.Callback(nil, false, errors.New("Can't PM to a channel message sender!"), 0)
	}
}

func (this *MessageCtx) ForwardAsync(m data.OForward, handler data.ResponseHandler) {
	m.ChatId = this.Msg.Chat.Id
	m.SourceMessageId = this.Msg.Id
	this.Bot.Remote.ForwardMessageAsync(m, handler)
}

func (this *MessageCtx) EditTextAsync(m data.OMessageEdit, handler data.ResponseHandler) {
	m.ChatId = this.Msg.Chat.Id
	m.SourceMessageId = this.Msg.Id
	this.Bot.Remote.EditMessageTextAsync(m, handler)
}

func (this *MessageCtx) DeleteAsync(handler data.ResponseHandler) {
	this.Bot.Remote.DeleteMessageAsync(data.ODelete{SourceChatId: this.Msg.Chat.Id, SourceMessageId: this.Msg.Id}, handler)
}

func (this *MessageCtx) KickSenderAsync(handler data.ResponseHandler) {
	if this.Msg.Chat.Type == data.Group || this.Msg.Chat.Type == data.Supergroup {
		this.Bot.Remote.KickMemberAsync(data.OChatMember{TargetData: data.TargetData{ChatId: this.Msg.Chat.Id}, UserId: this.Msg.From.Id}, handler)
	} else if handler != nil {
		handler.Callback(nil, false, errors.New("Tried to kick message sender from channel or PM"), 0)
	}
}

func (this *MessageCtx) MemberAsync(handler data.ResponseHandler) {
	if this.Msg.Chat.Type == data.Group || this.Msg.Chat.Type == data.Supergroup {
		this.Bot.Remote.GetChatMemberAsync(data.OChatMember{TargetData: data.TargetData{ChatId: this.Msg.Chat.Id}, UserId: this.Msg.From.Id}, handler)
	} else if handler != nil {
		handler.Callback(nil, false, errors.New("Tried to fetch chat info for sender from channel or PM"), 0)
	}
}

func (this *InlineCtx) Answer(o data.OInlineQueryAnswer) (error) {
	o.Id = this.Query.Id
	return this.Bot.Remote.AnswerInlineQuery(o)
}

func (this *InlineCtx) AnswerAsync(o data.OInlineQueryAnswer, handler data.ResponseHandler) {
	o.Id = this.Query.Id
	this.Bot.Remote.AnswerInlineQueryAsync(o, handler)
}

func (this *CallbackCtx) Answer(o data.OCallback) (error) {
	if !this.answered {
		this.answered = true
		o.Id = this.Cb.Id
		return this.Bot.Remote.AnswerCallbackQuery(o)
	} else {
		return errors.New("callback query already answered")
	}
}

func (this *CallbackCtx) AnswerAsync(o data.OCallback, handler data.ResponseHandler) {
	if !this.answered {
		this.answered = true
		o.Id = this.Cb.Id
		this.Bot.Remote.AnswerCallbackQueryAsync(o, handler)
	} else if handler != nil {
		handler.Callback(nil, false, errors.New("callback query already answered"), 0)
	}
}
