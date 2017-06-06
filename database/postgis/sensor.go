package postgis

import (
	"database/sql"
	"errors"
	"fmt"

	gostErrors "github.com/geodan/gost/errors"
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/odata"
)

func sensorParamFactory(values map[string]interface{}) (entities.Entity, error) {
	s := &entities.Sensor{}
	for as, value := range values {
		if value == nil {
			continue
		}

		if as == asMappings[entities.EntityTypeSensor][sensorID] {
			s.ID = value
		} else if as == asMappings[entities.EntityTypeSensor][sensorName] {
			s.Name = value.(string)
		} else if as == asMappings[entities.EntityTypeSensor][sensorDescription] {
			s.Description = value.(string)
		} else if as == asMappings[entities.EntityTypeSensor][sensorEncodingType] {
			encodingType := value.(int64)
			if encodingType != 0 {
				s.EncodingType = entities.EncodingValues[encodingType].Value
			}
		} else if as == asMappings[entities.EntityTypeSensor][sensorMetadata] {
			s.Metadata = value.(string)
		}
	}

	return s, nil
}

// GetSensor return a sensor by id
func (gdb *GostDatabase) GetSensor(id interface{}, qo *odata.QueryOptions) (*entities.Sensor, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Sensor does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Sensor{}, nil, intID, qo)
	sensor, err := processSensor(gdb.Db, query, qi)

	if err != nil {
		return nil, err
	}

	return sensor, nil
}

// GetSensorByDatastream retrieves a sensor by given datastream
func (gdb *GostDatabase) GetSensorByDatastream(id interface{}, qo *odata.QueryOptions) (*entities.Sensor, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Sensor{}, &entities.Datastream{}, intID, qo)
	sensor, err := processSensor(gdb.Db, query, qi)
	if err != nil {
		return nil, err
	}

	return sensor, nil
}

// GetSensors retrieves all sensors based on the QueryOptions
func (gdb *GostDatabase) GetSensors(qo *odata.QueryOptions) ([]*entities.Sensor, int, error) {
	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Sensor{}, nil, nil, qo)
	countSQL := gdb.QueryBuilder.CreateCountQuery(&entities.Sensor{}, nil, nil, qo)
	return processSensors(gdb.Db, query, qi, countSQL)
}

func processSensor(db *sql.DB, sql string, qi *QueryParseInfo) (*entities.Sensor, error) {
	sensors, _, err := processSensors(db, sql, qi, "")
	if err != nil {
		return nil, err
	}

	if len(sensors) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Sensor not found"))
	}

	return sensors[0], nil
}

func processSensors(db *sql.DB, sql string, qi *QueryParseInfo, countSQL string) ([]*entities.Sensor, int, error) {
	data, err := ExecuteSelect(db, qi, sql)
	if err != nil {
		return nil, 0, fmt.Errorf("Error executing query %v", err)
	}

	sensors := make([]*entities.Sensor, 0)
	for _, d := range data {
		entity := d.(*entities.Sensor)
		sensors = append(sensors, entity)
	}

	var count int
	if len(countSQL) > 0 {
		count, err = ExecuteSelectCount(db, countSQL)
		if err != nil {
			return nil, 0, fmt.Errorf("Error executing count %v", err)
		}
	}

	return sensors, count, nil
}

// PostSensor posts a sensor to the database
func (gdb *GostDatabase) PostSensor(sensor *entities.Sensor) (*entities.Sensor, error) {
	var sensorID int
	encoding, err1 := entities.CreateEncodingType(sensor.EncodingType)
	if err1 != nil {
		return nil, err1
	}

	sql := fmt.Sprintf("INSERT INTO %s.sensor (name, description, encodingtype, metadata) VALUES ($1, $2, $3, $4) RETURNING id", gdb.Schema)
	err2 := gdb.Db.QueryRow(sql, sensor.Name, sensor.Description, encoding.Code, sensor.Metadata).Scan(&sensorID)
	if err2 != nil {
		return nil, err2
	}

	sensor.ID = sensorID
	return sensor, nil
}

// SensorExists checks if a sensor is present in the database based on a given id
func (gdb *GostDatabase) SensorExists(id int) bool {
	return EntityExists(gdb, id, "sensor")
}

// PatchSensor updates a sensor in the database
func (gdb *GostDatabase) PatchSensor(id interface{}, s *entities.Sensor) (*entities.Sensor, error) {
	var err error
	var ok bool
	var intID int
	updates := make(map[string]interface{})

	if intID, ok = ToIntID(id); !ok || !gdb.SensorExists(intID) {
		return nil, gostErrors.NewRequestNotFound(errors.New("Sensor does not exist"))
	}

	if len(s.Name) > 0 {
		updates["name"] = s.Name
	}

	if len(s.Description) > 0 {
		updates["description"] = s.Description
	}

	if len(s.Metadata) > 0 {
		updates["metadata"] = s.Metadata
	}

	if len(s.EncodingType) > 0 {
		encoding, _ := entities.CreateEncodingType(s.EncodingType)
		updates["encodingtype"] = encoding.Code
	}

	if err = gdb.updateEntityColumns("sensor", updates, intID); err != nil {
		return nil, err
	}

	ns, _ := gdb.GetSensor(intID, nil)
	return ns, nil
}

// PutSensor receives a Sensor entity and changes it in the database
// returns the Sensor
func (gdb *GostDatabase) PutSensor(id interface{}, sensor *entities.Sensor) (*entities.Sensor, error) {
	return gdb.PatchSensor(id, sensor)
}

// DeleteSensor tries to delete a Sensor by the given id
func (gdb *GostDatabase) DeleteSensor(id interface{}) error {
	return DeleteEntity(gdb, id, "sensor")
}
