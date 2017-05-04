package rest

import (
	"fmt"

	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
)

func createObservedPropertiesEndpoint(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "ObservedProperties",
		EntityType: entities.EntityTypeObservedProperty,
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "ObservedProperties")),
		SupportedExpandParams: []string{
			"datastreams",
		},
		SupportedSelectParams: []string{
			"id",
			"name",
			"definition",
			"description",
			"datastreams",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/observedproperties", HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/observedproperties{id}", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/observedproperty", HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/observedproperty/{params}", HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/observedproperties{id}/{params}", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/observedproperties{id}/{params}/$value", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/observedproperties/{params}", HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/observedproperties/{params}/$value", HandleGetObservedProperties},

			{models.HTTPOperationPost, "/v1.0/observedproperties", HandlePostObservedProperty},
			{models.HTTPOperationDelete, "/v1.0/observedproperties{id}", HandleDeleteObservedProperty},
			{models.HTTPOperationPatch, "/v1.0/observedproperties{id}", HandlePatchObservedProperty},
			{models.HTTPOperationPut, "/v1.0/observedproperties{id}", HandlePutObservedProperty},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties", HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties{id}", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/observedproperty", HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/observedproperty/{params}", HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties{id}/{params}", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties{id}/{params}/$value", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties/{params}", HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties/{params}/$value", HandleGetObservedProperties},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/observedproperties", HandlePostObservedProperty},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/observedproperties{id}", HandleDeleteObservedProperty},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/observedproperties{id}", HandlePatchObservedProperty},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/observedproperties{id}", HandlePutObservedProperty},
		},
	}
}
