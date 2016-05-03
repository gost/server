package postgis

import (
	"fmt"
	"github.com/geodan/gost/sensorthings/entities"
	"strconv"
)

// GetObservedProperty todo
func (gdb *GostDatabase) GetObservedProperty(id string) (*entities.ObservedProperty, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var opID int
	var name string
	var definition string
	var description string
	sql := fmt.Sprintf("select id, name, definition, description FROM %s.observedproperty where id = $1", gdb.Schema)
	err2 := gdb.Db.QueryRow(sql, intID).Scan(&opID, &name, &definition, &description)

	if err2 != nil {
		return nil, err
	}

	op := entities.ObservedProperty{}
	op.ID = strconv.Itoa(opID)
	op.Name = name
	op.Definition = definition
	op.Description = description

	return &op, nil
}

// GetObservedProperties todo
func (gdb *GostDatabase) GetObservedProperties() ([]*entities.ObservedProperty, error) {
	sql := fmt.Sprintf("select id, name, definition, description FROM %s.observedproperty", gdb.Schema)
	rows, err := gdb.Db.Query(sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var observedproperties = []*entities.ObservedProperty{}

	for rows.Next() {
		op := entities.ObservedProperty{}

		var id int
		var name string
		var definition string
		var description string
		err2 := rows.Scan(&id, &name, &definition, &description)
		if err2 != nil {
			return nil, err2
		}

		op.ID = strconv.Itoa(id)
		op.Name = name
		op.Definition = definition
		op.Description = description

		observedproperties = append(observedproperties, &op)
	}

	return observedproperties, nil
}

// PostObservedProperty todo
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
