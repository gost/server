package config

import (
	"fmt"

	"github.com/gost/server/sensorthings/entities"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/rest/endpoint"
	"github.com/gost/server/sensorthings/rest/handlers"
)

// CreateFeaturesOfInterestEndpoint constructs the featuresOfInterest endpoint configuration
func CreateFeaturesOfInterestEndpoint(externalURL string) *endpoint.Endpoint {
	return &endpoint.Endpoint{
		Name:       "FeaturesOfInterest",
		EntityType: entities.EntityTypeFeatureOfInterest,
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "FeaturesOfInterest")),
		SupportedExpandParams: []string{
			"observations",
		},
		SupportedSelectParams: []string{
			"id",
			"name",
			"description",
			"encodingtype",
			"feature",
			"observations",
		},
		Operations: []models.EndpointOperation{
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/featuresofinterest", Handler: handlers.HandleGetFeatureOfInterests},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/featuresofinterest{id}", Handler: handlers.HandleGetFeatureOfInterest},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/featuresofinterest{id}/{params}", Handler: handlers.HandleGetFeatureOfInterest},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/featuresofinterest{id}/{params}/$value", Handler: handlers.HandleGetFeatureOfInterest},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/featuresofinterest/{params}", Handler: handlers.HandleGetFeatureOfInterests},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/observations{id}/featureofinterest", Handler: handlers.HandleGetFeatureOfInterestByObservation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/observations{id}/featureofinterest/{params}", Handler: handlers.HandleGetFeatureOfInterestByObservation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/observations{id}/featureofinterest/{params}/$value", Handler: handlers.HandleGetFeatureOfInterestByObservation},
			{OperationType: models.HTTPOperationPost, Path: "/v1.0/featuresofinterest", Handler: handlers.HandlePostFeatureOfInterest},
			{OperationType: models.HTTPOperationDelete, Path: "/v1.0/featuresofinterest{id}", Handler: handlers.HandleDeleteFeatureOfInterest},
			{OperationType: models.HTTPOperationPatch, Path: "/v1.0/featuresofinterest{id}", Handler: handlers.HandlePatchFeatureOfInterest},
			{OperationType: models.HTTPOperationPut, Path: "/v1.0/featuresofinterest{id}", Handler: handlers.HandlePutFeatureOfInterest},

			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/featuresofinterest", Handler: handlers.HandleGetFeatureOfInterests},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/featuresofinterest{id}", Handler: handlers.HandleGetFeatureOfInterest},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/featuresofinterest{id}/{params}", Handler: handlers.HandleGetFeatureOfInterest},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/featuresofinterest{id}/{params}/$value", Handler: handlers.HandleGetFeatureOfInterest},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/observations{id}/featureofinterest", Handler: handlers.HandleGetFeatureOfInterestByObservation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/observations{id}/featureofinterest/{params}", Handler: handlers.HandleGetFeatureOfInterestByObservation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/observations{id}/featureofinterest/{params}/$value", Handler: handlers.HandleGetFeatureOfInterestByObservation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/featuresofinterest/{params}", Handler: handlers.HandleGetFeatureOfInterests},

			{OperationType: models.HTTPOperationPost, Path: "/v1.0/{c:.*}/featuresofinterest", Handler: handlers.HandlePostFeatureOfInterest},
			{OperationType: models.HTTPOperationDelete, Path: "/v1.0/{c:.*}/featuresofinterest{id}", Handler: handlers.HandleDeleteFeatureOfInterest},
			{OperationType: models.HTTPOperationPatch, Path: "/v1.0/{c:.*}/featuresofinterest{id}", Handler: handlers.HandlePatchFeatureOfInterest},
			{OperationType: models.HTTPOperationPut, Path: "/v1.0/{c:.*}/featuresofinterest{id}", Handler: handlers.HandlePutFeatureOfInterest},
		},
	}
}
