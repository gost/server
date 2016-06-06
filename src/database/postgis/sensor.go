package postgis

import (
	"fmt"
	"github.com/geodan/gost/src/sensorthings/entities"
	"strconv"

	"database/sql"
	"errors"
	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// GetSensor todo
func (gdb *GostDatabase) GetSensor(id string, qo *odata.QueryOptions) (*entities.Sensor, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("select id, description, encodingtype, metadata from %s.sensor where id = $1", gdb.Schema)
	sensor, err := processSensor(gdb.Db, sql, intID)
	if err != nil {
		return nil, err
	}

	return sensor, nil
}

// GetSensorByDatastream todo
func (gdb *GostDatabase) GetSensorByDatastream(id string, qo *odata.QueryOptions) (*entities.Sensor, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("select id, description, encodingtype, metadata from %s.sensor inner join %s.datastream on datastream.sensor_id = sensor.id where datastream.id = $1", gdb.Schema, gdb.Schema)
	sensor, err := processSensor(gdb.Db, sql, intID)
	if err != nil {
		return nil, err
	}

	return sensor, nil
}

// GetSensors todo
func (gdb *GostDatabase) GetSensors(qo *odata.QueryOptions) ([]*entities.Sensor, error) {
	sql := fmt.Sprintf("select id, description, encodingtype, metadata FROM %s.sensor", gdb.Schema)
	return processSensors(gdb.Db, sql)
}

func processSensor(db *sql.DB, sql string, args ...interface{}) (*entities.Sensor, error) {
	sensors, err := processSensors(db, sql, args...)
	if err != nil {
		return nil, err
	}

	if len(sensors) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Sensor not found"))
	}

	return sensors[0], nil
}

func processSensors(db *sql.DB, sql string, args ...interface{}) ([]*entities.Sensor, error) {
	rows, err := db.Query(sql, args...)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var sensors = []*entities.Sensor{}

	for rows.Next() {
		var id, encodingtype int
		var description, metadata string
		err = rows.Scan(&id, &description, &encodingtype, &metadata)
		if err != nil {
			return nil, err
		}

		sensor := entities.Sensor{}
		sensor.ID = strconv.Itoa(id)
		sensor.Description = description
		sensor.Metadata = metadata
		sensor.EncodingType = entities.EncodingValues[encodingtype].Value

		sensors = append(sensors, &sensor)
	}

	return sensors, nil
}

// PostSensor todo
func (gdb *GostDatabase) PostSensor(sensor *entities.Sensor) (*entities.Sensor, error) {
	var sensorID int
	encoding, _ := entities.CreateEncodingType(sensor.EncodingType)
	sql := fmt.Sprintf("INSERT INTO %s.sensor (description, encodingtype, metadata) VALUES ($1, $2, $3) RETURNING id", gdb.Schema)
	err := gdb.Db.QueryRow(sql, sensor.Description, encoding.Code, sensor.Metadata).Scan(&sensorID)
	if err != nil {
		return nil, err
	}

	sensor.ID = strconv.Itoa(sensorID)
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

// DeleteSensor tries to delete a Sensor by the given id
func (gdb *GostDatabase) DeleteSensor(id string) error {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
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
