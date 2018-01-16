package postgis

import (
	"encoding/json"
	"fmt"

	entities "github.com/gost/core"

	"database/sql"
	"errors"

	"github.com/gost/godata"
	gostErrors "github.com/gost/server/errors"
	"github.com/gost/server/sensorthings/odata"
)

func locationParamFactory(values map[string]interface{}) (entities.Entity, error) {
	l := &entities.Location{}
	for as, value := range values {
		if value == nil {
			continue
		}

		if as == asMappings[entities.EntityTypeLocation][locationID] {
			l.ID = value
		} else if as == asMappings[entities.EntityTypeLocation][locationName] {
			l.Name = value.(string)
		} else if as == asMappings[entities.EntityTypeLocation][locationDescription] {
			l.Description = value.(string)
		} else if as == asMappings[entities.EntityTypeLocation][locationEncodingType] {
			encodingType := value.(int64)
			if encodingType != 0 {
				l.EncodingType = entities.EncodingValues[encodingType].Value
			}
		} else if as == asMappings[entities.EntityTypeLocation][locationLocation] {
			t := value.(string)
			locationMap, err := JSONToMap(&t)
			if err != nil {
				return nil, err
			}

			l.Location = locationMap
		}
	}

	return l, nil
}

// GetLocation retrieves the location for the given id from the database
func (gdb *GostDatabase) GetLocation(id interface{}, qo *odata.QueryOptions) (*entities.Location, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Location does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Location{}, nil, intID, qo)
	return processLocation(gdb.Db, query, qi)
}

// GetLocations retrieves all locations
func (gdb *GostDatabase) GetLocations(qo *odata.QueryOptions) ([]*entities.Location, int, bool, error) {
	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Location{}, nil, nil, qo)
	countSQL := gdb.QueryBuilder.CreateCountQuery(&entities.Location{}, nil, nil, qo)
	return processLocations(gdb.Db, query, qo, qi, countSQL, false)
}

// GetLocationsByHistoricalLocation retrieves all locations linked to the given HistoricalLocation
func (gdb *GostDatabase) GetLocationsByHistoricalLocation(hlID interface{}, qo *odata.QueryOptions) ([]*entities.Location, int, bool, error) {
	intID, ok := ToIntID(hlID)
	if !ok {
		return nil, 0, false, gostErrors.NewRequestNotFound(errors.New("HistoricaLocation does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Location{}, &entities.HistoricalLocation{}, intID, qo)
	countSQL := gdb.QueryBuilder.CreateCountQuery(&entities.Location{}, &entities.HistoricalLocation{}, intID, qo)
	return processLocations(gdb.Db, query, qo, qi, countSQL, true)
}

// GetLocationByDatastreamID returns a location linked to an observation
// todo fix staticcheck error: 'argument qo is overwritten before first use'
// remove qo parameter? or change function
func (gdb *GostDatabase) GetLocationByDatastreamID(datastreamID interface{}, qo *odata.QueryOptions) (*entities.Location, error) {
	intID, ok := ToIntID(datastreamID)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("datastream does not exist"))
	}

	qo = &odata.QueryOptions{}
	tq := godata.GoDataTopQuery(-1)
	qo.Top = &tq

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Location{}, &entities.Datastream{}, intID, qo)
	return processLocation(gdb.Db, query, qi)
}

// GetLocationsByThing retrieves all locations linked to the given thing
func (gdb *GostDatabase) GetLocationsByThing(thingID interface{}, qo *odata.QueryOptions) ([]*entities.Location, int, bool, error) {
	intID, ok := ToIntID(thingID)
	if !ok {
		return nil, 0, false, gostErrors.NewRequestNotFound(errors.New("Thing does not exist"))
	}

	tq := godata.GoDataTopQuery(1)
	qo.Top = &tq

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Location{}, &entities.Thing{}, intID, qo)
	countSQL := gdb.QueryBuilder.CreateCountQuery(&entities.Location{}, &entities.Thing{}, intID, qo)
	return processLocations(gdb.Db, query, qo, qi, countSQL, true)
}

func processLocation(db *sql.DB, sql string, qi *QueryParseInfo) (*entities.Location, error) {
	locations, _, _, err := processLocations(db, sql, nil, qi, "", false)
	if err != nil {
		return nil, err
	}

	if len(locations) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Location not found"))
	}

	return locations[0], nil
}

func processLocations(db *sql.DB, sql string, qo *odata.QueryOptions, qi *QueryParseInfo, countSQL string, disableNextLink bool) ([]*entities.Location, int, bool, error) {
	data, hasNext, err := ExecuteSelect(db, qi, sql, qo)
	if err != nil {
		return nil, 0, hasNext, fmt.Errorf("Error executing query %v", err)
	}

	locations := make([]*entities.Location, 0)
	for _, d := range data {
		entity := d.(*entities.Location)
		locations = append(locations, entity)
	}

	var count int
	if len(countSQL) > 0 {
		count, err = ExecuteSelectCount(db, countSQL)
		if err != nil {
			return nil, 0, hasNext, fmt.Errorf("error executing count %v", err)
		}
	}

	if disableNextLink {
		hasNext = false
	}

	return locations, count, hasNext, nil
}

// PostLocation receives a posted location entity and adds it to the database
// returns the created Location including the generated id
func (gdb *GostDatabase) PostLocation(location *entities.Location) (*entities.Location, error) {
	var locationID int
	locationBytes, _ := json.Marshal(location.Location)
	encoding, _ := entities.CreateEncodingType(location.EncodingType)

	sql2 := fmt.Sprintf("INSERT INTO %s.location (name, description, encodingtype, location) VALUES ($1, $2, $3, ST_SetSRID(ST_GeomFromGeoJSON('%s'),4326)) RETURNING id", gdb.Schema, string(locationBytes[:]))
	err := gdb.Db.QueryRow(sql2, location.Name, location.Description, encoding.Code).Scan(&locationID)
	if err != nil {
		return nil, err
	}

	location.ID = locationID
	return location, nil
}

// LocationExists checks if a location is present in the database based on a given id
func (gdb *GostDatabase) LocationExists(id interface{}) bool {
	return EntityExists(gdb, id, "location")
}

// PatchLocation updates a Location in the database
func (gdb *GostDatabase) PatchLocation(id interface{}, l *entities.Location) (*entities.Location, error) {
	var err error
	var ok bool
	var intID int
	updates := make(map[string]interface{})

	if intID, ok = ToIntID(id); !ok || !gdb.LocationExists(intID) {
		return nil, gostErrors.NewRequestNotFound(errors.New("Location does not exist"))
	}

	if len(l.Name) > 0 {
		updates["name"] = l.Name
	}

	if len(l.Description) > 0 {
		updates["description"] = l.Description
	}

	if len(l.Location) > 0 {
		locationBytes, _ := json.Marshal(l.Location)
		updates["location"] = fmt.Sprintf("ST_SetSRID(ST_GeomFromGeoJSON('%s'),4326)", string(locationBytes[:]))
	}

	if len(l.EncodingType) > 0 {
		encoding, _ := entities.CreateEncodingType(l.EncodingType)
		updates["encodingtype"] = encoding.Code
	}

	if err = gdb.updateEntityColumns("location", updates, intID); err != nil {
		return nil, err
	}

	ns, _ := gdb.GetLocation(intID, nil)
	return ns, nil
}

// DeleteLocation removes a given location from the database
func (gdb *GostDatabase) DeleteLocation(id interface{}) error {
	return DeleteEntity(gdb, id, "location")
}

// PutLocation receives a Location entity and changes it in the database
// returns the adapted Location
func (gdb *GostDatabase) PutLocation(id interface{}, location *entities.Location) (*entities.Location, error) {
	return gdb.PatchLocation(id, location)
}

// LinkLocation links a thing with a location
// fails when a thing or location cannot be found for the given id's
func (gdb *GostDatabase) LinkLocation(thingID interface{}, locationID interface{}) error {
	tid, ok := ToIntID(thingID)
	if !ok || !gdb.ThingExists(tid) {
		return gostErrors.NewRequestNotFound(errors.New("Thing does not exist"))
	}

	lid, ok := ToIntID(locationID)
	if !ok || !gdb.LocationExists(lid) {
		return gostErrors.NewRequestNotFound(errors.New("Location does not exist"))
	}

	sql2 := fmt.Sprintf("INSERT INTO %s.thing_to_location (thing_id, location_id) VALUES ($1, $2)", gdb.Schema)
	_, err3 := gdb.Db.Exec(sql2, tid, lid)

	return err3
}
