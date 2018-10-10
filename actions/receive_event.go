package actions

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Replay/broker"
	"github.com/ONSBR/Plataforma-Replay/db"
	"github.com/ONSBR/Plataforma-Replay/policy"
	"github.com/ONSBR/Plataforma-Replay/tape"
	"github.com/labstack/gommon/log"
)

var createTapePolicy policy.CreateTapePolicy

func ReceiveEvent(event *domain.Event) error {
	dumpAlreadyDone, err := createTapePolicy.DatabaseBackupAlreadyCreated(event.SystemID)
	if err != nil {
		log.Error(err)
		return err
	}
	currentTape := tape.GetTape(event.SystemID)
	if !dumpAlreadyDone {
		dbManager := db.GetDB()
		if err := dbManager.Backup(event.SystemID, currentTape); err != nil {
			log.Error(err)
		}
	}
	if err := currentTape.Record(event); err != nil {
		log.Error(err)
	}
	brk := broker.GetBroker()
	return brk.Publish(event)
}
