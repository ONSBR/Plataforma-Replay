package db

//DB is a interface to backup and restore data from database
type DB interface {
	Backup()
	Restore()
}
