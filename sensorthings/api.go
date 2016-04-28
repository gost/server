package sensorthings

import (
	"github.com/geodan/gost/configuration"
)

const (
	// APIPrefix for V1.0 endpoint
	APIPrefix string = "v1.0"
)

// API describes all request and responses to fulfill the SensorThings API standard
type API interface {
	GetConfig() *configuration.Config

	GetVersionInfo() *VersionInfo
	GetBasePathInfo() *ArrayResponse
	GetEndpoints() *[]Endpoint

	GetThing(id string, qo *QueryOptions) (*Thing, error)
	GetThings(qo *QueryOptions) (*ArrayResponse, error)
	PostThing(thing Thing) (*Thing, []error)
	DeleteThing(id string)
	PatchThing(thing Thing)

	GetLocation(id string) *Location
	GetLocations() *ArrayResponse
	PostLocation(location Location, thingID string) (*Location, []error)
	DeleteLocation(id string)
	PatchLocation(id string)

	PostHistoricalLocation(thingID string, locationID string) error
	LinkLocation(thingID string, locationID string) error
}

// APIv1 is the default implementation of SensorThingsApi, API needs a database
// provider, config, endpoint information to setup te needed services
type APIv1 struct {
	db        Database
	config    configuration.Config
	endPoints []Endpoint
	//mqtt      mqtt.MQTTServer
}

// NewAPI Initialise a new SensorThings API
func NewAPI(database Database, config configuration.Config) API {
	return &APIv1{
		db: database,
		//mqtt:   mqtt,
		config: config,
	}
}

// GetConfig return the current configuration.Config set for the api
func (a *APIv1) GetConfig() *configuration.Config {
	return &a.config
}

// GetVersionInfo retrieves the version info of the current supported SensorThings API Version and running server version
func (a *APIv1) GetVersionInfo() *VersionInfo {
	versionInfo := VersionInfo{
		GostServerVersion: GostServerVersion{configuration.ServerVersion},
		APIVersion:        APIVersion{configuration.SensorThingsAPIVersion},
	}

	return &versionInfo
}

// GetBasePathInfo when navigating to the base resource path will return a JSON array of the available SensorThings resource endpoints.
func (a *APIv1) GetBasePathInfo() *ArrayResponse {
	var rp interface{} = a.GetEndpoints()
	basePathInfo := ArrayResponse{
		Data: &rp,
	}

	return &basePathInfo
}

// GetEndpoints returns all configured endpoints for the HTTP server
func (a *APIv1) GetEndpoints() *[]Endpoint {
	if a.endPoints == nil {
		a.endPoints = CreateEndPoints(a.config.GetExternalServerURI())
	}

	return &a.endPoints
}

// GetThing returns a thing entity based on the given id and QueryOptions
// returns an error when the entity cannot be found
func (a *APIv1) GetThing(id string, qo *QueryOptions) (*Thing, error) {
	t, err := a.db.GetThing(id)
	if err != nil {
		return nil, err
	}

	t.SetLinks(a.config.GetExternalServerURI())
	return t, nil
}

// GetThings returns an array of thing entities based on the QueryOptions
func (a *APIv1) GetThings(qo *QueryOptions) (*ArrayResponse, error) {
	things, err := a.db.GetThings()
	if err != nil {
		return nil, err
	}

	uri := a.config.GetExternalServerURI()
	for idx, item := range things {
		i := *item
		item.SetLinks(uri)
		things[idx] = &i
	}

	var data interface{} = things

	var count = len(things)
	response := ArrayResponse{
		Count: &count,
		Data:  &data,
	}

	return &response, nil
}

// PostThing checks if a posted thing entity is valid and adds it to the database
// a posted thing can also contain Locations and DataStreams
func (a *APIv1) PostThing(thing Thing) (*Thing, []error) {
	_, err := thing.ContainsMandatoryPostParams()
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
				_, err3 := a.PostLocation(location, nt.ID)
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
					return nil, []error{err5}
				}
			}
		}
	}

	nt.SetLinks(a.config.GetExternalServerURI())
	//push to mqtt
	return nt, nil
}

// DeleteThing todo
func (a *APIv1) DeleteThing(id string) {
}

// PatchThing todo
func (a *APIv1) PatchThing(thing Thing) {

}

// GetLocation todo
func (a *APIv1) GetLocation(id string) *Location {
	return nil
}

// GetLocations todo
func (a *APIv1) GetLocations() *ArrayResponse {
	return nil
}

// PatchLocation todo
func (a *APIv1) PatchLocation(id string) {
	//	return nil
}

// PostLocation checks if the given location entity is valid and adds it to the database
// the new location will be linked to a thing if needed
func (a *APIv1) PostLocation(location Location, thingID string) (*Location, []error) {
	_, err := location.ContainsMandatoryPostParams()
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
			return nil, []error{err4}
		}
	}

	return l, nil
}

// DeleteLocation todo
func (a *APIv1) DeleteLocation(id string) {

}

// LinkLocation links a thing with a location in the database
func (a *APIv1) LinkLocation(thingID string, locationID string) error {
	err3 := a.db.LinkLocation(thingID, locationID)
	if err3 != nil {
		return err3
	}

	return nil
}

// PostHistoricalLocation is triggered by code and cannot be used from any endpoint PostHistoricalLocation
// adds a HistoricalLocation into the database
func (a *APIv1) PostHistoricalLocation(thingID string, locationID string) error {
	err := a.db.PostHistoricalLocation(thingID, locationID)
	if err != nil {
		return err
	}

	return nil
}
