package bot

import (
	"github.com/thewug/gogram"
)

type MessageCtx struct {
	Msg          *gogram.TMessage
	Edited        bool
	Cmd           bot.CommandData
	CmdParseError error
	Bot          *Bot
}

func (this *MessageCtx) Respond(m OMessage) (MessageCtx, error) {
	msg, err := m.Bot.Remote().SendMessage(this.Msg.Chat.Id, m.Text, m.ReplyTo, m.ParseMode, m.ReplyMarkup, !m.EnableWebPreview)
	return MessageCtx {
		Msg: msg,
		Bot: this.Bot,
	}, err
}

func (this *MessageCtx) Reply(m OMessage) (MessageCtx, error) {
	msg, err := m.Bot.Remote().SendMessage(this.Msg.Chat.Id, m.Text, &this.Msg.Message_id, m.ParseMode, m.ReplyMarkup, !m.EnableWebPreview)
	return MessageCtx {
		Msg: msg,
		Bot: this.Bot,
	}, err
}

func (this *MessageCtx) Delete() (error) {
	return m.Bot.Remote().DeleteMessage(this.Msg.Chat.Id, this.Msg.Message_id)
}

func (this *MessageCtx) RespondAsync(m OMessage, handler gogram.ResponseHandler) (MessageCtx, error) {
	m.Bot.Remote().SendMessageAsync(this.Msg.Chat.Id, m.Text, m.ReplyTo, m.ParseMode, m.ReplyMarkup, !m.EnableWebPreview, handler)
}

func (this *MessageCtx) ReplyAsync(m OMessage, handler gogram.ResponseHandler) (MessageCtx, error) {
	m.Bot.Remote().SendMessageAsync(this.Msg.Chat.Id, m.Text, &this.Msg.Message_id, m.ParseMode, m.ReplyMarkup, !m.EnableWebPreview, handler)
}

func (this *MessageCtx) DeleteAsync(handler gogram.ResponseHandler) (error) {
	return m.Bot.Remote().DeleteMessageAsync(this.Msg.Chat.Id, this.Msg.Message_id, handler)
}
