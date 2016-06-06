package api

import (
	"errors"
	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// GetHistoricalLocation retrieves a single HistoricalLocation by id
func (a *APIv1) GetHistoricalLocation(id interface{}, qo *odata.QueryOptions) (*entities.HistoricalLocation, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.HistoricalLocation{})
	if err != nil {
		return nil, err
	}

	hl, err := a.db.GetHistoricalLocation(id, qo)
	if err != nil {
		return nil, err
	}

	hl.SetLinks(a.config.GetExternalServerURI())
	return hl, nil
}

// GetHistoricalLocations retrieves all HistoricalLocations
func (a *APIv1) GetHistoricalLocations(qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.HistoricalLocation{})
	if err != nil {
		return nil, err
	}

	hl, err := a.db.GetHistoricalLocations(qo)
	return processHistoricalLocations(a, hl, err)
}

// GetHistoricalLocationsByLocation retrieves all HistoricalLocations linked to a given location
func (a *APIv1) GetHistoricalLocationsByLocation(locationID interface{}, qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.HistoricalLocation{})
	if err != nil {
		return nil, err
	}

	hl, err := a.db.GetHistoricalLocationsByLocation(locationID, qo)
	return processHistoricalLocations(a, hl, err)
}

// GetHistoricalLocationsByThing retrieves all HistoricalLocations linked to a given thing
func (a *APIv1) GetHistoricalLocationsByThing(thingID interface{}, qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.HistoricalLocation{})
	if err != nil {
		return nil, err
	}

	hl, err := a.db.GetHistoricalLocationsByThing(thingID, qo)
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
func (a *APIv1) PostHistoricalLocation(thingID interface{}, locationID interface{}) []error {
	err := a.db.PostHistoricalLocation(thingID, locationID)
	if err != nil {
		return []error{err}
	}

	return nil
}

// PatchHistoricalLocation todo
func (a *APIv1) PatchHistoricalLocation(id interface{}, hl *entities.HistoricalLocation) (*entities.HistoricalLocation, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteHistoricalLocation deletes a given HistoricalLocation from the database
func (a *APIv1) DeleteHistoricalLocation(id interface{}) error {
	return a.db.DeleteHistoricalLocation(id)
}
