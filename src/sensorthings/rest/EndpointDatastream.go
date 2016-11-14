package rest

import (
	"fmt"

	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

func createDatastreamsEndpoint(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "Datastreams",
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "Datastreams")),
		SupportedQueryOptions: []odata.QueryOptionType{
			odata.QueryOptionTop, odata.QueryOptionSkip, odata.QueryOptionOrderBy, odata.QueryOptionCount, odata.QueryOptionResultFormat,
			odata.QueryOptionExpand, odata.QueryOptionSelect, odata.QueryOptionFilter,
		},
		SupportedExpandParams: []string{
			"Thing",
			"Sensor",
			"Observedproperty",
			"Observations",
		},
		SupportedSelectParams: []string{
			"description",
			"unitOfMeasurement",
			"observationType",
			"observedArea",
			"phenomenonTime",
			"resultTime",
			"Thing",
			"Sensor",
			"ObservedProperty",
			"Observations",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/datastreams", HandleGetDatastreams},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/observedProperty", HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/observations", HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/sensor", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/thing", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/observedproperty/{params}", HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/observations/{params}", HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/sensor/{params}", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/sensor/{params}/$value", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/thing/{params}", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/thing/{params}/$value", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/{params}", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/{params}/$value", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams/{params}", HandleGetDatastreams},

			{models.HTTPOperationPost, "/v1.0/datastreams", HandlePostDatastream},
			{models.HTTPOperationPost, "/v1.0/things{id}/datastreams", HandlePostDatastreamByThing},
			{models.HTTPOperationDelete, "/v1.0/datastreams{id}", HandleDeleteDatastream},
			{models.HTTPOperationPatch, "/v1.0/datastreams{id}", HandlePatchDatastream},
			{models.HTTPOperationPut, "/v1.0/datastreams{id}", HandlePutDatastream},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams", HandleGetDatastreams},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/observedproperty", HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/observations", HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/sensor", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/thing", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/observedproperty/{params}", HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/observations/{params}", HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/sensor/{params}", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/sensor/{params}/$value", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/thing/{params}", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/thing/{params}/$value", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/{params}", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/{params}/$value", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams/{params}", HandleGetDatastreams},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/datastreams", HandlePostDatastream},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/datastreams{id}", HandleDeleteDatastream},
			{models.HTTPOperationPost, "/v1.0/{c:.*}/things{id}/datastreams", HandlePostDatastreamByThing},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/datastreams{id}", HandlePatchDatastream},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/datastreams{id}", HandlePutDatastream},
		},
	}
}
