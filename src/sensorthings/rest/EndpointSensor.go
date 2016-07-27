package rest

import (
	"fmt"

	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

func createSensorsEndpoint(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "Sensors",
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "Sensors")),
		SupportedQueryOptions: []odata.QueryOptionType{
			odata.QueryOptionTop, odata.QueryOptionSkip, odata.QueryOptionOrderBy, odata.QueryOptionCount, odata.QueryOptionResultFormat,
			odata.QueryOptionExpand, odata.QueryOptionSelect, odata.QueryOptionFilter,
		},
		SupportedExpandParams: []string{
			"Datastream",
		},
		SupportedSelectParams: []string{
			"description",
			"encodingType",
			"metadata",
			"Datastreams",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/Sensors", HandleGetSensors},
			{models.HTTPOperationGet, "/v1.0/Sensors{id}", HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/Sensors{id}/Datastreams", HandleGetDatastreamsBySensor},
			{models.HTTPOperationGet, "/v1.0/Sensors{id}/Datastreams/{params}", HandleGetDatastreamsBySensor},
			{models.HTTPOperationGet, "/v1.0/Sensors{id}/{params}", HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/Sensors{id}/{params}/$value", HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/Sensors/{params}", HandleGetSensors},

			{models.HTTPOperationPost, "/v1.0/Sensors", HandlePostSensors},
			{models.HTTPOperationDelete, "/v1.0/Sensors{id}", HandleDeleteSensor},
			{models.HTTPOperationPatch, "/v1.0/Sensors{id}", HandlePatchSensor},
			{models.HTTPOperationPut, "/v1.0/Sensors{id}", HandlePutSensor},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/Sensors", HandleGetSensors},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Sensors{id}", HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Sensors{id}/Datastreams", HandleGetDatastreamsBySensor},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Sensors{id}/Datastreams/{params}", HandleGetDatastreamsBySensor},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Sensors{id}/{params}", HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Sensors{id}/{params}/$value", HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Sensors/{params}", HandleGetSensors},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/Sensors", HandlePostSensors},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/Sensors{id}", HandleDeleteSensor},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/Sensors{id}", HandlePatchSensor},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/Sensors{id}", HandlePutSensor},
		},
	}
}
