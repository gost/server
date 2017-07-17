package config

import (
	"fmt"

	"github.com/gost/server/sensorthings/entities"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/rest/endpoint"
	"github.com/gost/server/sensorthings/rest/handlers"
)

func CreateSensorsEndpoint(externalURL string) *endpoint.Endpoint {
	return &endpoint.Endpoint{
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
			{models.HTTPOperationGet, "/v1.0/sensors", handlers.HandleGetSensors},
			{models.HTTPOperationGet, "/v1.0/sensors{id}", handlers.HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/sensor", handlers.HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/sensor/{params}", handlers.HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/sensor/{params}/$value", handlers.HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/sensors{id}/{params}", handlers.HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/sensors{id}/{params}/$value", handlers.HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/sensors/{params}", handlers.HandleGetSensors},

			{models.HTTPOperationPost, "/v1.0/sensors", handlers.HandlePostSensors},
			{models.HTTPOperationDelete, "/v1.0/sensors{id}", handlers.HandleDeleteSensor},
			{models.HTTPOperationPatch, "/v1.0/sensors{id}", handlers.HandlePatchSensor},
			{models.HTTPOperationPut, "/v1.0/sensors{id}", handlers.HandlePutSensor},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/sensors", handlers.HandleGetSensors},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/sensors{id}", handlers.HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/sensor", handlers.HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/sensor/{params}", handlers.HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/sensor/{params}/$value", handlers.HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/sensors{id}/{params}", handlers.HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/sensors{id}/{params}/$value", handlers.HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/sensors/{params}", handlers.HandleGetSensors},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/sensors", handlers.HandlePostSensors},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/sensors{id}", handlers.HandleDeleteSensor},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/sensors{id}", handlers.HandlePatchSensor},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/sensors{id}", handlers.HandlePutSensor},
		},
	}
}
