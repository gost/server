package postgis

import (
	"fmt"
	"github.com/geodan/gost/src/sensorthings/entities"
	"strconv"

	"database/sql"
	"encoding/json"
	"errors"
	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// GetLocation retrieves the location for the given id from the database
func (gdb *GostDatabase) GetLocation(id interface{}, qo *odata.QueryOptions) (*entities.Location, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Location does not exist"))
	}

	sql := fmt.Sprintf("select id, description, encodingtype, public.ST_AsGeoJSON(location) AS location from %s.location where id = $1", gdb.Schema)
	return processLocation(gdb.Db, sql, intID)
}

// GetLocations retrieves all locations
func (gdb *GostDatabase) GetLocations(qo *odata.QueryOptions) ([]*entities.Location, error) {
	sql := fmt.Sprintf("select id, description, encodingtype, public.ST_AsGeoJSON(location) AS location from %s.location", gdb.Schema)
	return processLocations(gdb.Db, sql)
}

// GetLocationsByHistoricalLocation retrieves all locations linked to the given HistoricalLocation
func (gdb *GostDatabase) GetLocationsByHistoricalLocation(hlID interface{}, qo *odata.QueryOptions) ([]*entities.Location, error) {
	intID, ok := ToIntID(hlID)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("HistoricaLocation does not exist"))
	}

	sql := fmt.Sprintf("select location.id, location.description, location.encodingtype, public.ST_AsGeoJSON(location.location) AS location from %s.location inner join %s.historicallocation on historicallocation.location_id = location.id where historicallocation.id = $1 limit 1", gdb.Schema, gdb.Schema)
	return processLocations(gdb.Db, sql, intID)
}

// GetLocationsByThing retrieves all locations linked to the given thing
func (gdb *GostDatabase) GetLocationsByThing(thingID interface{}, qo *odata.QueryOptions) ([]*entities.Location, error) {
	intID, ok := ToIntID(thingID)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Thing does not exist"))
	}

	sql := fmt.Sprintf("select location.id, location.description, location.encodingtype, public.ST_AsGeoJSON(location.location) AS location from %s.location inner join %s.thing_to_location on thing_to_location.location_id = location.id where thing_to_location.thing_id = $1 limit 1", gdb.Schema, gdb.Schema)
	return processLocations(gdb.Db, sql, intID)
}

func processLocation(db *sql.DB, sql string, args ...interface{}) (*entities.Location, error) {
	locations, err := processLocations(db, sql, args...)
	if err != nil {
		return nil, err
	}

	if len(locations) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Location not found"))
	}

	return locations[0], nil
}

func processLocations(db *sql.DB, sql string, args ...interface{}) ([]*entities.Location, error) {
	rows, err := db.Query(sql, args...)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var locations = []*entities.Location{}
	for rows.Next() {
		var sensorID, encodingtype int
		var description, location string
		err = rows.Scan(&sensorID, &description, &encodingtype, &location)
		if err != nil {
			return nil, err
		}

		locationMap, err := JSONToMap(&location)
		if err != nil {
			return nil, err
		}

		l := entities.Location{}
		l.ID = strconv.Itoa(sensorID)
		l.Description = description
		l.Location = locationMap
		l.EncodingType = entities.EncodingValues[encodingtype].Value

		locations = append(locations, &l)
	}

	return locations, nil
}

// PostLocation receives a posted location entity and adds it to the database
// returns the created Location including the generated id
func (gdb *GostDatabase) PostLocation(location *entities.Location) (*entities.Location, error) {
	var locationID int
	locationBytes, _ := json.Marshal(location.Location)
	encoding, _ := entities.CreateEncodingType(location.EncodingType)
	sql := fmt.Sprintf("INSERT INTO %s.location (description, encodingtype, location) VALUES ($1, $2, public.ST_GeomFromGeoJSON('%s')) RETURNING id", gdb.Schema, string(locationBytes[:]))
	err := gdb.Db.QueryRow(sql, location.Description, encoding.Code).Scan(&locationID)
	if err != nil {
		return nil, err
	}

	location.ID = strconv.Itoa(locationID)
	return location, nil
}

// LocationExists checks if a location is present in the database based on a given id
func (gdb *GostDatabase) LocationExists(locationID interface{}) bool {
	var result bool
	sql := fmt.Sprintf("SELECT exists (SELECT 1 FROM  %s.location WHERE id = $1 LIMIT 1)", gdb.Schema)
	err := gdb.Db.QueryRow(sql, locationID).Scan(&result)
	if err != nil {
		return false
	}

	return result
}

// DeleteLocation removes a given location from the database
func (gdb *GostDatabase) DeleteLocation(locationID interface{}) error {
	intID, ok := ToIntID(locationID)
	if !ok {
		return gostErrors.NewRequestNotFound(errors.New("Location does not exist"))
	}

	sql := fmt.Sprintf("DELETE FROM %s.location WHERE id = $1", gdb.Schema)
	r, err := gdb.Db.Exec(sql, intID)
	if err != nil {
		return err
	}

	if c, _ := r.RowsAffected(); c == 0 {
		return gostErrors.NewRequestNotFound(errors.New("Location not found"))
	}

	return nil
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

	sql := fmt.Sprintf("INSERT INTO %s.thing_to_location (thing_id, location_id) VALUES ($1, $2)", gdb.Schema)
	_, err3 := gdb.Db.Exec(sql, tid, lid)
	if err3 != nil {
		return err3
	}

	return nil
}
