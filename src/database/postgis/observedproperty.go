package postgis

import (
	"fmt"

	"github.com/geodan/gost/src/sensorthings/entities"

	"database/sql"
	"errors"

	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// GetTotalObservedProperties returns the total ObservedProperties count in the database
func (gdb *GostDatabase) GetTotalObservedProperties() int {
	var count int
	sql := fmt.Sprintf("SELECT Count(*) from %s.observedproperty", gdb.Schema)
	gdb.Db.QueryRow(sql).Scan(&count)
	return count
}

// GetObservedProperty returns an ObservedProperty by id
func (gdb *GostDatabase) GetObservedProperty(id interface{}, qo *odata.QueryOptions) (*entities.ObservedProperty, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("ObservedProperty does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.ObservedProperty{}, qo, "", "", nil)+" FROM %s.observedproperty where id = %v", gdb.Schema, intID)
	observedProperty, err := processObservedProperty(gdb.Db, sql, qo)
	if err != nil {
		return nil, err
	}

	return observedProperty, nil
}

// GetObservedPropertyByDatastream returns an ObservedProperty by id
func (gdb *GostDatabase) GetObservedPropertyByDatastream(id interface{}, qo *odata.QueryOptions) (*entities.ObservedProperty, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.ObservedProperty{}, qo, "observedproperty.", "", nil)+" FROM %s.observedproperty inner join %s.datastream on datastream.observedproperty_id = observedproperty.id where datastream.id = %v", gdb.Schema, gdb.Schema, intID)
	observedProperty, err := processObservedProperty(gdb.Db, sql, qo)
	if err != nil {
		return nil, err
	}

	return observedProperty, nil
}

// GetObservedProperties returns all observed properties
func (gdb *GostDatabase) GetObservedProperties(qo *odata.QueryOptions) ([]*entities.ObservedProperty, int, error) {
	sql := fmt.Sprintf("select "+CreateSelectString(&entities.ObservedProperty{}, qo, "", "", nil)+" FROM %s.observedproperty order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema)
	countSQL := fmt.Sprintf("select COUNT(*) FROM %s.observedproperty", gdb.Schema)
	return processObservedProperties(gdb.Db, sql, qo, countSQL)
}

func processObservedProperty(db *sql.DB, sql string, qo *odata.QueryOptions) (*entities.ObservedProperty, error) {
	observedProperties, _, err := processObservedProperties(db, sql, qo, "")
	if err != nil {
		return nil, err
	}

	if len(observedProperties) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("ObservedProperty not found"))
	}

	return observedProperties[0], nil
}

func processObservedProperties(db *sql.DB, sql string, qo *odata.QueryOptions, countSQL string) ([]*entities.ObservedProperty, int, error) {
	rows, err := db.Query(sql)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	var observedProperties = []*entities.ObservedProperty{}

	for rows.Next() {
		var opID interface{}
		var name string
		var definition string
		var description string

		var params []interface{}
		var qp []string
		if qo == nil || qo.QuerySelect == nil || len(qo.QuerySelect.Params) == 0 {
			op := &entities.ObservedProperty{}
			qp = op.GetPropertyNames()
		} else {
			qp = qo.QuerySelect.Params
		}

		for _, p := range qp {
			if p == "id" {
				params = append(params, &opID)
			}
			if p == "name" {
				params = append(params, &name)
			}
			if p == "definition" {
				params = append(params, &definition)
			}
			if p == "description" {
				params = append(params, &description)
			}
		}

		err = rows.Scan(params...)

		if err != nil {
			return nil, 0, err
		}
		op := entities.ObservedProperty{}
		op.ID = opID
		op.Name = name
		op.Definition = definition
		op.Description = description

		observedProperties = append(observedProperties, &op)
	}

	var count int
	if len(countSQL) > 0 {
		db.QueryRow(countSQL).Scan(&count)
	}

	return observedProperties, count, nil
}

// PostObservedProperty adds an ObservedProperty to the database
func (gdb *GostDatabase) PostObservedProperty(op *entities.ObservedProperty) (*entities.ObservedProperty, error) {
	var opID int
	sql := fmt.Sprintf("INSERT INTO %s.observedproperty (name, definition, description) VALUES ($1, $2, $3) RETURNING id", gdb.Schema)
	err := gdb.Db.QueryRow(sql, op.Name, op.Definition, op.Description).Scan(&opID)
	if err != nil {
		return nil, err
	}

	op.ID = opID
	return op, nil
}

// ObservedPropertyExists checks if a ObservedProperty is present in the database based on a given id.
func (gdb *GostDatabase) ObservedPropertyExists(thingID interface{}) bool {
	var result bool
	sql := fmt.Sprintf("SELECT exists (SELECT 1 FROM %s.observedproperty WHERE id = $1 LIMIT 1)", gdb.Schema)
	err := gdb.Db.QueryRow(sql, thingID).Scan(&result)
	if err != nil {
		return false
	}

	return result
}

// PatchObservedProperty updates a ObservedProperty in the database
func (gdb *GostDatabase) PatchObservedProperty(id interface{}, op *entities.ObservedProperty) (*entities.ObservedProperty, error) {
	var err error
	var ok bool
	var intID int
	updates := make(map[string]interface{})

	if intID, ok = ToIntID(id); !ok || !gdb.ObservedPropertyExists(intID) {
		return nil, gostErrors.NewRequestNotFound(errors.New("ObservedProperty does not exist"))
	}

	if len(op.Description) > 0 {
		updates["description"] = op.Description
	}

	if len(op.Definition) > 0 {
		updates["definition"] = op.Definition
	}

	if len(op.Name) > 0 {
		updates["name"] = op.Name
	}

	if err = gdb.updateEntityColumns("observedproperty", updates, intID); err != nil {
		return nil, err
	}

	ns, _ := gdb.GetObservedProperty(intID, nil)
	return ns, nil
}

// DeleteObservedProperty tries to delete a ObservedProperty by the given id
func (gdb *GostDatabase) DeleteObservedProperty(id interface{}) error {
	return DeleteEntity(gdb, id, "observedproperty")
}
