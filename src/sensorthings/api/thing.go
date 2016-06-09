package api

import (
	"errors"
	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// GetThing returns a thing entity based on the given id and QueryOptions
// returns an error when the entity cannot be found
func (a *APIv1) GetThing(id interface{}, qo *odata.QueryOptions, path string) (*entities.Thing, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Thing{})
	if err != nil {
		return nil, err
	}

	t, err := a.db.GetThing(id, qo)
	if err != nil {
		return nil, err
	}

	a.ProcessGetRequest(t, qo)
	return t, nil
}

// GetThingByDatastream returns a thing entity based on the given datastream id and QueryOptions
func (a *APIv1) GetThingByDatastream(id interface{}, qo *odata.QueryOptions, path string) (*entities.Thing, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Thing{})
	if err != nil {
		return nil, err
	}

	t, err := a.db.GetThingByDatastream(id, qo)
	if err != nil {
		return nil, err
	}

	a.ProcessGetRequest(t, qo)
	return t, nil
}

// GetThingsByLocation returns things based on the given location id and QueryOptions
func (a *APIv1) GetThingsByLocation(id interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Thing{})
	if err != nil {
		return nil, err
	}
	things, err := a.db.GetThingsByLocation(id, qo)
	return processThings(a, things, qo, path, err)
}

// GetThingByHistoricalLocation returns a thing entity based on the given HistoricalLocation id and QueryOptions
func (a *APIv1) GetThingByHistoricalLocation(id interface{}, qo *odata.QueryOptions, path string) (*entities.Thing, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Thing{})
	if err != nil {
		return nil, err
	}
	t, err := a.db.GetThingByHistoricalLocation(id, qo)
	if err != nil {
		return nil, err
	}

	a.ProcessGetRequest(t, qo)
	return t, nil
}

// GetThings returns an array of thing entities based on the QueryOptions
func (a *APIv1) GetThings(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Thing{})
	if err != nil {
		return nil, err
	}
	things, err := a.db.GetThings(qo)
	return processThings(a, things, qo, path, err)
}

func processThings(a *APIv1, observations []*entities.Thing, qo *odata.QueryOptions, path string, err error) (*models.ArrayResponse, error) {
	if err != nil {
		return nil, err
	}

	for idx, item := range observations {
		i := *item
		a.ProcessGetRequest(&i, qo)
		observations[idx] = &i
	}

	var data interface{} = observations
	total := a.db.GetTotalThings()
	return &models.ArrayResponse{
		Count:    total,
		NextLink: a.CreateNextLink(total, path, qo),
		Data:     &data,
	}, nil
}

// PostThing checks if a posted thing entity is valid and adds it to the database
// a posted thing can also contain Locations and DataStreams
func (a *APIv1) PostThing(thing *entities.Thing) (*entities.Thing, []error) {
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
			if l.ID == nil { //Id is null so a new location is posted
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

	nt.SetAllLinks(a.config.GetExternalServerURI())
	//push to mqtt
	return nt, nil
}

// DeleteThing deletes a given Thing from the database
func (a *APIv1) DeleteThing(id interface{}) error {
	return a.db.DeleteThing(id)
}

// PatchThing todo
func (a *APIv1) PatchThing(id interface{}, thing *entities.Thing) (*entities.Thing, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("patch thing not implemented yet"))
}
