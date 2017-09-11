package api

import (
	"errors"

	"fmt"

	entities "github.com/gost/core"
	gostErrors "github.com/gost/server/errors"
	"github.com/gost/server/sensorthings/odata"
)

// GetThing returns a thing entity based on the given id and QueryOptions
// returns an error when the entity cannot be found
func (a *APIv1) GetThing(id interface{}, qo *odata.QueryOptions, path string) (*entities.Thing, error) {
	t, err := a.db.GetThing(id, qo)
	if err != nil {
		return nil, err
	}

	a.SetLinks(t, qo)
	return t, nil
}

// GetThingByDatastream returns a thing entity based on the given datastream id and QueryOptions
func (a *APIv1) GetThingByDatastream(id interface{}, qo *odata.QueryOptions, path string) (*entities.Thing, error) {
	t, err := a.db.GetThingByDatastream(id, qo)
	if err != nil {
		return nil, err
	}

	a.SetLinks(t, qo)
	return t, nil
}

// GetThingsByLocation returns things based on the given location id and QueryOptions
func (a *APIv1) GetThingsByLocation(id interface{}, qo *odata.QueryOptions, path string) (*entities.ArrayResponse, error) {
	things, count, err := a.db.GetThingsByLocation(id, qo)
	return processThings(a, things, qo, path, count, err)
}

// GetThingByHistoricalLocation returns a thing entity based on the given HistoricalLocation id and QueryOptions
func (a *APIv1) GetThingByHistoricalLocation(id interface{}, qo *odata.QueryOptions, path string) (*entities.Thing, error) {
	t, err := a.db.GetThingByHistoricalLocation(id, qo)
	if err != nil {
		return nil, err
	}

	a.SetLinks(t, qo)
	return t, nil
}

// GetThings returns an array of thing entities based on the QueryOptions
func (a *APIv1) GetThings(qo *odata.QueryOptions, path string) (*entities.ArrayResponse, error) {
	things, count, err := a.db.GetThings(qo)
	return processThings(a, things, qo, path, count, err)
}

func processThings(a *APIv1, things []*entities.Thing, qo *odata.QueryOptions, path string, count int, err error) (*entities.ArrayResponse, error) {
	if err != nil {
		return nil, err
	}

	for idx, item := range things {
		i := *item
		a.SetLinks(&i, qo)
		things[idx] = &i
	}

	var data interface{} = things
	return a.createArrayResponse(count, path, qo, data), nil
}

// PostThing checks if a posted thing entity is valid and adds it to the database
// a posted thing can also contain Locations and DataStreams
func (a *APIv1) PostThing(thing *entities.Thing) (*entities.Thing, []error) {
	var err []error
	var err2 error
	_, err = containsMandatoryParams(thing)
	if len(err) > 0 {
		return nil, err
	}

	nt, err2 := a.db.PostThing(thing)
	if err2 != nil {
		return nil, []error{err2}
	}

	var postedLocations []*entities.Location
	var postedDatastreams []*entities.Datastream

	// Handle deep insert locations
	if thing.Locations != nil {
		for _, l := range thing.Locations {
			// New location posted
			if l.ID == nil { //Id is null so a new location is posted
				var nl *entities.Location
				if nl, err = a.PostLocationByThing(nt.ID, l); len(err) > 0 {
					a.reverseInserts(nt, postedLocations, postedDatastreams)
					err = append(err, gostErrors.NewConflictRequestError(errors.New("Location deep insert went wrong")))
					return nil, err
				}
				postedLocations = append(postedLocations, nl)
			} else { // posted id: link
				if err2 = a.LinkLocation(nt.ID, l.ID); len(err) > 0 {
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

				if _, err = a.PostHistoricalLocation(hl); len(err) > 0 {
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
				var nd *entities.Datastream
				if nd, err = a.PostDatastreamByThing(nt.ID, d); err != nil {
					a.reverseInserts(nt, postedLocations, postedDatastreams)
					err = append(err, gostErrors.NewConflictRequestError(errors.New("Creating Datastrean went wrong")))
					return nil, err
				}
				postedDatastreams = append(postedDatastreams, nd)
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

// PatchThing updates the given thing in the database
func (a *APIv1) PatchThing(id interface{}, thing *entities.Thing) (*entities.Thing, error) {
	if thing.Datastreams != nil || thing.HistoricalLocations != nil || isDeepPatchLocations(thing.Locations) {
		return nil, gostErrors.NewBadRequestError(errors.New("Unable to deep patch Thing"))
	}

	return a.db.PatchThing(id, thing)
}

// PutThing updates the given thing in the database
func (a *APIv1) PutThing(id interface{}, thing *entities.Thing) (*entities.Thing, []error) {
	var err error
	putthing, err := a.db.PutThing(id, thing)
	if err != nil {
		return nil, []error{err}
	}

	putthing.SetAllLinks(a.config.GetExternalServerURI())
	return putthing, nil
}

func isDeepPatchLocations(locations []*entities.Location) bool {
	if locations != nil {
		for _, l := range locations {
			if l.ID == nil || len(fmt.Sprintf("%v", l.ID)) == 0 {
				return true
			}
		}
	}

	return false
}
