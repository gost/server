package config

import (
	"fmt"

	entities "github.com/gost/core"
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
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/observedproperties", Handler: handlers.HandleGetObservedProperties},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/observedproperties{id}", Handler: handlers.HandleGetObservedProperty},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams{id}/observedproperty", Handler: handlers.HandleGetObservedPropertyByDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams{id}/observedproperty/{params}", Handler: handlers.HandleGetObservedPropertyByDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/observedproperties{id}/{params}", Handler: handlers.HandleGetObservedProperty},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/observedproperties{id}/{params}/$value", Handler: handlers.HandleGetObservedProperty},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/observedproperties/{params}", Handler: handlers.HandleGetObservedProperties},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/observedproperties/{params}/$value", Handler: handlers.HandleGetObservedProperties},

			{OperationType: models.HTTPOperationPost, Path: "/v1.0/observedproperties", Handler: handlers.HandlePostObservedProperty},
			{OperationType: models.HTTPOperationDelete, Path: "/v1.0/observedproperties{id}", Handler: handlers.HandleDeleteObservedProperty},
			{OperationType: models.HTTPOperationPatch, Path: "/v1.0/observedproperties{id}", Handler: handlers.HandlePatchObservedProperty},
			{OperationType: models.HTTPOperationPut, Path: "/v1.0/observedproperties{id}", Handler: handlers.HandlePutObservedProperty},

			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/observedproperties", Handler: handlers.HandleGetObservedProperties},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/observedproperties{id}", Handler: handlers.HandleGetObservedProperty},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/datastreams{id}/observedproperty", Handler: handlers.HandleGetObservedPropertyByDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/datastreams{id}/observedproperty/{params}", Handler: handlers.HandleGetObservedPropertyByDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/observedproperties{id}/{params}", Handler: handlers.HandleGetObservedProperty},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/observedproperties{id}/{params}/$value", Handler: handlers.HandleGetObservedProperty},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/observedproperties/{params}", Handler: handlers.HandleGetObservedProperties},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/observedproperties/{params}/$value", Handler: handlers.HandleGetObservedProperties},

			{OperationType: models.HTTPOperationPost, Path: "/v1.0/{c:.*}/observedproperties", Handler: handlers.HandlePostObservedProperty},
			{OperationType: models.HTTPOperationDelete, Path: "/v1.0/{c:.*}/observedproperties{id}", Handler: handlers.HandleDeleteObservedProperty},
			{OperationType: models.HTTPOperationPatch, Path: "/v1.0/{c:.*}/observedproperties{id}", Handler: handlers.HandlePatchObservedProperty},
			{OperationType: models.HTTPOperationPut, Path: "/v1.0/{c:.*}/observedproperties{id}", Handler: handlers.HandlePutObservedProperty},
		},
	}
}
