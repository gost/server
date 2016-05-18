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
	ds, err := a.db.GetDatastream(id)
	if err != nil {
		return nil, err
	}

	ds.SetLinks(a.config.GetExternalServerURI())
	return ds, nil
}

// GetDatastreams retrieves an array of sensors based on the given query
func (a *APIv1) GetDatastreams(qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	datastreams, err := a.db.GetDatastreams()
	return processDatastreams(a, datastreams, err)

}

// GetDatastreamsByThing returns all datastreams linked to the given thing
func (a *APIv1) GetDatastreamsByThing(thingID string, qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	datastreams, err := a.db.GetDatastreamsByThing(thingID)
	return processDatastreams(a, datastreams, err)
}

// GetDatastreamsBySensor todo
func (a *APIv1) GetDatastreamsBySensor(thingID string, qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
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
	_, err := datastream.ContainsMandatoryParams()
	if err != nil {
		return nil, err
	}

	ns, err2 := a.db.PostDatastream(datastream)
	if err2 != nil {
		return nil, []error{err2}
	}

	return ns, nil
}

// PostDatastreamByThing todo
func (a *APIv1) PostDatastreamByThing(thingID string, datastream *entities.Datastream) (*entities.Datastream, []error) {
	datastream.Thing = &entities.Thing{ID: thingID}
	return a.PostDatastream(datastream)
}

// PatchDatastream todo
func (a *APIv1) PatchDatastream(id string, datastream *entities.Datastream) (*entities.Datastream, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteDatastream todo
func (a *APIv1) DeleteDatastream(id string) error {
	return gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}
