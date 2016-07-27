package rest

import (
	"fmt"

	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

func createThingsEndpoint(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "Things",
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "Things")),
		SupportedQueryOptions: []odata.QueryOptionType{
			odata.QueryOptionTop, odata.QueryOptionSkip, odata.QueryOptionOrderBy, odata.QueryOptionCount, odata.QueryOptionResultFormat,
			odata.QueryOptionExpand, odata.QueryOptionSelect, odata.QueryOptionFilter,
		},
		SupportedExpandParams: []string{
			"Locations",
			"Datastreams",
			"HistoricalLocations",
		},
		SupportedSelectParams: []string{
			"properties",
			"description",
			"Locations",
			"Datastreams",
			"HistoricalLocations",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/Things", HandleGetThings},
			{models.HTTPOperationGet, "/v1.0/Things{id}", HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/Things{id}/Datastreams", HandleGetDatastreamsByThing},
			{models.HTTPOperationGet, "/v1.0/Things{id}/HistoricalLocations", HandleGetHistoricalLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/Things{id}/Locations", HandleGetLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/Things{id}/Datastreams/{params}", HandleGetDatastreamsByThing},
			{models.HTTPOperationGet, "/v1.0/Things{id}/HistoricalLocations/{params}", HandleGetHistoricalLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/Things{id}/Locations/{params}", HandleGetLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/Things{id}/{params}", HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/Things{id}/{params}/$value", HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/Things/{params}", HandleGetThings},

			{models.HTTPOperationPost, "/v1.0/Things", HandlePostThing},
			{models.HTTPOperationDelete, "/v1.0/Things{id}", HandleDeleteThing},
			{models.HTTPOperationPatch, "/v1.0/Things{id}", HandlePatchThing},
			{models.HTTPOperationPut, "/v1.0/Things{id}", HandlePutThing},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things", HandleGetThings},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things{id}", HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things{id}/Datastreams", HandleGetDatastreamsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things{id}/HistoricalLocations", HandleGetHistoricalLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things{id}/Locations", HandleGetLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things{id}/Datastreams/{params}", HandleGetDatastreamsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things{id}/HistoricalLocations/{params}", HandleGetHistoricalLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things{id}/Locations/{params}", HandleGetLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things{id}/{params}", HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things{id}/{params}/$value", HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things/{params}", HandleGetThings},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/Things", HandlePostThing},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/Things{id}", HandleDeleteThing},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/Things{id}", HandlePatchThing},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/Things{id}", HandlePutThing},
		},
	}
}
