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

	observations, count, err := a.db.GetObservations(qo)
	return processObservations(a, observations, qo, path, count, err)
}

// GetObservationsByFeatureOfInterest returns all observation by given FeatureOfInterest end QueryOptions
func (a *APIv1) GetObservationsByFeatureOfInterest(foiID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Observation{})
	if err != nil {
		return nil, err
	}

	observations, count, err := a.db.GetObservationsByFeatureOfInterest(foiID, qo)
	return processObservations(a, observations, qo, path, count, err)
}

// GetObservationsByDatastream returns all observations by given Datastream and QueryOptions
func (a *APIv1) GetObservationsByDatastream(datastreamID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Observation{})
	if err != nil {
		return nil, err
	}

	observations, count, err := a.db.GetObservationsByDatastream(datastreamID, qo)
	return processObservations(a, observations, qo, path, count, err)
}

func processObservations(a *APIv1, observations []*entities.Observation, qo *odata.QueryOptions, path string, count int, err error) (*models.ArrayResponse, error) {
	if err != nil {
		return nil, err
	}

	for idx, item := range observations {
		i := *item
		a.ProcessGetRequest(&i, qo)
		observations[idx] = &i
	}

	var data interface{} = observations
	return &models.ArrayResponse{
		Count:    count,
		NextLink: a.CreateNextLink(count, path, qo),
		Data:     &data,
	}, nil
}

// ConvertLocationToFoi converts a location to FOI
func ConvertLocationToFoi(l *entities.Location) *entities.FeatureOfInterest {
	foi := &entities.FeatureOfInterest{}
	foi.Name = l.Name
	foi.Description = l.Description
	foi.EncodingType = l.EncodingType
	foi.Feature = l.Location
	foi.OriginalLocationID = l.ID
	return foi
}

// CopyLocationToFoi copies the location of the thing to the FeatureOfInterest table. If it already
// exist, returns only the existing FeatureOfInterest ID
func CopyLocationToFoi(gdb *models.Database, datastreamID interface{}) (string, error) {
	db := *gdb
	var result string
	var l *entities.Location
	var err error

	if l, err = db.GetLocationByDatastreamID(datastreamID, nil); err != nil || l == nil {
		return "", gostErrors.NewConflictRequestError(errors.New("No location found for datastream.Thing"))
	}

	var featureOfInterestID interface{}

	// now check if the locationid already exists in featureofinterest.orginal_location id
	featureOfInterestID, _ = db.GetFeatureOfInterestIDByLocationID(l.ID)
	if featureOfInterestID == nil {
		// if the FeatureOfInterest does not exist already, create it now
		NewFeatureOfInterest := ConvertLocationToFoi(l)
		CreatedFeatureOfInterest, err := db.PostFeatureOfInterest(NewFeatureOfInterest)
		if err != nil {
			return "", err
		}
		result = toStringID(CreatedFeatureOfInterest.ID)
	} else {
		result = toStringID(featureOfInterestID)
	}

	return result, nil
}

// PostObservation checks for correctness of the observation and calls PostObservation on the database
func (a *APIv1) PostObservation(observation *entities.Observation) (*entities.Observation, []error) {
	_, err := observation.ContainsMandatoryParams()
	if err != nil {
		return nil, err
	}

	datastreamID := observation.Datastream.ID

	// there is no foi posted: try to copy it from thing.location...
	if observation.FeatureOfInterest == nil {
		foiID, err := CopyLocationToFoi(&a.db, datastreamID)

		if err != nil {
			errorMessage := "Unable to copy location of thing to featureofinterest."
			return nil, []error{gostErrors.NewBadRequestError(errors.New(errorMessage))}
		}

		observation.FeatureOfInterest = &entities.FeatureOfInterest{}
		observation.FeatureOfInterest.ID = foiID
	} else if observation.FeatureOfInterest != nil && observation.FeatureOfInterest.ID == nil {
		var foi *entities.FeatureOfInterest
		if foi, err = a.PostFeatureOfInterest(observation.FeatureOfInterest); err != nil {
			return nil, []error{gostErrors.NewConflictRequestError(errors.New("Unable to create deep inserted FeatureOfInterest"))}
		}
		observation.FeatureOfInterest = foi
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

func toStringID(id interface{}) string {
	return fmt.Sprintf("%v", id)
}

// PostObservationByDatastream creates an Observation with a linked datastream by given datastream id and calls PostObservation on the database
func (a *APIv1) PostObservationByDatastream(datastreamID interface{}, observation *entities.Observation) (*entities.Observation, []error) {
	d := &entities.Datastream{}
	d.ID = datastreamID
	observation.Datastream = d
	return a.PostObservation(observation)
}

// PatchObservation updates the given observation in the database
func (a *APIv1) PatchObservation(id interface{}, observation *entities.Observation) (*entities.Observation, error) {
	if observation.Datastream != nil || observation.FeatureOfInterest != nil {
		return nil, gostErrors.NewBadRequestError(errors.New("Unable to deep patch Observation"))
	}

	return a.db.PatchObservation(id, observation)
}

// PutObservation updates the given observation in the database
func (a *APIv1) PutObservation(id interface{}, observation *entities.Observation) (*entities.Observation, []error) {
	obs, err2 := a.db.PutObservation(id, observation)
	if err2 != nil {
		return nil, []error{err2}
	}
	return obs, nil
}

// DeleteObservation deletes a given Observation from the database
func (a *APIv1) DeleteObservation(id interface{}) error {
	return a.db.DeleteObservation(id)
}
