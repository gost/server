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
			{models.HTTPOperationGet, "/v1.0/observedproperties", HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/observedproperties{id}", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/observedproperties{id}/datastreams", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/observedproperties{id}/datastreams/{params}", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/observedproperties{id}/datastreams/{params}/$value", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/observedproperties{id}/{params}", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/observedproperties{id}/{params}/$value", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/observedproperties/{params}", HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/observedproperties/{params}/$value", HandleGetObservedProperties},

			{models.HTTPOperationPost, "/v1.0/observedproperties", HandlePostObservedProperty},
			{models.HTTPOperationDelete, "/v1.0/observedproperties{id}", HandleDeleteObservedProperty},
			{models.HTTPOperationPatch, "/v1.0/observedproperties{id}", HandlePatchObservedProperty},
			{models.HTTPOperationPut, "/v1.0/observedproperties{id}", HandlePutObservedProperty},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties", HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties{id}", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties{id}/datastreams", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties{id}/datastreams/{params}", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties{id}/datastreams/{params}/$value", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties{id}/{params}", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties{id}/{params}/$value", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties/{params}", HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties/{params}/$value", HandleGetObservedProperties},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/observedproperties", HandlePostObservedProperty},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/observedproperties{id}", HandleDeleteObservedProperty},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/observedproperties{id}", HandlePatchObservedProperty},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/observedproperties{id}", HandlePutObservedProperty},
		},
	}
}
