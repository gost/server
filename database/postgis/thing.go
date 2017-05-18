package postgis

import (
	"encoding/json"
	"fmt"

	"github.com/geodan/gost/sensorthings/entities"

	"database/sql"
	"errors"

	gostErrors "github.com/geodan/gost/errors"
	"github.com/geodan/gost/sensorthings/odata"
)

func thingParamFactory(values map[string]interface{}) (entities.Entity, error) {
	t := &entities.Thing{}
	for as, value := range values {
		if value == nil {
			continue
		}

		if as == asMappings[entities.EntityTypeThing][thingID] {
			t.ID = value
		} else if as == asMappings[entities.EntityTypeThing][thingName] {
			t.Name = value.(string)
		} else if as == asMappings[entities.EntityTypeThing][thingDescription] {
			t.Description = value.(string)
		} else if as == asMappings[entities.EntityTypeThing][thingProperties] {
			p := value.(string)
			propertiesMap, err := JSONToMap(&p)
			if err != nil {
				return nil, err
			}

			t.Properties = propertiesMap
		}
	}

	return t, nil
}

// GetThing returns a thing entity based on id and query
func (gdb *GostDatabase) GetThing(id interface{}, qo *odata.QueryOptions) (*entities.Thing, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Thing does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Thing{}, nil, intID, qo)
	return processThing(gdb.Db, query, qi)
}

//GetThingByDatastream retrieves the thing linked to a datastream
func (gdb *GostDatabase) GetThingByDatastream(id interface{}, qo *odata.QueryOptions) (*entities.Thing, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Thing{}, &entities.Datastream{}, intID, qo)
	return processThing(gdb.Db, query, qi)
}

//GetThingsByLocation retrieves the thing linked to a location
func (gdb *GostDatabase) GetThingsByLocation(id interface{}, qo *odata.QueryOptions) ([]*entities.Thing, int, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, 0, gostErrors.NewRequestNotFound(errors.New("Location does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Thing{}, &entities.Location{}, intID, qo)
	countSQL := gdb.QueryBuilder.CreateCountQuery(&entities.Thing{}, &entities.Location{}, intID, qo)
	return processThings(gdb.Db, query, qi, countSQL)
}

//GetThingByHistoricalLocation retrieves the thing linked to a HistoricalLocation
func (gdb *GostDatabase) GetThingByHistoricalLocation(id interface{}, qo *odata.QueryOptions) (*entities.Thing, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("HistoricalLocation does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Thing{}, &entities.HistoricalLocation{}, intID, qo)
	return processThing(gdb.Db, query, qi)
}

// GetThings returns an array of things
func (gdb *GostDatabase) GetThings(qo *odata.QueryOptions) ([]*entities.Thing, int, error) {
	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Thing{}, nil, nil, qo)
	countSQL := gdb.QueryBuilder.CreateCountQuery(&entities.Thing{}, nil, nil, qo)
	return processThings(gdb.Db, query, qi, countSQL)
}

func processThing(db *sql.DB, sql string, qi *QueryParseInfo) (*entities.Thing, error) {
	things, _, err := processThings(db, sql, qi, "")
	if err != nil {
		return nil, err
	}

	if len(things) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Thing not found"))
	}

	return things[0], nil
}

func processThings(db *sql.DB, sql string, qi *QueryParseInfo, countSQL string) ([]*entities.Thing, int, error) {
	data, err := ExecuteSelect(db, qi, sql)
	if err != nil {
		return nil, 0, fmt.Errorf("Error executing query %v", err)
	}

	things := make([]*entities.Thing, 0)
	for _, d := range data {
		entity := d.(*entities.Thing)
		things = append(things, entity)
	}

	var count int
	if len(countSQL) > 0 {
		count, err = ExecuteSelectCount(db, countSQL)
		if err != nil {
			return nil, 0, fmt.Errorf("Error executing count %v", err)
		}
	}

	return things, count, nil
}

// PostThing receives a posted thing entity and adds it to the database
// returns the created Thing including the generated id
func (gdb *GostDatabase) PostThing(thing *entities.Thing) (*entities.Thing, error) {
	jsonProperties, _ := json.Marshal(thing.Properties)
	var thingID int
	query := fmt.Sprintf("INSERT INTO %s.thing (name, description, properties) VALUES ($1, $2, $3) RETURNING id", gdb.Schema)
	err := gdb.Db.QueryRow(query, thing.Name, thing.Description, jsonProperties).Scan(&thingID)
	if err != nil {
		return nil, err
	}

	thing.ID = thingID
	return thing, nil
}

// PutThing receives a Thing entity and changes it in the database
// returns the adapted Thing
func (gdb *GostDatabase) PutThing(id interface{}, thing *entities.Thing) (*entities.Thing, error) {
	return gdb.PatchThing(id, thing)
}

// PatchThing receives a to be patched Thing entity and changes it in the database
// returns the patched Thing
func (gdb *GostDatabase) PatchThing(id interface{}, thing *entities.Thing) (*entities.Thing, error) {
	var err error
	var ok bool
	var intID int
	updates := make(map[string]interface{})

	thing.ID = id
	if intID, ok = ToIntID(id); !ok || !gdb.ThingExists(intID) {
		return nil, gostErrors.NewRequestNotFound(errors.New("Thing does not exist"))
	}

	if len(thing.Name) > 0 {
		updates[thingName] = thing.Name
	}

	if len(thing.Description) > 0 {
		updates[thingDescription] = thing.Description
	}

	if len(thing.Properties) > 0 {
		jsonProperties, _ := json.Marshal(thing.Properties)
		updates[thingProperties] = string(jsonProperties[:])
	}

	if thing.Locations != nil {
		if len(thing.Locations) > 0 {
			for _, l := range thing.Locations {
				location, _ := gdb.GetLocation(l.ID, nil)

				// todo: check if location exist
				if location != nil {
					query := fmt.Sprintf("update %s.thing_to_location set location_id  = $1 where thing_id= $2", gdb.Schema)
					res, err := gdb.Db.Exec(query, location.ID, intID)
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

					hl := &entities.HistoricalLocation{
						Thing:     thing,
						Locations: []*entities.Location{location},
					}

					hl.ContainsMandatoryParams()
					gdb.PostHistoricalLocation(hl)
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
func (gdb *GostDatabase) ThingExists(id interface{}) bool {
	return EntityExists(gdb, id, "thing")
}

// DeleteThing tries to delete a Thing by the given id
func (gdb *GostDatabase) DeleteThing(id interface{}) error {
	return DeleteEntity(gdb, id, "thing")
}
