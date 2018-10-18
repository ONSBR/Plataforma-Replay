package recorder

import (
	"bufio"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Replay/db"
	"github.com/ONSBR/Plataforma-Replay/policy"
	"github.com/ONSBR/Plataforma-Replay/tape"
	"github.com/labstack/gommon/log"
)

var createTapePolicy policy.CreateTapePolicy

const rootPath = "./"

type DefaultRecorder struct {
}

func (d *DefaultRecorder) Rec(event *domain.Event) error {
	currentTape, err := tape.GetOrCreateTape(event.SystemID, rootPath)
	if err != nil {
		return err
	}
	dumpAlreadyDone := createTapePolicy.DatabaseBackupAlreadyCreated(currentTape)
	if !dumpAlreadyDone {
		dbManager := db.GetDB()
		var reader *bufio.Reader
		reader, err := dbManager.Backup(event.SystemID)
		if err != nil {
			log.Error(err)
		} else {
			currentTape.RecordReader("dump.sql", "dump", reader)
		}
	}
	if err := currentTape.RecordEvent(event); err != nil {
		log.Error(err)
	}
	return nil
}

func (d *DefaultRecorder) Eject(systemID string) (*tape.Tape, error) {
	currentTape, err := tape.GetOrCreateTape(systemID, rootPath)
	if err != nil {
		return nil, err
	}
	return currentTape, currentTape.Close()
}

func (d *DefaultRecorder) Insert(tape *tape.Tape) error {
	return nil
}

func (d *DefaultRecorder) IsClosed() (bool, error) {
	return false, nil
}
