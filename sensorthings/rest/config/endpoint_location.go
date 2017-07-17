package config

import (
	"fmt"

	"github.com/gost/server/sensorthings/entities"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/rest/endpoint"
	"github.com/gost/server/sensorthings/rest/handlers"
)

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
			{models.HTTPOperationGet, "/v1.0/locations", handlers.HandleGetLocations},
			{models.HTTPOperationGet, "/v1.0/locations{id}", handlers.HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/locations", handlers.HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/locations/{params}", handlers.HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/locations{id}/{params}", handlers.HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/locations{id}/{params}/$value", handlers.HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/things{id}/locations", handlers.HandleGetLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/things{id}/locations/{params}", handlers.HandleGetLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/locations/{params}", handlers.HandleGetLocations},

			{models.HTTPOperationPost, "/v1.0/locations", handlers.HandlePostLocation},
			{models.HTTPOperationPost, "/v1.0/things{id}/locations", handlers.HandlePostLocationByThing},
			{models.HTTPOperationDelete, "/v1.0/locations{id}", handlers.HandleDeleteLocation},
			{models.HTTPOperationPatch, "/v1.0/locations{id}", handlers.HandlePatchLocation},
			{models.HTTPOperationPut, "/v1.0/locations{id}", handlers.HandlePutLocation},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/locations", handlers.HandleGetLocations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/locations{id}", handlers.HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}/locations", handlers.HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/things{id}/locations", handlers.HandleGetLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/things{id}/locations/{params}", handlers.HandleGetLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}/locations/{params}", handlers.HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/locations{id}/{params}", handlers.HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/locations{id}/{params}/$value", handlers.HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/locations/{params}", handlers.HandleGetLocations},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/locations", handlers.HandlePostLocation},
			{models.HTTPOperationPost, "/v1.0/{c:.*}/things{id}/locations", handlers.HandlePostLocationByThing},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/locations{id}", handlers.HandleDeleteLocation},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/locations{id}", handlers.HandlePatchLocation},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/locations{id}", handlers.HandlePutLocation},
		},
	}
}
