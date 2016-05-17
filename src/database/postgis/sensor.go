package postgis

import (
	"fmt"
	"github.com/geodan/gost/src/sensorthings/entities"
	"strconv"

	gostErrors "github.com/geodan/gost/src/errors"
)

// GetSensor todo
func (gdb *GostDatabase) GetSensor(id string) (*entities.Sensor, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var sensorID int
	var description, metadata string
	sql := fmt.Sprintf("select id, description, metadata from %s.sensor where id = $1", gdb.Schema)
	err = gdb.Db.QueryRow(sql, intID).Scan(&sensorID, &description, &metadata)

	if err != nil {
		return nil, gostErrors.NewRequestNotFound(fmt.Errorf("Sensors(%s) does not exist", id))
	}

	sensor := entities.Sensor{
		ID:          strconv.Itoa(sensorID),
		Description: description,
		Metadata:    metadata,
	}

	return &sensor, nil
}

// GetSensors todo
func (gdb *GostDatabase) GetSensors() ([]*entities.Sensor, error) {
	sql := fmt.Sprintf("select id, description, metadata FROM %s.sensor", gdb.Schema)
	rows, err := gdb.Db.Query(sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var sensors = []*entities.Sensor{}

	for rows.Next() {
		var id int
		var description, metadata string
		err = rows.Scan(&id, &description, &metadata)
		if err != nil {
			return nil, err
		}

		sensor := entities.Sensor{
			ID:          strconv.Itoa(id),
			Description: description,
			Metadata:    metadata,
		}
		sensors = append(sensors, &sensor)
	}

	return sensors, nil
}

// PostSensor todo
func (gdb *GostDatabase) PostSensor(sensor entities.Sensor) (*entities.Sensor, error) {
	var sensorID int
	sql := fmt.Sprintf("INSERT INTO %s.sensor (description, encodingtype, metadata) VALUES ($1, $2, $3) RETURNING id", gdb.Schema)
	err := gdb.Db.QueryRow(sql, sensor.Description, 1, sensor.Metadata).Scan(&sensorID)
	if err != nil {
		return nil, err
	}

	sensor.ID = strconv.Itoa(sensorID)
	return &sensor, nil
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
