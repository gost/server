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

var lMapping = map[string]string{"location": "public.ST_AsGeoJSON(location.location)"}

// GetLocation retrieves the location for the given id from the database
func (gdb *GostDatabase) GetLocation(id interface{}, qo *odata.QueryOptions) (*entities.Location, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Location does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Location{}, qo, "", "", lMapping)+" AS location from %s.location where id = %v", gdb.Schema, intID)
	return processLocation(gdb.Db, sql, qo)
}

// GetLocations retrieves all locations
func (gdb *GostDatabase) GetLocations(qo *odata.QueryOptions) ([]*entities.Location, int, error) {
	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Location{}, qo, "", "", lMapping)+" AS location from %s.location order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema)
	countSQL := fmt.Sprintf("select COUNT(*) FROM %s.location", gdb.Schema)
	return processLocations(gdb.Db, sql, qo, countSQL)
}

// GetLocationsByHistoricalLocation retrieves all locations linked to the given HistoricalLocation
func (gdb *GostDatabase) GetLocationsByHistoricalLocation(hlID interface{}, qo *odata.QueryOptions) ([]*entities.Location, int, error) {
	intID, ok := ToIntID(hlID)
	if !ok {
		return nil, 0, gostErrors.NewRequestNotFound(errors.New("HistoricaLocation does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Location{}, qo, "location.", "", lMapping)+" AS location from %s.location inner join %s.historicallocation on historicallocation.location_id = location.id where historicallocation.id = %v order by id desc limit 1", gdb.Schema, gdb.Schema, intID)
	countSQL := fmt.Sprintf("select COUNT(*) FROM %s.location inner join %s.historicallocation on historicallocation.location_id = location.id where historicallocation.id = %v", gdb.Schema, gdb.Schema, intID)
	return processLocations(gdb.Db, sql, qo, countSQL)
}

// GetLocationByDatastreamID returns a location linked to an observation
func (gdb *GostDatabase) GetLocationByDatastreamID(datastreamID interface{}) (*entities.Location, error) {
	intID, ok := ToIntID(datastreamID)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	sql := fmt.Sprintf("SELECT "+CreateSelectString(&entities.Location{}, nil, "location.", "", lMapping)+" FROM %s.location INNER JOIN %s.thing_to_location on location.id = thing_to_location.location_id INNER JOIN %s.datastream on thing_to_location.thing_id = datastream.thing_id WHERE datastream.id = %v ORDER BY location.id DESC LIMIT 1", gdb.Schema, gdb.Schema, gdb.Schema, intID)
	return processLocation(gdb.Db, sql, nil)
}

// GetLocationsByThing retrieves all locations linked to the given thing
func (gdb *GostDatabase) GetLocationsByThing(thingID interface{}, qo *odata.QueryOptions) ([]*entities.Location, int, error) {
	intID, ok := ToIntID(thingID)
	if !ok {
		return nil, 0, gostErrors.NewRequestNotFound(errors.New("Thing does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Location{}, qo, "location.", "", lMapping)+" AS location from %s.location inner join %s.thing_to_location on thing_to_location.location_id = location.id where thing_to_location.thing_id = %v order by id desc limit 1", gdb.Schema, gdb.Schema, intID)
	countSQL := fmt.Sprintf("select COUNT(*) FROM %s.location inner join %s.thing_to_location on thing_to_location.location_id = location.id where thing_to_location.thing_id = %v", gdb.Schema, gdb.Schema, intID)
	return processLocations(gdb.Db, sql, qo, countSQL)
}

func processLocation(db *sql.DB, sql string, qo *odata.QueryOptions) (*entities.Location, error) {
	locations, _, err := processLocations(db, sql, qo, "")
	if err != nil {
		return nil, err
	}

	if len(locations) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Location not found"))
	}

	return locations[0], nil
}

func processLocations(db *sql.DB, sql string, qo *odata.QueryOptions, countSQL string) ([]*entities.Location, int, error) {
	rows, err := db.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, 0, err
	}

	var locations = []*entities.Location{}
	for rows.Next() {
		var sensorID interface{}
		var encodingType int
		var name, description, location string

		var params []interface{}
		var qp []string
		if qo == nil || qo.QuerySelect == nil || len(qo.QuerySelect.Params) == 0 {
			s := &entities.Location{}
			qp = s.GetPropertyNames()
		} else {
			qp = qo.QuerySelect.Params
		}

		for _, p := range qp {
			if p == "id" {
				params = append(params, &sensorID)
			}
			if p == "encodingType" {
				params = append(params, &encodingType)
			}
			if p == "name" {
				params = append(params, &name)
			}
			if p == "description" {
				params = append(params, &description)
			}
			if p == "location" {
				params = append(params, &location)
			}
		}

		err = rows.Scan(params...)
		if err != nil {
			return nil, 0, err
		}

		locationMap, err := JSONToMap(&location)
		if err != nil {
			return nil, 0, err
		}

		l := entities.Location{}
		l.ID = sensorID
		l.Name = name
		l.Description = description
		l.Location = locationMap
		if encodingType != 0 {
			l.EncodingType = entities.EncodingValues[encodingType].Value
		}
		locations = append(locations, &l)
	}

	var count int
	if len(countSQL) > 0 {
		db.QueryRow(countSQL).Scan(&count)
	}

	return locations, count, nil
}

// PostLocation receives a posted location entity and adds it to the database
// returns the created Location including the generated id
func (gdb *GostDatabase) PostLocation(location *entities.Location) (*entities.Location, error) {
	var locationID int
	locationBytes, _ := json.Marshal(location.Location)
	encoding, _ := entities.CreateEncodingType(location.EncodingType)

	sql := fmt.Sprintf("INSERT INTO %s.location (name, description, encodingtype, location) VALUES ($1, $2, $3, ST_SetSRID(ST_GeomFromGeoJSON('%s'),4326)) RETURNING id", gdb.Schema, string(locationBytes[:]))
	err := gdb.Db.QueryRow(sql, location.Name, location.Description, encoding.Code).Scan(&locationID)
	if err != nil {
		return nil, err
	}

	location.ID = locationID
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
