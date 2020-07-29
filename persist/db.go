package persist

import (
	"github.com/thewug/gogram"
	"github.com/thewug/gogram/data"

	"database/sql"
	"log"
	"time"
)

type StateType string

type PersistentStateFactory func (jstr []byte, sbp StateBasePersistent) (gogram.State)

type PersistContext struct {
	db_pool *sql.DB
	table_name string
	stateFactories map[StateType]PersistentStateFactory
}

func InitStatePersistence(database *sql.DB, table string) *PersistContext {
	var p PersistContext
	p.db_pool = database
	p.table_name = table
	p.stateFactories = make(map[StateType]PersistentStateFactory)
	return &p
}

func (this *PersistContext) WriteStateExited(sender data.Sender) {
	query := "DELETE FROM " + this.table_name + " WHERE state_user = $1 AND state_channel = $2"
	_, err := this.db_pool.Exec(query, sender.User, sender.Channel)
	if err != nil { log.Println(err.Error()) }
}

func (this *PersistContext) WriteStateEntered(sender data.Sender, tag StateType, timestamp time.Time, jstr []byte) {
	query := "INSERT INTO " + this.table_name + " (state_user, state_channel, state_ts, state_persist, state_type) VALUES ($1, $2, $3, $4, $5)"
	_, err := this.db_pool.Exec(query, sender.User, sender.Channel, timestamp, jstr, tag)
	if err != nil { log.Println(err.Error()) }
}

func (this *PersistContext) LoadAllStates(machine *gogram.MessageStateMachine) error {
	query := "SELECT state_user, state_channel, state_ts, state_persist, state_type FROM " + this.table_name
	rows, err := this.db_pool.Query(query)

	if err != nil { return err }

	for rows.Next() {
		var sender data.Sender
		var sbp StateBasePersistent
		sbp.Ctx = this
		var jstr []byte

		err := rows.Scan(&sender.User, &sender.Channel, &sbp.Timestamp, &jstr, &sbp.Tag)
		if err != nil { return err }

		factory := this.stateFactories[sbp.Tag]
		if factory != nil {
			sbp.StateMachine = machine
			state := factory(jstr, sbp)
			machine.SetStateDirect(sender, state)
		}
	}

	return nil
}
