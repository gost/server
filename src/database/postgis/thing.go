package postgis

import (
	"encoding/json"
	"fmt"
	"github.com/geodan/gost/src/sensorthings/entities"

	"database/sql"
	"errors"
	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/odata"
)

var totalThings int

// GetTotalThings returns the total things count in the database
func (gdb *GostDatabase) GetTotalThings() int {
	return totalThings
}

// InitThings Initialises the datastream repository, setting totalThings on startup
func (gdb *GostDatabase) InitThings() {
	sql := fmt.Sprintf("SELECT Count(*) from %s.thing", gdb.Schema)
	gdb.Db.QueryRow(sql).Scan(&totalThings)
}

// GetThing returns a thing entity based on id and query
func (gdb *GostDatabase) GetThing(id interface{}, qo *odata.QueryOptions) (*entities.Thing, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Thing does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Thing{}, qo, "", "", nil)+" from %s.thing where id = %v", gdb.Schema, intID)
	datastream, err := processThing(gdb.Db, sql, qo)
	if err != nil {
		return nil, err
	}

	return datastream, nil
}

//GetThingByDatastream retrieves the thing linked to a datastream
func (gdb *GostDatabase) GetThingByDatastream(id interface{}, qo *odata.QueryOptions) (*entities.Thing, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Thing{}, qo, "thing.", "", nil)+" from %s.thing INNER JOIN %s.datastream ON datastream.thing_id = thing.id WHERE datastream.id = %v;", gdb.Schema, gdb.Schema, intID)
	return processThing(gdb.Db, sql, qo)
}

//GetThingsByLocation retrieves the thing linked to a location
func (gdb *GostDatabase) GetThingsByLocation(id interface{}, qo *odata.QueryOptions) ([]*entities.Thing, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Location does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Thing{}, qo, "thing.", "", nil)+" from %s.thing INNER JOIN %s.thing_to_location ON thing.id = thing_to_location.thing_id INNER JOIN %s.location ON thing_to_location.location_id = location.id WHERE location.id = %v  "+CreateTopSkipQueryString(qo), gdb.Schema, gdb.Schema, gdb.Schema, intID)
	return processThings(gdb.Db, sql, qo)
}

//GetThingByHistoricalLocation retrieves the thing linked to a HistoricalLocation
func (gdb *GostDatabase) GetThingByHistoricalLocation(id interface{}, qo *odata.QueryOptions) (*entities.Thing, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("HistoricalLocation does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Thing{}, qo, "thing.", "", nil)+" from %s.thing INNER JOIN %s.historicallocation ON historicallocation.thing_id = thing.id WHERE historicallocation.id = %v;", gdb.Schema, gdb.Schema, intID)
	return processThing(gdb.Db, sql, qo)
}

// GetThings returns an array of things
func (gdb *GostDatabase) GetThings(qo *odata.QueryOptions) ([]*entities.Thing, error) {
	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Thing{}, qo, "", "", nil)+" FROM %s.thing "+CreateTopSkipQueryString(qo), gdb.Schema)
	return processThings(gdb.Db, sql, qo)
}

func processThing(db *sql.DB, sql string, qo *odata.QueryOptions) (*entities.Thing, error) {
	observations, err := processThings(db, sql, qo)
	if err != nil {
		return nil, err
	}

	if len(observations) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Thing not found"))
	}

	return observations[0], nil
}

func processThings(db *sql.DB, sql string, qo *odata.QueryOptions) ([]*entities.Thing, error) {
	rows, err := db.Query(sql)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var things = []*entities.Thing{}
	for rows.Next() {
		var thingID interface{}
		var description string
		var properties *string

		var params []interface{}
		var qp []string
		if qo == nil || qo.QuerySelect == nil || len(qo.QuerySelect.Params) == 0 {
			t := &entities.Thing{}
			qp = t.GetPropertyNames()
		} else {
			qp = qo.QuerySelect.Params
		}

		for _, p := range qp {
			if p == "id" {
				params = append(params, &thingID)
			}
			if p == "description" {
				params = append(params, &description)
			}
			if p == "properties" {
				params = append(params, &properties)
			}
		}

		err = rows.Scan(params...)
		if err != nil {
			return nil, err
		}

		propMap, err := JSONToMap(properties)
		if err != nil {
			return nil, err
		}

		thing := entities.Thing{}
		thing.ID = thingID
		thing.Description = description
		thing.Properties = propMap

		things = append(things, &thing)
	}

	return things, nil
}

// PostThing receives a posted thing entity and adds it to the database
// returns the created Thing including the generated id
func (gdb *GostDatabase) PostThing(thing *entities.Thing) (*entities.Thing, error) {
	jsonProperties, _ := json.Marshal(thing.Properties)
	var thingID int
	sql := fmt.Sprintf("INSERT INTO %s.thing (description, properties) VALUES ($1, $2) RETURNING id", gdb.Schema)
	err := gdb.Db.QueryRow(sql, thing.Description, jsonProperties).Scan(&thingID)
	if err != nil {
		return nil, err
	}

	thing.ID = thingID
	totalThings++
	return thing, nil
}

// ThingExists checks if a thing is present in the database based on a given id
func (gdb *GostDatabase) ThingExists(thingID interface{}) bool {
	var result bool
	sql := fmt.Sprintf("SELECT exists (SELECT 1 FROM %s.thing WHERE id = $1 LIMIT 1)", gdb.Schema)
	err := gdb.Db.QueryRow(sql, thingID).Scan(&result)
	if err != nil {
		return false
	}

	return result
}

// DeleteThing tries to delete a Thing by the given id
func (gdb *GostDatabase) DeleteThing(id interface{}) error {
	intID, ok := ToIntID(id)
	if !ok {
		return gostErrors.NewRequestNotFound(errors.New("Thing does not exist"))
	}

	r, err := gdb.Db.Exec(fmt.Sprintf("DELETE FROM %s.thing WHERE id = $1", gdb.Schema), intID)
	if err != nil {
		return err
	}

	if c, _ := r.RowsAffected(); c == 0 {
		return gostErrors.NewRequestNotFound(errors.New("Thing not found"))
	}

	totalThings--
	return nil
}
