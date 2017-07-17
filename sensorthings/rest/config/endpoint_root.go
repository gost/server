package config

import (
	"fmt"

	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/rest/endpoint"
	"github.com/gost/server/sensorthings/rest/handlers"
)

func CreateRootEndpoint(externalURL string) *endpoint.Endpoint {
	return &endpoint.Endpoint{
		Name:       "Root",
		OutputInfo: false,
		URL:        fmt.Sprintf("%s/%s", externalURL, "v1.0"),
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0", handlers.HandleAPIRoot},
		},
	}
}
