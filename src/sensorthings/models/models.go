package models

import (
	"net/http"

	"github.com/geodan/gost/src/configuration"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/odata"
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
	GetTopics() *[]Topic

	GetThing(id string, qo *odata.QueryOptions) (*entities.Thing, error)
	GetThingByDatastream(id string, qo *odata.QueryOptions) (*entities.Thing, error)
	GetThings(qo *odata.QueryOptions) (*ArrayResponse, error)
	PostThing(thing *entities.Thing) (*entities.Thing, []error)
	PatchThing(id string, thing *entities.Thing) (*entities.Thing, error)
	DeleteThing(id string) error

	GetLocation(id string, qo *odata.QueryOptions) (*entities.Location, error)
	GetLocations(qo *odata.QueryOptions) (*ArrayResponse, error)
	GetLocationsByThing(thingID string, qo *odata.QueryOptions) (*ArrayResponse, error)
	PostLocation(location *entities.Location) (*entities.Location, []error)
	PostLocationByThing(thingID string, location *entities.Location) (*entities.Location, []error)
	PatchLocation(id string, location *entities.Location) (*entities.Location, error)
	DeleteLocation(id string) error

	GetHistoricalLocation(id string, qo *odata.QueryOptions) (*entities.HistoricalLocation, error)
	GetHistoricalLocations(qo *odata.QueryOptions) (*ArrayResponse, error)
	GetHistoricalLocationsByThing(thingID string, qo *odata.QueryOptions) (*ArrayResponse, error)
	PatchHistoricalLocation(id string, hl *entities.HistoricalLocation) (*entities.HistoricalLocation, error)
	DeleteHistoricalLocation(id string) error

	GetDatastream(id string, qo *odata.QueryOptions) (*entities.Datastream, error)
	GetDatastreams(qo *odata.QueryOptions) (*ArrayResponse, error)
	GetDatastreamsByThing(thingID string, qo *odata.QueryOptions) (*ArrayResponse, error)
	GetDatastreamsBySensor(sensorID string, qo *odata.QueryOptions) (*ArrayResponse, error)
	PostDatastream(datastream *entities.Datastream) (*entities.Datastream, []error)
	PostDatastreamByThing(thingID string, datastream *entities.Datastream) (*entities.Datastream, []error)
	PatchDatastream(id string, datastream *entities.Datastream) (*entities.Datastream, error)
	DeleteDatastream(id string) error

	GetFeatureOfInterest(id string, qo *odata.QueryOptions) (*entities.FeatureOfInterest, error)
	GetFeatureOfInterests(qo *odata.QueryOptions) (*ArrayResponse, error)
	PostFeatureOfInterest(foi *entities.FeatureOfInterest) (*entities.FeatureOfInterest, []error)
	PatchFeatureOfInterest(id string, foi *entities.FeatureOfInterest) (*entities.FeatureOfInterest, error)
	DeleteFeatureOfInterest(id string) error

	GetObservation(id string, qo *odata.QueryOptions) (*entities.Observation, error)
	GetObservations(qo *odata.QueryOptions) (*ArrayResponse, error)
	GetObservationsByDatastream(datastreamID string, qo *odata.QueryOptions) (*ArrayResponse, error)
	PostObservation(observation *entities.Observation) (*entities.Observation, []error)
	PostObservationByDatastream(datastreamID string, observation *entities.Observation) (*entities.Observation, []error)
	PatchObservation(id string, observation *entities.Observation) (*entities.Observation, error)
	DeleteObservation(id string) error

	GetObservedProperty(id string, qo *odata.QueryOptions) (*entities.ObservedProperty, error)
	GetObservedProperties(qo *odata.QueryOptions) (*ArrayResponse, error)
	GetObservedPropertiesByDatastream(datastreamID string, qo *odata.QueryOptions) (*ArrayResponse, error)
	PostObservedProperty(op *entities.ObservedProperty) (*entities.ObservedProperty, []error)
	PatchObservedProperty(id string, op *entities.ObservedProperty) (*entities.ObservedProperty, error)
	DeleteObservedProperty(id string) error

	GetSensor(id string, qo *odata.QueryOptions) (*entities.Sensor, error)
	GetSensors(qo *odata.QueryOptions) (*ArrayResponse, error)
	PostSensor(sensor *entities.Sensor) (*entities.Sensor, []error)
	PatchSensor(id string, sensor *entities.Sensor) (*entities.Sensor, error)
	DeleteSensor(id string) error

	PostHistoricalLocation(thingID string, locationID string) []error

	LinkLocation(thingID string, locationID string) error
}

// Database specifies the operations that the database provider needs to support
type Database interface {
	Start()
	CreateSchema(location string) error

	GetThing(string) (*entities.Thing, error)
	GetThingByDatastream(string) (*entities.Thing, error)
	GetThings() ([]*entities.Thing, error)
	PostThing(*entities.Thing) (*entities.Thing, error)

	GetLocation(string) (*entities.Location, error)
	GetLocations() ([]*entities.Location, error)
	GetLocationsByThing(string) ([]*entities.Location, error)
	PostLocation(*entities.Location) (*entities.Location, error)
	LinkLocation(string, locationID string) error

	GetObservedProperty(string) (*entities.ObservedProperty, error)
	GetObservedProperties() ([]*entities.ObservedProperty, error)
	PostObservedProperty(*entities.ObservedProperty) (*entities.ObservedProperty, error)

	GetSensor(string) (*entities.Sensor, error)
	GetSensors() ([]*entities.Sensor, error)
	PostSensor(*entities.Sensor) (*entities.Sensor, error)

	GetDatastream(string) (*entities.Datastream, error)
	GetDatastreams() ([]*entities.Datastream, error)
	GetDatastreamsByThing(string) ([]*entities.Datastream, error)
	PostDatastream(*entities.Datastream) (*entities.Datastream, error)

	GetFeatureOfInterest(string) (*entities.FeatureOfInterest, error)
	GetFeatureOfInterests() ([]*entities.FeatureOfInterest, error)
	PostFeatureOfInterest(*entities.FeatureOfInterest) (*entities.FeatureOfInterest, error)

	GetObservation(string) (*entities.Observation, error)
	GetObservations() ([]*entities.Observation, error)
	GetObservationsByDatastream(string) ([]*entities.Observation, error)
	PostObservation(*entities.Observation) (*entities.Observation, error)

	GetHistoricalLocation(string) (*entities.HistoricalLocation, error)
	GetHistoricalLocations() ([]*entities.HistoricalLocation, error)
	GetHistoricalLocationsByThing(string) ([]*entities.HistoricalLocation, error)
	PostHistoricalLocation(thingID string, locationID string) error

	ThingExists(thingID int) bool
	LocationExists(thingID int) bool
}

// MQTTClient interface defines the needed MQTT client operations
type MQTTClient interface {
	Start(*API)
	Stop()
	Publish(string, string, byte) //topic, message, qos
}

// Endpoint defines the rest endpoint options
type Endpoint interface {
	GetName() string
	GetURL() string
	GetOperations() []EndpointOperation
	GetSupportedQueryOptions() []odata.QueryOptionType
	GetSupportedExpandParams() []string
	GetSupportedSelectParams() []string
	SupportsQueryOptionType(queryOptionType odata.QueryOptionType) bool
	AreQueryOptionsSupported(queryOptions *odata.QueryOptions) (bool, []error)
}

// HTTPHandler func defines the format of the handler to process the incoming request
type HTTPHandler func(w http.ResponseWriter, r *http.Request, e *Endpoint, a *API)

// EndpointOperation contains the needed information to create an endpoint in the HTTP.Router
type EndpointOperation struct {
	OperationType HTTPOperation
	Path          string //relative path to the endpoint for example: /v1.0/myendpoint/
	Handler       HTTPHandler
}

// HTTPOperation describes the HTTP operation such as GET POST DELETE.
type HTTPOperation string

// HTTPOperation is a "enumeration" of the HTTP operations needed for all endpoints.
const (
	HTTPOperationGet    HTTPOperation = "GET"
	HTTPOperationPost   HTTPOperation = "POST"
	HTTPOperationPatch  HTTPOperation = "PATCH"
	HTTPOperationDelete HTTPOperation = "DELETE"
)

// Topic defines the MQTT PUBLISH topics
type Topic struct {
	Path    string
	Handler MQTTHandler
}

// MQTTHandler func defines the format of the handler to process the incoming MQTT publish message
type MQTTHandler func(a *API, topic string, message []byte)

// MQTTInternalHandler func defines the format of the handler to process the incoming MQTT publish message
type MQTTInternalHandler func(a *API, message []byte, id string)

// VersionInfo describes the version info for the GOST server version and supported SensorThings API version
type VersionInfo struct {
	GostServerVersion GostServerVersion `json:"gostServerVersion"`
	APIVersion        APIVersion        `json:"sensorThingsApiVersion"`
}

// GostServerVersion contains version information on the GOST server
type GostServerVersion struct {
	Version string `json:"version"`
}

// APIVersion contains version information on the supported SensorThings API version
type APIVersion struct {
	Version string `json:"version"`
}

// ArrayResponse is the default response format for sending content back
type ArrayResponse struct {
	Count int          `json:"count,omitempty"`
	Data  *interface{} `json:"value"`
}

// ErrorResponse is the default response format for sending errors back
type ErrorResponse struct {
	Error ErrorContent `json:"error"`
}

// ErrorContent holds information on the error that occurred
type ErrorContent struct {
	StatusText string   `json:"status"`
	StatusCode int      `json:"code"`
	Messages   []string `json:"message"`
}
