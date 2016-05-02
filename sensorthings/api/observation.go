package api

import (
	"errors"
	gostErrors "github.com/geodan/gost/errors"
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
)

// GetObservation todo
func (a *APIv1) GetObservation(id string, qo *odata.QueryOptions) (*entities.Observation, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// GetObservations todo
func (a *APIv1) GetObservations(qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// PostObservation todo
func (a *APIv1) PostObservation(observation entities.Observation, x string) (*entities.Observation, []error) {
	return nil, []error{gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))}
}

// PatchObservation todo
func (a *APIv1) PatchObservation(id string, observation entities.Observation) (*entities.Observation, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteObservation todo
func (a *APIv1) DeleteObservation(id string) error {
	return gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}
