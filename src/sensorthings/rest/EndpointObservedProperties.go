package rest

import (
	"fmt"

	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

func createObservedPropertiesEndpoint(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "ObservedProperties",
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "ObservedProperties")),
		SupportedQueryOptions: []odata.QueryOptionType{
			odata.QueryOptionTop, odata.QueryOptionSkip, odata.QueryOptionOrderBy, odata.QueryOptionCount, odata.QueryOptionResultFormat,
			odata.QueryOptionExpand, odata.QueryOptionSelect, odata.QueryOptionFilter,
		},
		SupportedExpandParams: []string{
			"Datastreams",
		},
		SupportedSelectParams: []string{
			"name",
			"definition",
			"description",
			"Datastreams",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/ObservedProperties", HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/ObservedProperties{id}", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/ObservedProperties{id}/Datastreams", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/ObservedProperties{id}/Datastreams/{params}", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/ObservedProperties{id}/Datastreams/{params}/$value", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/ObservedProperties{id}/{params}", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/ObservedProperties{id}/{params}/$value", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/ObservedProperties/{params}", HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/ObservedProperties/{params}/$value", HandleGetObservedProperties},

			{models.HTTPOperationPost, "/v1.0/ObservedProperties", HandlePostObservedProperty},
			{models.HTTPOperationDelete, "/v1.0/ObservedProperties{id}", HandleDeleteObservedProperty},
			{models.HTTPOperationPatch, "/v1.0/ObservedProperties{id}", HandlePatchObservedProperty},
			{models.HTTPOperationPut, "/v1.0/ObservedProperties{id}", HandlePatchObservedProperty},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/ObservedProperties", HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/ObservedProperties{id}", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/ObservedProperties{id}/Datastreams", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/ObservedProperties{id}/Datastreams/{params}", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/ObservedProperties{id}/Datastreams/{params}/$value", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/ObservedProperties{id}/{params}", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/ObservedProperties{id}/{params}/$value", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/ObservedProperties/{params}", HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/ObservedProperties/{params}/$value", HandleGetObservedProperties},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/ObservedProperties", HandlePostObservedProperty},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/ObservedProperties{id}", HandleDeleteObservedProperty},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/ObservedProperties{id}", HandlePatchObservedProperty},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/ObservedProperties{id}", HandlePatchObservedProperty},
		},
	}
}
