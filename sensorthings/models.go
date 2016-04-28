package sensorthings

import "github.com/geodan/gost/sensorthings/entities"

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
