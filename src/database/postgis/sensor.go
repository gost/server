package postgis

import (
	"database/sql"
	"errors"
	"fmt"

	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/odata"
	"strings"
)

// GetTotalSensors returns the total sensors count in the database
func (gdb *GostDatabase) GetTotalSensors() int {
	var count int
	sql := fmt.Sprintf("SELECT Count(*) from %s.sensor", gdb.Schema)
	gdb.Db.QueryRow(sql).Scan(&count)
	return count
}

// GetSensor return a sensor by id
func (gdb *GostDatabase) GetSensor(id interface{}, qo *odata.QueryOptions) (*entities.Sensor, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Sensor does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Sensor{}, qo, "", "", nil)+" from %s.sensor where id = %v", gdb.Schema, intID)
	sensor, err := processSensor(gdb.Db, sql, qo)
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

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Sensor{}, qo, "sensor.", "", nil)+" from %s.sensor inner join %s.datastream on datastream.sensor_id = sensor.id where datastream.id = %v", gdb.Schema, gdb.Schema, intID)
	sensor, err := processSensor(gdb.Db, sql, qo)
	if err != nil {
		return nil, err
	}

	return sensor, nil
}

// GetSensors retrieves all sensors based on the QueryOptions
func (gdb *GostDatabase) GetSensors(qo *odata.QueryOptions) ([]*entities.Sensor, int, error) {
	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Sensor{}, qo, "", "", nil)+" FROM %s.sensor order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema)
	countSQL := fmt.Sprintf("select COUNT(*) FROM %s.sensor", gdb.Schema)
	return processSensors(gdb.Db, sql, qo, countSQL)
}

func processSensor(db *sql.DB, sql string, qo *odata.QueryOptions) (*entities.Sensor, error) {
	sensors, _, err := processSensors(db, sql, qo, "")
	if err != nil {
		return nil, err
	}

	if len(sensors) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Sensor not found"))
	}

	return sensors[0], nil
}

func processSensors(db *sql.DB, sql string, qo *odata.QueryOptions, countSQL string) ([]*entities.Sensor, int, error) {
	rows, err := db.Query(sql)
	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	var sensors = []*entities.Sensor{}

	for rows.Next() {
		var id interface{}
		var encodingType int
		var name, description, metadata string

		var params []interface{}
		var qp []string
		if qo == nil || qo.QuerySelect == nil || len(qo.QuerySelect.Params) == 0 {
			s := &entities.Sensor{}
			qp = s.GetPropertyNames()
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
			if p == "encodingtype" {
				params = append(params, &encodingType)
			}
			if p == "description" {
				params = append(params, &description)
			}
			if p == "metadata" {
				params = append(params, &metadata)
			}
		}

		err = rows.Scan(params...)
		if err != nil {
			return nil, 0, err
		}

		sensor := entities.Sensor{}
		sensor.ID = id
		sensor.Name = name
		sensor.Description = description
		sensor.Metadata = metadata
		if encodingType != 0 {
			sensor.EncodingType = entities.EncodingValues[encodingType].Value
		}

		sensors = append(sensors, &sensor)
	}

	var count int
	if len(countSQL) > 0 {
		db.QueryRow(countSQL).Scan(&count)
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
	var intID int
	var ok bool

	if intID, ok = ToIntID(id); !ok || !gdb.SensorExists(intID) {
		return nil, gostErrors.NewRequestNotFound(errors.New("Sensor does not exist"))
	}
	encoding, err1 := entities.CreateEncodingType(sensor.EncodingType)
	if err1 != nil {
		return nil, err1
	}

	// INSERT INTO %s.sensor (name, description, encodingtype, metadata) VALUES ($1, $2, $3, $4)
	sql := fmt.Sprintf("update %s.sensor set name=$1, description=$2,encodingtype=$3, metadata=$4 where id=$5", gdb.Schema)
	_, err := gdb.Db.Exec(sql, sensor.Name, sensor.Description, encoding.Code, sensor.Metadata, intID)
	if err != nil {
		return nil, err
	}

	ns, _ := gdb.GetSensor(intID, nil)
	return ns, nil
}

// DeleteSensor tries to delete a Sensor by the given id
func (gdb *GostDatabase) DeleteSensor(id interface{}) error {
	return DeleteEntity(gdb, id, "sensor")
}
