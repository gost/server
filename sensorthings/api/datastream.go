package api

import (
	"errors"
	gostErrors "github.com/geodan/gost/errors"
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
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

// GetDatastreamsByThing todo
func (a *APIv1) GetDatastreamsByThing(thingID string, qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// GetDatastreamsBySensor todo
func (a *APIv1) GetDatastreamsBySensor(thingID string, qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// PostDatastream todo
func (a *APIv1) PostDatastream(datastream entities.Datastream) (*entities.Datastream, []error) {
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
func (a *APIv1) PostDatastreamByThing(thingID string, datastream entities.Datastream) (*entities.Datastream, []error) {
	datastream.Thing = &entities.Thing{ID: thingID}
	return a.PostDatastream(datastream)
}

// PatchDatastream todo
func (a *APIv1) PatchDatastream(id string, datastream entities.Datastream) (*entities.Datastream, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteDatastream todo
func (a *APIv1) DeleteDatastream(id string) error {
	return gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}
