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
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// GetSensors todo
func (a *APIv1) GetSensors(qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// PostSensor todo
func (a *APIv1) PostSensor(sensor entities.Sensor) (*entities.Sensor, []error) {
	return nil, []error{gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))}
}

// PatchSensor todo
func (a *APIv1) PatchSensor(id string, sensor entities.Sensor) (*entities.Sensor, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteSensor todo
func (a *APIv1) DeleteSensor(id string) error {
	return gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}
