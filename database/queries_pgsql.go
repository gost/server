package database

import (
	"fmt"
	"io/ioutil"
)

// GetCreateDatabaseQuery returns the database creation script for PostgreSQL
func GetCreateDatabaseQuery(schema string) (string, error) {
	content, err := ioutil.ReadFile("./scripts/createdb.sql")
	if err != nil {
    return "", err
	}

	return fmt.Sprintf(string(content), schema, schema), nil
}
