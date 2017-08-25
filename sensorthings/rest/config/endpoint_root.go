package config

import (
	"fmt"

	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/rest/endpoint"
	"github.com/gost/server/sensorthings/rest/handlers"
)

// CreateRootEndpoint creates the Root endpoint configuration
func CreateRootEndpoint(externalURL string) *endpoint.Endpoint {
	return &endpoint.Endpoint{
		Name:       "Root",
		OutputInfo: false,
		URL:        fmt.Sprintf("%s/%s", externalURL, "v1.0"),
		Operations: []models.EndpointOperation{
			{OperationType: models.HTTPOperationGet, Path: "/v1.0", Handler: handlers.HandleAPIRoot},
		},
	}
}
