package database

import (
	"fmt"
	"io/ioutil"
)

// GetCreateDatabaseQuery returns the database creation script for PostgreSQL
func GetCreateDatabaseQuery(location string, schema string) (*string, error) {
	bytes, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}

	content := string(bytes[:])
	formatted := fmt.Sprintf(content, schema, schema)
	return &formatted, nil
}
