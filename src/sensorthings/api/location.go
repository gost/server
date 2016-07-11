package api

import (
	"errors"
	"log"

	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// PostLocation tries to add a new location
func (a *APIv1) PostLocation(location *entities.Location) (*entities.Location, []error) {
	_, err := location.ContainsMandatoryParams()
	if err != nil {
		return nil, err
	}

	supported, err2 := entities.CheckEncodingSupported(location, location.EncodingType)
	if !supported || err2 != nil {
		return nil, []error{err2}
	}

	l, err2 := a.db.PostLocation(location)
	if err2 != nil {
		return nil, []error{err2}
	}
	l.SetAllLinks(a.config.GetExternalServerURI())
	return l, nil
}

// PostLocationByThing checks if the given location entity is valid and adds it to the database
// the new location will be linked to a thing if needed
func (a *APIv1) PostLocationByThing(thingID interface{}, location *entities.Location) (*entities.Location, []error) {
	l, err := a.PostLocation(location)
	if err != nil {
		return nil, err
	}

	if thingID != nil {
		err2 := a.LinkLocation(thingID, l.ID)
		if err2 != nil {
			err3 := a.DeleteLocation(l.ID)
			if err3 != nil {
				log.Printf("Error rolling back location %v", err3)
			}

			return nil, []error{err2}
		}

		hl := &entities.HistoricalLocation{
			Thing:     &entities.Thing{},
			Locations: []*entities.Location{l},
		}

		hl.Thing.ID = thingID
		hl.ContainsMandatoryParams()

		hl, err = a.PostHistoricalLocation(hl)
		if err != nil {
			err3 := a.DeleteHistoricalLocation(l.ID)
			if err3 != nil {
				log.Printf("Error rolling back location %v", err3)
			}

			return nil, []error{err2}
		}

	}

	l.SetAllLinks(a.config.GetExternalServerURI())

	return l, nil
}

// GetLocation retrieves a single location by id
func (a *APIv1) GetLocation(id interface{}, qo *odata.QueryOptions, path string) (*entities.Location, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Location{})
	if err != nil {
		return nil, err
	}

	l, err := a.db.GetLocation(id, qo)
	if err != nil {
		return nil, err
	}

	a.ProcessGetRequest(l, qo)
	return l, nil
}

// GetLocations retrieves all locations from the database and returns it as and ArrayResponse
func (a *APIv1) GetLocations(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Location{})
	if err != nil {
		return nil, err
	}

	locations, err := a.db.GetLocations(qo)
	return processLocations(a, locations, qo, path, err)
}

// GetLocationsByHistoricalLocation retrieves the latest locations linked to a HistoricalLocation
func (a *APIv1) GetLocationsByHistoricalLocation(hlID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Location{})
	if err != nil {
		return nil, err
	}

	locations, err := a.db.GetLocationsByHistoricalLocation(hlID, qo)
	return processLocations(a, locations, qo, path, err)
}

// GetLocationsByThing retrieves the latest locations linked to a thing
func (a *APIv1) GetLocationsByThing(thingID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Location{})
	if err != nil {
		return nil, err
	}

	locations, err := a.db.GetLocationsByThing(thingID, qo)
	return processLocations(a, locations, qo, path, err)
}

func processLocations(a *APIv1, locations []*entities.Location, qo *odata.QueryOptions, path string, err error) (*models.ArrayResponse, error) {
	for idx, item := range locations {
		i := *item
		a.ProcessGetRequest(&i, qo)
		locations[idx] = &i
	}

	var data interface{} = locations
	return &models.ArrayResponse{
		Count:    len(locations),
		NextLink: a.CreateNextLink(a.db.GetTotalLocations(), path, qo),
		Data:     &data,
	}, nil
}

// PatchLocation todo
func (a *APIv1) PatchLocation(id interface{}, location *entities.Location) (*entities.Location, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteLocation deletes a given Location from the database
func (a *APIv1) DeleteLocation(id interface{}) error {
	return a.db.DeleteLocation(id)
}

// LinkLocation links a thing with a location in the database
func (a *APIv1) LinkLocation(thingID interface{}, locationID interface{}) error {
	err3 := a.db.LinkLocation(thingID, locationID)
	if err3 != nil {
		return gostErrors.NewBadRequestError(err3)
	}
	return nil
}
