package api

import (
	"errors"

	gostErrors "github.com/gost/server/errors"
	entities "github.com/gost/core"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/odata"
)

// GetDatastream retrieves a sensor by id and given query
func (a *APIv1) GetDatastream(id interface{}, qo *odata.QueryOptions, path string) (*entities.Datastream, error) {
	ds, err := a.db.GetDatastream(id, qo)
	if err != nil {
		return nil, err
	}

	a.SetLinks(ds, qo)
	return ds, nil
}

// GetDatastreams retrieves an array of sensors based on the given query
func (a *APIv1) GetDatastreams(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	datastreams, count, err := a.db.GetDatastreams(qo)
	return processDatastreams(a, datastreams, qo, path, count, err)
}

// GetDatastreamsByThing returns all datastreams linked to the given thing
func (a *APIv1) GetDatastreamsByThing(thingID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	datastreams, count, err := a.db.GetDatastreamsByThing(thingID, qo)
	return processDatastreams(a, datastreams, qo, path, count, err)
}

// GetDatastreamByObservation returns a datastream linked to the given observation
func (a *APIv1) GetDatastreamByObservation(observationID interface{}, qo *odata.QueryOptions, path string) (*entities.Datastream, error) {
	ds, err := a.db.GetDatastreamByObservation(observationID, qo)
	if err != nil {
		return nil, err
	}

	a.SetLinks(ds, qo)
	return ds, nil
}

// GetDatastreamsBySensor returns all datastreams linked to the given sensor
func (a *APIv1) GetDatastreamsBySensor(sensorID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	datastreams, count, err := a.db.GetDatastreamsBySensor(sensorID, qo)
	return processDatastreams(a, datastreams, qo, path, count, err)
}

// GetDatastreamsByObservedProperty returns all datastreams linked to the given ObservedProperty
func (a *APIv1) GetDatastreamsByObservedProperty(oID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	datastreams, count, err := a.db.GetDatastreamsByObservedProperty(oID, qo)
	return processDatastreams(a, datastreams, qo, path, count, err)
}

func processDatastreams(a *APIv1, datastreams []*entities.Datastream, qo *odata.QueryOptions, path string, count int, err error) (*models.ArrayResponse, error) {
	if err != nil {
		return nil, err
	}

	for idx, item := range datastreams {
		i := *item
		a.SetLinks(&i, qo)
		datastreams[idx] = &i
	}

	var data interface{} = datastreams
	return a.createArrayResponse(count, path, qo, data), nil
}

// PostDatastream adds a new datastream to the database
func (a *APIv1) PostDatastream(datastream *entities.Datastream) (*entities.Datastream, []error) {
	var errors []error
	var err error
	var postedObservedProperty *entities.ObservedProperty
	var postedSensor *entities.Sensor
	var postedObservations []*entities.Observation

	_, errors = containsMandatoryParams(datastream)
	if len(errors) > 0 {
		a.revertPostDatastream(postedObservedProperty, postedSensor, postedObservations)
		return nil, errors
	}

	// Check if ObservedProperty is deep inserted
	if datastream.ObservedProperty != nil && datastream.ObservedProperty.ID == nil {
		var op *entities.ObservedProperty
		if op, err = a.db.PostObservedProperty(datastream.ObservedProperty); err != nil {
			a.revertPostDatastream(postedObservedProperty, postedSensor, postedObservations)
			return nil, []error{err}
		}

		datastream.ObservedProperty = op
		postedObservedProperty = op
	}

	// Check if Sensor is deep inserted
	if datastream.Sensor != nil && datastream.Sensor.ID == nil {
		var s *entities.Sensor
		if s, err = a.db.PostSensor(datastream.Sensor); err != nil {
			a.revertPostDatastream(postedObservedProperty, postedSensor, postedObservations)
			return nil, []error{err}
		}

		datastream.Sensor = s
		postedSensor = s
	}

	ns, err := a.db.PostDatastream(datastream)
	if err != nil {
		a.revertPostDatastream(postedObservedProperty, postedSensor, postedObservations)
		return nil, []error{err}
	}

	// Check if Observations are deep inserted
	if datastream.Observations != nil {
		for _, observation := range datastream.Observations {
			ds := &entities.Datastream{}
			ds.ID = datastream.ID
			observation.Datastream = ds

			var o *entities.Observation
			if o, errors = a.PostObservation(observation); len(errors) > 0 {
				a.revertPostDatastream(postedObservedProperty, postedSensor, postedObservations)
				return nil, errors
			}

			postedObservations = append(postedObservations, o)
		}
	}

	ns.SetAllLinks(a.config.GetExternalServerURI())
	return ns, nil
}

func (a *APIv1) revertPostDatastream(op *entities.ObservedProperty, sensor *entities.Sensor, observations []*entities.Observation) {
	if op != nil {
		a.DeleteObservedProperty(op.ID)
	}

	if sensor != nil {
		a.DeleteSensor(sensor.ID)
	}

	for _, observation := range observations {
		a.DeleteObservation(observation.ID)
	}
}

// PostDatastreamByThing adds a new datastream by given thing ID
func (a *APIv1) PostDatastreamByThing(thingID interface{}, datastream *entities.Datastream) (*entities.Datastream, []error) {
	t := &entities.Thing{}
	t.ID = thingID
	datastream.Thing = t
	return a.PostDatastream(datastream)
}

// PatchDatastream updates the given datastream in the database
func (a *APIv1) PatchDatastream(id interface{}, datastream *entities.Datastream) (*entities.Datastream, error) {
	if datastream.Observations != nil || datastream.Sensor != nil || datastream.ObservedProperty != nil || datastream.Thing != nil {
		return nil, gostErrors.NewBadRequestError(errors.New("Deep patch datastream not supported."))
	}

	return a.db.PatchDatastream(id, datastream)
}

// PutDatastream updates the given thing in the database
func (a *APIv1) PutDatastream(id interface{}, datastream *entities.Datastream) (*entities.Datastream, []error) {
	var err2 error
	putdatastream, err2 := a.db.PutDatastream(id, datastream)
	if err2 != nil {
		return nil, []error{err2}
	}

	// putdatastream.SetAllLinks(a.config.GetExternalServerURI())
	return putdatastream, nil
}

// DeleteDatastream deletes a datastream from the database
func (a *APIv1) DeleteDatastream(id interface{}) error {
	return a.db.DeleteDatastream(id)
}
