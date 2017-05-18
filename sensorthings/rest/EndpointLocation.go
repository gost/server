package rest

import (
	"fmt"

	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
)

func createLocationsEndpoint(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "Locations",
		EntityType: entities.EntityTypeLocation,
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "Locations")),
		SupportedQueryOptions: []odata.QueryOptionType{
			odata.QueryOptionTop, odata.QueryOptionSkip, odata.QueryOptionOrderBy, odata.QueryOptionCount, odata.QueryOptionResultFormat,
			odata.QueryOptionExpand, odata.QueryOptionSelect, odata.QueryOptionFilter,
		},
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
			{models.HTTPOperationGet, "/v1.0/locations", HandleGetLocations},
			{models.HTTPOperationGet, "/v1.0/locations{id}", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/locations", HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/locations/{params}", HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/locations{id}/{params}", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/locations{id}/{params}/$value", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/things{id}/locations", HandleGetLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/things{id}/locations/{params}", HandleGetLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/locations/{params}", HandleGetLocations},

			{models.HTTPOperationPost, "/v1.0/locations", HandlePostLocation},
			{models.HTTPOperationPost, "/v1.0/things{id}/locations", HandlePostLocationByThing},
			{models.HTTPOperationDelete, "/v1.0/locations{id}", HandleDeleteLocation},
			{models.HTTPOperationPatch, "/v1.0/locations{id}", HandlePatchLocation},
			{models.HTTPOperationPut, "/v1.0/locations{id}", HandlePutLocation},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/locations", HandleGetLocations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/locations{id}", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}/locations", HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/things{id}/locations", HandleGetLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/things{id}/locations/{params}", HandleGetLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}/locations/{params}", HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/locations{id}/{params}", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/locations{id}/{params}/$value", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/locations/{params}", HandleGetLocations},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/locations", HandlePostLocation},
			{models.HTTPOperationPost, "/v1.0/{c:.*}/things{id}/locations", HandlePostLocationByThing},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/locations{id}", HandleDeleteLocation},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/locations{id}", HandlePatchLocation},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/locations{id}", HandlePutLocation},
		},
	}
}
