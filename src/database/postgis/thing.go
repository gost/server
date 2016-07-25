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
func (gdb *GostDatabase) GetThingsByLocation(id interface{}, qo *odata.QueryOptions) ([]*entities.Thing, int, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, 0, gostErrors.NewRequestNotFound(errors.New("Location does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Thing{}, qo, "thing.", "", nil)+" from %s.thing INNER JOIN %s.thing_to_location ON thing.id = thing_to_location.thing_id INNER JOIN %s.location ON thing_to_location.location_id = location.id WHERE location.id = %v order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema, gdb.Schema, gdb.Schema, intID)
	countSQL := fmt.Sprintf("SELECT Count(*) from %s.thing INNER JOIN %s.historicallocation ON historicallocation.thing_id = thing.id WHERE historicallocation.id = %v;", gdb.Schema, gdb.Schema, intID)
	return processThings(gdb.Db, sql, qo, countSQL)
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
func (gdb *GostDatabase) GetThings(qo *odata.QueryOptions) ([]*entities.Thing, int, error) {
	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Thing{}, qo, "", "", nil)+" FROM %s.thing order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema)
	countSQL := fmt.Sprintf("select COUNT(*) as count FROM %s.thing", gdb.Schema)
	return processThings(gdb.Db, sql, qo, countSQL)
}

func processThing(db *sql.DB, sql string, qo *odata.QueryOptions) (*entities.Thing, error) {
	observations, _, err := processThings(db, sql, qo, "")
	if err != nil {
		return nil, err
	}

	if len(observations) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Thing not found"))
	}

	return observations[0], nil
}

func processThings(db *sql.DB, sql string, qo *odata.QueryOptions, countSQL string) ([]*entities.Thing, int, error) {
	rows, err := db.Query(sql)
	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}

	var things = []*entities.Thing{}
	for rows.Next() {
		var thingID interface{}
		var name, description string
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
			if p == "name" {
				params = append(params, &name)
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
			return nil, 0, err
		}

		propMap, err := JSONToMap(properties)
		if err != nil {
			return nil, 0, err
		}

		thing := entities.Thing{}
		thing.ID = thingID
		thing.Name = name
		thing.Description = description
		thing.Properties = propMap

		things = append(things, &thing)
	}

	var count int
	if len(countSQL) > 0 {
		db.QueryRow(countSQL).Scan(&count)
	}

	return things, count, nil
}

// PostThing receives a posted thing entity and adds it to the database
// returns the created Thing including the generated id
func (gdb *GostDatabase) PostThing(thing *entities.Thing) (*entities.Thing, error) {
	jsonProperties, _ := json.Marshal(thing.Properties)
	var thingID int
	sql := fmt.Sprintf("INSERT INTO %s.thing (name, description, properties) VALUES ($1, $2, $3) RETURNING id", gdb.Schema)
	err := gdb.Db.QueryRow(sql, thing.Name, thing.Description, jsonProperties).Scan(&thingID)
	if err != nil {
		return nil, err
	}

	thing.ID = thingID
	return thing, nil
}

// PatchThing receives a to be patched Thing entity and changes it in the database
// returns the patched Thing
func (gdb *GostDatabase) PatchThing(id interface{}, thing *entities.Thing) (*entities.Thing, error) {
	var err error
	var ok bool
	var intID int
	updates := make(map[string]interface{})

	if intID, ok = ToIntID(id); !ok || !gdb.ThingExists(intID) {
		return nil, gostErrors.NewRequestNotFound(errors.New("Thing does not exist"))
	}

	if len(thing.Name) > 0 {
		updates["name"] = thing.Name
	}

	if len(thing.Description) > 0 {
		updates["description"] = thing.Description
	}

	if len(thing.Properties) > 0 {
		jsonProperties, _ := json.Marshal(thing.Properties)
		updates["properties"] = string(jsonProperties[:])
	}

	if thing.Locations != nil {
		if len(thing.Locations) > 0 {
			for _, l := range thing.Locations {
				location, _ := gdb.GetLocation(l.ID, nil)

				// todo: check if location exist
				if location != nil {
					sql := fmt.Sprintf("update %s.thing_to_location set location_id  = $1 where thing_id= $2", gdb.Schema)
					res, err := gdb.Db.Exec(sql, location.ID, intID)
					if err != nil {
						return nil, err
					}
					if c, _ := res.RowsAffected(); c == 0 {
						sqlInsert := fmt.Sprintf("insert into %s.thing_to_location (location_id,thing_id) values ($1, $2)", gdb.Schema)
						_, err := gdb.Db.Exec(sqlInsert, location.ID, intID)
						if err != nil {
							return nil, err
						}
					}
				}
			}
		}
	}

	if err = gdb.updateEntityColumns("thing", updates, intID); err != nil {
		return nil, err
	}

	nt, _ := gdb.GetThing(intID, nil)
	return nt, nil
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
	return DeleteEntity(gdb, id, "thing")
}
