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
			"description",
			"encodingType",
			"feature",
			"Observations",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/FeaturesOfInterest", HandleGetFeatureOfInterests},
			{models.HTTPOperationGet, "/v1.0/FeaturesOfInterest{id}", HandleGetFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/FeatureOfInterest{id}/Observations", HandleGetObservationsByFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/FeatureOfInterest{id}/Observations/{params}", HandleGetObservationsByFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/FeaturesOfInterest{id}/{params}", HandleGetFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/FeaturesOfInterest{id}/{params}/$value", HandleGetFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/FeaturesOfInterest/{params}", HandleGetFeatureOfInterests},

			{models.HTTPOperationPost, "/v1.0/FeaturesOfInterest", HandlePostFeatureOfInterest},
			{models.HTTPOperationDelete, "/v1.0/FeaturesOfInterest{id}", HandleDeleteFeatureOfInterest},
			{models.HTTPOperationPatch, "/v1.0/FeaturesOfInterest{id}", HandlePatchFeatureOfInterest},
			{models.HTTPOperationPut, "/v1.0/FeaturesOfInterest{id}", HandlePutFeatureOfInterest},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/FeaturesOfInterest", HandleGetFeatureOfInterests},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/FeaturesOfInterest{id}", HandleGetFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/FeatureOfInterest{id}/Observations", HandleGetObservationsByFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/FeatureOfInterest{id}/Observations/{params}", HandleGetObservationsByFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/FeaturesOfInterest{id}/{params}", HandleGetFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/FeaturesOfInterest{id}/{params}/$value", HandleGetFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/FeaturesOfInterest/{params}", HandleGetFeatureOfInterests},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/FeaturesOfInterest", HandlePostFeatureOfInterest},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/FeaturesOfInterest{id}", HandleDeleteFeatureOfInterest},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/FeaturesOfInterest{id}", HandlePatchFeatureOfInterest},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/FeaturesOfInterest{id}", HandlePutFeatureOfInterest},
		},
	}
}
