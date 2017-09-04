package config

import (
	"fmt"

	entities "github.com/gost/core"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/rest/endpoint"
	"github.com/gost/server/sensorthings/rest/handlers"
)

// CreateDatastreamsEndpoint constructs the Datastreams endpoint configuration
func CreateDatastreamsEndpoint(externalURL string) *endpoint.Endpoint {
	return &endpoint.Endpoint{
		Name:       "Datastreams",
		EntityType: entities.EntityTypeDatastream,
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "Datastreams")),
		SupportedExpandParams: []string{
			"thing",
			"sensor",
			"observedproperty",
			"observations",
		},
		SupportedSelectParams: []string{
			"id",
			"name",
			"description",
			"unitofmeasurement",
			"observationtype",
			"observedarea",
			"phenomenontime",
			"resulttime",
			"thing",
			"sensor",
			"observedproperty",
			"observations",
		},
		Operations: []models.EndpointOperation{
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams", Handler: handlers.HandleGetDatastreams},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams{id}", Handler: handlers.HandleGetDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/observedproperties{id}/datastreams", Handler: handlers.HandleGetDatastreamsByObservedProperty},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/observedproperties{id}/datastreams/{params}", Handler: handlers.HandleGetDatastreamsByObservedProperty},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/observedproperties{id}/datastreams/{params}/$value", Handler: handlers.HandleGetDatastreamsByObservedProperty},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/observations{id}/datastream", Handler: handlers.HandleGetDatastreamByObservation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/sensors{id}/datastreams", Handler: handlers.HandleGetDatastreamsBySensor},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/sensors{id}/datastreams/{params}", Handler: handlers.HandleGetDatastreamsBySensor},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/things{id}/datastreams", Handler: handlers.HandleGetDatastreamsByThing},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/observations{id}/datastream/{params}", Handler: handlers.HandleGetDatastreamByObservation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/observations{id}/datastream/{params}/$value", Handler: handlers.HandleGetDatastreamByObservation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/things{id}/datastreams/{params}", Handler: handlers.HandleGetDatastreamsByThing},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams{id}/{params}", Handler: handlers.HandleGetDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams{id}/{params}/$value", Handler: handlers.HandleGetDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/datastreams/{params}", Handler: handlers.HandleGetDatastreams},

			{OperationType: models.HTTPOperationPost, Path: "/v1.0/datastreams", Handler: handlers.HandlePostDatastream},
			{OperationType: models.HTTPOperationPost, Path: "/v1.0/things{id}/datastreams", Handler: handlers.HandlePostDatastreamByThing},
			{OperationType: models.HTTPOperationDelete, Path: "/v1.0/datastreams{id}", Handler: handlers.HandleDeleteDatastream},
			{OperationType: models.HTTPOperationPatch, Path: "/v1.0/datastreams{id}", Handler: handlers.HandlePatchDatastream},
			{OperationType: models.HTTPOperationPut, Path: "/v1.0/datastreams{id}", Handler: handlers.HandlePutDatastream},

			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/datastreams", Handler: handlers.HandleGetDatastreams},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/datastreams{id}", Handler: handlers.HandleGetDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/observedproperties{id}/datastreams", Handler: handlers.HandleGetDatastreamsByObservedProperty},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/observedproperties{id}/datastreams/{params}", Handler: handlers.HandleGetDatastreamsByObservedProperty},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/observedproperties{id}/datastreams/{params}/$value", Handler: handlers.HandleGetDatastreamsByObservedProperty},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/observations{id}/datastream", Handler: handlers.HandleGetDatastreamByObservation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/observations{id}/datastream/{params}", Handler: handlers.HandleGetDatastreamByObservation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/observations{id}/datastream/{params}/$value", Handler: handlers.HandleGetDatastreamByObservation},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/sensors{id}/datastreams", Handler: handlers.HandleGetDatastreamsBySensor},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/sensors{id}/datastreams/{params}", Handler: handlers.HandleGetDatastreamsBySensor},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/things{id}/datastreams", Handler: handlers.HandleGetDatastreamsByThing},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/things{id}/datastreams/{params}", Handler: handlers.HandleGetDatastreamsByThing},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/datastreams{id}/{params}", Handler: handlers.HandleGetDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/datastreams{id}/{params}/$value", Handler: handlers.HandleGetDatastream},
			{OperationType: models.HTTPOperationGet, Path: "/v1.0/{c:.*}/datastreams/{params}", Handler: handlers.HandleGetDatastreams},

			{OperationType: models.HTTPOperationPost, Path: "/v1.0/{c:.*}/datastreams", Handler: handlers.HandlePostDatastream},
			{OperationType: models.HTTPOperationDelete, Path: "/v1.0/{c:.*}/datastreams{id}", Handler: handlers.HandleDeleteDatastream},
			{OperationType: models.HTTPOperationPost, Path: "/v1.0/{c:.*}/things{id}/datastreams", Handler: handlers.HandlePostDatastreamByThing},
			{OperationType: models.HTTPOperationPatch, Path: "/v1.0/{c:.*}/datastreams{id}", Handler: handlers.HandlePatchDatastream},
			{OperationType: models.HTTPOperationPut, Path: "/v1.0/{c:.*}/datastreams{id}", Handler: handlers.HandlePutDatastream},
		},
	}
}
