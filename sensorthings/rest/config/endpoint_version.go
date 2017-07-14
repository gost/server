package config

import (
	"fmt"

	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/rest/endpoint"
	"github.com/geodan/gost/sensorthings/rest/handlers"
)

func CreateVersionEndpoint(externalURL string) *endpoint.Endpoint {
	return &endpoint.Endpoint{
		Name:       "Version",
		OutputInfo: false,
		URL:        fmt.Sprintf("%s/%s", externalURL, "Version"),
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/version", handlers.HandleVersion},
		},
	}
}
