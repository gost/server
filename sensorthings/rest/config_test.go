package rest

import (
	"strings"
	"testing"

	"github.com/geodan/gost/configuration"
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
	"github.com/stretchr/testify/assert"
)

func TestCreateEndPoints(t *testing.T) {
	//arrange
	endpoints := CreateEndPoints("http://test.com")

	//assert
	assert.Equal(t, 10, len(endpoints))
}

func TestCreateEndPointVersion(t *testing.T) {
	//arrange
	ve := createVersionEndpoint("http://test.com")

	//assert
	containsVersionPath := containsEndpoint("version", ve.Operations)
	assert.Equal(t, true, containsVersionPath, "Version endpoint needs to contain an endpoint containing the path Version")
}

func containsEndpoint(epName string, eps []models.EndpointOperation) bool {
	for _, o := range eps {
		if strings.Contains(o.Path, epName) {
			return true
		}
	}

	return false
}

type MockAPI struct {
}

// NewAPI Initialise a new SensorThings API
func NewMockAPI() models.API {
	api := &MockAPI{}
	return api
}

func (a *MockAPI) Start()                                                  {}
func (a *MockAPI) GetConfig() *configuration.Config                        { return nil }
func (a *MockAPI) GetAcceptedPaths() []string                              { return []string{} }
func (a *MockAPI) GetBasePathInfo() *models.ArrayResponse                  { return nil }
func (a *MockAPI) GetEndpoints() *map[entities.EntityType]models.Endpoint  { return nil }
func (a *MockAPI) initRest()                                               {}
func (a *MockAPI) GetTopics() *[]models.Topic                              { return nil }
func (a *MockAPI) SetLinks(entity entities.Entity, qo *odata.QueryOptions) {}
func (a *MockAPI) CreateNextLink(count int, incomingURL string, qo *odata.QueryOptions) string {
	return ""
}

func (a *MockAPI) GetThing(id interface{}, qo *odata.QueryOptions, path string) (*entities.Thing, error) {
	return nil, nil
}
func (a *MockAPI) GetThingByDatastream(id interface{}, qo *odata.QueryOptions, path string) (*entities.Thing, error) {
	return nil, nil
}
func (a *MockAPI) GetThingsByLocation(id interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return nil, nil
}
func (a *MockAPI) GetThingByHistoricalLocation(id interface{}, qo *odata.QueryOptions, path string) (*entities.Thing, error) {
	return nil, nil
}
func (a *MockAPI) GetThings(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return nil, nil
}
func (a *MockAPI) PostThing(thing *entities.Thing) (*entities.Thing, []error) { return nil, nil }
func (a *MockAPI) PatchThing(id interface{}, thing *entities.Thing) (*entities.Thing, error) {
	return nil, nil
}
func (a *MockAPI) PutThing(id interface{}, thing *entities.Thing) (*entities.Thing, []error) {
	return nil, nil
}
func (a *MockAPI) DeleteThing(id interface{}) error { return nil }

func (a *MockAPI) GetLocation(id interface{}, qo *odata.QueryOptions, path string) (*entities.Location, error) {
	return nil, nil
}
func (a *MockAPI) GetLocations(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return nil, nil
}
func (a *MockAPI) GetLocationsByHistoricalLocation(hlID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return nil, nil
}
func (a *MockAPI) GetLocationsByThing(thingID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return nil, nil
}
func (a *MockAPI) PostLocation(location *entities.Location) (*entities.Location, []error) {
	return nil, nil
}
func (a *MockAPI) PostLocationByThing(thingID interface{}, location *entities.Location) (*entities.Location, []error) {
	return nil, nil
}
func (a *MockAPI) PatchLocation(id interface{}, location *entities.Location) (*entities.Location, error) {
	return nil, nil
}
func (a *MockAPI) PutLocation(id interface{}, location *entities.Location) (*entities.Location, []error) {
	return nil, nil
}
func (a *MockAPI) DeleteLocation(id interface{}) error { return nil }

func (a *MockAPI) GetHistoricalLocation(id interface{}, qo *odata.QueryOptions, path string) (*entities.HistoricalLocation, error) {
	return nil, nil
}
func (a *MockAPI) GetHistoricalLocations(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return nil, nil
}
func (a *MockAPI) GetHistoricalLocationsByLocation(locationID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return nil, nil
}
func (a *MockAPI) GetHistoricalLocationsByThing(thingID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return nil, nil
}
func (a *MockAPI) PostHistoricalLocation(hl *entities.HistoricalLocation) (*entities.HistoricalLocation, []error) {
	return nil, nil
}
func (a *MockAPI) PutHistoricalLocation(id interface{}, hl *entities.HistoricalLocation) (*entities.HistoricalLocation, []error) {
	return nil, nil
}
func (a *MockAPI) PatchHistoricalLocation(id interface{}, hl *entities.HistoricalLocation) (*entities.HistoricalLocation, error) {
	return nil, nil
}
func (a *MockAPI) DeleteHistoricalLocation(id interface{}) error { return nil }

func (a *MockAPI) GetDatastream(id interface{}, qo *odata.QueryOptions, path string) (*entities.Datastream, error) {
	return nil, nil
}
func (a *MockAPI) GetDatastreams(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return nil, nil
}
func (a *MockAPI) GetDatastreamByObservation(id interface{}, qo *odata.QueryOptions, path string) (*entities.Datastream, error) {
	return nil, nil
}
func (a *MockAPI) GetDatastreamsByThing(thingID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return nil, nil
}
func (a *MockAPI) GetDatastreamsBySensor(sensorID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return nil, nil
}
func (a *MockAPI) GetDatastreamsByObservedProperty(sensorID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return nil, nil
}
func (a *MockAPI) PostDatastream(datastream *entities.Datastream) (*entities.Datastream, []error) {
	return nil, nil
}
func (a *MockAPI) PostDatastreamByThing(thingID interface{}, datastream *entities.Datastream) (*entities.Datastream, []error) {
	return nil, nil
}
func (a *MockAPI) PatchDatastream(id interface{}, datastream *entities.Datastream) (*entities.Datastream, error) {
	return nil, nil
}
func (a *MockAPI) PutDatastream(id interface{}, datastream *entities.Datastream) (*entities.Datastream, []error) {
	return nil, nil
}
func (a *MockAPI) DeleteDatastream(id interface{}) error { return nil }

func (a *MockAPI) GetFeatureOfInterest(id interface{}, qo *odata.QueryOptions, path string) (*entities.FeatureOfInterest, error) {
	return nil, nil
}
func (a *MockAPI) GetFeatureOfInterestByObservation(id interface{}, qo *odata.QueryOptions, path string) (*entities.FeatureOfInterest, error) {
	return nil, nil
}
func (a *MockAPI) GetFeatureOfInterests(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return nil, nil
}
func (a *MockAPI) PostFeatureOfInterest(foi *entities.FeatureOfInterest) (*entities.FeatureOfInterest, []error) {
	return nil, nil
}
func (a *MockAPI) PatchFeatureOfInterest(id interface{}, foi *entities.FeatureOfInterest) (*entities.FeatureOfInterest, error) {
	return nil, nil
}
func (a *MockAPI) PutFeatureOfInterest(id interface{}, foi *entities.FeatureOfInterest) (*entities.FeatureOfInterest, []error) {
	return nil, nil
}
func (a *MockAPI) DeleteFeatureOfInterest(id interface{}) error { return nil }

func (a *MockAPI) GetObservation(id interface{}, qo *odata.QueryOptions, path string) (*entities.Observation, error) {
	return nil, nil
}
func (a *MockAPI) GetObservations(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return nil, nil
}
func (a *MockAPI) GetObservationsByDatastream(datastreamID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return nil, nil
}
func (a *MockAPI) GetObservationsByFeatureOfInterest(foiID interface{}, qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return nil, nil
}
func (a *MockAPI) PostObservation(observation *entities.Observation) (*entities.Observation, []error) {
	return nil, nil
}
func (a *MockAPI) PostObservationByDatastream(datastreamID interface{}, observation *entities.Observation) (*entities.Observation, []error) {
	return nil, nil
}
func (a *MockAPI) PatchObservation(id interface{}, observation *entities.Observation) (*entities.Observation, error) {
	return nil, nil
}
func (a *MockAPI) PutObservation(id interface{}, observation *entities.Observation) (*entities.Observation, []error) {
	return nil, nil
}
func (a *MockAPI) DeleteObservation(id interface{}) error { return nil }

func (a *MockAPI) GetObservedProperty(id interface{}, qo *odata.QueryOptions, path string) (*entities.ObservedProperty, error) {
	return nil, nil
}
func (a *MockAPI) GetObservedProperties(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return nil, nil
}
func (a *MockAPI) GetObservedPropertyByDatastream(datastreamID interface{}, qo *odata.QueryOptions, path string) (*entities.ObservedProperty, error) {
	return nil, nil
}
func (a *MockAPI) PostObservedProperty(op *entities.ObservedProperty) (*entities.ObservedProperty, []error) {
	return nil, nil
}
func (a *MockAPI) PatchObservedProperty(id interface{}, op *entities.ObservedProperty) (*entities.ObservedProperty, error) {
	return nil, nil
}
func (a *MockAPI) PutObservedProperty(id interface{}, op *entities.ObservedProperty) (*entities.ObservedProperty, []error) {
	return nil, nil
}
func (a *MockAPI) DeleteObservedProperty(id interface{}) error { return nil }

func (a *MockAPI) GetSensor(id interface{}, qo *odata.QueryOptions, path string) (*entities.Sensor, error) {
	return nil, nil
}
func (a *MockAPI) GetSensorByDatastream(id interface{}, qo *odata.QueryOptions, path string) (*entities.Sensor, error) {
	return nil, nil
}
func (a *MockAPI) GetSensors(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	return nil, nil
}
func (a *MockAPI) PostSensor(sensor *entities.Sensor) (*entities.Sensor, []error) { return nil, nil }
func (a *MockAPI) PatchSensor(id interface{}, sensor *entities.Sensor) (*entities.Sensor, error) {
	return nil, nil
}
func (a *MockAPI) DeleteSensor(id interface{}) error { return nil }
func (a *MockAPI) PutSensor(id interface{}, sensor *entities.Sensor) (*entities.Sensor, []error) {
	return nil, nil
}

func (a *MockAPI) LinkLocation(thingID interface{}, locationID interface{}) error { return nil }

func (a *MockAPI) GetVersionInfo() *models.VersionInfo {
	versionInfo := models.VersionInfo{
		GostServerVersion: models.GostServerVersion{Version: configuration.ServerVersion},
		APIVersion:        models.APIVersion{Version: configuration.SensorThingsAPIVersion},
	}

	return &versionInfo
}
