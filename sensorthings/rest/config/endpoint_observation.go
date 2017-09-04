package config

import (
	"fmt"

	entities "github.com/gost/core"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/rest/endpoint"
	"github.com/gost/server/sensorthings/rest/handlers"
)

// CreateObservationsEndpoint constructs the Observations endpoint configuration
func CreateObservationsEndpoint(externalURL string) *endpoint.Endpoint {
	return &endpoint.Endpoint{
		Name:       "Observations",
		EntityType: entities.EntityTypeObservation,
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "Observations")),
		SupportedExpandParams: []string{
			"datastream",
			"featureofinterest",
		},
		SupportedSelectParams: []string{
			"id",
			"result",
			"phenomenontime",
			"resulttime",
			"resultquality",
			"validtime",
			"parameters",
			"datastream",
			"featureofinterest",
		},
		Operations: []models.EndpointOperation{
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/observations", Handler: handlers.HandleGetObservations},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/observations{id}", Handler: handlers.HandleGetObservation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams{id}/observations", Handler: handlers.HandleGetObservationsByDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams{id}/observations/{params}", Handler: handlers.HandleGetObservationsByDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/featureofinterest{id}/observations", Handler: handlers.HandleGetObservationsByFeatureOfInterest},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/featuresofinterest{id}/observations", Handler: handlers.HandleGetObservationsByFeatureOfInterest},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/featureofinterest{id}/observations/{params}", Handler: handlers.HandleGetObservationsByFeatureOfInterest},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/featuresofinterest{id}/observations/{params}", Handler: handlers.HandleGetObservationsByFeatureOfInterest},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/observations{id}/{params}", Handler: handlers.HandleGetObservation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/observations{id}/{params}/$value", Handler: handlers.HandleGetObservation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/observations/{params}", Handler: handlers.HandleGetObservations},

			{OperationType: models.HTTPOperationPost, Path: "/v1.0/observations", Handler: handlers.HandlePostObservation},
			{OperationType: models.HTTPOperationPost, Path: "/v1.0/datastreams{id}/observations", Handler: handlers.HandlePostObservationByDatastream},
			{OperationType: models.HTTPOperationDelete, Path: "/v1.0/observations{id}", Handler: handlers.HandleDeleteObservation},
			{OperationType: models.HTTPOperationPatch, Path: "/v1.0/observations{id}", Handler: handlers.HandlePatchObservation},
			{OperationType: models.HTTPOperationPut, Path: "/v1.0/observations{id}", Handler: handlers.HandlePutObservation},

			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/observations", Handler: handlers.HandleGetObservations},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/observations{id}", Handler: handlers.HandleGetObservation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/featuresofinterest{id}/observations", Handler: handlers.HandleGetObservationsByFeatureOfInterest},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/featuresofinterest{id}/observations/{params}", Handler: handlers.HandleGetObservationsByFeatureOfInterest},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/datastreams{id}/observations", Handler: handlers.HandleGetObservationsByDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/datastreams{id}/observations/{params}", Handler: handlers.HandleGetObservationsByDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/observations{id}/{params}", Handler: handlers.HandleGetObservation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/observations{id}/{params}/$value", Handler: handlers.HandleGetObservation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/observations/{params}", Handler: handlers.HandleGetObservations},

			{OperationType: models.HTTPOperationPost, Path: "/v1.0/{c:.*}/observations", Handler: handlers.HandlePostObservation},
			{OperationType: models.HTTPOperationPost, Path: "/v1.0/{c:.*}/datastreams{id}/observations", Handler: handlers.HandlePostObservationByDatastream},
			{OperationType: models.HTTPOperationDelete, Path: "/v1.0/{c:.*}/observations{id}", Handler: handlers.HandleDeleteObservation},
			{OperationType: models.HTTPOperationPatch, Path: "/v1.0/{c:.*}/observations{id}", Handler: handlers.HandlePatchObservation},
			{OperationType: models.HTTPOperationPut, Path: "/v1.0/{c:.*}/observations{id}", Handler: handlers.HandlePutObservation},
		},
	}
}
