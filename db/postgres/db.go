package postgres

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	_ "github.com/lib/pq"

	"github.com/labstack/gommon/log"

	"github.com/ONSBR/Plataforma-Deployer/env"

	"github.com/PMoneda/whaler"
)

type DBPostgres struct {
}

func (db *DBPostgres) Backup(dbName, path string) error {
	log.Info("dumping database ", dbName, " to ", path)
	if r, err := whaler.RunCommand("postgres", "pg_dump", "-U", env.Get("POSTGRES_USER", "postgres"), "-d", dbName, "-f", path); err != nil {
		log.Error(err)
		return err
	} else {
		buf, _ := ioutil.ReadAll(r)
		log.Info("postgres output")
		log.Info(string(buf))
	}
	return nil
}

func (db *DBPostgres) Restore(dbName, path string) error {
	if err := db.createDatabaseIfNotExist(dbName); err != nil {
		log.Info(err.Error())
		return err
	}
	if buf, err := whaler.RunCommand("postgres", "psql", "-U", env.Get("POSTGRES_USER", "postgres"), "-d", dbName, "-f", path); err != nil {
		log.Info(err.Error())
		return err
	} else {
		if _, err := ioutil.ReadAll(buf); err != nil {
			return err
		}
	}
	return nil
}

func (db *DBPostgres) createDatabaseIfNotExist(name string) error {
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		env.Get("POSTGRES_HOST", "localhost"), env.Get("POSTGRES_PORT", "5432"), env.Get("POSTGRES_USER", "postgres"), env.Get("POSTGRES_PASSWORD", "postgres"), "postgres")

	dataConn, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Error(err)
		return err
	}
	defer dataConn.Close()
	dataConn.SetMaxOpenConns(1)
	dataConn.SetMaxIdleConns(1)
	log.Info("droping database")
	if _, err := dataConn.Exec(fmt.Sprintf("DROP DATABASE %s", name)); err != nil {
		log.Error(err)
	}
	log.Info("creating database ", name)
	if _, err := dataConn.Exec(fmt.Sprintf("CREATE DATABASE %s", name)); err != nil {
		return err
	}
	return nil
}
