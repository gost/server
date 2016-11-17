package rest

import (
	"fmt"

	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

func createFeaturesOfInterestEndpoint(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "FeaturesOfInterest",
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "FeaturesOfInterest")),
		SupportedQueryOptions: []odata.QueryOptionType{
			odata.QueryOptionTop, odata.QueryOptionSkip, odata.QueryOptionOrderBy, odata.QueryOptionCount, odata.QueryOptionResultFormat,
			odata.QueryOptionExpand, odata.QueryOptionSelect, odata.QueryOptionFilter,
		},
		SupportedExpandParams: []string{
			"Observation",
		},
		SupportedSelectParams: []string{
			"id",
			"name",
			"description",
			"encodingType",
			"feature",
			"Observations",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/featuresofinterest", HandleGetFeatureOfInterests},
			{models.HTTPOperationGet, "/v1.0/featuresofinterest{id}", HandleGetFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/featuresofinterest{id}/{params}", HandleGetFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/featuresofinterest{id}/{params}/$value", HandleGetFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/featuresofinterest/{params}", HandleGetFeatureOfInterests},
			{models.HTTPOperationGet, "/v1.0/observations{id}/featureofinterest", HandleGetFeatureOfInterestByObservation},
			{models.HTTPOperationGet, "/v1.0/observations{id}/featureofinterest/{params}", HandleGetFeatureOfInterestByObservation},
			{models.HTTPOperationGet, "/v1.0/observations{id}/featureofinterest/{params}/$value", HandleGetFeatureOfInterestByObservation},
			{models.HTTPOperationPost, "/v1.0/featuresofinterest", HandlePostFeatureOfInterest},
			{models.HTTPOperationDelete, "/v1.0/featuresofinterest{id}", HandleDeleteFeatureOfInterest},
			{models.HTTPOperationPatch, "/v1.0/featuresofinterest{id}", HandlePatchFeatureOfInterest},
			{models.HTTPOperationPut, "/v1.0/featuresofinterest{id}", HandlePutFeatureOfInterest},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/featuresofinterest", HandleGetFeatureOfInterests},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/featuresofinterest{id}", HandleGetFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/featuresofinterest{id}/{params}", HandleGetFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/featuresofinterest{id}/{params}/$value", HandleGetFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations{id}/featureofInterest", HandleGetFeatureOfInterestByObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations{id}/featureofinterest/{params}", HandleGetFeatureOfInterestByObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations{id}/featureofinterest/{params}/$value", HandleGetFeatureOfInterestByObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/featuresofinterest/{params}", HandleGetFeatureOfInterests},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/featuresofinterest", HandlePostFeatureOfInterest},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/featuresofinterest{id}", HandleDeleteFeatureOfInterest},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/featuresofinterest{id}", HandlePatchFeatureOfInterest},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/featuresofinterest{id}", HandlePutFeatureOfInterest},
		},
	}
}
