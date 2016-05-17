package postgis

import (
	"fmt"
	"github.com/geodan/gost/src/sensorthings/entities"
	"strconv"

	gostErrors "github.com/geodan/gost/src/errors"
)

// GetObservedProperty returns an ObservedProperty by id
func (gdb *GostDatabase) GetObservedProperty(id string) (*entities.ObservedProperty, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var opID int
	var name, definition, description string
	sql := fmt.Sprintf("select id, name, definition, description FROM %s.observedproperty where id = $1", gdb.Schema)
	err = gdb.Db.QueryRow(sql, intID).Scan(&opID, &name, &definition, &description)

	if err != nil {
		return nil, gostErrors.NewRequestNotFound(fmt.Errorf("ObservedProperties(%s) does not exist", id))
	}

	op := entities.ObservedProperty{
		ID:          strconv.Itoa(opID),
		Name:        name,
		Definition:  definition,
		Description: description,
	}

	return &op, nil
}

// GetObservedProperties returns all observed properties
func (gdb *GostDatabase) GetObservedProperties() ([]*entities.ObservedProperty, error) {
	sql := fmt.Sprintf("select id, name, definition, description FROM %s.observedproperty", gdb.Schema)
	rows, err := gdb.Db.Query(sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var observedProperties = []*entities.ObservedProperty{}

	for rows.Next() {
		var opID int
		var name string
		var definition string
		var description string
		err2 := rows.Scan(&opID, &name, &definition, &description)
		if err2 != nil {
			return nil, err2
		}

		op := entities.ObservedProperty{
			ID:          strconv.Itoa(opID),
			Name:        name,
			Definition:  definition,
			Description: description,
		}

		observedProperties = append(observedProperties, &op)
	}

	return observedProperties, nil
}

// PostObservedProperty adds an ObservedProperty to the database
func (gdb *GostDatabase) PostObservedProperty(op entities.ObservedProperty) (*entities.ObservedProperty, error) {
	var opID int
	sql := fmt.Sprintf("INSERT INTO %s.observedproperty (name, definition, description) VALUES ($1, $2, $3) RETURNING id", gdb.Schema)
	err := gdb.Db.QueryRow(sql, op.Name, op.Definition, op.Description).Scan(&opID)
	if err != nil {
		return nil, err
	}

	op.ID = strconv.Itoa(opID)
	return &op, nil
}

// ObservedPropertyExists checks if a ObservedProperty is present in the database based on a given id.
func (gdb *GostDatabase) ObservedPropertyExists(thingID int) bool {
	var result bool
	sql := fmt.Sprintf("SELECT exists (SELECT 1 FROM %s.observedproperty WHERE id = $1 LIMIT 1)", gdb.Schema)
	err := gdb.Db.QueryRow(sql, thingID).Scan(&result)
	if err != nil {
		return false
	}

	return result
}
