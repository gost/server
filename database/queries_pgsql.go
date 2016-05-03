package database

import (
	"fmt"
	"io/ioutil"
	"log"
)

// GetCreateDatabaseQuery returns the database creation script for PostgreSQL
func GetCreateDatabaseQuery(schema string) string {
	content, err := ioutil.ReadFile("./scripts/createdb.sql")
	if err != nil {
		log.Fatal("db create read error: ", err)
	}

	return fmt.Sprintf(string(content), schema, schema)
}
