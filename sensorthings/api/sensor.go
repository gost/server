package api

import (
	"errors"
	gostErrors "github.com/geodan/gost/errors"
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
)

// GetSensor todo
func (a *APIv1) GetSensor(id string, qo *odata.QueryOptions) (*entities.Sensor, error) {
	s, err := a.db.GetSensor(id)
	if err != nil {
		return nil, err
	}

	s.SetLinks(a.config.GetExternalServerURI())
	return s, nil
}

// GetSensors todo
func (a *APIv1) GetSensors(qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	sensors, err := a.db.GetSensors()
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

// PostSensor todo
func (a *APIv1) PostSensor(sensor entities.Sensor) (*entities.Sensor, []error) {
	_, err := sensor.ContainsMandatoryParams()
	if err != nil {
		return nil, err
	}

	ns, err2 := a.db.PostSensor(sensor)
	if err2 != nil {
		return nil, []error{err2}
	}

	return ns, nil
}

// PatchSensor todo
func (a *APIv1) PatchSensor(id string, sensor entities.Sensor) (*entities.Sensor, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteSensor todo
func (a *APIv1) DeleteSensor(id string) error {
	return gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}
