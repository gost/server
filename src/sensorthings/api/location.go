package api

import (
	"errors"
	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// PostLocation todo
func (a *APIv1) PostLocation(location entities.Location) (*entities.Location, []error) {
	_, err := location.ContainsMandatoryParams()
	if err != nil {
		return nil, err
	}

	l, err2 := a.db.PostLocation(location)
	if err2 != nil {
		return nil, []error{err2}
	}

	return l, nil
}

// PostLocationByThing checks if the given location entity is valid and adds it to the database
// the new location will be linked to a thing if needed
func (a *APIv1) PostLocationByThing(thingID string, location entities.Location) (*entities.Location, []error) {
	l, err := a.PostLocation(location)
	if err != nil {
		return nil, err
	}

	if len(thingID) != 0 {
		err2 := a.LinkLocation(thingID, l.ID)
		if err2 != nil {
			return nil, []error{err2}
		}

		err3 := a.PostHistoricalLocation(thingID, l.ID)
		if err3 != nil {
			return nil, err3
		}
	}

	return l, nil
}

// GetLocation retrieves a single location by id
func (a *APIv1) GetLocation(id string, qo *odata.QueryOptions) (*entities.Location, error) {
	l, err := a.db.GetLocation(id)
	if err != nil {
		return nil, err
	}

	l.SetLinks(a.config.GetExternalServerURI())
	return l, nil
}

// GetLocations retrieves all locations from the database and returns it as and ArrayResponse
func (a *APIv1) GetLocations(qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	locations, err := a.db.GetLocations()
	return processLocations(a, locations, err)
}

// GetLocationsByThing retrieves the latest locations linked to a thing
func (a *APIv1) GetLocationsByThing(thingID string, qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	locations, err := a.db.GetLocationsByThing(thingID)
	return processLocations(a, locations, err)
}

func processLocations(a *APIv1, locations []*entities.Location, err error) (*models.ArrayResponse, error) {
	uri := a.config.GetExternalServerURI()
	for idx, item := range locations {
		i := *item
		i.SetLinks(uri)
		locations[idx] = &i
	}

	var data interface{} = locations
	return &models.ArrayResponse{
		Count: len(locations),
		Data:  &data,
	}, nil
}

// PatchLocation todo
func (a *APIv1) PatchLocation(id string, location entities.Location) (*entities.Location, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteLocation todo
func (a *APIv1) DeleteLocation(id string) error {
	return gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// LinkLocation links a thing with a location in the database
func (a *APIv1) LinkLocation(thingID string, locationID string) error {
	err3 := a.db.LinkLocation(thingID, locationID)
	if err3 != nil {
		return err3
	}

	return nil
}
