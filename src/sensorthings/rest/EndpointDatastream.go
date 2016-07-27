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
			{models.HTTPOperationGet, "/v1.0/Datastreams", HandleGetDatastreams},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/ObservedProperty", HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Observations", HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Sensor", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Thing", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/ObservedProperty/{params}", HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Observations/{params}", HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Sensor/{params}", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Sensor/{params}/$value", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Thing/{params}", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Thing/{params}/$value", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/{params}", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/{params}/$value", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams/{params}", HandleGetDatastreams},

			{models.HTTPOperationPost, "/v1.0/Datastreams", HandlePostDatastream},
			{models.HTTPOperationPost, "/v1.0/Things{id}/Datastreams", HandlePostDatastreamByThing},
			{models.HTTPOperationDelete, "/v1.0/Datastreams{id}", HandleDeleteDatastream},
			{models.HTTPOperationPatch, "/v1.0/Datastreams{id}", HandlePatchDatastream},
			{models.HTTPOperationPut, "/v1.0/Datastreams{id}", HandlePutDatastream},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams", HandleGetDatastreams},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/ObservedProperty", HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/Observations", HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/Sensor", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/Thing", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/ObservedProperty/{params}", HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/Observations/{params}", HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/Sensor/{params}", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/Sensor/{params}/$value", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/Thing/{params}", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/Thing/{params}/$value", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/{params}", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/{params}/$value", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams/{params}", HandleGetDatastreams},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/Datastreams", HandlePostDatastream},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/Datastreams{id}", HandleDeleteDatastream},
			{models.HTTPOperationPost, "/v1.0/{c:.*}/Things{id}/Datastreams", HandlePostDatastreamByThing},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/Datastreams{id}", HandlePatchDatastream},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/Datastreams{id}", HandlePutDatastream},
		},
	}
}
