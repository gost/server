package config

import (
	"fmt"

	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/rest/endpoint"
	"github.com/geodan/gost/sensorthings/rest/handlers"
)

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
			{models.HTTPOperationGet, "/v1.0/datastreams", handlers.HandleGetDatastreams},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}", handlers.HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/observedproperties{id}/datastreams", handlers.HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/observedproperties{id}/datastreams/{params}", handlers.HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/observedproperties{id}/datastreams/{params}/$value", handlers.HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/observations{id}/datastream", handlers.HandleGetDatastreamByObservation},
			{models.HTTPOperationGet, "/v1.0/sensors{id}/datastreams", handlers.HandleGetDatastreamsBySensor},
			{models.HTTPOperationGet, "/v1.0/sensors{id}/datastreams/{params}", handlers.HandleGetDatastreamsBySensor},
			{models.HTTPOperationGet, "/v1.0/things{id}/datastreams", handlers.HandleGetDatastreamsByThing},
			{models.HTTPOperationGet, "/v1.0/observations{id}/datastream/{params}", handlers.HandleGetDatastreamByObservation},
			{models.HTTPOperationGet, "/v1.0/observations{id}/datastream/{params}/$value", handlers.HandleGetDatastreamByObservation},
			{models.HTTPOperationGet, "/v1.0/things{id}/datastreams/{params}", handlers.HandleGetDatastreamsByThing},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/{params}", handlers.HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/{params}/$value", handlers.HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams/{params}", handlers.HandleGetDatastreams},

			{models.HTTPOperationPost, "/v1.0/datastreams", handlers.HandlePostDatastream},
			{models.HTTPOperationPost, "/v1.0/things{id}/datastreams", handlers.HandlePostDatastreamByThing},
			{models.HTTPOperationDelete, "/v1.0/datastreams{id}", handlers.HandleDeleteDatastream},
			{models.HTTPOperationPatch, "/v1.0/datastreams{id}", handlers.HandlePatchDatastream},
			{models.HTTPOperationPut, "/v1.0/datastreams{id}", handlers.HandlePutDatastream},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams", handlers.HandleGetDatastreams},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}", handlers.HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties{id}/datastreams", handlers.HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties{id}/datastreams/{params}", handlers.HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties{id}/datastreams/{params}/$value", handlers.HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations{id}/datastream", handlers.HandleGetDatastreamByObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations{id}/datastream/{params}", handlers.HandleGetDatastreamByObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations{id}/datastream/{params}/$value", handlers.HandleGetDatastreamByObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/sensors{id}/datastreams", handlers.HandleGetDatastreamsBySensor},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/sensors{id}/datastreams/{params}", handlers.HandleGetDatastreamsBySensor},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/things{id}/datastreams", handlers.HandleGetDatastreamsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/things{id}/datastreams/{params}", handlers.HandleGetDatastreamsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/{params}", handlers.HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/{params}/$value", handlers.HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams/{params}", handlers.HandleGetDatastreams},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/datastreams", handlers.HandlePostDatastream},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/datastreams{id}", handlers.HandleDeleteDatastream},
			{models.HTTPOperationPost, "/v1.0/{c:.*}/things{id}/datastreams", handlers.HandlePostDatastreamByThing},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/datastreams{id}", handlers.HandlePatchDatastream},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/datastreams{id}", handlers.HandlePutDatastream},
		},
	}
}
