package rest

import (
	"fmt"

	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
)

func createThingsEndpoint(externalURL string) *Endpoint {
	return &Endpoint{
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
			{models.HTTPOperationGet, "/v1.0/things", HandleGetThings},
			{models.HTTPOperationGet, "/v1.0/things{id}", HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/thing", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/thing/{params}", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/thing/{params}/$value", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/thing", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/thing/{params}", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/thing/{params}/$value", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/locations{id}/things", HandleGetThingsByLocation},
			{models.HTTPOperationGet, "/v1.0/locations{id}/things/{params}", HandleGetThingsByLocation},
			{models.HTTPOperationGet, "/v1.0/things{id}/{params}", HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/things{id}/{params}/$value", HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/things/{params}", HandleGetThings},

			{models.HTTPOperationPost, "/v1.0/things", HandlePostThing},
			{models.HTTPOperationDelete, "/v1.0/things{id}", HandleDeleteThing},
			{models.HTTPOperationPatch, "/v1.0/things{id}", HandlePatchThing},
			{models.HTTPOperationPut, "/v1.0/things{id}", HandlePutThing},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/things", HandleGetThings},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/things{id}", HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/locations{id}/things", HandleGetThingsByLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/locations{id}/things/{params}", HandleGetThingsByLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/thing", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/thing/{params}", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/thing/{params}/$value", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}/thing", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}/thing/{params}", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}/thing/{params}/$value", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/things{id}/{params}", HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/things{id}/{params}/$value", HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/things/{params}", HandleGetThings},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/things", HandlePostThing},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/things{id}", HandleDeleteThing},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/things{id}", HandlePatchThing},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/things{id}", HandlePutThing},
		},
	}
}
