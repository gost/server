package postgis

import (
	"database/sql"
	"errors"
	"fmt"

	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/odata"
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
func (gdb *GostDatabase) GetSensors(qo *odata.QueryOptions) ([]*entities.Sensor, error) {
	sql := fmt.Sprintf("select "+CreateSelectString(&entities.Sensor{}, qo, "", "", nil)+" FROM %s.sensor order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema)
	return processSensors(gdb.Db, sql, qo)
}

func processSensor(db *sql.DB, sql string, qo *odata.QueryOptions) (*entities.Sensor, error) {
	sensors, err := processSensors(db, sql, qo)
	if err != nil {
		return nil, err
	}

	if len(sensors) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Sensor not found"))
	}

	return sensors[0], nil
}

func processSensors(db *sql.DB, sql string, qo *odata.QueryOptions) ([]*entities.Sensor, error) {
	rows, err := db.Query(sql)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var sensors = []*entities.Sensor{}

	for rows.Next() {
		var id interface{}
		var encodingType int
		var description, metadata string

		var params []interface{}
		var qp []string
		if qo == nil || qo.QuerySelect == nil || len(qo.QuerySelect.Params) == 0 {
			s := &entities.Sensor{}
			qp = s.GetPropertyNames()
		} else {
			qp = qo.QuerySelect.Params
		}

		for _, p := range qp {
			if p == "id" {
				params = append(params, &id)
			}
			if p == "encodingType" {
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
			return nil, err
		}

		sensor := entities.Sensor{}
		sensor.ID = id
		sensor.Description = description
		sensor.Metadata = metadata
		if encodingType != 0 {
			sensor.EncodingType = entities.EncodingValues[encodingType].Value
		}

		sensors = append(sensors, &sensor)
	}

	return sensors, nil
}

// PostSensor posts a sensor to the database
func (gdb *GostDatabase) PostSensor(sensor *entities.Sensor) (*entities.Sensor, error) {
	var sensorID int
	encoding, _ := entities.CreateEncodingType(sensor.EncodingType)
	sql := fmt.Sprintf("INSERT INTO %s.sensor (description, encodingtype, metadata) VALUES ($1, $2, $3) RETURNING id", gdb.Schema)
	err := gdb.Db.QueryRow(sql, sensor.Description, encoding.Code, sensor.Metadata).Scan(&sensorID)
	if err != nil {
		return nil, err
	}

	sensor.ID = sensorID
	return sensor, nil
}

// SensorExists checks if a sensor is present in the database based on a given id
func (gdb *GostDatabase) SensorExists(thingID int) bool {
	var result bool
	sql := fmt.Sprintf("SELECT exists (SELECT 1 FROM %s.sensor WHERE id = $1 LIMIT 1)", gdb.Schema)
	err := gdb.Db.QueryRow(sql, thingID).Scan(&result)
	if err != nil {
		return false
	}

	return result
}

// PatchSensor updates a sensor in the database
func (gdb *GostDatabase) PatchSensor(id interface{}, s *entities.Sensor) (*entities.Sensor, error) {
	var err error
	var ok bool
	var intID int

	if intID, ok = ToIntID(id); !ok || !gdb.SensorExists(intID) {
		return nil, gostErrors.NewRequestNotFound(errors.New("Sensor does not exist"))
	}

	if len(s.Description) > 0 {
		if err = gdb.updateEntityColumn("sensor", "description", s.Description, intID); err != nil {
			return nil, err
		}
	}

	if len(s.Metadata) > 0 {
		if err = gdb.updateEntityColumn("sensor", "metadata", s.Metadata, intID); err != nil {
			return nil, err
		}
	}

	if len(s.EncodingType) > 0 {
		encoding, _ := entities.CreateEncodingType(s.EncodingType)
		if err = gdb.updateEntityColumn("sensor", "encodingtype", encoding.Code, intID); err != nil {
			return nil, err
		}
	}

	ns, _ := gdb.GetSensor(intID, nil)
	return ns, nil
}

// DeleteSensor tries to delete a Sensor by the given id
func (gdb *GostDatabase) DeleteSensor(id interface{}) error {
	intID, ok := ToIntID(id)
	if !ok {
		return gostErrors.NewRequestNotFound(errors.New("Sensor does not exist"))
	}

	r, err := gdb.Db.Exec(fmt.Sprintf("DELETE FROM %s.sensor WHERE id = $1", gdb.Schema), intID)
	if err != nil {
		return err
	}

	if c, _ := r.RowsAffected(); c == 0 {
		return gostErrors.NewRequestNotFound(errors.New("Sensor not found"))
	}

	return nil
}
