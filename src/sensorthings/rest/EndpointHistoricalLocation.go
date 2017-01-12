package rest

import (
	"fmt"

	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

func createHistoricalLocationsEndpoint(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "HistoricalLocations",
		EntityType: entities.EntityTypeHistoricalLocation,
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "HistoricalLocations")),
		SupportedQueryOptions: []odata.QueryOptionType{
			odata.QueryOptionTop, odata.QueryOptionSkip, odata.QueryOptionOrderBy, odata.QueryOptionCount, odata.QueryOptionResultFormat,
			odata.QueryOptionExpand, odata.QueryOptionSelect, odata.QueryOptionFilter,
		},
		SupportedExpandParams: []string{
			"locations",
			"thing",
		},
		SupportedSelectParams: []string{
			"id",
			"time",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/historicallocations", HandleGetHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/things{id}/historicallocations", HandleGetHistoricalLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/things{id}/historicallocations/{params}", HandleGetHistoricalLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/locations{id}/historicallocations", HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/locations{id}/historicallocations/{params}", HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/locations{id}/historicallocations/{params}/$value", HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/{params}", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/{params}/$value", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/historicallocations/{params}", HandleGetHistoricalLocations},

			{models.HTTPOperationPost, "/v1.0/historicallocations", HandlePostHistoricalLocation},
			{models.HTTPOperationDelete, "/v1.0/historicallocations{id}", HandleDeleteHistoricalLocations},
			{models.HTTPOperationPatch, "/v1.0/historicallocations{id}", HandlePatchHistoricalLocations},
			{models.HTTPOperationPut, "/v1.0/historicallocations{id}", HandlePutHistoricalLocation},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations", HandleGetHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/things{id}/historicallocations", HandleGetHistoricalLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/things{id}/historicallocations/{params}", HandleGetHistoricalLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/locations{id}/historicallocations", HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/locations{id}/historicallocations/{params}", HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/locations{id}/historicallocations/{params}/$value", HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}/{params}", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}/{params}/$value", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations/{params}", HandleGetHistoricalLocations},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/historicallocations", HandlePostHistoricalLocation},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/historicallocations{id}", HandleDeleteHistoricalLocations},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/historicallocations{id}", HandlePatchHistoricalLocations},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/historicallocations{id}", HandlePutHistoricalLocation},
		},
	}
}
