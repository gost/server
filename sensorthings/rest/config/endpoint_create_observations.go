package config

import (
	"fmt"

	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/rest/endpoint"
	"github.com/gost/server/sensorthings/rest/handlers"
)

// CreateCreateObservationsEndpoint constructs the CreateObservations endpoint configuration
func CreateCreateObservationsEndpoint(externalURL string) *endpoint.Endpoint {
	return &endpoint.Endpoint{
		Name:       "CreateObservations",
		OutputInfo: false,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "CreateObservations")),
		Operations: []models.EndpointOperation{
			{OperationType: models.HTTPOperationPost, Path: "/v1.0/createobservations", Handler: handlers.HandlePostCreateObservations},
		},
	}
}
