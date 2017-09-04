package config

import (
	"fmt"

	entities "github.com/gost/core"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/rest/endpoint"
	"github.com/gost/server/sensorthings/rest/handlers"
)

// CreateHistoricalLocationsEndpoint constructs the HistoricalLocations endpoint configuration
func CreateHistoricalLocationsEndpoint(externalURL string) *endpoint.Endpoint {
	return &endpoint.Endpoint{
		Name:       "HistoricalLocations",
		EntityType: entities.EntityTypeHistoricalLocation,
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "HistoricalLocations")),
		SupportedExpandParams: []string{
			"locations",
			"thing",
		},
		SupportedSelectParams: []string{
			"id",
			"time",
		},
		Operations: []models.EndpointOperation{
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/historicallocations", Handler: handlers.HandleGetHistoricalLocations},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/historicallocations{id}", Handler: handlers.HandleGetHistoricalLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/things{id}/historicallocations", Handler: handlers.HandleGetHistoricalLocationsByThing},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/things{id}/historicallocations/{params}", Handler: handlers.HandleGetHistoricalLocationsByThing},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/locations{id}/historicallocations", Handler: handlers.HandleGetHistoricalLocationsByLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/locations{id}/historicallocations/{params}", Handler: handlers.HandleGetHistoricalLocationsByLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/locations{id}/historicallocations/{params}/$value", Handler: handlers.HandleGetHistoricalLocationsByLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/historicallocations{id}/{params}", Handler: handlers.HandleGetHistoricalLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/historicallocations{id}/{params}/$value", Handler: handlers.HandleGetHistoricalLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/historicallocations/{params}", Handler: handlers.HandleGetHistoricalLocations},

			{OperationType: models.HTTPOperationPost, Path: "/v1.0/historicallocations", Handler: handlers.HandlePostHistoricalLocation},
			{OperationType: models.HTTPOperationDelete, Path: "/v1.0/historicallocations{id}", Handler: handlers.HandleDeleteHistoricalLocations},
			{OperationType: models.HTTPOperationPatch, Path: "/v1.0/historicallocations{id}", Handler: handlers.HandlePatchHistoricalLocations},
			{OperationType: models.HTTPOperationPut, Path: "/v1.0/historicallocations{id}", Handler: handlers.HandlePutHistoricalLocation},

			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/historicallocations", Handler: handlers.HandleGetHistoricalLocations},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/historicallocations{id}", Handler: handlers.HandleGetHistoricalLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/things{id}/historicallocations", Handler: handlers.HandleGetHistoricalLocationsByThing},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/things{id}/historicallocations/{params}", Handler: handlers.HandleGetHistoricalLocationsByThing},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/locations{id}/historicallocations", Handler: handlers.HandleGetHistoricalLocationsByLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/locations{id}/historicallocations/{params}", Handler: handlers.HandleGetHistoricalLocationsByLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/locations{id}/historicallocations/{params}/$value", Handler: handlers.HandleGetHistoricalLocationsByLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/historicallocations{id}/{params}", Handler: handlers.HandleGetHistoricalLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/historicallocations{id}/{params}/$value", Handler: handlers.HandleGetHistoricalLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/historicallocations/{params}", Handler: handlers.HandleGetHistoricalLocations},

			{OperationType: models.HTTPOperationPost, Path: "/v1.0/{c:.*}/historicallocations", Handler: handlers.HandlePostHistoricalLocation},
			{OperationType: models.HTTPOperationDelete, Path: "/v1.0/{c:.*}/historicallocations{id}", Handler: handlers.HandleDeleteHistoricalLocations},
			{OperationType: models.HTTPOperationPatch, Path: "/v1.0/{c:.*}/historicallocations{id}", Handler: handlers.HandlePatchHistoricalLocations},
			{OperationType: models.HTTPOperationPut, Path: "/v1.0/{c:.*}/historicallocations{id}", Handler: handlers.HandlePutHistoricalLocation},
		},
	}
}
