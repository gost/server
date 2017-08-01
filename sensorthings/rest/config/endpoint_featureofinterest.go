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
			{models.HTTPOperationGet, "/v1.0/featuresofinterest", handlers.HandleGetFeatureOfInterests},
			{models.HTTPOperationGet, "/v1.0/featuresofinterest{id}", handlers.HandleGetFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/featuresofinterest{id}/{params}", handlers.HandleGetFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/featuresofinterest{id}/{params}/$value", handlers.HandleGetFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/featuresofinterest/{params}", handlers.HandleGetFeatureOfInterests},
			{models.HTTPOperationGet, "/v1.0/observations{id}/featureofinterest", handlers.HandleGetFeatureOfInterestByObservation},
			{models.HTTPOperationGet, "/v1.0/observations{id}/featureofinterest/{params}", handlers.HandleGetFeatureOfInterestByObservation},
			{models.HTTPOperationGet, "/v1.0/observations{id}/featureofinterest/{params}/$value", handlers.HandleGetFeatureOfInterestByObservation},
			{models.HTTPOperationPost, "/v1.0/featuresofinterest", handlers.HandlePostFeatureOfInterest},
			{models.HTTPOperationDelete, "/v1.0/featuresofinterest{id}", handlers.HandleDeleteFeatureOfInterest},
			{models.HTTPOperationPatch, "/v1.0/featuresofinterest{id}", handlers.HandlePatchFeatureOfInterest},
			{models.HTTPOperationPut, "/v1.0/featuresofinterest{id}", handlers.HandlePutFeatureOfInterest},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/featuresofinterest", handlers.HandleGetFeatureOfInterests},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/featuresofinterest{id}", handlers.HandleGetFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/featuresofinterest{id}/{params}", handlers.HandleGetFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/featuresofinterest{id}/{params}/$value", handlers.HandleGetFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations{id}/featureofInterest", handlers.HandleGetFeatureOfInterestByObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations{id}/featureofinterest/{params}", handlers.HandleGetFeatureOfInterestByObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations{id}/featureofinterest/{params}/$value", handlers.HandleGetFeatureOfInterestByObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/featuresofinterest/{params}", handlers.HandleGetFeatureOfInterests},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/featuresofinterest", handlers.HandlePostFeatureOfInterest},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/featuresofinterest{id}", handlers.HandleDeleteFeatureOfInterest},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/featuresofinterest{id}", handlers.HandlePatchFeatureOfInterest},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/featuresofinterest{id}", handlers.HandlePutFeatureOfInterest},
		},
	}
}
