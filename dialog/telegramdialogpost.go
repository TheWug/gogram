package dialog

import (
	"storage"

	"github.com/thewug/gogram"
	"github.com/thewug/gogram/data"

	"errors"
	"time"
)

var ErrDialogTypeMismatch error = errors.New("Mismatched Dialog ID")

type Dialog interface {
	JSON() (string, error)
	ID() data.DialogID
}

type TelegramDialogPost struct {
	msg_id    data.MsgID
	msg_ts    time.Time
	chat_id   data.ChatID
	dialog_id data.DialogID

	dialog    Dialog
}

func (this *TelegramDialogPost) FirstSave(settings storage.UpdaterSettings, msg_id data.MsgID, chat_id data.ChatID, ts time.Time, d Dialog) (error) {
	*this = TelegramDialogPost{
		msg_id: msg_id,
		msg_ts: ts,
		chat_id: chat_id,
		dialog_id: d.ID(),
		dialog: d,
	}
	return this.Save(settings)
}

func (this *TelegramDialogPost) Save(settings storage.UpdaterSettings) error {
	json, err := this.dialog.JSON()
	if err != nil { return err }

	return storage.WriteDialogPost(settings, this.dialog_id, this.msg_id, this.chat_id, json, this.msg_ts)
}

func (this *TelegramDialogPost) Delete(settings storage.UpdaterSettings) error {
	return storage.EraseDialogPost(settings, this.msg_id, this.chat_id)
}

func (this *TelegramDialogPost) Load(found *storage.DialogPost, id data.DialogID, dlg Dialog) {
	this.msg_id = found.MsgId
	this.msg_ts = found.MsgTs
	this.chat_id = found.ChatId
	this.dialog_id = id
	this.dialog = dlg
}

func (this *TelegramDialogPost) Ctx(bot *gogram.TelegramBot) *gogram.MessageCtx {
	if this.msg_id == 0 || this.chat_id == 0 {
		return nil
	}

	return gogram.NewMessageCtx(&data.TMessage{
		Id: this.msg_id,
		Chat: data.TChat{
			Id: this.chat_id,
		},
	}, false, bot)
}

func (this *TelegramDialogPost) IsUnset() bool {
	return (this.dialog_id == "")
}

func NewTelegramDialogPost(msg_id data.MsgID, chat_id data.ChatID, dialog_id data.DialogID, msg_ts time.Time, dialog Dialog) TelegramDialogPost {
	return TelegramDialogPost{
		msg_id: msg_id,
		chat_id: chat_id,
		dialog_id: dialog_id,
		msg_ts: msg_ts,
		dialog: dialog,
	}
}
