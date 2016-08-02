package rest

import (
	"fmt"

	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

func createObservationsEndpoint(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "Observations",
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "Observations")),
		SupportedQueryOptions: []odata.QueryOptionType{
			odata.QueryOptionTop, odata.QueryOptionSkip, odata.QueryOptionOrderBy, odata.QueryOptionCount, odata.QueryOptionResultFormat,
			odata.QueryOptionExpand, odata.QueryOptionSelect, odata.QueryOptionFilter,
		},
		SupportedExpandParams: []string{
			"Datastream",
			"FeatureOfInterest",
		},
		SupportedSelectParams: []string{
			"description",
			"encodingType",
			"feature",
			"Observations",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/Observations", HandleGetObservations},
			{models.HTTPOperationGet, "/v1.0/Observations{id}", HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/Observations{id}/Datastream", HandleGetDatastreamByObservation},
			{models.HTTPOperationGet, "/v1.0/Observations{id}/FeatureOfInterest", HandleGetFeatureOfInterestByObservation},
			{models.HTTPOperationGet, "/v1.0/Observations{id}/Datastream/{params}", HandleGetDatastreamByObservation},
			{models.HTTPOperationGet, "/v1.0/Observations{id}/Datastream/{params}/$value", HandleGetDatastreamByObservation},
			{models.HTTPOperationGet, "/v1.0/Observations{id}/FeatureOfInterest/{params}", HandleGetFeatureOfInterestByObservation},
			{models.HTTPOperationGet, "/v1.0/Observations{id}/FeatureOfInterest/{params}/$value", HandleGetFeatureOfInterestByObservation},
			{models.HTTPOperationGet, "/v1.0/Observations{id}/{params}", HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/Observations{id}/{params}/$value", HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/Observations/{params}", HandleGetObservations},

			{models.HTTPOperationPost, "/v1.0/Observations", HandlePostObservation},
			{models.HTTPOperationPost, "/v1.0/Datastreams{id}/Observations", HandlePostObservationByDatastream},
			{models.HTTPOperationDelete, "/v1.0/Observations{id}", HandleDeleteObservation},
			{models.HTTPOperationPatch, "/v1.0/Observations{id}", HandlePatchObservation},
			{models.HTTPOperationPut, "/v1.0/Observations{id}", HandlePutObservation},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/Observations", HandleGetObservations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Observations{id}", HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Observations{id}/Datastream", HandleGetDatastreamByObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Observations{id}/FeatureOfInterest", HandleGetFeatureOfInterestByObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Observations{id}/Datastream/{params}", HandleGetDatastreamByObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Observations{id}/Datastream/{params}/$value", HandleGetDatastreamByObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Observations{id}/FeatureOfInterest/{params}", HandleGetFeatureOfInterestByObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Observations{id}/FeatureOfInterest/{params}/$value", HandleGetFeatureOfInterestByObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Observations{id}/{params}", HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Observations{id}/{params}/$value", HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Observations/{params}", HandleGetObservations},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/Observations", HandlePostObservation},
			{models.HTTPOperationPost, "/v1.0/{c:.*}/Datastreams{id}/Observations", HandlePostObservationByDatastream},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/Observations{id}", HandleDeleteObservation},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/Observations{id}", HandlePatchObservation},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/Observations{id}", HandlePutObservation},
		},
	}
}
