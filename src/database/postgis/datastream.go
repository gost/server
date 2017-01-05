package postgis

import (
	"encoding/json"
	"errors"
	"fmt"

	"database/sql"

	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/odata"
	"strings"
)

var dsMapping = map[string]string{"observedArea": "public.ST_AsGeoJSON(datastream.observedarea) AS observedarea"}

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
			ds.ResultTime = value.(string)
		} else if as == asMappings[entities.EntityTypeDatastream][datastreamObservationType] {
			obs, _ := entities.GetObservationTypeByID(value.(int64))
			ds.ObservationType = obs.Value
		} else if as == asMappings[entities.EntityTypeDatastream][datastreamPhenomenonTime] {
			ds.PhenomenonTime = value.(string)
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
	sql := fmt.Sprintf(sqlString, gdb.Schema, gdb.Schema, id)
	rows, err := gdb.Db.Query(sql)
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

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Datastream{}, qo, "", "", dsMapping)+" FROM %s.datastream where id = %v", gdb.Schema, intID)
	datastream, err := processDatastream(gdb.Db, sql, qo)
	if err != nil {
		return nil, err
	}

	if qo != nil {
		hasSelectQuery := (qo.QuerySelect != nil)
		var containsObservedArea = true
		if hasSelectQuery {
			containsObservedArea = ContainsToLower(qo.QuerySelect.Params, "observedArea")
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
	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Datastream{}, qo, "", "", dsMapping)+" FROM %s.datastream order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema)
	countSQL := fmt.Sprintf("select COUNT(*) FROM %s.datastream", gdb.Schema)
	return processDatastreams(gdb.Db, sql, qo, countSQL)
}

// expand

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
func (gdb *GostDatabase) GetDatastreamsByThing(thingID interface{}, qo *odata.QueryOptions) ([]*entities.Datastream, int, error) {
	intID, ok := ToIntID(thingID)
	if !ok {
		return nil, 0, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Datastream{}, qo, "datastream.", "", dsMapping)+" FROM %s.datastream inner join %s.thing on thing.id = datastream.thing_id where thing.id = %v order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema, gdb.Schema, intID)
	countSQL := fmt.Sprintf("select COUNT(*) FROM %s.datastream inner join %s.thing on thing.id = datastream.thing_id where thing.id = %v", gdb.Schema, gdb.Schema, intID)
	return processDatastreams(gdb.Db, sql, qo, countSQL)
}

// GetDatastreamsBySensor retrieves all datastreams linked to the given sensor
func (gdb *GostDatabase) GetDatastreamsBySensor(sensorID interface{}, qo *odata.QueryOptions) ([]*entities.Datastream, int, error) {
	intID, ok := ToIntID(sensorID)
	if !ok {
		return nil, 0, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Datastream{}, qo, "datastream.", "", dsMapping)+" FROM %s.datastream inner join %s.sensor on sensor.id = datastream.sensor_id where sensor.id = %v order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema, gdb.Schema, intID)
	countSQL := fmt.Sprintf("select COUNT(*) FROM %s.datastream inner join %s.sensor on sensor.id = datastream.sensor_id where sensor.id = %v", gdb.Schema, gdb.Schema, intID)
	return processDatastreams(gdb.Db, sql, qo, countSQL)
}

// GetDatastreamsByObservedProperty retrieves all datastreams linked to the given ObservedProerty
func (gdb *GostDatabase) GetDatastreamsByObservedProperty(oID interface{}, qo *odata.QueryOptions) ([]*entities.Datastream, int, error) {
	intID, ok := ToIntID(oID)
	if !ok {
		return nil, 0, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Datastream{}, qo, "datastream.", "", dsMapping)+" FROM %s.datastream inner join %s.observedproperty on observedproperty.id = datastream.observedproperty_id where observedproperty.id = %v order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema, gdb.Schema, intID)
	CountSQL := fmt.Sprintf("select COUNT(*) FROM %s.datastream inner join %s.observedproperty on observedproperty.id = datastream.observedproperty_id where observedproperty.id = %v", gdb.Schema, gdb.Schema, intID)
	return processDatastreams(gdb.Db, sql, qo, CountSQL)
}

func processDatastream(db *sql.DB, sql string, qo *odata.QueryOptions) (*entities.Datastream, error) {
	datastreams, _, err := processDatastreams(db, sql, qo, "")
	if err != nil {
		return nil, err
	}

	if len(datastreams) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	return datastreams[0], nil
}

func processDatastreams(db *sql.DB, sql string, qo *odata.QueryOptions, countSQL string) ([]*entities.Datastream, int, error) {
	rows, err := db.Query(sql)
	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}

	var datastreams = []*entities.Datastream{}
	for rows.Next() {
		var id interface{}
		var name, description, unitofmeasurement string
		var observedarea *string
		var phenomenonTime string
		var resultTime string
		var ot int64

		var params []interface{}
		var qp []string
		if qo == nil || qo.QuerySelect == nil || len(qo.QuerySelect.Params) == 0 {
			d := &entities.Datastream{}
			qp = d.GetPropertyNames()
		} else {
			qp = qo.QuerySelect.Params
		}

		for _, p := range qp {
			p = strings.ToLower(p)
			if p == "id" {
				params = append(params, &id)
			}
			if p == "name" {
				params = append(params, &name)
			}
			if p == "description" {
				params = append(params, &description)
			}
			if p == "unitofmeasurement" {
				params = append(params, &unitofmeasurement)
			}
			if p == "observationtype" {
				params = append(params, &ot)
			}
			if p == "observedarea" {
				params = append(params, &observedarea)
			}
			if p == "phenomenontime" {
				params = append(params, &phenomenonTime)
			}
			if p == "resulttime" {
				params = append(params, &resultTime)
			}
		}

		err = rows.Scan(params...)

		unitOfMeasurementMap, err := JSONToMap(&unitofmeasurement)
		if err != nil {
			return nil, 0, err
		}

		observedAreaMap, err := JSONToMap(observedarea)
		if err != nil {
			return nil, 0, err
		}

		datastream := entities.Datastream{}
		datastream.ID = id
		datastream.Name = name
		datastream.Description = description
		datastream.UnitOfMeasurement = unitOfMeasurementMap
		datastream.PhenomenonTime = phenomenonTime
		datastream.ResultTime = resultTime
		datastream.ObservedArea = observedAreaMap
		if ot != 0 {
			obs, _ := entities.GetObservationTypeByID(ot)
			datastream.ObservationType = obs.Value
		}

		datastreams = append(datastreams, &datastream)
	}

	var count int
	if len(countSQL) > 0 {
		db.QueryRow(countSQL).Scan(&count)
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

// PostDatastream todo
// TODO: !!!!ADD phenomenonTime SUPPORT!!!!
// TODO: !!!!ADD resulttime SUPPORT!!!!
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
		period := ParseTMPeriod(d.PhenomenonTime)
		phenomenonTime = "'" + ToPostgresPeriodFormat(period) + "'"
	}
	// get the ObservationType id in the lookup table
	observationType, err := entities.GetObservationTypeByValue(d.ObservationType)

	if err != nil {
		return nil, gostErrors.NewBadRequestError(errors.New("ObservationType does not exist"))
	}

	sql := fmt.Sprintf("INSERT INTO %s.datastream (name, description, unitofmeasurement, observedarea, thing_id, sensor_id, observedproperty_id, observationtype, phenomenonTime) VALUES ($1, $2, $3, %s, $4, $5, $6, $7, %s) RETURNING id", gdb.Schema, geom, phenomenonTime)
	err = gdb.Db.QueryRow(sql, d.Name, d.Description, unitOfMeasurement, tID, sID, oID, observationType.Code).Scan(&dsID)
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

	/*err := CheckDatastreamRelationsExist(gdb, datastream)
	if err != nil {
		return nil, err
	}
	tID, _ := ToIntID(datastream.Thing.ID)
	sID, _ := ToIntID(datastream.Sensor.ID)
	oID, _ := ToIntID(datastream.ObservedProperty.ID)

	var intID int
	var ok bool
	if intID, ok = ToIntID(id); !ok || !gdb.DatastreamExists(intID) {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	unitOfMeasurement, _ := json.Marshal(datastream.UnitOfMeasurement)
	geom := "NULL"
	if len(datastream.ObservedArea) != 0 {
		observedAreaBytes, _ := json.Marshal(datastream.ObservedArea)
		geom = fmt.Sprintf("ST_SetSRID(ST_GeomFromGeoJSON('%s'),4326)", string(observedAreaBytes[:]))
	}

	// get the ObservationType id in the lookup table
	observationType, err := entities.GetObservationTypeByValue(datastream.ObservationType)

	if err != nil {
		return nil, gostErrors.NewBadRequestError(errors.New("ObservationType does not exist"))
	}

	sql := fmt.Sprintf("update %s.datastream set name=$1, description=$2, unitofmeasurement=$3, observedarea=%s, thing_id=$4, sensor_id=$5, observedproperty_id=$6, observationtype=$7 where id = $8", gdb.Schema, geom)
	_, err2 := gdb.Db.Exec(sql, datastream.Name, datastream.Description, unitOfMeasurement, tID, sID, oID, observationType.Code, intID)
	if err2 != nil {
		return nil, err
	}
	datastream.ID = intID

	// clear inner entities to serves links upon response
	datastream.Thing = nil
	datastream.Sensor = nil
	datastream.ObservedProperty = nil

	return datastream, nil*/
}

// DeleteDatastream tries to delete a Datastream by the given id
func (gdb *GostDatabase) DeleteDatastream(id interface{}) error {
	return DeleteEntity(gdb, id, "datastream")
}

// DatastreamExists checks if a Datastream is present in the database based on a given id
func (gdb *GostDatabase) DatastreamExists(id int) bool {
	return EntityExists(gdb, id, "datastream")
}
