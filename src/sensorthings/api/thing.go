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

func processThings(a *APIv1, things []*entities.Thing, qo *odata.QueryOptions, path string, err error) (*models.ArrayResponse, error) {
	if err != nil {
		return nil, err
	}

	for idx, item := range things {
		i := *item
		a.ProcessGetRequest(&i, qo)
		things[idx] = &i
	}

	var data interface{} = things
	total := len(things)
	return &models.ArrayResponse{
		Count:    total,
		NextLink: a.CreateNextLink(total, path, qo),
		Data:     &data,
	}, nil
}

// PostThing checks if a posted thing entity is valid and adds it to the database
// a posted thing can also contain Locations and DataStreams
func (a *APIv1) PostThing(thing *entities.Thing) (*entities.Thing, []error) {
	var err []error
	var err2 error
	_, err = thing.ContainsMandatoryParams()
	if err != nil {
		return nil, err
	}

	nt, err2 := a.db.PostThing(thing)
	if err2 != nil {
		return nil, []error{err2}
	}

	postedLocations := make([]*entities.Location, 0)
	postedDatastreams := make([]*entities.Datastream, 0)

	// Handle deep insert locations
	if thing.Locations != nil {
		for _, l := range thing.Locations {
			// New location posted
			if l.ID == nil { //Id is null so a new location is posted
				if nl, err := a.PostLocationByThing(nt.ID, l); err != nil {
					a.reverseInserts(nt, postedLocations, postedDatastreams)
					err = append(err, gostErrors.NewConflictRequestError(errors.New("Location deep inserted went wrong")))
					return nil, err
				} else {
					postedLocations = append(postedLocations, nl)
				}
			} else { // posted id: link
				if err2 = a.LinkLocation(nt.ID, l.ID); err != nil {
					a.reverseInserts(nt, postedLocations, postedDatastreams)
					err = append(err, gostErrors.NewConflictRequestError(errors.New("Location linking went wrong")))
					err = append(err, err2)
					return nil, err
				}

				hl := &entities.HistoricalLocation{
					Thing:     nt,
					Locations: []*entities.Location{l},
				}

				hl.ContainsMandatoryParams()

				if hl, err = a.PostHistoricalLocation(hl); err != nil {
					a.reverseInserts(nt, postedLocations, postedDatastreams)
					err = append(err, gostErrors.NewConflictRequestError(errors.New("Creating Historical Location went wrong")))
					return nil, err
				}
			}
		}
	}

	// Handle deep insert datastreams
	if thing.Datastreams != nil {
		for _, d := range thing.Datastreams {
			// New location posted
			if d.ID == nil { //Id is null so a new datastream is posted
				if nd, err := a.PostDatastreamByThing(nt.ID, d); err != nil {
					a.reverseInserts(nt, postedLocations, postedDatastreams)
					err = append(err, gostErrors.NewConflictRequestError(errors.New("Creating Datastrean went wrong")))
					return nil, err
				} else {
					postedDatastreams = append(postedDatastreams, nd)
				}
			} else {
				a.reverseInserts(nt, postedLocations, postedDatastreams)
				err = append(err, gostErrors.NewConflictRequestError(errors.New("ID found for deep inserted datastream, linking to an existing Datastream is not allowed")))
				return nil, err
			}
		}
	}

	nt.SetAllLinks(a.config.GetExternalServerURI())
	//push to mqtt
	return nt, nil
}

func (a *APIv1) reverseInserts(thing *entities.Thing, locations []*entities.Location, datastreams []*entities.Datastream) {
	for _, datastream := range datastreams {
		a.DeleteDatastream(datastream.ID)
	}

	for _, location := range locations {
		a.DeleteLocation(location.ID)
	}

	a.DeleteThing(thing.ID)
}

// DeleteThing deletes a given Thing from the database
func (a *APIv1) DeleteThing(id interface{}) error {
	return a.db.DeleteThing(id)
}

// PatchThing todo
func (a *APIv1) PatchThing(id interface{}, thing *entities.Thing) (*entities.Thing, error) {
	nt, _ := a.db.PatchThing(id, thing)
	return nt, nil
}
