package rest

import (
	"fmt"

	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
)

func createDatastreamsEndpoint(externalURL string) *Endpoint {
	return &Endpoint{
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
			{models.HTTPOperationGet, "/v1.0/datastreams", HandleGetDatastreams},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/observedproperties{id}/datastreams", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/observedproperties{id}/datastreams/{params}", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/observedproperties{id}/datastreams/{params}/$value", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/observations{id}/datastream", HandleGetDatastreamByObservation},
			{models.HTTPOperationGet, "/v1.0/sensors{id}/datastreams", HandleGetDatastreamsBySensor},
			{models.HTTPOperationGet, "/v1.0/sensors{id}/datastreams/{params}", HandleGetDatastreamsBySensor},
			{models.HTTPOperationGet, "/v1.0/things{id}/datastreams", HandleGetDatastreamsByThing},
			{models.HTTPOperationGet, "/v1.0/observations{id}/datastream/{params}", HandleGetDatastreamByObservation},
			{models.HTTPOperationGet, "/v1.0/observations{id}/datastream/{params}/$value", HandleGetDatastreamByObservation},
			{models.HTTPOperationGet, "/v1.0/things{id}/datastreams/{params}", HandleGetDatastreamsByThing},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/{params}", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams{id}/{params}/$value", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/datastreams/{params}", HandleGetDatastreams},

			{models.HTTPOperationPost, "/v1.0/datastreams", HandlePostDatastream},
			{models.HTTPOperationPost, "/v1.0/things{id}/datastreams", HandlePostDatastreamByThing},
			{models.HTTPOperationDelete, "/v1.0/datastreams{id}", HandleDeleteDatastream},
			{models.HTTPOperationPatch, "/v1.0/datastreams{id}", HandlePatchDatastream},
			{models.HTTPOperationPut, "/v1.0/datastreams{id}", HandlePutDatastream},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams", HandleGetDatastreams},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties{id}/datastreams", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties{id}/datastreams/{params}", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observedproperties{id}/datastreams/{params}/$value", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations{id}/datastream", HandleGetDatastreamByObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations{id}/datastream/{params}", HandleGetDatastreamByObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/observations{id}/datastream/{params}/$value", HandleGetDatastreamByObservation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/sensors{id}/datastreams", HandleGetDatastreamsBySensor},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/sensors{id}/datastreams/{params}", HandleGetDatastreamsBySensor},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/things{id}/datastreams", HandleGetDatastreamsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/things{id}/datastreams/{params}", HandleGetDatastreamsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/{params}", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams{id}/{params}/$value", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/datastreams/{params}", HandleGetDatastreams},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/datastreams", HandlePostDatastream},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/datastreams{id}", HandleDeleteDatastream},
			{models.HTTPOperationPost, "/v1.0/{c:.*}/things{id}/datastreams", HandlePostDatastreamByThing},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/datastreams{id}", HandlePatchDatastream},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/datastreams{id}", HandlePutDatastream},
		},
	}
}
