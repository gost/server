package rest

import (
	"fmt"

	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
)

func createSensorsEndpoint(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "Sensors",
		EntityType: entities.EntityTypeSensor,
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "Sensors")),
		SupportedExpandParams: []string{
			"datastreams",
		},
		SupportedSelectParams: []string{
			"id",
			"name",
			"description",
			"encodingtype",
			"metadata",
			"datastreams",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/sensors", HandleGetSensors},
			{models.HTTPOperationGet, "/v1.0/sensors{id}", HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/sensor", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/sensor/{params}", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/sensor/{params}/$value", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/sensors{id}/{params}", HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/sensors{id}/{params}/$value", HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/sensors/{params}", HandleGetSensors},

			{models.HTTPOperationPost, "/v1.0/sensors", HandlePostSensors},
			{models.HTTPOperationDelete, "/v1.0/sensors{id}", HandleDeleteSensor},
			{models.HTTPOperationPatch, "/v1.0/sensors{id}", HandlePatchSensor},
			{models.HTTPOperationPut, "/v1.0/sensors{id}", HandlePutSensor},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/sensors", HandleGetSensors},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/sensors{id}", HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/sensor", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/sensor/{params}", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/sensor/{params}/$value", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/sensors{id}/{params}", HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/sensors{id}/{params}/$value", HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/sensors/{params}", HandleGetSensors},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/sensors", HandlePostSensors},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/sensors{id}", HandleDeleteSensor},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/sensors{id}", HandlePatchSensor},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/sensors{id}", HandlePutSensor},
		},
	}
}
