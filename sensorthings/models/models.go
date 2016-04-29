package models

import (
	"github.com/geodan/gost/configuration"
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/odata"
	"net/http"
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

	GetThing(id string, qo *odata.QueryOptions) (*entities.Thing, error)
	GetThings(qo *odata.QueryOptions) (*ArrayResponse, error)
	PostThing(thing entities.Thing) (*entities.Thing, []error)
	DeleteThing(id string)
	PatchThing(thing entities.Thing)

	GetLocation(id string) *entities.Location
	GetLocations() *ArrayResponse
	PostLocation(location entities.Location, thingID string) (*entities.Location, []error)
	DeleteLocation(id string)
	PatchLocation(id string)

	PostHistoricalLocation(thingID string, locationID string) error
	LinkLocation(thingID string, locationID string) error
}

// Database specifies the operations that the database provider needs to support
type Database interface {
	Start()
	CreateSchema()

	GetThing(string) (*entities.Thing, error)
	GetThings() ([]*entities.Thing, error)
	PostThing(thing entities.Thing) (*entities.Thing, error)

	PostLocation(location entities.Location) (*entities.Location, error)
	LinkLocation(thingID string, locationID string) error

	PostHistoricalLocation(thingID string, locationID string) error

	ThingExists(thingID int) bool
	LocationExists(thingID int) bool
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
	Count *int         `json:"count,omitempty"`
	Data  *interface{} `json:"value"`
}
