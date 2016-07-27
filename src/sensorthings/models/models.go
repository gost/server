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
	Start()
	GetConfig() *configuration.Config

	GetVersionInfo() *VersionInfo
	GetBasePathInfo() *ArrayResponse
	GetEndpoints() *[]Endpoint
	GetTopics() *[]Topic

	GetThing(id interface{}, qo *odata.QueryOptions, path string) (*entities.Thing, error)
	GetThingByDatastream(id interface{}, qo *odata.QueryOptions, path string) (*entities.Thing, error)
	GetThingsByLocation(id interface{}, qo *odata.QueryOptions, path string) (*ArrayResponse, error)
	GetThingByHistoricalLocation(id interface{}, qo *odata.QueryOptions, path string) (*entities.Thing, error)
	GetThings(qo *odata.QueryOptions, path string) (*ArrayResponse, error)
	PostThing(thing *entities.Thing) (*entities.Thing, []error)
	PatchThing(id interface{}, thing *entities.Thing) (*entities.Thing, error)
	PutThing(id interface{}, thing *entities.Thing) (*entities.Thing, []error)
	DeleteThing(id interface{}) error

	GetLocation(id interface{}, qo *odata.QueryOptions, path string) (*entities.Location, error)
	GetLocations(qo *odata.QueryOptions, path string) (*ArrayResponse, error)
	GetLocationsByHistoricalLocation(hlID interface{}, qo *odata.QueryOptions, path string) (*ArrayResponse, error)
	GetLocationsByThing(thingID interface{}, qo *odata.QueryOptions, path string) (*ArrayResponse, error)
	PostLocation(location *entities.Location) (*entities.Location, []error)
	PostLocationByThing(thingID interface{}, location *entities.Location) (*entities.Location, []error)
	PatchLocation(id interface{}, location *entities.Location) (*entities.Location, error)
	PutLocation(id interface{}, location *entities.Location) (*entities.Location, []error)
	DeleteLocation(id interface{}) error

	GetHistoricalLocation(id interface{}, qo *odata.QueryOptions, path string) (*entities.HistoricalLocation, error)
	GetHistoricalLocations(qo *odata.QueryOptions, path string) (*ArrayResponse, error)
	GetHistoricalLocationsByLocation(locationID interface{}, qo *odata.QueryOptions, path string) (*ArrayResponse, error)
	GetHistoricalLocationsByThing(thingID interface{}, qo *odata.QueryOptions, path string) (*ArrayResponse, error)
	PostHistoricalLocation(hl *entities.HistoricalLocation) (*entities.HistoricalLocation, []error)
	PatchHistoricalLocation(id interface{}, hl *entities.HistoricalLocation) (*entities.HistoricalLocation, error)
	DeleteHistoricalLocation(id interface{}) error

	GetDatastream(id interface{}, qo *odata.QueryOptions, path string) (*entities.Datastream, error)
	GetDatastreams(qo *odata.QueryOptions, path string) (*ArrayResponse, error)
	GetDatastreamByObservation(id interface{}, qo *odata.QueryOptions, path string) (*entities.Datastream, error)
	GetDatastreamsByThing(thingID interface{}, qo *odata.QueryOptions, path string) (*ArrayResponse, error)
	GetDatastreamsBySensor(sensorID interface{}, qo *odata.QueryOptions, path string) (*ArrayResponse, error)
	GetDatastreamsByObservedProperty(sensorID interface{}, qo *odata.QueryOptions, path string) (*ArrayResponse, error)
	PostDatastream(datastream *entities.Datastream) (*entities.Datastream, []error)
	PostDatastreamByThing(thingID interface{}, datastream *entities.Datastream) (*entities.Datastream, []error)
	PatchDatastream(id interface{}, datastream *entities.Datastream) (*entities.Datastream, error)
	PutDatastream(id interface{}, datastream *entities.Datastream) (*entities.Datastream, []error)
	DeleteDatastream(id interface{}) error

	GetFeatureOfInterest(id interface{}, qo *odata.QueryOptions, path string) (*entities.FeatureOfInterest, error)
	GetFeatureOfInterestByObservation(id interface{}, qo *odata.QueryOptions, path string) (*entities.FeatureOfInterest, error)
	GetFeatureOfInterests(qo *odata.QueryOptions, path string) (*ArrayResponse, error)
	PostFeatureOfInterest(foi *entities.FeatureOfInterest) (*entities.FeatureOfInterest, []error)
	PatchFeatureOfInterest(id interface{}, foi *entities.FeatureOfInterest) (*entities.FeatureOfInterest, error)
	DeleteFeatureOfInterest(id interface{}) error

	GetObservation(id interface{}, qo *odata.QueryOptions, path string) (*entities.Observation, error)
	GetObservations(qo *odata.QueryOptions, path string) (*ArrayResponse, error)
	GetObservationsByDatastream(datastreamID interface{}, qo *odata.QueryOptions, path string) (*ArrayResponse, error)
	GetObservationsByFeatureOfInterest(foiID interface{}, qo *odata.QueryOptions, path string) (*ArrayResponse, error)
	PostObservation(observation *entities.Observation) (*entities.Observation, []error)
	PostObservationByDatastream(datastreamID interface{}, observation *entities.Observation) (*entities.Observation, []error)
	PatchObservation(id interface{}, observation *entities.Observation) (*entities.Observation, error)
	DeleteObservation(id interface{}) error

	GetObservedProperty(id interface{}, qo *odata.QueryOptions, path string) (*entities.ObservedProperty, error)
	GetObservedProperties(qo *odata.QueryOptions, path string) (*ArrayResponse, error)
	GetObservedPropertyByDatastream(datastreamID interface{}, qo *odata.QueryOptions, path string) (*entities.ObservedProperty, error)
	PostObservedProperty(op *entities.ObservedProperty) (*entities.ObservedProperty, []error)
	PatchObservedProperty(id interface{}, op *entities.ObservedProperty) (*entities.ObservedProperty, error)
	DeleteObservedProperty(id interface{}) error

	GetSensor(id interface{}, qo *odata.QueryOptions, path string) (*entities.Sensor, error)
	GetSensorByDatastream(id interface{}, qo *odata.QueryOptions, path string) (*entities.Sensor, error)
	GetSensors(qo *odata.QueryOptions, path string) (*ArrayResponse, error)
	PostSensor(sensor *entities.Sensor) (*entities.Sensor, []error)
	PatchSensor(id interface{}, sensor *entities.Sensor) (*entities.Sensor, error)
	DeleteSensor(id interface{}) error

	LinkLocation(thingID interface{}, locationID interface{}) error
}

// Database specifies the operations that the database provider needs to support
type Database interface {
	Start()
	CreateSchema(location string) error

	GetThing(id interface{}, qo *odata.QueryOptions) (*entities.Thing, error)
	GetThingByDatastream(id interface{}, qo *odata.QueryOptions) (t *entities.Thing, e error)
	GetThingsByLocation(id interface{}, qo *odata.QueryOptions) (t []*entities.Thing, count int, e error)
	GetThingByHistoricalLocation(id interface{}, qo *odata.QueryOptions) (t *entities.Thing, e error)
	GetThings(qo *odata.QueryOptions) (t []*entities.Thing, count int, e error)
	PostThing(*entities.Thing) (*entities.Thing, error)
	PatchThing(interface{}, *entities.Thing) (*entities.Thing, error)
	PutThing(interface{}, *entities.Thing) (*entities.Thing, error)
	DeleteThing(id interface{}) error

	GetLocation(id interface{}, qo *odata.QueryOptions) (*entities.Location, error)
	GetLocations(qo *odata.QueryOptions) (l []*entities.Location, count int, e error)
	GetLocationsByHistoricalLocation(id interface{}, qo *odata.QueryOptions) (l []*entities.Location, count int, e error)
	GetLocationsByThing(id interface{}, qo *odata.QueryOptions) (l []*entities.Location, count int, e error)
	GetLocationByDatastreamID(id interface{}) (*entities.Location, error)
	PostLocation(*entities.Location) (*entities.Location, error)
	LinkLocation(id interface{}, locationID interface{}) error
	PatchLocation(interface{}, *entities.Location) (*entities.Location, error)
	DeleteLocation(id interface{}) error
	PutLocation(interface{}, *entities.Location) (*entities.Location, error)

	GetObservedProperty(id interface{}, qo *odata.QueryOptions) (*entities.ObservedProperty, error)
	GetObservedPropertyByDatastream(id interface{}, qo *odata.QueryOptions) (*entities.ObservedProperty, error)
	GetObservedProperties(qo *odata.QueryOptions) (o []*entities.ObservedProperty, count int, e error)
	PostObservedProperty(*entities.ObservedProperty) (*entities.ObservedProperty, error)
	PatchObservedProperty(interface{}, *entities.ObservedProperty) (*entities.ObservedProperty, error)
	DeleteObservedProperty(id interface{}) error

	GetSensor(id interface{}, qo *odata.QueryOptions) (*entities.Sensor, error)
	GetSensorByDatastream(id interface{}, qo *odata.QueryOptions) (*entities.Sensor, error)
	GetSensors(qo *odata.QueryOptions) (s []*entities.Sensor, count int, e error)
	PostSensor(*entities.Sensor) (*entities.Sensor, error)
	PatchSensor(interface{}, *entities.Sensor) (*entities.Sensor, error)
	DeleteSensor(id interface{}) error

	GetDatastream(id interface{}, qo *odata.QueryOptions) (*entities.Datastream, error)
	GetDatastreams(qo *odata.QueryOptions) (d []*entities.Datastream, count int, e error)
	GetDatastreamByObservation(id interface{}, qo *odata.QueryOptions) (*entities.Datastream, error)
	GetDatastreamsByThing(id interface{}, qo *odata.QueryOptions) (d []*entities.Datastream, count int, e error)
	GetDatastreamsBySensor(id interface{}, qo *odata.QueryOptions) (d []*entities.Datastream, count int, e error)
	GetDatastreamsByObservedProperty(id interface{}, qo *odata.QueryOptions) (d []*entities.Datastream, count int, e error)
	PostDatastream(*entities.Datastream) (*entities.Datastream, error)
	PatchDatastream(interface{}, *entities.Datastream) (*entities.Datastream, error)
	DeleteDatastream(id interface{}) error
	DatastreamExists(int) bool
	PutDatastream(interface{}, *entities.Datastream) (*entities.Datastream, error)

	GetFeatureOfInterest(id interface{}, qo *odata.QueryOptions) (*entities.FeatureOfInterest, error)
	GetFeatureOfInterestByLocationID(id interface{}) (*entities.FeatureOfInterest, error)
	GetFeatureOfInterestByObservation(id interface{}, qo *odata.QueryOptions) (*entities.FeatureOfInterest, error)
	GetFeatureOfInterests(qo *odata.QueryOptions) (f []*entities.FeatureOfInterest, count int, e error)
	PostFeatureOfInterest(*entities.FeatureOfInterest) (*entities.FeatureOfInterest, error)
	PatchFeatureOfInterest(interface{}, *entities.FeatureOfInterest) (*entities.FeatureOfInterest, error)
	DeleteFeatureOfInterest(id interface{}) error

	GetObservation(id interface{}, qo *odata.QueryOptions) (*entities.Observation, error)
	GetObservations(qo *odata.QueryOptions) (o []*entities.Observation, count int, e error)
	GetObservationsByDatastream(id interface{}, qo *odata.QueryOptions) (o []*entities.Observation, count int, e error)
	GetObservationsByFeatureOfInterest(id interface{}, qo *odata.QueryOptions) (o []*entities.Observation, count int, e error)
	PostObservation(*entities.Observation) (*entities.Observation, error)
	PatchObservation(interface{}, *entities.Observation) (*entities.Observation, error)
	DeleteObservation(id interface{}) error

	GetHistoricalLocation(id interface{}, qo *odata.QueryOptions) (*entities.HistoricalLocation, error)
	GetHistoricalLocations(qo *odata.QueryOptions) (h []*entities.HistoricalLocation, count int, e error)
	GetHistoricalLocationsByLocation(id interface{}, qo *odata.QueryOptions) (h []*entities.HistoricalLocation, count int, e error)
	GetHistoricalLocationsByThing(id interface{}, qo *odata.QueryOptions) (h []*entities.HistoricalLocation, count int, e error)
	PostHistoricalLocation(*entities.HistoricalLocation) (*entities.HistoricalLocation, error)
	PatchHistoricalLocation(interface{}, *entities.HistoricalLocation) (*entities.HistoricalLocation, error)
	DeleteHistoricalLocation(id interface{}) error

	ThingExists(thingID interface{}) bool
	LocationExists(thingID interface{}) bool
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
	ShowOutputInfo() bool
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
	HTTPOperationPut    HTTPOperation = "PUT"
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
	Count    int          `json:"count,omitempty"`
	NextLink string       `json:"@iot.nextLink,omitempty"`
	Data     *interface{} `json:"value"`
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
