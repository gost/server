package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	entities "github.com/gost/core"
	"github.com/gost/server/configuration"
	gostErrors "github.com/gost/server/errors"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/odata"
	"github.com/gost/server/sensorthings/rest/endpoint"
	"github.com/stretchr/testify/assert"
)

func TestVersionResponse(t *testing.T) {
	// act
	r := request("GET", "/version", nil)
	//HandleVersion()
	version := models.VersionInfo{}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &version)

	// assert
	assertStatusCode(http.StatusOK, r, t)
	assert.Equal(t, configuration.SensorThingsAPIVersion, version.APIVersion.Version)
	assert.Equal(t, configuration.ServerVersion, version.GostServerVersion.Version)
}

func assertStatusCode(expectedStatusCode int, r *http.Response, t *testing.T) {
	assert.Equal(t, expectedStatusCode, r.StatusCode)
}

func request(method, url string, body interface{}) *http.Response {
	var reader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reader = bytes.NewReader(b)
	}

	client := &http.Client{}
	request, _ := http.NewRequest(method, getServer().URL+url, reader)
	r, _ := client.Do(request)

	return r
}

var testServer *httptest.Server

func getRouter() *mux.Router {
	a := newMockAPI()
	eps := endpoint.EndpointsToSortedList(a.GetEndpoints())
	router := mux.NewRouter().StrictSlash(false)

	for _, e := range eps {
		op := e
		operation := op.Operation
		method := fmt.Sprintf("%s", operation.OperationType)
		router.Methods(method).
			Path(operation.Path).
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				operation.Handler(w, r, &op.Endpoint, &a)
			})
	}

	return router
}

func getServer() *httptest.Server {
	if testServer == nil {
		router := getRouter()
		testServer = httptest.NewServer(router)
	}

	return testServer
}

func newMockThing(id int) *entities.Thing {
	thing := &entities.Thing{Name: fmt.Sprintf("thing %v", id), Description: fmt.Sprintf("description of thing %v", id), Properties: map[string]interface{}{"type": "none"}}
	thing.ID = id
	return thing
}

func newMockSensor(id int) *entities.Sensor {
	sensor := &entities.Sensor{Name: fmt.Sprintf("sensor %v", id), Description: fmt.Sprintf("description of sensor %v", id), EncodingType: "PDF", Metadata: "none"}
	sensor.ID = id
	return sensor
}

func newMockLocation(id int) *entities.Location {
	location := &entities.Location{Name: fmt.Sprintf("location %v", id),
		Description:  fmt.Sprintf("description of location %v", id),
		EncodingType: "application/vnd.geo+json",
		Location:     map[string]interface{}{"coordinates": "test"}}
	location.ID = id
	return location
}

func newMockHistoricalLocation(id int) *entities.HistoricalLocation {
	historicalLocation := &entities.HistoricalLocation{
		Time: "2017-07-17T07:03:09.194Z",
	}
	historicalLocation.ID = id
	return historicalLocation
}

func newMockObservedProperty(id int) *entities.ObservedProperty {
	op := &entities.ObservedProperty{Name: fmt.Sprintf("sensor %v", id), Description: fmt.Sprintf("description of sensor %v", id), Definition: "none"}
	op.ID = id
	return op
}

func newMockObservation(id int) *entities.Observation {
	op := &entities.Observation{Result: 35, PhenomenonTime: "2017-07-17T07:03:09.194Z", ResultQuality: "high"}
	op.ID = id
	return op
}

func newMockFeatureOfInterest(id int) *entities.FeatureOfInterest {
	foi := &entities.FeatureOfInterest{Name: fmt.Sprintf("foi %v", id), Description: fmt.Sprintf("description of foi %v", id), EncodingType: "application/vnd.geo+json"}
	foi.ID = id
	return foi
}

func newMockDatastream(id int) *entities.Datastream {
	ds := &entities.Datastream{Name: fmt.Sprintf("datastream %v", id), Description: fmt.Sprintf("description of datastream %v", id)}
	ds.ID = id
	return ds
}

type MockAPI struct {
	config *configuration.Config
}

// NewAPI Initialise a new SensorThings API
func newMockAPI() models.API {
	api := MockAPI{}
	return &api
}

func (a *MockAPI) Start() {}
func (a *MockAPI) GetConfig() *configuration.Config {
	if a.config != nil {
		return a.config
	}

	a.config = &configuration.Config{
		Server: configuration.ServerConfig{
			IndentedJSON:      true,
			MaxEntityResponse: 200,
			ExternalURI:       "localhost",
			HTTPS:             false,
			HTTPSCert:         "",
			HTTPSKey:          "",
			Host:              "0.0.0.0",
			Name:              "GOST Server",
			Port:              8080,
		},
	}

	return a.config
}

func (a *MockAPI) GetAcceptedPaths() []string { return []string{} }
func (a *MockAPI) GetBasePathInfo() *models.ArrayResponse {
	bpi := []models.Endpoint{}
	ep := *a.GetEndpoints()
	for _, e := range ep {
		if e.ShowOutputInfo() {
			bpi = append(bpi, e)
		}
	}

	var i interface{} = bpi
	basePathInfo := models.ArrayResponse{
		Data: &i,
	}

	return &basePathInfo
}

func (a *MockAPI) GetEndpoints() *map[entities.EntityType]models.Endpoint {
	eps := createEndPoints("localhost")
	return &eps
}

func (a *MockAPI) initRest()                                               {}
func (a *MockAPI) GetTopics() *[]models.Topic                              { return nil }
func (a *MockAPI) SetLinks(entity entities.Entity, qo *odata.QueryOptions) {}
func (a *MockAPI) CreateNextLink(count int, incomingURL string, qo *odata.QueryOptions) string {
	return ""
}

func (a *MockAPI) GetThing(id interface{}, qo *odata.QueryOptions, path string) (*entities.Thing, error) {
	return getMockThing(id)
}
func (a *MockAPI) GetThings(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return getMockThings()
}
func (a *MockAPI) GetThingByDatastream(id interface{}, qo *odata.QueryOptions, path string) (*entities.Thing, error) {
	return getMockThing(id)
}
func (a *MockAPI) GetThingsByLocation(id interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return getMockThings()
}
func (a *MockAPI) GetThingByHistoricalLocation(id interface{}, qo *odata.QueryOptions, path string) (*entities.Thing, error) {
	return getMockThing(id)
}

func getMockThing(id interface{}) (*entities.Thing, error) {
	intID, ok := toIntID(id)
	if !ok || intID != 1 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Thing does not exist"))
	}
	return newMockThing(intID), nil
}

func getMockLocation(id interface{}) (*entities.Location, error) {
	intID, ok := toIntID(id)
	if !ok || intID != 1 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Location does not exist"))
	}
	return newMockLocation(intID), nil
}

func getMockHistoricalLocation(id interface{}) (*entities.HistoricalLocation, error) {
	intID, ok := toIntID(id)
	if !ok || intID != 1 {
		return nil, gostErrors.NewRequestNotFound(errors.New("HsitoricalLocation does not exist"))
	}
	return newMockHistoricalLocation(intID), nil
}

func getMockFeatureOfInterest(id interface{}) (*entities.FeatureOfInterest, error) {
	intID, ok := toIntID(id)
	if !ok || intID != 1 {
		return nil, gostErrors.NewRequestNotFound(errors.New("featureOfInterest does not exist"))
	}
	return newMockFeatureOfInterest(intID), nil
}

func getMockSensor(id interface{}) (*entities.Sensor, error) {
	intID, ok := toIntID(id)
	if !ok || intID != 1 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Sensor does not exist"))
	}
	return newMockSensor(intID), nil
}

func getMockDatastream(id interface{}) (*entities.Datastream, error) {
	intID, ok := toIntID(id)
	if !ok || intID != 1 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}
	return newMockDatastream(intID), nil
}

func getMockObservedProperty(id interface{}) (*entities.ObservedProperty, error) {
	intID, ok := toIntID(id)
	if !ok || intID != 1 {
		return nil, gostErrors.NewRequestNotFound(errors.New("ObservedProperty does not exist"))
	}
	return newMockObservedProperty(intID), nil
}

func getMockObservation(id interface{}) (*entities.Observation, error) {
	intID, ok := toIntID(id)
	if !ok || intID != 1 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Observation does not exist"))
	}
	return newMockObservation(intID), nil
}

func getMockThings() (*models.ArrayResponse, error) {
	var data interface{} = []*entities.Thing{newMockThing(1), newMockThing(2)}
	return &models.ArrayResponse{
		Count: 2,
		Data:  &data,
	}, nil
}

func getMockLocations() (*models.ArrayResponse, error) {
	var data interface{} = []*entities.Location{newMockLocation(1), newMockLocation(2)}
	return &models.ArrayResponse{
		Count: 2,
		Data:  &data,
	}, nil
}

func getMockHistoricalLocations() (*models.ArrayResponse, error) {
	var data interface{} = []*entities.HistoricalLocation{newMockHistoricalLocation(1), newMockHistoricalLocation(2)}
	return &models.ArrayResponse{
		Count: 2,
		Data:  &data,
	}, nil
}

func getMockfeaturesOfInterest() (*models.ArrayResponse, error) {
	var data interface{} = []*entities.FeatureOfInterest{newMockFeatureOfInterest(1), newMockFeatureOfInterest(2)}
	return &models.ArrayResponse{
		Count: 2,
		Data:  &data,
	}, nil
}

func getMockSensors() (*models.ArrayResponse, error) {
	var data interface{} = []*entities.Sensor{newMockSensor(1), newMockSensor(2)}
	return &models.ArrayResponse{
		Count: 2,
		Data:  &data,
	}, nil
}

func getMockObservedProperties() (*models.ArrayResponse, error) {
	var data interface{} = []*entities.ObservedProperty{newMockObservedProperty(1), newMockObservedProperty(2)}
	return &models.ArrayResponse{
		Count: 2,
		Data:  &data,
	}, nil
}

func getMockObservations() (*models.ArrayResponse, error) {
	var data interface{} = []*entities.Observation{newMockObservation(1), newMockObservation(2)}
	return &models.ArrayResponse{
		Count: 2,
		Data:  &data,
	}, nil
}

func getMockDatastreams() (*models.ArrayResponse, error) {
	var data interface{} = []*entities.Datastream{newMockDatastream(1), newMockDatastream(2)}
	return &models.ArrayResponse{
		Count: 2,
		Data:  &data,
	}, nil
}

func (a *MockAPI) PostThing(thing *entities.Thing) (*entities.Thing, []error) {
	return thing, nil
}

func (a *MockAPI) PatchThing(id interface{}, thing *entities.Thing) (*entities.Thing, error) {
	return thing, nil
}
func (a *MockAPI) PutThing(id interface{}, thing *entities.Thing) (*entities.Thing, []error) {
	return thing, nil
}
func (a *MockAPI) DeleteThing(id interface{}) error { return nil }

func (a *MockAPI) GetLocation(id interface{}, qo *odata.QueryOptions, path string) (*entities.Location, error) {
	return getMockLocation(id)
}
func (a *MockAPI) GetLocations(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return getMockLocations()
}
func (a *MockAPI) GetLocationsByHistoricalLocation(hlID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return getMockLocations()
}
func (a *MockAPI) GetLocationsByThing(thingID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return getMockLocations()
}
func (a *MockAPI) PostLocation(location *entities.Location) (*entities.Location, []error) {
	return location, nil
}
func (a *MockAPI) PostLocationByThing(thingID interface{}, location *entities.Location) (*entities.Location, []error) {
	return location, nil
}
func (a *MockAPI) PatchLocation(id interface{}, location *entities.Location) (*entities.Location, error) {
	return location, nil
}
func (a *MockAPI) PutLocation(id interface{}, location *entities.Location) (*entities.Location, []error) {
	return location, nil
}
func (a *MockAPI) DeleteLocation(id interface{}) error { return nil }

func (a *MockAPI) GetHistoricalLocation(id interface{}, qo *odata.QueryOptions, path string) (*entities.HistoricalLocation, error) {
	return getMockHistoricalLocation(id)
}
func (a *MockAPI) GetHistoricalLocations(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return getMockHistoricalLocations()
}
func (a *MockAPI) GetHistoricalLocationsByLocation(locationID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return getMockHistoricalLocations()
}
func (a *MockAPI) GetHistoricalLocationsByThing(thingID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return getMockHistoricalLocations()
}
func (a *MockAPI) PostHistoricalLocation(hl *entities.HistoricalLocation) (*entities.HistoricalLocation, []error) {
	return hl, nil
}
func (a *MockAPI) PutHistoricalLocation(id interface{}, hl *entities.HistoricalLocation) (*entities.HistoricalLocation, []error) {
	return hl, nil
}
func (a *MockAPI) PatchHistoricalLocation(id interface{}, hl *entities.HistoricalLocation) (*entities.HistoricalLocation, error) {
	return hl, nil
}
func (a *MockAPI) DeleteHistoricalLocation(id interface{}) error { return nil }

func (a *MockAPI) GetDatastream(id interface{}, qo *odata.QueryOptions, path string) (*entities.Datastream, error) {
	return getMockDatastream(id)
}
func (a *MockAPI) GetDatastreams(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return getMockDatastreams()
}
func (a *MockAPI) GetDatastreamByObservation(id interface{}, qo *odata.QueryOptions, path string) (*entities.Datastream, error) {
	return getMockDatastream(id)
}
func (a *MockAPI) GetDatastreamsByThing(thingID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return getMockDatastreams()
}
func (a *MockAPI) GetDatastreamsBySensor(sensorID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return getMockDatastreams()
}
func (a *MockAPI) GetDatastreamsByObservedProperty(sensorID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return getMockDatastreams()
}
func (a *MockAPI) PostDatastream(datastream *entities.Datastream) (*entities.Datastream, []error) {
	return datastream, nil
}
func (a *MockAPI) PostDatastreamByThing(thingID interface{}, datastream *entities.Datastream) (*entities.Datastream, []error) {
	return datastream, nil
}
func (a *MockAPI) PatchDatastream(id interface{}, datastream *entities.Datastream) (*entities.Datastream, error) {
	return datastream, nil
}
func (a *MockAPI) PutDatastream(id interface{}, datastream *entities.Datastream) (*entities.Datastream, []error) {
	return datastream, nil
}
func (a *MockAPI) DeleteDatastream(id interface{}) error { return nil }

func (a *MockAPI) GetFeatureOfInterest(id interface{}, qo *odata.QueryOptions, path string) (*entities.FeatureOfInterest, error) {
	return getMockFeatureOfInterest(id)
}
func (a *MockAPI) GetFeatureOfInterestByObservation(id interface{}, qo *odata.QueryOptions, path string) (*entities.FeatureOfInterest, error) {
	return getMockFeatureOfInterest(1)
}
func (a *MockAPI) GetFeatureOfInterests(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return getMockfeaturesOfInterest()
}
func (a *MockAPI) PostFeatureOfInterest(foi *entities.FeatureOfInterest) (*entities.FeatureOfInterest, []error) {
	return foi, nil
}
func (a *MockAPI) PatchFeatureOfInterest(id interface{}, foi *entities.FeatureOfInterest) (*entities.FeatureOfInterest, error) {
	return foi, nil
}
func (a *MockAPI) PutFeatureOfInterest(id interface{}, foi *entities.FeatureOfInterest) (*entities.FeatureOfInterest, []error) {
	return foi, nil
}
func (a *MockAPI) DeleteFeatureOfInterest(id interface{}) error { return nil }

func (a *MockAPI) GetObservation(id interface{}, qo *odata.QueryOptions, path string) (*entities.Observation, error) {
	return getMockObservation(id)
}
func (a *MockAPI) GetObservations(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return getMockObservations()
}
func (a *MockAPI) GetObservationsByDatastream(datastreamID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return getMockObservations()
}
func (a *MockAPI) GetObservationsByFeatureOfInterest(foiID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return getMockObservations()
}
func (a *MockAPI) PostObservation(observation *entities.Observation) (*entities.Observation, []error) {
	return observation, nil
}
func (a *MockAPI) PostObservationByDatastream(datastreamID interface{}, observation *entities.Observation) (*entities.Observation, []error) {
	return observation, nil
}
func (a *MockAPI) PatchObservation(id interface{}, observation *entities.Observation) (*entities.Observation, error) {
	return observation, nil
}
func (a *MockAPI) PutObservation(id interface{}, observation *entities.Observation) (*entities.Observation, []error) {
	return observation, nil
}
func (a *MockAPI) DeleteObservation(id interface{}) error { return nil }

func (a *MockAPI) GetObservedProperty(id interface{}, qo *odata.QueryOptions, path string) (*entities.ObservedProperty, error) {
	return getMockObservedProperty(id)
}
func (a *MockAPI) GetObservedProperties(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return getMockObservedProperties()
}
func (a *MockAPI) GetObservedPropertyByDatastream(datastreamID interface{}, qo *odata.QueryOptions, path string) (*entities.ObservedProperty, error) {
	return getMockObservedProperty(datastreamID)
}
func (a *MockAPI) PostObservedProperty(op *entities.ObservedProperty) (*entities.ObservedProperty, []error) {
	return op, nil
}
func (a *MockAPI) PatchObservedProperty(id interface{}, op *entities.ObservedProperty) (*entities.ObservedProperty, error) {
	return op, nil
}
func (a *MockAPI) PutObservedProperty(id interface{}, op *entities.ObservedProperty) (*entities.ObservedProperty, []error) {
	return op, nil
}
func (a *MockAPI) DeleteObservedProperty(id interface{}) error { return nil }

func (a *MockAPI) GetSensor(id interface{}, qo *odata.QueryOptions, path string) (*entities.Sensor, error) {
	return getMockSensor(id)
}
func (a *MockAPI) GetSensorByDatastream(id interface{}, qo *odata.QueryOptions, path string) (*entities.Sensor, error) {
	return getMockSensor(id)
}
func (a *MockAPI) GetSensors(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return getMockSensors()
}
func (a *MockAPI) PostSensor(sensor *entities.Sensor) (*entities.Sensor, []error) { return sensor, nil }
func (a *MockAPI) PatchSensor(id interface{}, sensor *entities.Sensor) (*entities.Sensor, error) {
	return sensor, nil
}
func (a *MockAPI) DeleteSensor(id interface{}) error { return nil }
func (a *MockAPI) PutSensor(id interface{}, sensor *entities.Sensor) (*entities.Sensor, []error) {
	return sensor, nil
}

func (a *MockAPI) LinkLocation(thingID interface{}, locationID interface{}) error { return nil }

func (a *MockAPI) GetVersionInfo() *models.VersionInfo {
	versionInfo := models.VersionInfo{
		GostServerVersion: models.GostServerVersion{Version: configuration.ServerVersion},
		APIVersion:        models.APIVersion{Version: configuration.SensorThingsAPIVersion},
	}

	return &versionInfo
}

// ToIntID converts an interface to int id used for the id's in the database
func toIntID(id interface{}) (int, bool) {
	switch t := id.(type) {
	case string:
		intID, err := strconv.Atoi(t)
		if err != nil {
			return 0, false
		}
		return intID, true
	case float64:
		return int(t), true
	}

	intID, err := strconv.Atoi(fmt.Sprintf("%v", id))
	if err != nil {
		// why not return:  0, err
		return 0, false
	}

	// why not return: intID, nil
	return intID, true
}

func createEndPoints(externalURL string) map[entities.EntityType]models.Endpoint {
	return map[entities.EntityType]models.Endpoint{
		entities.EntityTypeVersion: &endpoint.Endpoint{
			Name:       "Version",
			OutputInfo: false,
			Operations: []models.EndpointOperation{
				{OperationType: models.HTTPOperationGet, Path: "/version", Handler: HandleVersion},
			},
		},
		entities.EntityTypeUnknown: &endpoint.Endpoint{
			Name:       "Root",
			OutputInfo: false,
			Operations: []models.EndpointOperation{
				{OperationType: models.HTTPOperationGet, Path: "/v1.0", Handler: HandleAPIRoot},
			},
		},
		entities.EntityTypeThing: &endpoint.Endpoint{
			Name:       "Things",
			EntityType: entities.EntityTypeThing,
			OutputInfo: true,
			Operations: []models.EndpointOperation{
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/things", Handler: HandleGetThings},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/things{id}", Handler: HandleGetThing},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/historicallocations{id}/thing", Handler: HandleGetThingByHistoricalLocation},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams{id}/thing", Handler: HandleGetThingByDatastream},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/locations{id}/things", Handler: HandleGetThingsByLocation},

				{OperationType: models.HTTPOperationPost, Path: "/v1.0/things", Handler: HandlePostThing},
				{OperationType: models.HTTPOperationDelete, Path: "/v1.0/things{id}", Handler: HandleDeleteThing},
				{OperationType: models.HTTPOperationPatch, Path: "/v1.0/things{id}", Handler: HandlePatchThing},
				{OperationType: models.HTTPOperationPut, Path: "/v1.0/things{id}", Handler: HandlePutThing},
			},
		},
		entities.EntityTypeDatastream: &endpoint.Endpoint{
			Name:       "Datastreams",
			EntityType: entities.EntityTypeDatastream,
			OutputInfo: true,
			Operations: []models.EndpointOperation{
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams", Handler: HandleGetDatastreams},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams{id}", Handler: HandleGetDatastream},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/observedproperties{id}/datastreams", Handler: HandleGetDatastreamsByObservedProperty},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/observations{id}/datastream", Handler: HandleGetDatastreamByObservation},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/sensors{id}/datastreams", Handler: HandleGetDatastreamsBySensor},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/things{id}/datastreams", Handler: HandleGetDatastreamsByThing},

				{OperationType: models.HTTPOperationPost, Path: "/v1.0/datastreams", Handler: HandlePostDatastream},
				{OperationType: models.HTTPOperationPost, Path: "/v1.0/things{id}/datastreams", Handler: HandlePostDatastreamByThing},
				{OperationType: models.HTTPOperationDelete, Path: "/v1.0/datastreams{id}", Handler: HandleDeleteDatastream},
				{OperationType: models.HTTPOperationPatch, Path: "/v1.0/datastreams{id}", Handler: HandlePatchDatastream},
				{OperationType: models.HTTPOperationPut, Path: "/v1.0/datastreams{id}", Handler: HandlePutDatastream},
			},
		},
		entities.EntityTypeObservedProperty: &endpoint.Endpoint{
			Name:       "ObservedProperties",
			EntityType: entities.EntityTypeObservedProperty,
			OutputInfo: true,
			Operations: []models.EndpointOperation{
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/observedproperties", Handler: HandleGetObservedProperties},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/observedproperties{id}", Handler: HandleGetObservedProperty},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams{id}/observedproperty", Handler: HandleGetObservedPropertyByDatastream},

				{OperationType: models.HTTPOperationPost, Path: "/v1.0/observedproperties", Handler: HandlePostObservedProperty},
				{OperationType: models.HTTPOperationDelete, Path: "/v1.0/observedproperties{id}", Handler: HandleDeleteObservedProperty},
				{OperationType: models.HTTPOperationPatch, Path: "/v1.0/observedproperties{id}", Handler: HandlePatchObservedProperty},
				{OperationType: models.HTTPOperationPut, Path: "/v1.0/observedproperties{id}", Handler: HandlePutObservedProperty},
			},
		},
		entities.EntityTypeLocation: &endpoint.Endpoint{
			Name:       "Locations",
			EntityType: entities.EntityTypeLocation,
			OutputInfo: true,
			Operations: []models.EndpointOperation{
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/locations", Handler: HandleGetLocations},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/locations{id}", Handler: HandleGetLocation},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/historicallocations{id}/locations", Handler: HandleGetLocationsByHistoricalLocations},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/things{id}/locations", Handler: HandleGetLocationsByThing},

				{OperationType: models.HTTPOperationPost, Path: "/v1.0/locations", Handler: HandlePostLocation},
				{OperationType: models.HTTPOperationPost, Path: "/v1.0/things{id}/locations", Handler: HandlePostLocationByThing},
				{OperationType: models.HTTPOperationDelete, Path: "/v1.0/locations{id}", Handler: HandleDeleteLocation},
				{OperationType: models.HTTPOperationPatch, Path: "/v1.0/locations{id}", Handler: HandlePatchLocation},
				{OperationType: models.HTTPOperationPut, Path: "/v1.0/locations{id}", Handler: HandlePutLocation},
			},
		},
		entities.EntityTypeSensor: &endpoint.Endpoint{
			Name:       "Sensors",
			EntityType: entities.EntityTypeSensor,
			OutputInfo: true,
			Operations: []models.EndpointOperation{
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/sensors", Handler: HandleGetSensors},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/sensors{id}", Handler: HandleGetSensor},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams{id}/sensor", Handler: HandleGetSensorByDatastream},

				{OperationType: models.HTTPOperationPost, Path: "/v1.0/sensors", Handler: HandlePostSensors},
				{OperationType: models.HTTPOperationDelete, Path: "/v1.0/sensors{id}", Handler: HandleDeleteSensor},
				{OperationType: models.HTTPOperationPatch, Path: "/v1.0/sensors{id}", Handler: HandlePatchSensor},
				{OperationType: models.HTTPOperationPut, Path: "/v1.0/sensors{id}", Handler: HandlePutSensor},
			},
		},
		entities.EntityTypeObservation: &endpoint.Endpoint{
			Name:       "Observations",
			EntityType: entities.EntityTypeObservation,
			OutputInfo: true,
			Operations: []models.EndpointOperation{
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/observations", Handler: HandleGetObservations},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/observations{id}", Handler: HandleGetObservation},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams{id}/observations", Handler: HandleGetObservationsByDatastream},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/featureofinterest{id}/observations", Handler: HandleGetObservationsByFeatureOfInterest},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/featuresofinterest{id}/observations", Handler: HandleGetObservationsByFeatureOfInterest},

				{OperationType: models.HTTPOperationPost, Path: "/v1.0/observations", Handler: HandlePostObservation},
				{OperationType: models.HTTPOperationPost, Path: "/v1.0/datastreams{id}/observations", Handler: HandlePostObservationByDatastream},
				{OperationType: models.HTTPOperationDelete, Path: "/v1.0/observations{id}", Handler: HandleDeleteObservation},
				{OperationType: models.HTTPOperationPatch, Path: "/v1.0/observations{id}", Handler: HandlePatchObservation},
				{OperationType: models.HTTPOperationPut, Path: "/v1.0/observations{id}", Handler: HandlePutObservation},
			},
		},
		entities.EntityTypeFeatureOfInterest: &endpoint.Endpoint{
			Name:       "FeaturesOfInterest",
			EntityType: entities.EntityTypeFeatureOfInterest,
			OutputInfo: true,
			Operations: []models.EndpointOperation{
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/featuresofinterest", Handler: HandleGetFeatureOfInterests},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/featuresofinterest{id}", Handler: HandleGetFeatureOfInterest},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/observations{id}/featureofinterest", Handler: HandleGetFeatureOfInterestByObservation},
				{OperationType: models.HTTPOperationPost, Path: "/v1.0/featuresofinterest", Handler: HandlePostFeatureOfInterest},
				{OperationType: models.HTTPOperationDelete, Path: "/v1.0/featuresofinterest{id}", Handler: HandleDeleteFeatureOfInterest},
				{OperationType: models.HTTPOperationPatch, Path: "/v1.0/featuresofinterest{id}", Handler: HandlePatchFeatureOfInterest},
				{OperationType: models.HTTPOperationPut, Path: "/v1.0/featuresofinterest{id}", Handler: HandlePutFeatureOfInterest},
			},
		},
		entities.EntityTypeHistoricalLocation: &endpoint.Endpoint{
			Name:       "HistoricalLocations",
			EntityType: entities.EntityTypeHistoricalLocation,
			OutputInfo: true,
			Operations: []models.EndpointOperation{
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/historicallocations", Handler: HandleGetHistoricalLocations},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/historicallocations{id}", Handler: HandleGetHistoricalLocation},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/things{id}/historicallocations", Handler: HandleGetHistoricalLocationsByThing},
				{OperationType: models.HTTPOperationGet, Path: "/v1.0/locations{id}/historicallocations", Handler: HandleGetHistoricalLocationsByLocation},

				{OperationType: models.HTTPOperationPost, Path: "/v1.0/historicallocations", Handler: HandlePostHistoricalLocation},
				{OperationType: models.HTTPOperationDelete, Path: "/v1.0/historicallocations{id}", Handler: HandleDeleteHistoricalLocations},
				{OperationType: models.HTTPOperationPatch, Path: "/v1.0/historicallocations{id}", Handler: HandlePatchHistoricalLocations},
				{OperationType: models.HTTPOperationPut, Path: "/v1.0/historicallocations{id}", Handler: HandlePutHistoricalLocation},
			},
		},
	}
}
