package api

import (
	"errors"
	gostErrors "github.com/geodan/gost/errors"
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
)

// PostLocation todo
func (a *APIv1) PostLocation(location entities.Location) (*entities.Location, []error) {
	return nil, []error{gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))}
}

// PostLocationByThing checks if the given location entity is valid and adds it to the database
// the new location will be linked to a thing if needed
func (a *APIv1) PostLocationByThing(thingID string, location entities.Location) (*entities.Location, []error) {
	_, err := location.ContainsMandatoryParams()
	if err != nil {
		return nil, err
	}

	l, err2 := a.db.PostLocation(location)
	if err2 != nil {
		return nil, []error{err2}
	}

	if len(thingID) != 0 {
		err3 := a.LinkLocation(thingID, l.ID)
		if err3 != nil {
			return nil, []error{err3}
		}

		err4 := a.PostHistoricalLocation(thingID, l.ID)
		if err4 != nil {
			return nil, err4
		}
	}

	return l, nil
}

// GetLocation todo
func (a *APIv1) GetLocation(id string, qo *odata.QueryOptions) (*entities.Location, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// GetLocations todo
func (a *APIv1) GetLocations(qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
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
