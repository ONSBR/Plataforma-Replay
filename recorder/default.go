package recorder

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Replay/db"
	"github.com/ONSBR/Plataforma-Replay/policy"
	"github.com/ONSBR/Plataforma-Replay/sdk"
	"github.com/ONSBR/Plataforma-Replay/tape"
	"github.com/labstack/gommon/log"
)

var createTapePolicy policy.CreateTapePolicy

type DefaultRecorder struct {
	systemID string
	path     string
	tape     *tape.Tape
}

func (d *DefaultRecorder) Rec(event *domain.Event) error {
	currentTape, err := tape.GetOrCreateTape(event.SystemID, d.path)
	if err != nil {
		return err
	}
	dbName, err := sdk.GetDBName(event.SystemID)
	if err != nil {
		return err
	}
	dumpAlreadyDone := createTapePolicy.DatabaseBackupAlreadyCreated(currentTape)
	if !dumpAlreadyDone {
		log.Info("creating database dump")
		dumpPath := fmt.Sprintf("%s/dump.sql", currentTape.Dest())
		if fd, err := os.Create(dumpPath); err != nil {
			log.Error(err)
			return err
		} else {
			fd.Close()
		}
		dbManager := db.GetDB()
		err := dbManager.Backup(dbName, dumpPath)
		if err != nil {
			log.Error(err)
			return err
		} else if err := currentTape.Record(&tape.Segment{FileName: "dump.sql", SegmentType: "dump", Timestamp: time.Now().Unix()}); err != nil {
			log.Error(err)
			return err
		}
	}
	if err := currentTape.RecordEvent(event); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

/*
Eject stops recording and eject the tape for a systemID, in this case if the user wants to record again
he should start the process from the begining
*/
func (d *DefaultRecorder) Eject() (*tape.Tape, error) {
	currentTape, err := tape.GetOrCreateTape(d.systemID, d.path)
	if err != nil {
		return nil, err
	}
	return currentTape, currentTape.Close()
}

func (d *DefaultRecorder) GetOrCreateTape(systemID string) (*tape.Tape, error) {
	return tape.GetOrCreateTape(systemID, d.path)
}

func (d *DefaultRecorder) GetTape() (*tape.Tape, error) {
	return tape.GetTape(d.systemID, d.path)
}

func (d *DefaultRecorder) InsertTape(tape *tape.Tape) {
	d.tape = tape
}

func (d *DefaultRecorder) IsRecording() bool {
	tp := tape.Tape{
		SystemID: d.systemID,
		Path:     d.path,
	}
	return tp.Exist()
}

func (d *DefaultRecorder) AvailableTapesToDownload(systemID string) ([]string, error) {
	fd, err := os.Open(d.path)
	if err != nil {
		return nil, err
	}
	files, errF := fd.Readdirnames(0)
	if errF != nil {
		return nil, errF
	}
	availables := make([]string, 0)
	for _, name := range files {
		if strings.HasPrefix(name, systemID) && strings.HasSuffix(name, ".rec") {
			availables = append(availables, name)
		}
	}
	return availables, nil
}

func (d *DefaultRecorder) TapePath(name string) string {
	return fmt.Sprintf("%s/%s", d.path, name)
}
