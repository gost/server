package config

import (
	"fmt"

	entities "github.com/gost/core"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/rest/endpoint"
	"github.com/gost/server/sensorthings/rest/handlers"
)

// CreateLocationsEndpoint constructs the Locations endpoint configuration
func CreateLocationsEndpoint(externalURL string) *endpoint.Endpoint {
	return &endpoint.Endpoint{
		Name:       "Locations",
		EntityType: entities.EntityTypeLocation,
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "Locations")),
		SupportedExpandParams: []string{
			"things",
			"historicallocations",
		},
		SupportedSelectParams: []string{
			"id",
			"name",
			"description",
			"encodingtype",
			"location",
			"things",
			"historicallocations",
		},
		Operations: []models.EndpointOperation{
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/locations", Handler: handlers.HandleGetLocations},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/locations{id}", Handler: handlers.HandleGetLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/historicallocations{id}/locations", Handler: handlers.HandleGetLocationsByHistoricalLocations},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/historicallocations{id}/locations/{params}", Handler: handlers.HandleGetLocationsByHistoricalLocations},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/locations{id}/{params}", Handler: handlers.HandleGetLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/locations{id}/{params}/$value", Handler: handlers.HandleGetLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/things{id}/locations", Handler: handlers.HandleGetLocationsByThing},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/things{id}/locations/{params}", Handler: handlers.HandleGetLocationsByThing},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/locations/{params}", Handler: handlers.HandleGetLocations},

			{OperationType: models.HTTPOperationPost, Path: "/v1.0/locations", Handler: handlers.HandlePostLocation},
			{OperationType: models.HTTPOperationPost, Path: "/v1.0/things{id}/locations", Handler: handlers.HandlePostLocationByThing},
			{OperationType: models.HTTPOperationDelete, Path: "/v1.0/locations{id}", Handler: handlers.HandleDeleteLocation},
			{OperationType: models.HTTPOperationPatch, Path: "/v1.0/locations{id}", Handler: handlers.HandlePatchLocation},
			{OperationType: models.HTTPOperationPut, Path: "/v1.0/locations{id}", Handler: handlers.HandlePutLocation},

			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/locations", Handler: handlers.HandleGetLocations},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/locations{id}", Handler: handlers.HandleGetLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/historicallocations{id}/locations", Handler: handlers.HandleGetLocationsByHistoricalLocations},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/things{id}/locations", Handler: handlers.HandleGetLocationsByThing},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/things{id}/locations/{params}", Handler: handlers.HandleGetLocationsByThing},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/historicallocations{id}/locations/{params}", Handler: handlers.HandleGetLocationsByHistoricalLocations},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/locations{id}/{params}", Handler: handlers.HandleGetLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/locations{id}/{params}/$value", Handler: handlers.HandleGetLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/locations/{params}", Handler: handlers.HandleGetLocations},

			{OperationType: models.HTTPOperationPost, Path: "/v1.0/{c:.*}/locations", Handler: handlers.HandlePostLocation},
			{OperationType: models.HTTPOperationPost, Path: "/v1.0/{c:.*}/things{id}/locations", Handler: handlers.HandlePostLocationByThing},
			{OperationType: models.HTTPOperationDelete, Path: "/v1.0/{c:.*}/locations{id}", Handler: handlers.HandleDeleteLocation},
			{OperationType: models.HTTPOperationPatch, Path: "/v1.0/{c:.*}/locations{id}", Handler: handlers.HandlePatchLocation},
			{OperationType: models.HTTPOperationPut, Path: "/v1.0/{c:.*}/locations{id}", Handler: handlers.HandlePutLocation},
		},
	}
}
