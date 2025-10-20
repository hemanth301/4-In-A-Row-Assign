package services

import (
	"database/sql"
	"io/ioutil"
)

// RunMigrations runs the schema.sql file against the given DB
func RunMigrations(db *sql.DB, schemaPath string) error {
	b, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(b))
	return err
}
