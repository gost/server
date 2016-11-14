package rest

import (
	"fmt"

	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

func createHistoricalLocationsEndpoint(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "HistoricalLocations",
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
			"time",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/historicallocations", HandleGetHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/locations", HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/thing", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/thing/{params}", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/thing/{params}/$value", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/locations/{params}", HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/{params}", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/historicallocations{id}/{params}/$value", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/historicallocations/{params}", HandleGetHistoricalLocations},

			{models.HTTPOperationPost, "/v1.0/historicallocations", HandlePostHistoricalLocation},
			{models.HTTPOperationDelete, "/v1.0/historicallocations{id}", HandleDeleteHistoricalLocations},
			{models.HTTPOperationPatch, "/v1.0/historicallocations{id}", HandlePatchHistoricalLocations},
			{models.HTTPOperationPut, "/v1.0/historicallocations{id}", HandlePutHistoricalLocation},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations", HandleGetHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}/locations", HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}/thing", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}/thing/{params}", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}/thing/{params}/$value", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/historicallocations{id}/locations/{params}", HandleGetLocationsByHistoricalLocations},
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
