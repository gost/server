package rest

import (
	"fmt"

	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
)

func createObservationsEndpoint(externalURL string) *Endpoint {
	return &Endpoint{
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
			{models.HTTPOperationGet, "/v1.0/observations", HandleGetObservations},
			{models.HTTPOperationGet, "/v1.0/observations{id}", HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/observations", HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/observations/{params}", HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/featureofinterest{id}/observations", HandleGetObservationsByFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/featuresofinterest{id}/observations", HandleGetObservationsByFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/featureofinterest{id}/observations/{params}", HandleGetObservationsByFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/featuresofinterest{id}/observations/{params}", HandleGetObservationsByFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/observations{id}/{params}", HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/observations{id}/{params}/$value", HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/observations/{params}", HandleGetObservations},

			{models.HTTPOperationPost, "/v1.0/observations", HandlePostObservation},
			{models.HTTPOperationPost, "/v1.0/datastreams{id}/observations", HandlePostObservationByDatastream},
			{models.HTTPOperationDelete, "/v1.0/observations{id}", HandleDeleteObservation},
			{models.HTTPOperationPatch, "/v1.0/observations{id}", HandlePatchObservation},
			{models.HTTPOperationPut, "/v1.0/observations{id}", HandlePutObservation},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations", HandleGetObservations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations{id}", HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/featuresofinterest{id}/observations", HandleGetObservationsByFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/featuresofinterest{id}/observations/{params}", HandleGetObservationsByFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/observations", HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/observations/{params}", HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations{id}/{params}", HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations{id}/{params}/$value", HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations/{params}", HandleGetObservations},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/observations", HandlePostObservation},
			{models.HTTPOperationPost, "/v1.0/{c:.*}/datastreams{id}/observations", HandlePostObservationByDatastream},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/observations{id}", HandleDeleteObservation},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/observations{id}", HandlePatchObservation},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/observations{id}", HandlePutObservation},
		},
	}
}
