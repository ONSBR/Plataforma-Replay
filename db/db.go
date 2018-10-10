package db

import (
	"github.com/ONSBR/Plataforma-Replay/db/postgres"
	"github.com/ONSBR/Plataforma-Replay/tape"
)

//DB is a interface to backup and restore data from database
type DB interface {
	Backup(systemID string, tape *tape.Tape) error
	Restore(systemID string, tape *tape.Tape) error
}

func GetDB() DB {
	return new(postgres.DBPostgres)
}
