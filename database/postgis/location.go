package postgis

import (
	"fmt"
	"github.com/geodan/gost/sensorthings/entities"
	"strconv"

	"database/sql"
	"encoding/json"
	"errors"
	gostErrors "github.com/geodan/gost/errors"
)

// GetLocation retrieves the location for the given id from the database
func (gdb *GostDatabase) GetLocation(id string) (*entities.Location, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	sql := "select id, description, encodingtype, public.ST_AsGeoJSON(location) AS location from location where id = $1"
	return processLocation(gdb.Db, sql, intID)
}

// GetLocations retrieves all locations
func (gdb *GostDatabase) GetLocations() ([]*entities.Location, error) {
	sql := "select id, description, encodingtype, public.ST_AsGeoJSON(location) AS location from location"
	return processLocations(gdb.Db, sql)

}

// GetLocationsByThing retrieves all locations linked to the given thing
func (gdb *GostDatabase) GetLocationsByThing(thingID string) ([]*entities.Location, error) {
	intID, err := strconv.Atoi(thingID)
	if err != nil {
		return nil, err
	}

	sql := "select location.id, location.description, location.encodingtype, public.ST_AsGeoJSON(location.location) AS location from location inner join thing_to_location on thing_to_location.location_id = location.id where thing_to_location.thing_id = $1 limit 1"
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

		l := entities.Location{
			ID:          strconv.Itoa(sensorID),
			Description: description,
			Location:    locationMap,
		}

		locations = append(locations, &l)
	}

	return locations, nil
}

// PostLocation receives a posted location entity and adds it to the database
// returns the created Location including the generated id
//TODO: ENCODINGTYPE
func (gdb *GostDatabase) PostLocation(location entities.Location) (*entities.Location, error) {
	var locationID int
	locationBytes, _ := json.Marshal(location.Location)
	sql := fmt.Sprintf("INSERT INTO location (description, encodingtype, location) VALUES ($1, $2, public.ST_GeomFromGeoJSON('%s')) RETURNING id", string(locationBytes[:]))
	fmt.Println(sql)
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
	sql := "SELECT exists (SELECT 1 FROM  location WHERE id = $1 LIMIT 1)"
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
		return fmt.Errorf("Thing(%s) does not exist", thingID)
	}

	lid, err2 := strconv.Atoi(thingID)
	if !gdb.ThingExists(lid) || err2 != nil {
		return fmt.Errorf("Location(%s) does not exist", locationID)
	}

	sql := "INSERT INTO thing_to_location (thing_id, location_id) VALUES ($1, $2)"
	_, err3 := gdb.Db.Exec(sql, tid, lid)
	if err3 != nil {
		return err3
	}

	return nil
}
