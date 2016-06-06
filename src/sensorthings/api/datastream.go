package api

import (
	"errors"
	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// GetDatastream retrieves a sensor by id and given query
func (a *APIv1) GetDatastream(id string, qo *odata.QueryOptions) (*entities.Datastream, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Datastream{})
	if err != nil {
		return nil, err
	}

	ds, err := a.db.GetDatastream(id, qo)
	if err != nil {
		return nil, err
	}

	ds.SetLinks(a.config.GetExternalServerURI())
	return ds, nil
}

// GetDatastreams retrieves an array of sensors based on the given query
func (a *APIv1) GetDatastreams(qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Datastream{})
	if err != nil {
		return nil, err
	}

	datastreams, err := a.db.GetDatastreams(qo)
	return processDatastreams(a, datastreams, err)

}

// GetDatastreamsByThing returns all datastreams linked to the given thing
func (a *APIv1) GetDatastreamsByThing(thingID string, qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Datastream{})
	if err != nil {
		return nil, err
	}

	datastreams, err := a.db.GetDatastreamsByThing(thingID, qo)
	return processDatastreams(a, datastreams, err)
}

// GetDatastreamByObservation returns a datastream linked to the given observation
func (a *APIv1) GetDatastreamByObservation(observationID string, qo *odata.QueryOptions) (*entities.Datastream, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Datastream{})
	if err != nil {
		return nil, err
	}

	ds, err := a.db.GetDatastreamByObservation(observationID, qo)
	if err != nil {
		return nil, err
	}

	ds.SetLinks(a.config.GetExternalServerURI())
	return ds, nil
}

// GetDatastreamsBySensor returns all datastreams linked to the given sensor
func (a *APIv1) GetDatastreamsBySensor(sensorID string, qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Datastream{})
	if err != nil {
		return nil, err
	}

	datastreams, err := a.db.GetDatastreamsBySensor(sensorID, qo)
	return processDatastreams(a, datastreams, err)
}

// GetDatastreamsByObservedProperty returns all datastreams linked to the given ObservedProperty
func (a *APIv1) GetDatastreamsByObservedProperty(oID string, qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Datastream{})
	if err != nil {
		return nil, err
	}

	datastreams, err := a.db.GetDatastreamsByObservedProperty(oID, qo)
	return processDatastreams(a, datastreams, err)
}

func processDatastreams(a *APIv1, datastreams []*entities.Datastream, err error) (*models.ArrayResponse, error) {
	if err != nil {
		return nil, err
	}

	uri := a.config.GetExternalServerURI()
	for idx, item := range datastreams {
		i := *item
		i.SetLinks(uri)
		datastreams[idx] = &i
	}

	var data interface{} = datastreams
	return &models.ArrayResponse{
		Count: len(datastreams),
		Data:  &data,
	}, nil
}

// PostDatastream todo
func (a *APIv1) PostDatastream(datastream *entities.Datastream) (*entities.Datastream, []error) {
	/*
		_, err := thing.ContainsMandatoryParams()
		if err != nil {
			return nil, err
		}

		nt, err2 := a.db.PostThing(thing)
		if err2 != nil {
			return nil, []error{err2}
		}

		// Handle locations
		if thing.Locations != nil {
			for _, l := range thing.Locations {
				// New location posted
				if len(l.ID) == 0 { //Id is null so a new location is posted
					_, err3 := a.PostLocationByThing(nt.ID, l)
					if err3 != nil {
						return nil, err3
					}
				} else { // posted id: link
					err4 := a.LinkLocation(nt.ID, l.ID)
					if err4 != nil {
						// todo: thing is posted, delete it
						return nil, []error{err4}
					}

					err5 := a.PostHistoricalLocation(nt.ID, l.ID)
					if err5 != nil {
						// todo: things is posted, delete it
						return nil, err5
					}
				}
			}
		}

		nt.SetLinks(a.config.GetExternalServerURI())
		//push to mqtt
		return nt, nil
	*/

	_, err := datastream.ContainsMandatoryParams()
	if err != nil {
		return nil, err
	}

	ns, err2 := a.db.PostDatastream(datastream)
	if err2 != nil {
		return nil, []error{err2}
	}

	ns.SetLinks(a.config.GetExternalServerURI())

	return ns, nil
}

// PostDatastreamByThing todo
func (a *APIv1) PostDatastreamByThing(thingID string, datastream *entities.Datastream) (*entities.Datastream, []error) {
	t := &entities.Thing{}
	t.ID = thingID
	datastream.Thing = t
	return a.PostDatastream(datastream)
}

// PatchDatastream todo
func (a *APIv1) PatchDatastream(id string, datastream *entities.Datastream) (*entities.Datastream, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteDatastream deletes a datastream from the database
func (a *APIv1) DeleteDatastream(id string) error {
	return a.db.DeleteDatastream(id)
}
