package config

import (
	"fmt"

	entities "github.com/gost/core"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/rest/endpoint"
	"github.com/gost/server/sensorthings/rest/handlers"
)

// CreateSensorsEndpoint creates the Sensors endpoint configuration
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
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/sensors", Handler: handlers.HandleGetSensors},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/sensors{id}", Handler: handlers.HandleGetSensor},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams{id}/sensor", Handler: handlers.HandleGetSensorByDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams{id}/sensor/{params}", Handler: handlers.HandleGetSensorByDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams{id}/sensor/{params}/$value", Handler: handlers.HandleGetSensorByDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/sensors{id}/{params}", Handler: handlers.HandleGetSensor},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/sensors{id}/{params}/$value", Handler: handlers.HandleGetSensor},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/sensors/{params}", Handler: handlers.HandleGetSensors},

			{OperationType: models.HTTPOperationPost, Path: "/v1.0/sensors", Handler: handlers.HandlePostSensors},
			{OperationType: models.HTTPOperationDelete, Path: "/v1.0/sensors{id}", Handler: handlers.HandleDeleteSensor},
			{OperationType: models.HTTPOperationPatch, Path: "/v1.0/sensors{id}", Handler: handlers.HandlePatchSensor},
			{OperationType: models.HTTPOperationPut, Path: "/v1.0/sensors{id}", Handler: handlers.HandlePutSensor},

			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/sensors", Handler: handlers.HandleGetSensors},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/sensors{id}", Handler: handlers.HandleGetSensor},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/datastreams{id}/sensor", Handler: handlers.HandleGetSensorByDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/datastreams{id}/sensor/{params}", Handler: handlers.HandleGetSensorByDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/datastreams{id}/sensor/{params}/$value", Handler: handlers.HandleGetSensorByDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/sensors{id}/{params}", Handler: handlers.HandleGetSensor},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/sensors{id}/{params}/$value", Handler: handlers.HandleGetSensor},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/sensors/{params}", Handler: handlers.HandleGetSensors},

			{OperationType: models.HTTPOperationPost, Path: "/v1.0/{c:.*}/sensors", Handler: handlers.HandlePostSensors},
			{OperationType: models.HTTPOperationDelete, Path: "/v1.0/{c:.*}/sensors{id}", Handler: handlers.HandleDeleteSensor},
			{OperationType: models.HTTPOperationPatch, Path: "/v1.0/{c:.*}/sensors{id}", Handler: handlers.HandlePatchSensor},
			{OperationType: models.HTTPOperationPut, Path: "/v1.0/{c:.*}/sensors{id}", Handler: handlers.HandlePutSensor},
		},
	}
}
