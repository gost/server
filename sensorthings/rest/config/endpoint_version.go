package config

import (
	"fmt"

	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/rest/endpoint"
	"github.com/gost/server/sensorthings/rest/handlers"
)

// CreateVersionEndpoint creates the Version endpoint configuration
func CreateVersionEndpoint(externalURL string) *endpoint.Endpoint {
	return &endpoint.Endpoint{
		Name:       "Version",
		OutputInfo: false,
		URL:        fmt.Sprintf("%s/%s", externalURL, "Version"),
		Operations: []models.EndpointOperation{
			{OperationType: models.HTTPOperationGet, Path: "/version", Handler: handlers.HandleVersion},
		},
	}
}
