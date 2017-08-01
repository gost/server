package config

import (
	"fmt"

	"github.com/gost/server/sensorthings/entities"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/rest/endpoint"
	"github.com/gost/server/sensorthings/rest/handlers"
)

// CreateObservedPropertiesEndpoint constructs the ObservedProperties endpoint configuration
func CreateObservedPropertiesEndpoint(externalURL string) *endpoint.Endpoint {
	return &endpoint.Endpoint{
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
			{models.HTTPOperationGet, "/v1.0/observedproperties", handlers.HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/observedproperties{id}", handlers.HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/observedproperty", handlers.HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/observedproperty/{params}", handlers.HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/observedproperties{id}/{params}", handlers.HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/observedproperties{id}/{params}/$value", handlers.HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/observedproperties/{params}", handlers.HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/observedproperties/{params}/$value", handlers.HandleGetObservedProperties},

			{models.HTTPOperationPost, "/v1.0/observedproperties", handlers.HandlePostObservedProperty},
			{models.HTTPOperationDelete, "/v1.0/observedproperties{id}", handlers.HandleDeleteObservedProperty},
			{models.HTTPOperationPatch, "/v1.0/observedproperties{id}", handlers.HandlePatchObservedProperty},
			{models.HTTPOperationPut, "/v1.0/observedproperties{id}", handlers.HandlePutObservedProperty},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties", handlers.HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties{id}", handlers.HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/observedproperty", handlers.HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/observedproperty/{params}", handlers.HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties{id}/{params}", handlers.HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties{id}/{params}/$value", handlers.HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties/{params}", handlers.HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties/{params}/$value", handlers.HandleGetObservedProperties},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/observedproperties", handlers.HandlePostObservedProperty},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/observedproperties{id}", handlers.HandleDeleteObservedProperty},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/observedproperties{id}", handlers.HandlePatchObservedProperty},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/observedproperties{id}", handlers.HandlePutObservedProperty},
		},
	}
}
