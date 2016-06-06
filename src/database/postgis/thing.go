package postgis

import (
	"encoding/json"
	"fmt"
	"github.com/geodan/gost/src/sensorthings/entities"
	"strconv"

	"database/sql"
	"errors"
	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// GetThing returns a thing entity based on id and query
func (gdb *GostDatabase) GetThing(id string, qo *odata.QueryOptions) (*entities.Thing, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("select id, description, properties from %s.thing where id = $1", gdb.Schema)
	datastream, err := processThing(gdb.Db, sql, intID)
	if err != nil {
		return nil, err
	}

	return datastream, nil
}

//GetThingByDatastream retrieves the thing linked to a datastream
func (gdb *GostDatabase) GetThingByDatastream(id string, qo *odata.QueryOptions) (*entities.Thing, error) {
	tID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("select thing.id, thing.description, thing.properties from %s.thing INNER JOIN %s.datastream ON datastream.thing_id = thing.id WHERE datastream.id = $1;", gdb.Schema, gdb.Schema)
	return processThing(gdb.Db, sql, tID)
}

//GetThingsByLocation retrieves the thing linked to a location
func (gdb *GostDatabase) GetThingsByLocation(id string, qo *odata.QueryOptions) ([]*entities.Thing, error) {
	tID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("select thing.id, thing.description, thing.properties from %s.thing INNER JOIN %s.thing_to_location ON thing.id = thing_to_location.thing_id	INNER JOIN %s.location ON thing_to_location.location_id = location.id WHERE location.id = $1;", gdb.Schema, gdb.Schema, gdb.Schema)
	return processThings(gdb.Db, sql, tID)
}

//GetThingByHistoricalLocation retrieves the thing linked to a HistoricalLocation
func (gdb *GostDatabase) GetThingByHistoricalLocation(id string, qo *odata.QueryOptions) (*entities.Thing, error) {
	tID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("select thing.id, thing.description, thing.properties from %s.thing INNER JOIN %s.historicallocation ON historicallocation.thing_id = thing.id WHERE historicallocation.id = $1;", gdb.Schema, gdb.Schema)
	return processThing(gdb.Db, sql, tID)
}

// GetThings returns an array of things
func (gdb *GostDatabase) GetThings(qo *odata.QueryOptions) ([]*entities.Thing, error) {
	sql := fmt.Sprintf("select id, description, properties FROM %s.thing", gdb.Schema)
	return processThings(gdb.Db, sql)
}

func processThing(db *sql.DB, sql string, args ...interface{}) (*entities.Thing, error) {
	observations, err := processThings(db, sql, args...)
	if err != nil {
		return nil, err
	}

	if len(observations) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Thing not found"))
	}

	return observations[0], nil
}

func processThings(db *sql.DB, sql string, args ...interface{}) ([]*entities.Thing, error) {
	rows, err := db.Query(sql, args...)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var things = []*entities.Thing{}
	for rows.Next() {
		var thingID int
		var description string
		var properties *string

		err = rows.Scan(&thingID, &description, &properties)
		if err != nil {
			return nil, err
		}

		propMap, err := JSONToMap(properties)
		if err != nil {
			return nil, err
		}

		thing := entities.Thing{}
		thing.ID = strconv.Itoa(thingID)
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

	thing.ID = strconv.Itoa(thingID)
	return thing, nil
}

// ThingExists checks if a thing is present in the database based on a given id
func (gdb *GostDatabase) ThingExists(thingID int) bool {
	var result bool
	sql := fmt.Sprintf("SELECT exists (SELECT 1 FROM %s.thing WHERE id = $1 LIMIT 1)", gdb.Schema)
	err := gdb.Db.QueryRow(sql, thingID).Scan(&result)
	if err != nil {
		return false
	}

	return result
}

// DeleteThing tries to delete a Thing by the given id
func (gdb *GostDatabase) DeleteThing(id string) error {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	r, err := gdb.Db.Exec(fmt.Sprintf("DELETE FROM %s.thing WHERE id = $1", gdb.Schema), intID)
	if err != nil {
		return err
	}

	if c, _ := r.RowsAffected(); c == 0 {
		return gostErrors.NewRequestNotFound(errors.New("Thing not found"))
	}

	return nil
}
