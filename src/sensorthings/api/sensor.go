package api

import (
	"errors"
	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// GetSensor retrieves a sensor by id and given query
func (a *APIv1) GetSensor(id interface{}, qo *odata.QueryOptions) (*entities.Sensor, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Sensor{})
	if err != nil {
		return nil, err
	}

	s, err := a.db.GetSensor(id, qo)
	if err != nil {
		return nil, err
	}

	s.SetLinks(a.config.GetExternalServerURI())
	return s, nil
}

// GetSensorByDatastream retrieves a sensor by given datastream
func (a *APIv1) GetSensorByDatastream(id interface{}, qo *odata.QueryOptions) (*entities.Sensor, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Sensor{})
	if err != nil {
		return nil, err
	}

	s, err := a.db.GetSensorByDatastream(id, qo)
	if err != nil {
		return nil, err
	}

	s.SetLinks(a.config.GetExternalServerURI())
	return s, nil
}

// GetSensors retrieves an array of sensors based on the given query
func (a *APIv1) GetSensors(qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Sensor{})
	if err != nil {
		return nil, err
	}

	sensors, err := a.db.GetSensors(qo)
	if err != nil {
		return nil, err
	}

	uri := a.config.GetExternalServerURI()
	for idx, item := range sensors {
		i := *item
		i.SetLinks(uri)
		sensors[idx] = &i
	}

	var data interface{} = sensors
	return &models.ArrayResponse{
		Count: len(sensors),
		Data:  &data,
	}, nil
}

// PostSensor adds a new sensor to the database
func (a *APIv1) PostSensor(sensor *entities.Sensor) (*entities.Sensor, []error) {
	_, err := sensor.ContainsMandatoryParams()
	if err != nil {
		return nil, err
	}

	supported, err2 := entities.CheckEncodingSupported(sensor, sensor.EncodingType)
	if !supported || err2 != nil {
		return nil, []error{err2}
	}

	ns, err2 := a.db.PostSensor(sensor)
	if err2 != nil {
		return nil, []error{err2}
	}

	ns.SetLinks(a.config.GetExternalServerURI())

	return ns, nil
}

// PatchSensor updates a sensor in the database
func (a *APIv1) PatchSensor(id interface{}, sensor *entities.Sensor) (*entities.Sensor, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteSensor deletes a sensor from the database by given sensor id
func (a *APIv1) DeleteSensor(id interface{}) error {
	return a.db.DeleteSensor(id)
}
