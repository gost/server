package api

import (
	"errors"
	gostErrors "github.com/geodan/gost/errors"
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
)

// GetThing returns a thing entity based on the given id and QueryOptions
// returns an error when the entity cannot be found
func (a *APIv1) GetThing(id string, qo *odata.QueryOptions) (*entities.Thing, error) {
	t, err := a.db.GetThing(id)
	if err != nil {
		return nil, err
	}

	t.SetLinks(a.config.GetExternalServerURI())
	return t, nil
}

// GetThings returns an array of thing entities based on the QueryOptions
func (a *APIv1) GetThings(qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	things, err := a.db.GetThings()
	if err != nil {
		return nil, err
	}

	uri := a.config.GetExternalServerURI()
	for idx, item := range things {
		i := *item
		i.SetLinks(uri)
		things[idx] = &i
	}

	var data interface{} = things

	var count = len(things)
	response := models.ArrayResponse{
		Count: &count,
		Data:  &data,
	}

	return &response, nil
}

// PostThing checks if a posted thing entity is valid and adds it to the database
// a posted thing can also contain Locations and DataStreams
func (a *APIv1) PostThing(thing entities.Thing) (*entities.Thing, []error) {
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
			location := *l
			// New location posted
			if len(l.ID) == 0 { //Id is null so a new location is posted
				_, err3 := a.PostLocationByThing(nt.ID, location)
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
}

// DeleteThing todo
func (a *APIv1) DeleteThing(id string) error {
	return gostErrors.NewRequestNotImplemented(errors.New("delete thing not implemented yet"))
}

// PatchThing todo
func (a *APIv1) PatchThing(id string, thing entities.Thing) (*entities.Thing, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("patch thing not implemented yet"))
}
