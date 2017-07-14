package config

import (
	"fmt"

	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/rest/endpoint"
	"github.com/geodan/gost/sensorthings/rest/handlers"
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
