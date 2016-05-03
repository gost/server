package postgis

import (
	"fmt"
	"github.com/geodan/gost/sensorthings/entities"
	"strconv"

	"encoding/json"
	gostErrors "github.com/geodan/gost/errors"
)

// GetLocation retrieves the location for the given id from the database
func (gdb *GostDatabase) GetLocation(id string) (*entities.Location, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var sensorID, encodingtype int
	var description, location string
	sql := fmt.Sprintf("select id, description, encodingtype, ST_AsGeoJSON(location) AS location from %s.location where id = $1", gdb.Schema)
	err = gdb.Db.QueryRow(sql, intID).Scan(&sensorID, &description, &encodingtype, &location)

	if err != nil {
		return nil, gostErrors.NewRequestNotFound(fmt.Errorf("Locations(%s) does not exist", id))
	}

	locationMap, err := JSONToMap(location)
	if err != nil {
		return nil, err
	}

	l := entities.Location{
		ID:          strconv.Itoa(sensorID),
		Description: description,
		Location:    locationMap,
	}

	return &l, nil
}

// GetLocations todo
func (gdb *GostDatabase) GetLocations() ([]*entities.Location, error) {
	return nil, nil
}

// PostLocation receives a posted location entity and adds it to the database
// returns the created Location including the generated id
//TODO: ENCODINGTYPE
func (gdb *GostDatabase) PostLocation(location entities.Location) (*entities.Location, error) {
	var locationID int
	locationBytes, _ := json.Marshal(location.Location)
	jsonToGeom := fmt.Sprintf("ST_GeomFromGeoJSON('%s')", string(locationBytes[:]))
	sql := fmt.Sprintf("INSERT INTO %s.location (description, encodingtype, location) VALUES ($1, $2, %s) RETURNING id", gdb.Schema, jsonToGeom)
	err := gdb.Db.QueryRow(sql, location.Description, 1).Scan(&locationID)
	if err != nil {
		return nil, err
	}

	location.ID = strconv.Itoa(locationID)
	return &location, nil
}

// LocationExists checks if a location is present in the database based on a given id
func (gdb *GostDatabase) LocationExists(locationID int) bool {
	var result bool
	sql := fmt.Sprintf("SELECT exists (SELECT 1 FROM  %s.location WHERE id = $1 LIMIT 1)", gdb.Schema)
	err := gdb.Db.QueryRow(sql, locationID).Scan(&result)
	if err != nil {
		return false
	}

	return result
}

// LinkLocation links a thing with a location
// fails when a thing or location cannot be found for the given id's
func (gdb *GostDatabase) LinkLocation(thingID string, locationID string) error {
	tid, err := strconv.Atoi(thingID)
	if !gdb.ThingExists(tid) || err != nil {
		return fmt.Errorf("Thing(%v) does not exist", thingID)
	}

	lid, err2 := strconv.Atoi(thingID)
	if !gdb.ThingExists(lid) || err2 != nil {
		return fmt.Errorf("Location(%v) does not exist", locationID)
	}

	sql := fmt.Sprintf("INSERT INTO %s.thing_to_location (thing_id, location_id) VALUES ($1, $2)", gdb.Schema)
	_, err3 := gdb.Db.Exec(sql, tid, lid)
	if err3 != nil {
		return err3
	}

	return nil
}
