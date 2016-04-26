package sensorthings

import (
	"github.com/tebben/gost/configuration"
)

const (
	API_PREFIX string = "v1.0" // API Prefix for V1.0 endpoint
)

type SensorThingsApi interface {
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

type Api struct {
	db        Database
	config    configuration.Config
	endPoints []Endpoint
	//mqtt      mqtt.MQTTServer
}

// Initialise a new SensorThings API
func NewApi(database Database, config configuration.Config) SensorThingsApi {
	return &Api{
		db:     database,
		//mqtt:   mqtt,
		config: config,
	}
}

// Get the version info of the current supported SensorThings API Version and running server version
func (a *Api) GetVersionInfo() *VersionInfo {
	versionInfo := VersionInfo{
		GostServerVersion:      GostServerVersion{configuration.ServerVersion},
		SensorThingsApiVersion: SensorThingsApiVersion{configuration.SensorThingsAPIVersion},
	}

	return &versionInfo
}

// Navigating to the base resource path will return a JSON array of the available SensorThings resource endpoints.
func (a *Api) GetBasePathInfo() *ArrayResponse {
	var rp interface{} = a.GetEndpoints()
	basePathInfo := ArrayResponse{
		Data: &rp,
	}

	return &basePathInfo
}

func (a *Api) GetEndpoints() *[]Endpoint {
	if a.endPoints == nil {
		a.endPoints = CreateEndPoints(a.config.GetExternalServerUri())
	}

	return &a.endPoints
}

func (a *Api) GetThing(id string, qo *QueryOptions) (*Thing, error) {
	t, err := a.db.GetThing(id)
	if(err != nil){
		return nil, err
	}

	t.SetLinks(a.config.GetExternalServerUri())
	return t, nil
}

func (a *Api) GetThings(qo *QueryOptions) (*ArrayResponse, error) {
	things, err := a.db.GetThings()
	if(err != nil){
		return nil, err
	}

	uri := a.config.GetExternalServerUri()
	for idx, item := range(things){
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

func (a *Api) PostThing(thing Thing) (*Thing, []error) {
	_, err := thing.ContainsMandatoryPostParams()
	if(err != nil){
		return nil, err
	}

	nt, err2 := a.db.PostThing(thing)
	if(err2 != nil){
		return nil, []error{err2}
	}

	// Handle locations
	if(thing.Locations != nil){
		for _, l := range thing.Locations{
			location := *l
			// New location posted
			if(len(l.ID) == 0){ //Id is null so a new location is posted
				_, err3 := a.PostLocation(location, nt.ID)
				if(err3 != nil){
					return nil, err3
				}
			}else { // posted id: link
				err4 := a.LinkLocation(nt.ID, l.ID)
				if(err4 != nil){
					// todo: thing is posted, delete it
					return nil, []error{err4}
				}

				err5 := a.PostHistoricalLocation(nt.ID, l.ID)
				if(err5 != nil){
					// todo: things is posted, delete it
					return nil, []error{err5}
				}
			}
		}
	}

	nt.SetLinks(a.config.GetExternalServerUri())
	//push to mqtt
	return nt, nil
}

func (a *Api) DeleteThing(id string) {
}

func (a *Api) PatchThing(thing Thing) {

}

func (a *Api) GetLocation(id string) *Location{
	return nil
}

func (a *Api) GetLocations() *ArrayResponse{
	return nil
}

func (a *Api) PatchLocation(id string) {
//	return nil
}

func (a *Api) PostLocation(location Location, thingID string) (*Location, []error){
	_, err := location.ContainsMandatoryPostParams()
	if(err != nil){
		return nil, err
	}

	l, err2 := a.db.PostLocation(location)
	if(err2 != nil){
		return nil, []error{err2}
	}

	if(len(thingID) != 0) {
		err3 := a.LinkLocation(thingID, l.ID)
		if(err3 != nil){
			return nil, []error{err3}
		}

		err4 := a.PostHistoricalLocation(thingID, l.ID)
		if(err4 != nil){
			return nil, []error{err4}
		}
	}

	return l, nil
}

func (a *Api) DeleteLocation(id string){

}

func (a *Api) LinkLocation(thingID string, locationID string) error{
	err3 := a.db.LinkLocation(thingID, locationID);
	if err3 != nil {
		return err3
	}

	return nil
}

func (a *Api) PostHistoricalLocation(thingID string, locationID string) error{
	err := a.db.PostHistoricalLocation(thingID, locationID);
	if err != nil {
		return err
	}

	return nil
}