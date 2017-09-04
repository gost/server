package config

import (
	"fmt"

	entities "github.com/gost/core"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/rest/endpoint"
	"github.com/gost/server/sensorthings/rest/handlers"
)

// CreateThingsEndpoint creates the Things endpoint configuration
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
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/things", Handler: handlers.HandleGetThings},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/things{id}", Handler: handlers.HandleGetThing},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/historicallocations{id}/thing", Handler: handlers.HandleGetThingByHistoricalLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/historicallocations{id}/thing/{params}", Handler: handlers.HandleGetThingByHistoricalLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/historicallocations{id}/thing/{params}/$value", Handler: handlers.HandleGetThingByHistoricalLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams{id}/thing", Handler: handlers.HandleGetThingByDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams{id}/thing/{params}", Handler: handlers.HandleGetThingByDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams{id}/thing/{params}/$value", Handler: handlers.HandleGetThingByDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/locations{id}/things", Handler: handlers.HandleGetThingsByLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/locations{id}/things/{params}", Handler: handlers.HandleGetThingsByLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/things{id}/{params}", Handler: handlers.HandleGetThing},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/things{id}/{params}/$value", Handler: handlers.HandleGetThing},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/things/{params}", Handler: handlers.HandleGetThings},

			{OperationType: models.HTTPOperationPost, Path: "/v1.0/things", Handler: handlers.HandlePostThing},
			{OperationType: models.HTTPOperationDelete, Path: "/v1.0/things{id}", Handler: handlers.HandleDeleteThing},
			{OperationType: models.HTTPOperationPatch, Path: "/v1.0/things{id}", Handler: handlers.HandlePatchThing},
			{OperationType: models.HTTPOperationPut, Path: "/v1.0/things{id}", Handler: handlers.HandlePutThing},

			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/things", Handler: handlers.HandleGetThings},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/things{id}", Handler: handlers.HandleGetThing},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/locations{id}/things", Handler: handlers.HandleGetThingsByLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/locations{id}/things/{params}", Handler: handlers.HandleGetThingsByLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/datastreams{id}/thing", Handler: handlers.HandleGetThingByDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/datastreams{id}/thing/{params}", Handler: handlers.HandleGetThingByDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/datastreams{id}/thing/{params}/$value", Handler: handlers.HandleGetThingByDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/historicallocations{id}/thing", Handler: handlers.HandleGetThingByHistoricalLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/historicallocations{id}/thing/{params}", Handler: handlers.HandleGetThingByHistoricalLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/historicallocations{id}/thing/{params}/$value", Handler: handlers.HandleGetThingByHistoricalLocation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/things{id}/{params}", Handler: handlers.HandleGetThing},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/things{id}/{params}/$value", Handler: handlers.HandleGetThing},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/things/{params}", Handler: handlers.HandleGetThings},

			{OperationType: models.HTTPOperationPost, Path: "/v1.0/{c:.*}/things", Handler: handlers.HandlePostThing},
			{OperationType: models.HTTPOperationDelete, Path: "/v1.0/{c:.*}/things{id}", Handler: handlers.HandleDeleteThing},
			{OperationType: models.HTTPOperationPatch, Path: "/v1.0/{c:.*}/things{id}", Handler: handlers.HandlePatchThing},
			{OperationType: models.HTTPOperationPut, Path: "/v1.0/{c:.*}/things{id}", Handler: handlers.HandlePutThing},
		},
	}
}
