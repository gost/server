package config

import (
	"fmt"

	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/rest/endpoint"
	"github.com/geodan/gost/sensorthings/rest/handlers"
)

func CreateThingsEndpoint(externalURL string) *endpoint.Endpoint {
	return &endpoint.Endpoint{
		Name:       "Things",
		EntityType: entities.EntityTypeThing,
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "Things")),
		SupportedExpandParams: []string{
			"locations",
			"datastreams",
			"historicallocations",
		},
		SupportedSelectParams: []string{
			"id",
			"name",
			"properties",
			"description",
			"Locations",
			"datastreams",
			"historicallocations",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/things", handlers.HandleGetThings},
			{models.HTTPOperationGet, "/v1.0/things{id}", handlers.HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/thing", handlers.HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/thing/{params}", handlers.HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/thing/{params}/$value", handlers.HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/thing", handlers.HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/thing/{params}", handlers.HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/thing/{params}/$value", handlers.HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/locations{id}/things", handlers.HandleGetThingsByLocation},
			{models.HTTPOperationGet, "/v1.0/locations{id}/things/{params}", handlers.HandleGetThingsByLocation},
			{models.HTTPOperationGet, "/v1.0/things{id}/{params}", handlers.HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/things{id}/{params}/$value", handlers.HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/things/{params}", handlers.HandleGetThings},

			{models.HTTPOperationPost, "/v1.0/things", handlers.HandlePostThing},
			{models.HTTPOperationDelete, "/v1.0/things{id}", handlers.HandleDeleteThing},
			{models.HTTPOperationPatch, "/v1.0/things{id}", handlers.HandlePatchThing},
			{models.HTTPOperationPut, "/v1.0/things{id}", handlers.HandlePutThing},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/things", handlers.HandleGetThings},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/things{id}", handlers.HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/locations{id}/things", handlers.HandleGetThingsByLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/locations{id}/things/{params}", handlers.HandleGetThingsByLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/thing", handlers.HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/thing/{params}", handlers.HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/thing/{params}/$value", handlers.HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}/thing", handlers.HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}/thing/{params}", handlers.HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}/thing/{params}/$value", handlers.HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/things{id}/{params}", handlers.HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/things{id}/{params}/$value", handlers.HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/things/{params}", handlers.HandleGetThings},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/things", handlers.HandlePostThing},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/things{id}", handlers.HandleDeleteThing},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/things{id}", handlers.HandlePatchThing},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/things{id}", handlers.HandlePutThing},
		},
	}
}
