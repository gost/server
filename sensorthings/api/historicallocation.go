package api

import (
	"errors"
	gostErrors "github.com/geodan/gost/errors"
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
)

// GetHistoricalLocation todo
func (a *APIv1) GetHistoricalLocation(id string, qo *odata.QueryOptions) (*entities.HistoricalLocation, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// GetHistoricalLocations todo
func (a *APIv1) GetHistoricalLocations(qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// PostHistoricalLocation is triggered by code and cannot be used from any endpoint PostHistoricalLocation
// adds a HistoricalLocation into the database
func (a *APIv1) PostHistoricalLocation(thingID string, locationID string) []error {
	err := a.db.PostHistoricalLocation(thingID, locationID)
	if err != nil {
		return []error{err}
	}

	return nil
}

// PatchHistoricalLocation todo
func (a *APIv1) PatchHistoricalLocation(id string, hl entities.HistoricalLocation) (*entities.HistoricalLocation, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteHistoricalLocation todo
func (a *APIv1) DeleteHistoricalLocation(id string) error {
	return gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}
