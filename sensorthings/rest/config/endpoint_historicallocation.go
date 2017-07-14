package config

import (
	"fmt"

	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/rest/endpoint"
	"github.com/geodan/gost/sensorthings/rest/handlers"
)

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
			{models.HTTPOperationGet, "/v1.0/historicallocations", handlers.HandleGetHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}", handlers.HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/things{id}/historicallocations", handlers.HandleGetHistoricalLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/things{id}/historicallocations/{params}", handlers.HandleGetHistoricalLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/locations{id}/historicallocations", handlers.HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/locations{id}/historicallocations/{params}", handlers.HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/locations{id}/historicallocations/{params}/$value", handlers.HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/{params}", handlers.HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/{params}/$value", handlers.HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/historicallocations/{params}", handlers.HandleGetHistoricalLocations},

			{models.HTTPOperationPost, "/v1.0/historicallocations", handlers.HandlePostHistoricalLocation},
			{models.HTTPOperationDelete, "/v1.0/historicallocations{id}", handlers.HandleDeleteHistoricalLocations},
			{models.HTTPOperationPatch, "/v1.0/historicallocations{id}", handlers.HandlePatchHistoricalLocations},
			{models.HTTPOperationPut, "/v1.0/historicallocations{id}", handlers.HandlePutHistoricalLocation},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations", handlers.HandleGetHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}", handlers.HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/things{id}/historicallocations", handlers.HandleGetHistoricalLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/things{id}/historicallocations/{params}", handlers.HandleGetHistoricalLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/locations{id}/historicallocations", handlers.HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/locations{id}/historicallocations/{params}", handlers.HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/locations{id}/historicallocations/{params}/$value", handlers.HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}/{params}", handlers.HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}/{params}/$value", handlers.HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations/{params}", handlers.HandleGetHistoricalLocations},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/historicallocations", handlers.HandlePostHistoricalLocation},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/historicallocations{id}", handlers.HandleDeleteHistoricalLocations},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/historicallocations{id}", handlers.HandlePatchHistoricalLocations},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/historicallocations{id}", handlers.HandlePutHistoricalLocation},
		},
	}
}
