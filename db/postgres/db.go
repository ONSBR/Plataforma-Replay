package postgres

import "github.com/ONSBR/Plataforma-Replay/tape"

type DBPostgres struct {
}

func (db *DBPostgres) Backup(systemID string, tape *tape.Tape) error {
	return nil
}

func (db *DBPostgres) Restore(systemID string, tape *tape.Tape) error {
	return nil
}
