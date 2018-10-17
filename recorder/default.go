package recorder

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Replay/db"
	"github.com/ONSBR/Plataforma-Replay/policy"
	"github.com/ONSBR/Plataforma-Replay/tape"
	"github.com/labstack/gommon/log"
)

var createTapePolicy policy.CreateTapePolicy

type DefaultRecorder struct {
	currentTape *tape.Tape
}

func (d *DefaultRecorder) Rec(event *domain.Event) error {
	dumpAlreadyDone, err := createTapePolicy.DatabaseBackupAlreadyCreated(event.SystemID)
	if err != nil {
		log.Error(err)
		return err
	}
	if !dumpAlreadyDone {
		dbManager := db.GetDB()
		if err := dbManager.Backup(event.SystemID, d.currentTape.Dest()); err != nil {
			log.Error(err)
		}
	}
	if err := d.currentTape.Record(event); err != nil {
		log.Error(err)
	}
	return nil
}

func (d *DefaultRecorder) Eject(systemID string) (*tape.Tape, error) {
	return nil, nil
}

func (d *DefaultRecorder) Insert(tape *tape.Tape) error {
	return nil
}

func (d *DefaultRecorder) IsClosed() (bool, error) {
	return false, nil
}
