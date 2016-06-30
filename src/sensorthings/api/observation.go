package api

import (
	"encoding/json"
	"errors"
	"fmt"

	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// GetObservation returns an observation by id
func (a *APIv1) GetObservation(id interface{}, qo *odata.QueryOptions, path string) (*entities.Observation, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Observation{})
	if err != nil {
		return nil, err
	}

	o, err := a.db.GetObservation(id, qo)
	if err != nil {
		return nil, err
	}

	a.ProcessGetRequest(o, qo)
	return o, nil
}

// GetObservations return all observations by given QueryOptions
func (a *APIv1) GetObservations(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Observation{})
	if err != nil {
		return nil, err
	}

	observations, err := a.db.GetObservations(qo)
	return processObservations(a, observations, qo, path, err)
}

// GetObservationsByFeatureOfInterest returns all observation by given FeatureOfInterest end QueryOptions
func (a *APIv1) GetObservationsByFeatureOfInterest(foiID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Observation{})
	if err != nil {
		return nil, err
	}

	observations, err := a.db.GetObservationsByFeatureOfInterest(foiID, qo)
	return processObservations(a, observations, qo, path, err)
}

// GetObservationsByDatastream returns all observations by given Datastream and QueryOptions
func (a *APIv1) GetObservationsByDatastream(datastreamID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Observation{})
	if err != nil {
		return nil, err
	}

	observations, err := a.db.GetObservationsByDatastream(datastreamID, qo)
	return processObservations(a, observations, qo, path, err)
}

func processObservations(a *APIv1, observations []*entities.Observation, qo *odata.QueryOptions, path string, err error) (*models.ArrayResponse, error) {
	if err != nil {
		return nil, err
	}

	for idx, item := range observations {
		i := *item
		a.ProcessGetRequest(&i, qo)
		observations[idx] = &i
	}

	var numberOfObservations = len(observations)

	var data interface{} = observations
	return &models.ArrayResponse{
		Count:    numberOfObservations,
		NextLink: a.CreateNextLink(a.db.GetTotalObservations(), path, qo),
		Data:     &data,
	}, nil
}

// PostObservation checks for correctness of the observation and calls PostObservation on the database
func (a *APIv1) PostObservation(observation *entities.Observation) (*entities.Observation, []error) {
	_, err := observation.ContainsMandatoryParams()
	if err != nil {
		return nil, err
	}

	datastreamID := observation.Datastream.ID

	if observation.FeatureOfInterest == nil {
		foiID, err := a.foiRepository.GetFoiIDByDatastreamID(&a.db, toStringID(datastreamID))

		if err != nil {
			return nil, []error{gostErrors.NewBadRequestError(errors.New("Unable to link or create FeatureOfInterest for Observation. The linked Thing should have a location or supply (deep insert) a new FeatureOfInterest or link to a known FeatureOfInterest by id \"FeatureOfInterest\":{\"@iot.id\":30198}"))}
		}

		observation.FeatureOfInterest = &entities.FeatureOfInterest{}
		observation.FeatureOfInterest.ID = foiID
	} else if observation.FeatureOfInterest != nil && observation.FeatureOfInterest.ID == nil {
		if foi, err := a.PostFeatureOfInterest(observation.FeatureOfInterest); err != nil {
			return nil, []error{gostErrors.NewConflictRequestError(errors.New("Unable to create deep inserted FeatureOfInterest"))}
		} else {
			observation.FeatureOfInterest = foi
		}
	}

	no, err2 := a.db.PostObservation(observation)
	if err2 != nil {
		return nil, []error{err2}
	}

	no.SetAllLinks(a.config.GetExternalServerURI())

	json, _ := json.Marshal(no)
	s := string(json)

	//ToDo: MQTT TEST
	a.mqtt.Publish(fmt.Sprintf("Datastreams(%v)/Observations", datastreamID), s, 0)
	a.mqtt.Publish("Observations", s, 0)

	return no, nil
}

// PostObservationByDatastream creates an Observation with a linked datastream by given datastream id and calls PostObservation on the database
func (a *APIv1) PostObservationByDatastream(datastreamID interface{}, observation *entities.Observation) (*entities.Observation, []error) {
	d := &entities.Datastream{}
	d.ID = datastreamID
	observation.Datastream = d
	return a.PostObservation(observation)
}

// PatchObservation todo
func (a *APIv1) PatchObservation(id interface{}, observation *entities.Observation) (*entities.Observation, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteObservation deletes a given Observation from the database
func (a *APIv1) DeleteObservation(id interface{}) error {
	return a.db.DeleteObservation(id)
}
