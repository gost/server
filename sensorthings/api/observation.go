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

// GetObservationsByDatastream todo
func (a *APIv1) GetObservationsByDatastream(datastreamID string, qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// PostObservation todo
func (a *APIv1) PostObservation(observation entities.Observation) (*entities.Observation, []error) {
	_, err := observation.ContainsMandatoryParams()
	if err != nil {
		return nil, err
	}

	//ToDo check for linked featureofinterest
	no, err2 := a.db.PostObservation(observation)
	if err2 != nil {
		return nil, []error{err2}
	}

	no.SetLinks(a.config.GetExternalServerURI())
	return no, nil
}

// PostObservationByDatastream creates a Datastream with given id for the Observation and calls PostObservation
func (a *APIv1) PostObservationByDatastream(datastreamID string, observation entities.Observation) (*entities.Observation, []error) {
	observation.Datastream = &entities.Datastream{ID: datastreamID}
	return a.PostObservation(observation)
}

// PatchObservation todo
func (a *APIv1) PatchObservation(id string, observation entities.Observation) (*entities.Observation, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteObservation todo
func (a *APIv1) DeleteObservation(id string) error {
	return gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}
