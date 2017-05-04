package rest

import (
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
)

// Endpoint contains all information for creating and handling a main SensorThings endpoint.
// A SensorThings endpoint contains multiple EndpointOperations
// Endpoint can be marshalled to JSON for returning endpoint information requested
// by the user: http://www.sensorup.com/docs/#resource-path
type Endpoint struct {
	Name                  string                     `json:"name"` // Name of the endpoint
	URL                   string                     `json:"url"`  // External URL to the endpoint
	EntityType            entities.EntityType        `json:"-"`
	OutputInfo            bool                       `json:"-"` //Output when BasePathInfo is requested by the user
	Operations            []models.EndpointOperation `json:"-"`
	SupportedExpandParams []string                   `json:"-"`
	SupportedSelectParams []string                   `json:"-"`
}

// GetName returns the endpoint name
func (e *Endpoint) GetName() string {
	return e.Name
}

// ShowOutputInfo returns true if the endpoint should output his info when BasePathInfo is requested
func (e *Endpoint) ShowOutputInfo() bool {
	return e.OutputInfo
}

// GetURL returns the external url
func (e *Endpoint) GetURL() string {
	return e.URL
}

// GetOperations returns all operations for this endpoint such as GET, POST
func (e *Endpoint) GetOperations() []models.EndpointOperation {
	return e.Operations
}

// GetSupportedExpandParams returns which entities can be expanded
func (e *Endpoint) GetSupportedExpandParams() []string {
	return e.SupportedExpandParams
}

// GetSupportedSelectParams returns the supported select parameters for this endpoint
func (e *Endpoint) GetSupportedSelectParams() []string {
	return e.SupportedSelectParams
}
