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
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations", HandleGetHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}/Locations", HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}/Thing", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}/Thing/{params}", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}/Thing/{params}/$value", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}/Locations/{params}", HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}/{params}", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}/{params}/$value", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations/{params}", HandleGetHistoricalLocations},

			{models.HTTPOperationPost, "/v1.0/HistoricalLocations", HandlePostHistoricalLocation},
			{models.HTTPOperationDelete, "/v1.0/HistoricalLocations{id}", HandleDeleteHistoricalLocations},
			{models.HTTPOperationPatch, "/v1.0/HistoricalLocations{id}", HandlePatchHistoricalLocations},
			{models.HTTPOperationPut, "/v1.0/HistoricalLocations{id}", HandlePatchHistoricalLocations},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/HistoricalLocations", HandleGetHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/HistoricalLocations{id}", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/HistoricalLocations{id}/Locations", HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/HistoricalLocations{id}/Thing", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/HistoricalLocations{id}/Thing/{params}", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/HistoricalLocations{id}/Thing/{params}/$value", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/HistoricalLocations{id}/Locations/{params}", HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/HistoricalLocations{id}/{params}", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/HistoricalLocations{id}/{params}/$value", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/HistoricalLocations/{params}", HandleGetHistoricalLocations},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/HistoricalLocations", HandlePostHistoricalLocation},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/HistoricalLocations{id}", HandleDeleteHistoricalLocations},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/HistoricalLocations{id}", HandlePatchHistoricalLocations},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/HistoricalLocations{id}", HandlePatchHistoricalLocations},
		},
	}
}
