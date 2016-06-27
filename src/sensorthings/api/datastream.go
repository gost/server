package api

import (
	"errors"
	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// GetDatastream retrieves a sensor by id and given query
func (a *APIv1) GetDatastream(id interface{}, qo *odata.QueryOptions, path string) (*entities.Datastream, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Datastream{})
	if err != nil {
		return nil, err
	}

	ds, err := a.db.GetDatastream(id, qo)
	if err != nil {
		return nil, err
	}

	a.ProcessGetRequest(ds, qo)
	return ds, nil
}

// GetDatastreams retrieves an array of sensors based on the given query
func (a *APIv1) GetDatastreams(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Datastream{})
	if err != nil {
		return nil, err
	}

	datastreams, err := a.db.GetDatastreams(qo)
	return processDatastreams(a, datastreams, qo, path, err)

}

// GetDatastreamsByThing returns all datastreams linked to the given thing
func (a *APIv1) GetDatastreamsByThing(thingID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Datastream{})
	if err != nil {
		return nil, err
	}

	datastreams, err := a.db.GetDatastreamsByThing(thingID, qo)
	return processDatastreams(a, datastreams, qo, path, err)
}

// GetDatastreamByObservation returns a datastream linked to the given observation
func (a *APIv1) GetDatastreamByObservation(observationID interface{}, qo *odata.QueryOptions, path string) (*entities.Datastream, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Datastream{})
	if err != nil {
		return nil, err
	}

	ds, err := a.db.GetDatastreamByObservation(observationID, qo)
	if err != nil {
		return nil, err
	}

	a.ProcessGetRequest(ds, qo)
	return ds, nil
}

// GetDatastreamsBySensor returns all datastreams linked to the given sensor
func (a *APIv1) GetDatastreamsBySensor(sensorID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Datastream{})
	if err != nil {
		return nil, err
	}

	datastreams, err := a.db.GetDatastreamsBySensor(sensorID, qo)
	return processDatastreams(a, datastreams, qo, path, err)
}

// GetDatastreamsByObservedProperty returns all datastreams linked to the given ObservedProperty
func (a *APIv1) GetDatastreamsByObservedProperty(oID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Datastream{})
	if err != nil {
		return nil, err
	}

	datastreams, err := a.db.GetDatastreamsByObservedProperty(oID, qo)
	return processDatastreams(a, datastreams, qo, path, err)
}

func processDatastreams(a *APIv1, datastreams []*entities.Datastream, qo *odata.QueryOptions, path string, err error) (*models.ArrayResponse, error) {
	if err != nil {
		return nil, err
	}

	for idx, item := range datastreams {
		i := *item
		a.ProcessGetRequest(&i, qo)
		datastreams[idx] = &i
	}

	var data interface{} = datastreams
	return &models.ArrayResponse{
		Count:    a.db.GetTotalDatastreams(),
		NextLink: a.CreateNextLink(a.db.GetTotalDatastreams(), path, qo),
		Data:     &data,
	}, nil
}

// PostDatastream todo
func (a *APIv1) PostDatastream(datastream *entities.Datastream) (*entities.Datastream, []error) {
	var postedObservedProperty *entities.ObservedProperty
	var postedSensor *entities.Sensor
	postedObservations := make([]*entities.Observation, 0)

	// Check if ObservedProperty is deep inserted
	if datastream.ObservedProperty != nil && datastream.ObservedProperty.ID == nil {
		if op, err := a.db.PostObservedProperty(datastream.ObservedProperty); err != nil {
			a.revertPostDatastream(postedObservedProperty, postedSensor, postedObservations)
			return nil, []error{err}
		} else {
			datastream.ObservedProperty = op
			postedObservedProperty = op
		}
	}

	// Check if Sensor is deep inserted
	if datastream.Sensor != nil && datastream.Sensor.ID == nil {
		if s, err := a.db.PostSensor(datastream.Sensor); err != nil {
			a.revertPostDatastream(postedObservedProperty, postedSensor, postedObservations)
			return nil, []error{err}
		} else {
			datastream.Sensor = s
			postedSensor = s
		}
	}

	_, err := datastream.ContainsMandatoryParams()
	if err != nil {
		a.revertPostDatastream(postedObservedProperty, postedSensor, postedObservations)
		return nil, err
	}

	ns, err2 := a.db.PostDatastream(datastream)
	if err2 != nil {
		a.revertPostDatastream(postedObservedProperty, postedSensor, postedObservations)
		return nil, []error{err2}
	}

	// Check if Observations are deep inserted
	if datastream.Observations != nil {
		for _, observation := range datastream.Observations {
			ds := &entities.Datastream{}
			ds.ID = datastream.ID
			observation.Datastream = ds

			if o, err := a.PostObservation(observation); err != nil {
				a.revertPostDatastream(postedObservedProperty, postedSensor, postedObservations)
				return nil, err
			} else {
				postedObservations = append(postedObservations, o)
			}
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

// PatchDatastream todo
func (a *APIv1) PatchDatastream(id interface{}, datastream *entities.Datastream) (*entities.Datastream, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteDatastream deletes a datastream from the database
func (a *APIv1) DeleteDatastream(id interface{}) error {
	return a.db.DeleteDatastream(id)
}
