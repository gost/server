package api

import (
	"errors"
	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// GetHistoricalLocation retrieves a single HistoricalLocation by id
func (a *APIv1) GetHistoricalLocation(id string, qo *odata.QueryOptions) (*entities.HistoricalLocation, error) {
	hl, err := a.db.GetHistoricalLocation(id)
	if err != nil {
		return nil, err
	}

	hl.SetLinks(a.config.GetExternalServerURI())
	return hl, nil
}

// GetHistoricalLocations retrieves all HistoricalLocations
func (a *APIv1) GetHistoricalLocations(qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	hl, err := a.db.GetHistoricalLocations()
	return processHistoricalLocations(a, hl, err)
}

// GetHistoricalLocationsByLocation retrieves all HistoricalLocations linked to a given location
func (a *APIv1) GetHistoricalLocationsByLocation(locationID string, qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	hl, err := a.db.GetHistoricalLocationsByLocation(locationID)
	return processHistoricalLocations(a, hl, err)
}

// GetHistoricalLocationsByThing retrieves all HistoricalLocations linked to a given thing
func (a *APIv1) GetHistoricalLocationsByThing(thingID string, qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	hl, err := a.db.GetHistoricalLocationsByThing(thingID)
	return processHistoricalLocations(a, hl, err)
}

func processHistoricalLocations(a *APIv1, historicalLocations []*entities.HistoricalLocation, err error) (*models.ArrayResponse, error) {
	uri := a.config.GetExternalServerURI()
	for idx, item := range historicalLocations {
		i := *item
		i.SetLinks(uri)
		historicalLocations[idx] = &i
	}

	var data interface{} = historicalLocations
	return &models.ArrayResponse{
		Count: len(historicalLocations),
		Data:  &data,
	}, nil
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
func (a *APIv1) PatchHistoricalLocation(id string, hl *entities.HistoricalLocation) (*entities.HistoricalLocation, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteHistoricalLocation todo
func (a *APIv1) DeleteHistoricalLocation(id string) error {
	return gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}
