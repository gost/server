package postgis

import (
	"fmt"

	"github.com/geodan/gost/src/sensorthings/entities"

	"database/sql"
	"errors"

	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/odata"
)

func observedPropertyParamFactory(values map[string]interface{}) (entities.Entity, error) {
	op := &entities.ObservedProperty{}
	for as, value := range values {
		if value == nil {
			continue
		}

		if as == asMappings[entities.EntityTypeObservedProperty][observedPropertyID] {
			op.ID = value
		} else if as == asMappings[entities.EntityTypeObservedProperty][observedPropertyName] {
			op.Name = value.(string)
		} else if as == asMappings[entities.EntityTypeObservedProperty][observedPropertyDescription] {
			op.Description = value.(string)
		} else if as == asMappings[entities.EntityTypeObservedProperty][observedPropertyDefinition] {
			op.Definition = value.(string)
		}
	}

	return op, nil
}

// GetObservedProperty returns an ObservedProperty by id
func (gdb *GostDatabase) GetObservedProperty(id interface{}, qo *odata.QueryOptions) (*entities.ObservedProperty, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("ObservedProperty does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.ObservedProperty{}, nil, intID, qo)
	observedProperty, err := processObservedProperty(gdb.Db, query, qi)
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

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.ObservedProperty{}, &entities.Datastream{}, intID, qo)
	observedProperty, err := processObservedProperty(gdb.Db, query, qi)
	if err != nil {
		return nil, err
	}

	return observedProperty, nil
}

// GetObservedProperties returns all observed properties
func (gdb *GostDatabase) GetObservedProperties(qo *odata.QueryOptions) ([]*entities.ObservedProperty, int, error) {
	query, qi := gdb.QueryBuilder.CreateQuery(&entities.ObservedProperty{}, nil, nil, qo)
	countSQL := gdb.QueryBuilder.CreateCountQuery(&entities.ObservedProperty{}, nil, nil, qo)
	return processObservedProperties(gdb.Db, query, qi, countSQL)
}

func processObservedProperty(db *sql.DB, sql string, qi *QueryParseInfo) (*entities.ObservedProperty, error) {
	ops, _, err := processObservedProperties(db, sql, qi, "")
	if err != nil {
		return nil, err
	}

	if len(ops) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("ObservedProperty not found"))
	}

	return ops[0], nil
}

func processObservedProperties(db *sql.DB, sql string, qi *QueryParseInfo, countSQL string) ([]*entities.ObservedProperty, int, error) {
	data, err := ExecuteSelect(db, qi, sql)
	if err != nil {
		return nil, 0, fmt.Errorf("Error executing query %v", err)
	}

	obs := make([]*entities.ObservedProperty, 0)
	for _, d := range data {
		entity := d.(*entities.ObservedProperty)
		obs = append(obs, entity)
	}

	var count int
	if len(countSQL) > 0 {
		count, err = ExecuteSelectCount(db, countSQL)
		if err != nil {
			return nil, 0, fmt.Errorf("Error executing count %v", err)
		}
	}

	return obs, count, nil
}

// PostObservedProperty adds an ObservedProperty to the database
func (gdb *GostDatabase) PostObservedProperty(op *entities.ObservedProperty) (*entities.ObservedProperty, error) {
	var opID int
	query := fmt.Sprintf("INSERT INTO %s.observedproperty (name, definition, description) VALUES ($1, $2, $3) RETURNING id", gdb.Schema)
	err := gdb.Db.QueryRow(query, op.Name, op.Definition, op.Description).Scan(&opID)
	if err != nil {
		return nil, err
	}

	op.ID = opID
	return op, nil
}

// PutObservedProperty updates a ObservedProperty in the database
func (gdb *GostDatabase) PutObservedProperty(id interface{}, op *entities.ObservedProperty) (*entities.ObservedProperty, error) {
	return gdb.PatchObservedProperty(id, op)
}

// ObservedPropertyExists checks if a ObservedProperty is present in the database based on a given id.
func (gdb *GostDatabase) ObservedPropertyExists(id interface{}) bool {
	return EntityExists(gdb, id, "observedproperty")
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
