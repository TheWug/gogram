package persist

import (
	"github.com/thewug/gogram"
	"github.com/thewug/gogram/data"

	"encoding/json"
	"time"
)

type StateBasePersistent struct {
	gogram.StateBase

	Persist interface{}
	Timestamp time.Time
	Tag StateType

	Ctx *PersistContext
}

func (this *StateBasePersistent) StateExited(sender data.Sender) {
	this.Ctx.WriteStateExited(sender)
}

func (this *StateBasePersistent) StateEntered(sender data.Sender) {
	jstr, _ := json.Marshal(this.Persist)
	this.Ctx.WriteStateEntered(sender, this.Tag, time.Now(), jstr)
}

func Register(ctx *PersistContext, machine *gogram.MessageStateMachine, tag StateType, factory PersistentStateFactory) StateBasePersistent {
	ctx.stateFactories[tag] = factory

	var sbp StateBasePersistent
	sbp.StateBase = gogram.MakeBase(machine)

	sbp.Tag = tag
	sbp.Ctx = ctx
	return sbp
}
