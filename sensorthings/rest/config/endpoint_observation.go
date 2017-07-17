package config

import (
	"fmt"

	"github.com/gost/server/sensorthings/entities"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/rest/endpoint"
	"github.com/gost/server/sensorthings/rest/handlers"
)

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
			{models.HTTPOperationGet, "/v1.0/observations", handlers.HandleGetObservations},
			{models.HTTPOperationGet, "/v1.0/observations{id}", handlers.HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/observations", handlers.HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/observations/{params}", handlers.HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/featureofinterest{id}/observations", handlers.HandleGetObservationsByFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/featuresofinterest{id}/observations", handlers.HandleGetObservationsByFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/featureofinterest{id}/observations/{params}", handlers.HandleGetObservationsByFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/featuresofinterest{id}/observations/{params}", handlers.HandleGetObservationsByFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/observations{id}/{params}", handlers.HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/observations{id}/{params}/$value", handlers.HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/observations/{params}", handlers.HandleGetObservations},

			{models.HTTPOperationPost, "/v1.0/observations", handlers.HandlePostObservation},
			{models.HTTPOperationPost, "/v1.0/datastreams{id}/observations", handlers.HandlePostObservationByDatastream},
			{models.HTTPOperationDelete, "/v1.0/observations{id}", handlers.HandleDeleteObservation},
			{models.HTTPOperationPatch, "/v1.0/observations{id}", handlers.HandlePatchObservation},
			{models.HTTPOperationPut, "/v1.0/observations{id}", handlers.HandlePutObservation},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations", handlers.HandleGetObservations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations{id}", handlers.HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/featuresofinterest{id}/observations", handlers.HandleGetObservationsByFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/featuresofinterest{id}/observations/{params}", handlers.HandleGetObservationsByFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/observations", handlers.HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/observations/{params}", handlers.HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations{id}/{params}", handlers.HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations{id}/{params}/$value", handlers.HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations/{params}", handlers.HandleGetObservations},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/observations", handlers.HandlePostObservation},
			{models.HTTPOperationPost, "/v1.0/{c:.*}/datastreams{id}/observations", handlers.HandlePostObservationByDatastream},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/observations{id}", handlers.HandleDeleteObservation},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/observations{id}", handlers.HandlePatchObservation},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/observations{id}", handlers.HandlePutObservation},
		},
	}
}
