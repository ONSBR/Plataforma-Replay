package db

import (
	"github.com/ONSBR/Plataforma-Replay/db/postgres"
)

//DB is a interface to backup and restore data from database
type DB interface {
	Backup(systemID, path string) error
	Restore(systemID, path string) error
}

func GetDB() DB {
	return new(postgres.DBPostgres)
}
