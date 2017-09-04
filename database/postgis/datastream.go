package postgis

import (
	"encoding/json"
	"errors"
	"fmt"

	"database/sql"

	"github.com/gost/now"
	gostErrors "github.com/gost/server/errors"
	entities "github.com/gost/core"
	"github.com/gost/server/sensorthings/odata"
)

func datastreamParamFactory(values map[string]interface{}) (entities.Entity, error) {
	ds := &entities.Datastream{}
	for as, value := range values {
		if value == nil {
			continue
		}

		if as == asMappings[entities.EntityTypeDatastream][datastreamID] {
			ds.ID = value
		} else if as == asMappings[entities.EntityTypeDatastream][datastreamObservedArea] {
			t := value.(string)
			observedAreaMap, err := JSONToMap(&t)
			if err != nil {
				return nil, err
			}
			ds.ObservedArea = observedAreaMap
		} else if as == asMappings[entities.EntityTypeDatastream][datastreamName] {
			ds.Name = value.(string)
		} else if as == asMappings[entities.EntityTypeDatastream][datastreamDescription] {
			ds.Description = value.(string)
		} else if as == asMappings[entities.EntityTypeDatastream][datastreamResultTime] {
			ds.ResultTime = now.PostgresToIso8601Period(value.(string))
		} else if as == asMappings[entities.EntityTypeDatastream][datastreamObservationType] {
			obs, _ := entities.GetObservationTypeByID(value.(int64))
			ds.ObservationType = obs.Value
		} else if as == asMappings[entities.EntityTypeDatastream][datastreamPhenomenonTime] {
			ds.PhenomenonTime = now.PostgresToIso8601Period(value.(string))
		} else if as == asMappings[entities.EntityTypeDatastream][datastreamUnitOfMeasurement] {
			t := value.(string)
			unitOfMeasurementMap, err := JSONToMap(&t)
			if err != nil {
				return nil, err
			}

			ds.UnitOfMeasurement = unitOfMeasurementMap
		}
	}

	return ds, nil
}

// GetObservedArea returns the observed area of all observations of datastream
func (gdb *GostDatabase) GetObservedArea(id int) (map[string]interface{}, error) {
	sqlString := "select ST_AsGeoJSON(ST_ConvexHull(ST_Collect(feature))) as geom from %s.featureofinterest where id in (select distinct featureofinterest_id from %s.observation where stream_id=%v)"
	sql2 := fmt.Sprintf(sqlString, gdb.Schema, gdb.Schema, id)
	rows, err := gdb.Db.Query(sql2)
	var geom string
	var propMap map[string]interface{}
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&geom)

		if err == nil {
			propMap, _ = JSONToMap(&geom)
		}
	}
	return propMap, err
}

// GetDatastream retrieves a datastream by id
func (gdb *GostDatabase) GetDatastream(id interface{}, qo *odata.QueryOptions) (*entities.Datastream, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Datastream{}, nil, intID, qo)
	datastream, err := processDatastream(gdb.Db, query, qi)
	if err != nil {
		return nil, err
	}

	if qo != nil {
		hasSelectQuery := (qo.Select != nil)
		var containsObservedArea = true
		if hasSelectQuery {
			containsObservedArea = ContainsToLower(qo.Select.SelectItems, "observedArea")
		}

		// calculate observedArea on the fly when not present in database
		if containsObservedArea {
			if datastream.ObservedArea == nil {
				observedArea, _ := gdb.GetObservedArea(intID)
				datastream.ObservedArea = observedArea
			}
		}
	}

	return datastream, nil
}

// GetDatastreams retrieves all datastreams
func (gdb *GostDatabase) GetDatastreams(qo *odata.QueryOptions) ([]*entities.Datastream, int, error) {
	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Datastream{}, nil, nil, qo)
	countSQL := gdb.QueryBuilder.CreateCountQuery(&entities.Datastream{}, nil, nil, qo)
	return processDatastreams(gdb.Db, query, qi, countSQL)
}

// GetDatastreamByObservation retrieves a datastream linked to the given observation
func (gdb *GostDatabase) GetDatastreamByObservation(observationID interface{}, qo *odata.QueryOptions) (*entities.Datastream, error) {
	intID, ok := ToIntID(observationID)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Datastream{}, &entities.Observation{}, intID, qo)
	return processDatastream(gdb.Db, query, qi)
}

// GetDatastreamsByThing retrieves all datastreams linked to the given thing
func (gdb *GostDatabase) GetDatastreamsByThing(thingID interface{}, qo *odata.QueryOptions) ([]*entities.Datastream, int, error) {
	intID, ok := ToIntID(thingID)
	if !ok {
		return nil, 0, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Datastream{}, &entities.Thing{}, intID, qo)
	countSQL := gdb.QueryBuilder.CreateCountQuery(&entities.Datastream{}, &entities.Thing{}, intID, qo)
	return processDatastreams(gdb.Db, query, qi, countSQL)
}

// GetDatastreamsBySensor retrieves all datastreams linked to the given sensor
func (gdb *GostDatabase) GetDatastreamsBySensor(sensorID interface{}, qo *odata.QueryOptions) ([]*entities.Datastream, int, error) {
	intID, ok := ToIntID(sensorID)
	if !ok {
		return nil, 0, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Datastream{}, &entities.Sensor{}, intID, qo)
	countSQL := gdb.QueryBuilder.CreateCountQuery(&entities.Datastream{}, &entities.Sensor{}, intID, qo)
	return processDatastreams(gdb.Db, query, qi, countSQL)
}

// GetDatastreamsByObservedProperty retrieves all datastreams linked to the given ObservedProerty
func (gdb *GostDatabase) GetDatastreamsByObservedProperty(oID interface{}, qo *odata.QueryOptions) ([]*entities.Datastream, int, error) {
	intID, ok := ToIntID(oID)
	if !ok {
		return nil, 0, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Datastream{}, &entities.ObservedProperty{}, intID, qo)
	countSQL := gdb.QueryBuilder.CreateCountQuery(&entities.Datastream{}, &entities.ObservedProperty{}, intID, qo)
	return processDatastreams(gdb.Db, query, qi, countSQL)
}

func processDatastream(db *sql.DB, sql string, qi *QueryParseInfo) (*entities.Datastream, error) {
	datastreams, _, err := processDatastreams(db, sql, qi, "")
	if err != nil {
		return nil, err
	}

	if len(datastreams) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	return datastreams[0], nil
}

func processDatastreams(db *sql.DB, sql string, qi *QueryParseInfo, countSQL string) ([]*entities.Datastream, int, error) {
	data, err := ExecuteSelect(db, qi, sql)
	if err != nil {
		return nil, 0, fmt.Errorf("Error executing query %v", err)
	}

	datastreams := make([]*entities.Datastream, 0)
	for _, d := range data {
		entity := d.(*entities.Datastream)
		datastreams = append(datastreams, entity)
	}

	var count int
	if len(countSQL) > 0 {
		count, err = ExecuteSelectCount(db, countSQL)
		if err != nil {
			return nil, 0, fmt.Errorf("Error executing count %v", err)
		}
	}

	return datastreams, count, nil
}

// CheckDatastreamRelationsExist check if the related entities exist
func CheckDatastreamRelationsExist(gdb *GostDatabase, d *entities.Datastream) error {
	var tID, sID, oID int
	var ok bool

	if tID, ok = ToIntID(d.Thing.ID); !ok || !gdb.ThingExists(tID) {
		return gostErrors.NewBadRequestError(errors.New("Thing does not exist"))
	}

	if sID, ok = ToIntID(d.Sensor.ID); !ok || !gdb.SensorExists(sID) {
		return gostErrors.NewBadRequestError(errors.New("Sensor does not exist"))
	}

	if oID, ok = ToIntID(d.ObservedProperty.ID); !ok || !gdb.ObservedPropertyExists(oID) {
		return gostErrors.NewBadRequestError(errors.New("ObservedProperty does not exist"))
	}
	return nil
}

// PostDatastream posts a datastream
func (gdb *GostDatabase) PostDatastream(d *entities.Datastream) (*entities.Datastream, error) {
	err := CheckDatastreamRelationsExist(gdb, d)
	if err != nil {
		return nil, err
	}
	tID, _ := ToIntID(d.Thing.ID)
	sID, _ := ToIntID(d.Sensor.ID)
	oID, _ := ToIntID(d.ObservedProperty.ID)
	var dsID int

	unitOfMeasurement, _ := json.Marshal(d.UnitOfMeasurement)
	geom := "NULL"
	if len(d.ObservedArea) != 0 {
		observedAreaBytes, _ := json.Marshal(d.ObservedArea)
		geom = fmt.Sprintf("ST_SetSRID(ST_GeomFromGeoJSON('%s'),4326)", string(observedAreaBytes[:]))
	}

	phenomenonTime := "NULL"
	if len(d.PhenomenonTime) != 0 {
		phenomenonTime = "'" + now.Iso8601ToPostgresPeriod(d.PhenomenonTime) + "'"
	}

	resultTime := "NULL"
	if len(d.ResultTime) != 0 {
		resultTime = "'" + now.Iso8601ToPostgresPeriod(d.ResultTime) + "'"
	}
	// get the ObservationType id in the lookup table
	observationType, err := entities.GetObservationTypeByValue(d.ObservationType)

	if err != nil {
		return nil, gostErrors.NewBadRequestError(errors.New("ObservationType does not exist"))
	}

	sql2 := fmt.Sprintf("INSERT INTO %s.datastream (name, description, unitofmeasurement, observedarea, thing_id, sensor_id, observedproperty_id, observationtype, phenomenonTime, resulttime) VALUES ($1, $2, $3, %s, $4, $5, $6, $7, %s, %s) RETURNING id", gdb.Schema, geom, phenomenonTime, resultTime)
	err = gdb.Db.QueryRow(sql2, d.Name, d.Description, unitOfMeasurement, tID, sID, oID, observationType.Code).Scan(&dsID)
	if err != nil {
		return nil, err
	}

	d.ID = dsID

	// clear inner entities to serves links upon response
	d.Thing = nil
	d.Sensor = nil
	d.ObservedProperty = nil

	return d, nil
}

// PatchDatastream updates a Datastream in the database
func (gdb *GostDatabase) PatchDatastream(id interface{}, ds *entities.Datastream) (*entities.Datastream, error) {
	var err error
	var ok bool
	var intID int
	updates := make(map[string]interface{})

	if intID, ok = ToIntID(id); !ok || !gdb.DatastreamExists(intID) {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	if len(ds.Name) > 0 {
		updates["name"] = ds.Name
	}

	if len(ds.Description) > 0 {
		updates["description"] = ds.Description
	}

	if len(ds.ObservationType) > 0 {
		observationType, err := entities.GetObservationTypeByValue(ds.ObservationType)
		if err != nil {
			return nil, gostErrors.NewBadRequestError(errors.New("ObservationType does not exist"))
		}

		updates["observationtype"] = observationType.Code
	}

	if len(ds.UnitOfMeasurement) > 0 {
		j, _ := json.Marshal(ds.UnitOfMeasurement)
		updates["unitofmeasurement"] = string(j[:])
	}

	if len(ds.ObservedArea) > 0 {
		observedAreaBytes, _ := json.Marshal(ds.ObservedArea)
		updates["observedarea"] = fmt.Sprintf("ST_SetSRID(ST_GeomFromGeoJSON('%s'),4326)", string(observedAreaBytes[:]))
	}

	if len(ds.PhenomenonTime) > 0 {
		phenomenonTime := now.Iso8601ToPostgresPeriod(ds.PhenomenonTime)
		updates["phenomenontime"] = phenomenonTime
	}

	if len(ds.ResultTime) > 0 {
		resultTime := now.Iso8601ToPostgresPeriod(ds.ResultTime)
		updates["resulttime"] = resultTime
	}

	if err = gdb.updateEntityColumns("datastream", updates, intID); err != nil {
		return nil, err
	}

	nd, _ := gdb.GetDatastream(intID, nil)
	return nd, nil
}

// PutDatastream receives a Datastream entity and changes it in the database
// returns the adapted Datastream
func (gdb *GostDatabase) PutDatastream(id interface{}, datastream *entities.Datastream) (*entities.Datastream, error) {
	return gdb.PatchDatastream(id, datastream)
}

// DeleteDatastream tries to delete a Datastream by the given id
func (gdb *GostDatabase) DeleteDatastream(id interface{}) error {
	return DeleteEntity(gdb, id, "datastream")
}

// DatastreamExists checks if a Datastream is present in the database based on a given id
func (gdb *GostDatabase) DatastreamExists(id int) bool {
	return EntityExists(gdb, id, "datastream")
}
