package api

import (
	"encoding/json"
	"errors"
	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// GetObservation returns an observation by id
func (a *APIv1) GetObservation(id string, qo *odata.QueryOptions) (*entities.Observation, error) {
	o, err := a.db.GetObservation(id, qo)
	if err != nil {
		return nil, err
	}

	o.SetLinks(a.config.GetExternalServerURI())
	return o, nil
}

// GetObservations return all observations
func (a *APIv1) GetObservations(qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	observations, err := a.db.GetObservations(qo)
	return processObservations(a, observations, err)
}

// GetObservationsByFeatureOfInterest todo
func (a *APIv1) GetObservationsByFeatureOfInterest(foiID string, qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	observations, err := a.db.GetObservationsByFeatureOfInterest(foiID, qo)
	return processObservations(a, observations, err)
}

// GetObservationsByDatastream todo
func (a *APIv1) GetObservationsByDatastream(datastreamID string, qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	observations, err := a.db.GetObservationsByDatastream(datastreamID, qo)
	return processObservations(a, observations, err)
}

func processObservations(a *APIv1, observations []*entities.Observation, err error) (*models.ArrayResponse, error) {
	if err != nil {
		return nil, err
	}

	uri := a.config.GetExternalServerURI()
	for idx, item := range observations {
		i := *item
		i.SetLinks(uri)
		observations[idx] = &i
	}

	var data interface{} = observations
	return &models.ArrayResponse{
		Count: len(observations),
		Data:  &data,
	}, nil
}

// PostObservation todo
func (a *APIv1) PostObservation(observation *entities.Observation) (*entities.Observation, []error) {
	_, err := observation.ContainsMandatoryParams()
	if err != nil {
		return nil, err
	}

	//ToDo check for linked featureofinterest -> POST
	no, err2 := a.db.PostObservation(observation)
	if err2 != nil {
		return nil, []error{err2}
	}

	no.SetLinks(a.config.GetExternalServerURI())

	json, _ := json.Marshal(no)
	s := string(json)

	//ToDo: TEST
	a.mqtt.Publish("Datastreams(1)/Observations", s, 0)
	a.mqtt.Publish("Observations", s, 0)

	return no, nil
}

// PostObservationByDatastream creates a Datastream with given id for the Observation and calls PostObservation
func (a *APIv1) PostObservationByDatastream(datastreamID string, observation *entities.Observation) (*entities.Observation, []error) {
	d := &entities.Datastream{}
	d.ID = datastreamID
	observation.Datastream = d
	return a.PostObservation(observation)
}

// PatchObservation todo
func (a *APIv1) PatchObservation(id string, observation *entities.Observation) (*entities.Observation, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteObservation deletes a given Observation from the database
func (a *APIv1) DeleteObservation(id string) error {
	return a.db.DeleteObservation(id)
}
