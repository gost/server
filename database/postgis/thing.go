package postgis

import (
	"encoding/json"
	"fmt"
	"github.com/geodan/gost/sensorthings/entities"
	"strconv"
)

// GetThing returns a thing entity based on id and query
func (gdb *GostDatabase) GetThing(id string) (*entities.Thing, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var thingID int
	var description string
	var properties string
	sql := fmt.Sprintf("select id, description, properties from %s.thing where id = $1", gdb.Schema)
	err2 := gdb.Db.QueryRow(sql, intID).Scan(&thingID, &description, &properties)

	if err2 != nil {
		return nil, err
	}

	thing := entities.Thing{}
	thing.ID = strconv.Itoa(thingID)
	thing.Description = description

	var p map[string]string
	err3 := json.Unmarshal([]byte(properties), &p)
	if err3 != nil {
		return nil, err3
	}

	thing.Properties = p

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
		thing := entities.Thing{}

		var id int
		var description string
		var properties string
		err2 := rows.Scan(&id, &description, &properties)
		if err2 != nil {
			return nil, err2
		}

		thing.ID = strconv.Itoa(id)
		thing.Description = description

		var p map[string]string
		err3 := json.Unmarshal([]byte(properties), &p)
		if err3 != nil {
			return nil, err3
		}

		thing.Properties = p
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
	sql := fmt.Sprintf("SELECT exists (SELECT 1 FROM %s.thing WHERE id = $1 LIMIT 1)", gdb.Schema)
	err := gdb.Db.QueryRow(sql, thingID).Scan(&result)
	if err != nil {
		return false
	}

	return result
}
