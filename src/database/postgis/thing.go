package postgis

import (
	"encoding/json"
	"fmt"
	"github.com/geodan/gost/src/sensorthings/entities"
	"strconv"

	gostErrors "github.com/geodan/gost/src/errors"
)

// GetThing returns a thing entity based on id and query
func (gdb *GostDatabase) GetThing(id string) (*entities.Thing, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var thingID int
	var description string
	var properties *string

	sql := fmt.Sprintf("select id, description, properties from %s.thing where id = $1", gdb.Schema)
	err = gdb.Db.QueryRow(sql, intID).Scan(&thingID, &description, &properties)
	if err != nil {
		return nil, gostErrors.NewRequestNotFound(fmt.Errorf("Things(%s) does not exist", id))
	}

	propMap, err := JSONToMap(properties)
	if err != nil {
		return nil, err
	}

	thing := entities.Thing{
		ID:          strconv.Itoa(thingID),
		Description: description,
		Properties:  propMap,
	}

	return &thing, nil
}

//GetThingByDatastream retrieves the thing linked to a datastream
func (gdb *GostDatabase) GetThingByDatastream(id string) (*entities.Thing, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var thingID int
	var description string
	var properties *string

	sql := fmt.Sprintf("select thing.id, thing.description, thing.properties from %s.thing INNER JOIN %s.datastream ON datastream.thing_id = thing.id WHERE datastream.id = $1;", gdb.Schema, gdb.Schema)
	err = gdb.Db.QueryRow(sql, intID).Scan(&thingID, &description, &properties)
	if err != nil {
		return nil, gostErrors.NewRequestNotFound(fmt.Errorf("Datastreams(%s)/Thing does not exist", id))
	}

	propMap, err := JSONToMap(properties)
	if err != nil {
		return nil, err
	}

	thing := entities.Thing{
		ID:          strconv.Itoa(thingID),
		Description: description,
		Properties:  propMap,
	}

	return &thing, nil
}

// GetThings returns an array of things
func (gdb *GostDatabase) GetThings() ([]*entities.Thing, error) {
	sql := fmt.Sprintf("select id, description, properties FROM %s.thing", gdb.Schema)
	rows, err := gdb.Db.Query(sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
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

		thing := entities.Thing{
			ID:          strconv.Itoa(thingID),
			Description: description,
			Properties:  propMap,
		}

		things = append(things, &thing)
	}

	return things, nil
}

// PostThing receives a posted thing entity and adds it to the database
// returns the created Thing including the generated id
func (gdb *GostDatabase) PostThing(thing entities.Thing) (*entities.Thing, error) {
	jsonProperties, _ := json.Marshal(thing.Properties)
	var thingID int
	sql := fmt.Sprintf("INSERT INTO %s.thing (description, properties) VALUES ($1, $2) RETURNING id", gdb.Schema)
	err := gdb.Db.QueryRow(sql, thing.Description, jsonProperties).Scan(&thingID)
	if err != nil {
		return nil, err
	}

	thing.ID = strconv.Itoa(thingID)
	return &thing, nil
}

// ThingExists checks if a thing is present in the database based on a given id
func (gdb *GostDatabase) ThingExists(thingID int) bool {
	var result bool
	sql := fmt.Sprintf("SELECT %s.exists (SELECT 1 FROM thing WHERE id = $1 LIMIT 1)", gdb.Schema)
	err := gdb.Db.QueryRow(sql, thingID).Scan(&result)
	if err != nil {
		return false
	}

	return result
}
