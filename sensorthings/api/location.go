package api

import (
	"errors"
	"log"
	"time"

	gostErrors "github.com/geodan/gost/errors"
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
)

// PostLocation tries to add a new location
func (a *APIv1) PostLocation(location *entities.Location) (*entities.Location, []error) {
	_, err := location.ContainsMandatoryParams()
	if err != nil {
		return nil, err
	}

	if location.EncodingType != entities.EncodingGeoJSON.Value {
		err := errors.New("Encoding not supported. Supported encoding: " + entities.EncodingGeoJSON.Value)
		return nil, []error{err}
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
	var err []error
	var err2 error
	l, err := a.PostLocation(location)
	if len(err) > 0 {
		return nil, err
	}

	if thingID != nil {
		err2 = a.LinkLocation(thingID, l.ID)
		if err2 != nil {
			err3 := a.DeleteLocation(l.ID)
			if err3 != nil {
				log.Printf("Error rolling back location %v", err2)
			}

			return nil, []error{err2}
		}

		hl := &entities.HistoricalLocation{
			Thing:     &entities.Thing{},
			Locations: []*entities.Location{l},
		}

		hl.Thing.ID = thingID
		hl.Time = time.Now().UTC().Format(time.RFC3339Nano)
		hl.ContainsMandatoryParams()

		_, err = a.PostHistoricalLocation(hl)
		if len(err) > 0 {
			err2 := a.DeleteHistoricalLocation(l.ID)
			if err2 != nil {
				log.Printf("Error rolling back location %v", err2)
			}

			return nil, []error{err2}
		}
	}

	l.SetAllLinks(a.config.GetExternalServerURI())

	return l, nil
}

// GetLocation retrieves a single location by id
func (a *APIv1) GetLocation(id interface{}, qo *odata.QueryOptions, path string) (*entities.Location, error) {
	l, err := a.db.GetLocation(id, qo)
	if err != nil {
		return nil, err
	}

	a.SetLinks(l, qo)
	return l, nil
}

// GetLocations retrieves all locations from the database and returns it as and ArrayResponse
func (a *APIv1) GetLocations(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	locations, count, err := a.db.GetLocations(qo)
	return processLocations(a, locations, qo, path, count, err)
}

// GetLocationsByHistoricalLocation retrieves the latest locations linked to a HistoricalLocation
func (a *APIv1) GetLocationsByHistoricalLocation(hlID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	locations, count, err := a.db.GetLocationsByHistoricalLocation(hlID, qo)
	return processLocations(a, locations, qo, path, count, err)
}

// GetLocationsByThing retrieves the latest locations linked to a thing
func (a *APIv1) GetLocationsByThing(thingID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	locations, count, err := a.db.GetLocationsByThing(thingID, qo)
	return processLocations(a, locations, qo, path, count, err)
}

func processLocations(a *APIv1, locations []*entities.Location, qo *odata.QueryOptions, path string, count int, err error) (*models.ArrayResponse, error) {
	for idx, item := range locations {
		i := *item
		a.SetLinks(&i, qo)
		locations[idx] = &i
	}

	var data interface{} = locations
	return &models.ArrayResponse{
		Count:    count,
		NextLink: a.CreateNextLink(count, path, qo),
		Data:     &data,
	}, nil
}

// PatchLocation updates the given location in the database
func (a *APIv1) PatchLocation(id interface{}, location *entities.Location) (*entities.Location, error) {
	if location.HistoricalLocations != nil || location.Things != nil {
		return nil, gostErrors.NewBadRequestError(errors.New("Unable to deep patch Location"))
	}

	if len(location.EncodingType) != 0 {
		if location.EncodingType != entities.EncodingGeoJSON.Value {
			err := errors.New("Encoding not supported. Supported encoding: " + entities.EncodingGeoJSON.Value)
			return nil, err
		}
	}

	return a.db.PatchLocation(id, location)
}

// PutLocation updates the given thing in the database
func (a *APIv1) PutLocation(id interface{}, location *entities.Location) (*entities.Location, []error) {
	var err2 error
	putlocation, err2 := a.db.PutLocation(id, location)
	if err2 != nil {
		return nil, []error{err2}
	}

	putlocation.SetAllLinks(a.config.GetExternalServerURI())
	return putlocation, nil
}

// DeleteLocation deletes a given Location from the database
func (a *APIv1) DeleteLocation(id interface{}) error {
	return a.db.DeleteLocation(id)
}

// LinkLocation links a thing with a location in the database
func (a *APIv1) LinkLocation(thingID interface{}, locationID interface{}) error {
	err := a.db.LinkLocation(thingID, locationID)
	if err != nil {
		return gostErrors.NewBadRequestError(err)
	}
	return nil
}
