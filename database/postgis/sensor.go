package postgis

import (
	"fmt"
	"github.com/geodan/gost/sensorthings/entities"
	"strconv"
)

// GetSensor todo
func (gdb *GostDatabase) GetSensor(id string) (*entities.Sensor, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var sensorID int
	var description string
	var metadata string
	sql := fmt.Sprintf("select id, description, metadata from %s.sensor where id = $1", gdb.Schema)
	err2 := gdb.Db.QueryRow(sql, intID).Scan(&sensorID, &description, &metadata)

	if err2 != nil {
		return nil, err
	}

	sensor := entities.Sensor{}
	sensor.ID = strconv.Itoa(sensorID)
	sensor.Description = description
	sensor.Metadata = metadata

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
		sensor := entities.Sensor{}

		var id int
		var description string
		var metadata string
		err2 := rows.Scan(&id, &description, &metadata)
		if err2 != nil {
			return nil, err2
		}

		sensor.ID = strconv.Itoa(id)
		sensor.Description = description
		sensor.Metadata = metadata

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
