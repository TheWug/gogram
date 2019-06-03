package gogram

import (
	"github.com/thewug/gogram/data"
)

type MessageCtx struct {
	Msg          *data.TMessage
	Edited        bool
	Cmd           CommandData
	CmdParseError error
	Bot          *TelegramBot
	Machine      *MessageStateMachine
}

func NewMessageCtx(msg *data.TMessage, edited bool, bot *TelegramBot) (*MessageCtx) {
	ctx := MessageCtx{
		Msg: msg,
		Edited: edited,
		Bot: bot,
		Machine: bot.state_machine,
	}

	ctx.Cmd, ctx.CmdParseError = ParseCommand(ctx.Msg)
	return &ctx
}

func (this *MessageCtx) SetState(newstate State) {
	this.Machine.SetState(this.Msg.Sender(), newstate)
}

func (this *MessageCtx) GetState() (State) {
	state, _ := this.Machine.UserStates[this.Msg.Sender()]
	return state
}

func (this *MessageCtx) Respond(m data.OMessage) (MessageCtx, error) {
	msg, err := this.Bot.Remote().SendMessage(this.Msg.Chat.Id, m.Text, m.ReplyTo, m.ParseMode, m.ReplyMarkup, !m.EnableWebPreview)
	return MessageCtx {
		Msg: msg,
		Bot: this.Bot,
	}, err
}

func (this *MessageCtx) Reply(m data.OMessage) (MessageCtx, error) {
	msg, err := this.Bot.Remote().SendMessage(this.Msg.Chat.Id, m.Text, &this.Msg.Message_id, m.ParseMode, m.ReplyMarkup, !m.EnableWebPreview)
	return MessageCtx {
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
