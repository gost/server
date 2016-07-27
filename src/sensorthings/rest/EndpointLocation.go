package rest

import (
	"fmt"

	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

func createLocationsEndpoint(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "Locations",
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "Locations")),
		SupportedQueryOptions: []odata.QueryOptionType{
			odata.QueryOptionTop, odata.QueryOptionSkip, odata.QueryOptionOrderBy, odata.QueryOptionCount, odata.QueryOptionResultFormat,
			odata.QueryOptionExpand, odata.QueryOptionSelect, odata.QueryOptionFilter,
		},
		SupportedExpandParams: []string{
			"Things",
			"HistoricalLocations",
		},
		SupportedSelectParams: []string{
			"description",
			"encodingType",
			"location",
			"Things",
			"HistoricalLocations",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/Locations", HandleGetLocations},
			{models.HTTPOperationGet, "/v1.0/Locations{id}", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/Locations{id}/Things", HandleGetThingsByLocation},
			{models.HTTPOperationGet, "/v1.0/Locations{id}/HistoricalLocations", HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/Locations{id}/Things/{params}", HandleGetThingsByLocation},
			{models.HTTPOperationGet, "/v1.0/Locations{id}/HistoricalLocations/{params}", HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/Locations{id}/HistoricalLocations/{params}/$value", HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/Locations{id}/{params}", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/Locations{id}/{params}/$value", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/Locations/{params}", HandleGetLocations},

			{models.HTTPOperationPost, "/v1.0/Locations", HandlePostLocation},
			{models.HTTPOperationPost, "/v1.0/Things{id}/Locations", HandlePostLocationByThing},
			{models.HTTPOperationDelete, "/v1.0/Locations{id}", HandleDeleteLocation},
			{models.HTTPOperationPatch, "/v1.0/Locations{id}", HandlePatchLocation},
			{models.HTTPOperationPut, "/v1.0/Locations{id}", HandlePutLocation},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/Locations", HandleGetLocations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Locations{id}", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Locations{id}/Things", HandleGetThingsByLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Locations{id}/HistoricalLocations", HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Locations{id}/Things/{params}", HandleGetThingsByLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Locations{id}/HistoricalLocations/{params}", HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Locations{id}/HistoricalLocations/{params}/$value", HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Locations{id}/{params}", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Locations{id}/{params}/$value", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Locations/{params}", HandleGetLocations},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/Locations", HandlePostLocation},
			{models.HTTPOperationPost, "/v1.0/{c:.*}/Things{id}/Locations", HandlePostLocationByThing},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/Locations{id}", HandleDeleteLocation},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/Locations{id}", HandlePatchLocation},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/Locations{id}", HandlePutLocation},
		},
	}
}
