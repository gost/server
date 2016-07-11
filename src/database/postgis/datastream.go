package postgis

import (
	"encoding/json"
	"errors"
	"fmt"

	"database/sql"

	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/odata"
)

var totalDatastreams int
var dsMapping = map[string]string{"observedArea": "public.ST_AsGeoJSON(datastream.observedarea) AS observedarea"}

// GetTotalDatastreams returns the amount of datastreams in the database
func (gdb *GostDatabase) GetTotalDatastreams() int {
	return totalDatastreams
}

// GetDatastream retrieves a datastream by id
func (gdb *GostDatabase) GetDatastream(id interface{}, qo *odata.QueryOptions) (*entities.Datastream, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Datastream{}, qo, "", "", dsMapping)+" FROM %s.datastream where id = %v", gdb.Schema, intID)
	datastream, err := processDatastream(gdb.Db, sql, qo)
	if err != nil {
		return nil, err
	}

	return datastream, nil
}

// GetDatastreams retrieves all datastreams
func (gdb *GostDatabase) GetDatastreams(qo *odata.QueryOptions) ([]*entities.Datastream, error) {
	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Datastream{}, qo, "", "", dsMapping)+" FROM %s.datastream order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema)
	return processDatastreams(gdb.Db, sql, qo)
}

// GetDatastreamByObservation retrieves a datastream linked to the given observation
func (gdb *GostDatabase) GetDatastreamByObservation(observationID interface{}, qo *odata.QueryOptions) (*entities.Datastream, error) {
	tID, ok := ToIntID(observationID)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Datastream{}, qo, "datastream.", "", dsMapping)+" FROM %s.datastream inner join %s.observation on datastream.id = observation.stream_id where observation.id = %v", gdb.Schema, gdb.Schema, tID)
	return processDatastream(gdb.Db, sql, qo)
}

// GetDatastreamsByThing retrieves all datastreams linked to the given thing
func (gdb *GostDatabase) GetDatastreamsByThing(thingID interface{}, qo *odata.QueryOptions) ([]*entities.Datastream, error) {
	tID, ok := ToIntID(thingID)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Datastream{}, qo, "datastream.", "", dsMapping)+" FROM %s.datastream inner join %s.thing on thing.id = datastream.thing_id where thing.id = %v order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema, gdb.Schema, tID)
	return processDatastreams(gdb.Db, sql, qo)
}

// GetDatastreamsBySensor retrieves all datastreams linked to the given sensor
func (gdb *GostDatabase) GetDatastreamsBySensor(sensorID interface{}, qo *odata.QueryOptions) ([]*entities.Datastream, error) {
	tID, ok := ToIntID(sensorID)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Datastream{}, qo, "datastream.", "", dsMapping)+" FROM %s.datastream inner join %s.sensor on sensor.id = datastream.sensor_id where sensor.id = %v order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema, gdb.Schema, tID)
	return processDatastreams(gdb.Db, sql, qo)
}

// GetDatastreamsByObservedProperty retrieves all datastreams linked to the given ObservedProerty
func (gdb *GostDatabase) GetDatastreamsByObservedProperty(oID interface{}, qo *odata.QueryOptions) ([]*entities.Datastream, error) {
	tID, ok := ToIntID(oID)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Datastream{}, qo, "datastream.", "", dsMapping)+" FROM %s.datastream inner join %s.observedproperty on observedproperty.id = datastream.observedproperty_id where observedproperty.id = %v order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema, gdb.Schema, tID)
	return processDatastreams(gdb.Db, sql, qo)
}

func processDatastream(db *sql.DB, sql string, qo *odata.QueryOptions) (*entities.Datastream, error) {
	datastreams, err := processDatastreams(db, sql, qo)
	if err != nil {
		return nil, err
	}

	if len(datastreams) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	return datastreams[0], nil
}

func processDatastreams(db *sql.DB, sql string, qo *odata.QueryOptions) ([]*entities.Datastream, error) {
	rows, err := db.Query(sql)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var datastreams = []*entities.Datastream{}
	for rows.Next() {
		var id interface{}
		var description, unitofmeasurement string
		var observedarea *string
		var ot int

		var params []interface{}
		var qp []string
		if qo == nil || qo.QuerySelect == nil || len(qo.QuerySelect.Params) == 0 {
			d := &entities.Datastream{}
			qp = d.GetPropertyNames()
		} else {
			qp = qo.QuerySelect.Params
		}

		for _, p := range qp {
			if p == "id" {
				params = append(params, &id)
			}
			if p == "description" {
				params = append(params, &description)
			}
			if p == "unitOfMeasurement" {
				params = append(params, &unitofmeasurement)
			}
			if p == "observationType" {
				params = append(params, &ot)
			}
			if p == "observedArea" {
				params = append(params, &observedarea)
			}
		}

		err = rows.Scan(params...)

		unitOfMeasurementMap, err := JSONToMap(&unitofmeasurement)
		if err != nil {
			return nil, err
		}

		observedAreaMap, err := JSONToMap(observedarea)
		if err != nil {
			return nil, err
		}

		datastream := entities.Datastream{}
		datastream.ID = id
		datastream.Description = description
		datastream.UnitOfMeasurement = unitOfMeasurementMap
		datastream.ObservedArea = observedAreaMap
		if ot != 0 {
			obs, _ := entities.GetObservationTypeByID(ot)
			datastream.ObservationType = obs.Value
		}

		datastreams = append(datastreams, &datastream)
	}

	return datastreams, nil
}

// PostDatastream todo
// TODO: !!!!ADD phenomenonTime SUPPORT!!!!
// TODO: !!!!ADD resulttime SUPPORT!!!!
func (gdb *GostDatabase) PostDatastream(d *entities.Datastream) (*entities.Datastream, error) {
	var dsID int

	tID, ok := ToIntID(d.Thing.ID)
	if !ok || !gdb.ThingExists(tID) {
		return nil, gostErrors.NewBadRequestError(errors.New("Thing does not exist"))
	}

	sID, ok := ToIntID(d.Sensor.ID)
	if !ok || !gdb.SensorExists(sID) {
		return nil, gostErrors.NewBadRequestError(errors.New("Sensor does not exist"))
	}

	oID, ok := ToIntID(d.ObservedProperty.ID)
	if !ok || !gdb.ObservedPropertyExists(oID) {
		return nil, gostErrors.NewBadRequestError(errors.New("ObservedProperty does not exist"))
	}

	unitOfMeasurement, _ := json.Marshal(d.UnitOfMeasurement)
	geom := "NULL"
	if len(d.ObservedArea) != 0 {
		observedAreaBytes, _ := json.Marshal(d.ObservedArea)
		geom = fmt.Sprintf("ST_SetSRID(ST_GeomFromGeoJSON('%s'),4326)", string(observedAreaBytes[:]))
	}

	// get the ObservationType id in the lookup table
	observationType, err := entities.GetObservationTypeByValue(d.ObservationType)

	if err != nil {
		return nil, gostErrors.NewBadRequestError(errors.New("ObservationType does not exist"))
	}

	sql := fmt.Sprintf("INSERT INTO %s.datastream (description, unitofmeasurement, observedarea, thing_id, sensor_id, observedproperty_id, observationtype) VALUES ($1, $2, %s, $3, $4, $5, $6) RETURNING id", gdb.Schema, geom)
	err = gdb.Db.QueryRow(sql, d.Description, unitOfMeasurement, tID, sID, oID, observationType.Code).Scan(&dsID)
	if err != nil {
		return nil, err
	}

	d.ID = dsID

	// clear inner entities to serves links upon response
	d.Thing = nil
	d.Sensor = nil
	d.ObservedProperty = nil

	totalDatastreams++
	return d, nil
}

// PatchDatastream updates a Datastream in the database
func (gdb *GostDatabase) PatchDatastream(id interface{}, ds *entities.Datastream) (*entities.Datastream, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("Not implemented"))
}

// DeleteDatastream tries to delete a Datastream by the given id
func (gdb *GostDatabase) DeleteDatastream(id interface{}) error {
	intID, ok := ToIntID(id)
	if !ok {
		return gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	r, err := gdb.Db.Exec(fmt.Sprintf("DELETE FROM %s.datastream WHERE id = $1", gdb.Schema), intID)
	if err != nil {
		return err
	}

	if c, _ := r.RowsAffected(); c == 0 {
		return gostErrors.NewRequestNotFound(errors.New("Datastream not found"))
	}

	totalDatastreams--
	return nil
}

// DatastreamExists checks if a Datastream is present in the database based on a given id
func (gdb *GostDatabase) DatastreamExists(databaseID int) bool {
	var result bool
	sql := fmt.Sprintf("SELECT exists (SELECT 1 FROM %s.datastream WHERE id = $1 LIMIT 1)", gdb.Schema)
	err := gdb.Db.QueryRow(sql, databaseID).Scan(&result)
	if err != nil {
		return false
	}

	return result
}
