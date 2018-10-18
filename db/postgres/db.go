package postgres

import (
	"bufio"
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"

	"github.com/ONSBR/Plataforma-Deployer/env"

	"github.com/PMoneda/whaler"
)

type DBPostgres struct {
}

func (db *DBPostgres) Backup(systemID string) (*bufio.Reader, error) {
	if dbName, err := db.getDbName(systemID); err != nil {
		return nil, err
	} else {
		if output, err := whaler.RunCommand("postgres", "pg_dump", "-U", env.Get("POSTGRES_USER", "postgres"), dbName); err != nil {
			return nil, err
		} else {
			buf := make([]byte, 8)
			_, err = output.Read(buf) //remove binary encoded bytes
			if err != nil {
				return nil, err
			}
			return output, nil
		}
	}
}

func (db *DBPostgres) getDbName(systemID string) (string, error) {
	response := make([]map[string]interface{}, 0)
	apicore.Query(apicore.Filter{
		Entity: "installedApp",
		Map:    "core",
		Name:   "bySystemIdAndType",
		Params: []apicore.Param{
			apicore.Param{
				Key:   "systemId",
				Value: systemID,
			},
			apicore.Param{
				Key:   "type",
				Value: "domain",
			},
		},
	}, &response)
	if len(response) == 0 {
		return "", fmt.Errorf("Cannot find database name for system %s", systemID)
	}
	switch t := response[0]["name"].(type) {
	case string:
		return t, nil
	default:
		return "", fmt.Errorf("Invalid datatype for database name")
	}
}

func (db *DBPostgres) Restore(systemID, path string) error {
	return nil
}
