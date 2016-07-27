package rest

import (
	"fmt"

	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// CreateEndPoints creates the pre-defined endpoint config, the config contains all endpoint info
// describing the SupportedQueryOptions (if needed) and EndpointOperation for each endpoint
// parameter externalURL is the URL where the GOST service can be reached, main endpoint urls
// are generated based upon this URL
func CreateEndPoints(externalURL string) []models.Endpoint {
	endpoints := []models.Endpoint{
		createVersion(externalURL),
		createRoot(externalURL),
		createThings(externalURL),
		createDatastreams(externalURL),
		createObservedProperties(externalURL),
		createLocations(externalURL),
		createSensors(externalURL),
		createObservations(externalURL),
		createFeaturesOfInterest(externalURL),
		createHistoricalLocations(externalURL),
	}

	return endpoints
}

func createRoot(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "Root",
		OutputInfo: false,
		URL:        fmt.Sprintf("%s/%s", externalURL, "v1.0"),
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0", HandleAPIRoot},
		},
	}
}

func createVersion(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "Version",
		OutputInfo: false,
		URL:        fmt.Sprintf("%s/%s", externalURL, "Version"),
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/Version", HandleVersion},
		},
	}
}

func createThings(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "Things",
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "Things")),
		SupportedQueryOptions: []odata.QueryOptionType{
			odata.QueryOptionTop, odata.QueryOptionSkip, odata.QueryOptionOrderBy, odata.QueryOptionCount, odata.QueryOptionResultFormat,
			odata.QueryOptionExpand, odata.QueryOptionSelect, odata.QueryOptionFilter,
		},
		SupportedExpandParams: []string{
			"Locations",
			"Datastreams",
			"HistoricalLocations",
		},
		SupportedSelectParams: []string{
			"properties",
			"description",
			"Locations",
			"Datastreams",
			"HistoricalLocations",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/Things", HandleGetThings},
			{models.HTTPOperationGet, "/v1.0/Things{id}", HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/Things{id}/Datastreams", HandleGetDatastreamsByThing},
			{models.HTTPOperationGet, "/v1.0/Things{id}/HistoricalLocations", HandleGetHistoricalLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/Things{id}/Locations", HandleGetLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/Things{id}/Datastreams/{params}", HandleGetDatastreamsByThing},
			{models.HTTPOperationGet, "/v1.0/Things{id}/HistoricalLocations/{params}", HandleGetHistoricalLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/Things{id}/Locations/{params}", HandleGetLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/Things{id}/{params}", HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/Things{id}/{params}/$value", HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/Things/{params}", HandleGetThings},

			{models.HTTPOperationPost, "/v1.0/Things", HandlePostThing},
			{models.HTTPOperationDelete, "/v1.0/Things{id}", HandleDeleteThing},
			{models.HTTPOperationPatch, "/v1.0/Things{id}", HandlePatchThing},
			{models.HTTPOperationPut, "/v1.0/Things{id}", HandlePutThing},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things", HandleGetThings},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things{id}", HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things{id}/Datastreams", HandleGetDatastreamsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things{id}/HistoricalLocations", HandleGetHistoricalLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things{id}/Locations", HandleGetLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things{id}/Datastreams/{params}", HandleGetDatastreamsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things{id}/HistoricalLocations/{params}", HandleGetHistoricalLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things{id}/Locations/{params}", HandleGetLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things{id}/{params}", HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things{id}/{params}/$value", HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Things/{params}", HandleGetThings},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/Things", HandlePostThing},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/Things{id}", HandleDeleteThing},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/Things{id}", HandlePatchThing},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/Things{id}", HandlePutThing},
		},
	}
}

func createDatastreams(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "Datastreams",
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "Datastreams")),
		SupportedQueryOptions: []odata.QueryOptionType{
			odata.QueryOptionTop, odata.QueryOptionSkip, odata.QueryOptionOrderBy, odata.QueryOptionCount, odata.QueryOptionResultFormat,
			odata.QueryOptionExpand, odata.QueryOptionSelect, odata.QueryOptionFilter,
		},
		SupportedExpandParams: []string{
			"Thing",
			"Sensor",
			"Observedproperty",
			"Observations",
		},
		SupportedSelectParams: []string{
			"description",
			"unitOfMeasurement",
			"observationType",
			"observedArea",
			"phenomenonTime",
			"resultTime",
			"Thing",
			"Sensor",
			"ObservedProperty",
			"Observations",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/Datastreams", HandleGetDatastreams},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/ObservedProperty", HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Observations", HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Sensor", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Thing", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/ObservedProperty/{params}", HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Observations/{params}", HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Sensor/{params}", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Sensor/{params}/$value", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Thing/{params}", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Thing/{params}/$value", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/{params}", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/{params}/$value", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams/{params}", HandleGetDatastreams},

			{models.HTTPOperationPost, "/v1.0/Datastreams", HandlePostDatastream},
			{models.HTTPOperationPost, "/v1.0/Things{id}/Datastreams", HandlePostDatastreamByThing},
			{models.HTTPOperationDelete, "/v1.0/Datastreams{id}", HandleDeleteDatastream},
			{models.HTTPOperationPatch, "/v1.0/Datastreams{id}", HandlePatchDatastream},
			{models.HTTPOperationPut, "/v1.0/Datastreams{id}", HandlePatchDatastream},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams", HandleGetDatastreams},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/ObservedProperty", HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/Observations", HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/Sensor", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/Thing", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/ObservedProperty/{params}", HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/Observations/{params}", HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/Sensor/{params}", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/Sensor/{params}/$value", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/Thing/{params}", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/Thing/{params}/$value", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/{params}", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams{id}/{params}/$value", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Datastreams/{params}", HandleGetDatastreams},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/Datastreams", HandlePostDatastream},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/Datastreams{id}", HandleDeleteDatastream},
			{models.HTTPOperationPost, "/v1.0/{c:.*}/Things{id}/Datastreams", HandlePostDatastreamByThing},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/Datastreams{id}", HandlePatchDatastream},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/Datastreams{id}", HandlePatchDatastream},
		},
	}
}

func createObservedProperties(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "ObservedProperties",
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "ObservedProperties")),
		SupportedQueryOptions: []odata.QueryOptionType{
			odata.QueryOptionTop, odata.QueryOptionSkip, odata.QueryOptionOrderBy, odata.QueryOptionCount, odata.QueryOptionResultFormat,
			odata.QueryOptionExpand, odata.QueryOptionSelect, odata.QueryOptionFilter,
		},
		SupportedExpandParams: []string{
			"Datastreams",
		},
		SupportedSelectParams: []string{
			"name",
			"definition",
			"description",
			"Datastreams",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/ObservedProperties", HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/ObservedProperties{id}", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/ObservedProperties{id}/Datastreams", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/ObservedProperties{id}/Datastreams/{params}", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/ObservedProperties{id}/Datastreams/{params}/$value", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/ObservedProperties{id}/{params}", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/ObservedProperties{id}/{params}/$value", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/ObservedProperties/{params}", HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/ObservedProperties/{params}/$value", HandleGetObservedProperties},

			{models.HTTPOperationPost, "/v1.0/ObservedProperties", HandlePostObservedProperty},
			{models.HTTPOperationDelete, "/v1.0/ObservedProperties{id}", HandleDeleteObservedProperty},
			{models.HTTPOperationPatch, "/v1.0/ObservedProperties{id}", HandlePatchObservedProperty},
			{models.HTTPOperationPut, "/v1.0/ObservedProperties{id}", HandlePatchObservedProperty},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/ObservedProperties", HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/ObservedProperties{id}", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/ObservedProperties{id}/Datastreams", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/ObservedProperties{id}/Datastreams/{params}", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/ObservedProperties{id}/Datastreams/{params}/$value", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/ObservedProperties{id}/{params}", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/ObservedProperties{id}/{params}/$value", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/ObservedProperties/{params}", HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/ObservedProperties/{params}/$value", HandleGetObservedProperties},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/ObservedProperties", HandlePostObservedProperty},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/ObservedProperties{id}", HandleDeleteObservedProperty},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/ObservedProperties{id}", HandlePatchObservedProperty},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/ObservedProperties{id}", HandlePatchObservedProperty},
		},
	}
}

func createLocations(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "Locations",
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "Locations")),
		SupportedQueryOptions: []odata.QueryOptionType{
			odata.QueryOptionTop, odata.QueryOptionSkip, odata.QueryOptionOrderBy, odata.QueryOptionCount, odata.QueryOptionResultFormat,
			odata.QueryOptionExpand, odata.QueryOptionSelect, odata.QueryOptionFilter,
		},
		SupportedExpandParams: []string{
			"Things",
			"HistoricalLocations",
		},
		SupportedSelectParams: []string{
			"description",
			"encodingType",
			"location",
			"Things",
			"HistoricalLocations",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/Locations", HandleGetLocations},
			{models.HTTPOperationGet, "/v1.0/Locations{id}", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/Locations{id}/Things", HandleGetThingsByLocation},
			{models.HTTPOperationGet, "/v1.0/Locations{id}/HistoricalLocations", HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/Locations{id}/Things/{params}", HandleGetThingsByLocation},
			{models.HTTPOperationGet, "/v1.0/Locations{id}/HistoricalLocations/{params}", HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/Locations{id}/HistoricalLocations/{params}/$value", HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/Locations{id}/{params}", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/Locations{id}/{params}/$value", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/Locations/{params}", HandleGetLocations},

			{models.HTTPOperationPost, "/v1.0/Locations", HandlePostLocation},
			{models.HTTPOperationPost, "/v1.0/Things{id}/Locations", HandlePostLocationByThing},
			{models.HTTPOperationDelete, "/v1.0/Locations{id}", HandleDeleteLocation},
			{models.HTTPOperationPatch, "/v1.0/Locations{id}", HandlePatchLocation},
			{models.HTTPOperationPut, "/v1.0/Locations{id}", HandlePatchLocation},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/Locations", HandleGetLocations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Locations{id}", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Locations{id}/Things", HandleGetThingsByLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Locations{id}/HistoricalLocations", HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Locations{id}/Things/{params}", HandleGetThingsByLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Locations{id}/HistoricalLocations/{params}", HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Locations{id}/HistoricalLocations/{params}/$value", HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Locations{id}/{params}", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Locations{id}/{params}/$value", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Locations/{params}", HandleGetLocations},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/Locations", HandlePostLocation},
			{models.HTTPOperationPost, "/v1.0/{c:.*}/Things{id}/Locations", HandlePostLocationByThing},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/Locations{id}", HandleDeleteLocation},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/Locations{id}", HandlePatchLocation},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/Locations{id}", HandlePutLocation},
		},
	}
}

func createSensors(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "Sensors",
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "Sensors")),
		SupportedQueryOptions: []odata.QueryOptionType{
			odata.QueryOptionTop, odata.QueryOptionSkip, odata.QueryOptionOrderBy, odata.QueryOptionCount, odata.QueryOptionResultFormat,
			odata.QueryOptionExpand, odata.QueryOptionSelect, odata.QueryOptionFilter,
		},
		SupportedExpandParams: []string{
			"Datastream",
		},
		SupportedSelectParams: []string{
			"description",
			"encodingType",
			"metadata",
			"Datastreams",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/Sensors", HandleGetSensors},
			{models.HTTPOperationGet, "/v1.0/Sensors{id}", HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/Sensors{id}/Datastreams", HandleGetDatastreamsBySensor},
			{models.HTTPOperationGet, "/v1.0/Sensors{id}/Datastreams/{params}", HandleGetDatastreamsBySensor},
			{models.HTTPOperationGet, "/v1.0/Sensors{id}/{params}", HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/Sensors{id}/{params}/$value", HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/Sensors/{params}", HandleGetSensors},

			{models.HTTPOperationPost, "/v1.0/Sensors", HandlePostSensors},
			{models.HTTPOperationDelete, "/v1.0/Sensors{id}", HandleDeleteSensor},
			{models.HTTPOperationPatch, "/v1.0/Sensors{id}", HandlePatchSensor},
			{models.HTTPOperationPut, "/v1.0/Sensors{id}", HandlePatchSensor},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/Sensors", HandleGetSensors},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Sensors{id}", HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Sensors{id}/Datastreams", HandleGetDatastreamsBySensor},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Sensors{id}/Datastreams/{params}", HandleGetDatastreamsBySensor},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Sensors{id}/{params}", HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Sensors{id}/{params}/$value", HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/Sensors/{params}", HandleGetSensors},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/Sensors", HandlePostSensors},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/Sensors{id}", HandleDeleteSensor},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/Sensors{id}", HandlePatchSensor},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/Sensors{id}", HandlePatchSensor},
		},
	}
}

func createObservations(externalURL string) *Endpoint {
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
			{models.HTTPOperationPut, "/v1.0/Observations{id}", HandlePatchObservation},

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
			{models.HTTPOperationPut, "/v1.0/{c:.*}/Observations{id}", HandlePatchObservation},
		},
	}
}

func createFeaturesOfInterest(externalURL string) *Endpoint {
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
			{models.HTTPOperationPut, "/v1.0/FeaturesOfInterest{id}", HandlePatchFeatureOfInterest},

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
			{models.HTTPOperationPut, "/v1.0/{c:.*}/FeaturesOfInterest{id}", HandlePatchFeatureOfInterest},
		},
	}
}

func createHistoricalLocations(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "HistoricalLocations",
		OutputInfo: true,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "HistoricalLocations")),
		SupportedQueryOptions: []odata.QueryOptionType{
			odata.QueryOptionTop, odata.QueryOptionSkip, odata.QueryOptionOrderBy, odata.QueryOptionCount, odata.QueryOptionResultFormat,
			odata.QueryOptionExpand, odata.QueryOptionSelect, odata.QueryOptionFilter,
		},
		SupportedExpandParams: []string{
			"locations",
			"thing",
		},
		SupportedSelectParams: []string{
			"time",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations", HandleGetHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}/Locations", HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}/Thing", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}/Thing/{params}", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}/Thing/{params}/$value", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}/Locations/{params}", HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}/{params}", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}/{params}/$value", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations/{params}", HandleGetHistoricalLocations},

			{models.HTTPOperationPost, "/v1.0/HistoricalLocations", HandlePostHistoricalLocation},
			{models.HTTPOperationDelete, "/v1.0/HistoricalLocations{id}", HandleDeleteHistoricalLocations},
			{models.HTTPOperationPatch, "/v1.0/HistoricalLocations{id}", HandlePatchHistoricalLocations},
			{models.HTTPOperationPut, "/v1.0/HistoricalLocations{id}", HandlePatchHistoricalLocations},

			{models.HTTPOperationGet, "/v1.0/{c:.*}/HistoricalLocations", HandleGetHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/HistoricalLocations{id}", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/HistoricalLocations{id}/Locations", HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/HistoricalLocations{id}/Thing", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/HistoricalLocations{id}/Thing/{params}", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/HistoricalLocations{id}/Thing/{params}/$value", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/HistoricalLocations{id}/Locations/{params}", HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/HistoricalLocations{id}/{params}", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/HistoricalLocations{id}/{params}/$value", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/{c:.*}/HistoricalLocations/{params}", HandleGetHistoricalLocations},

			{models.HTTPOperationPost, "/v1.0/{c:.*}/HistoricalLocations", HandlePostHistoricalLocation},
			{models.HTTPOperationDelete, "/v1.0/{c:.*}/HistoricalLocations{id}", HandleDeleteHistoricalLocations},
			{models.HTTPOperationPatch, "/v1.0/{c:.*}/HistoricalLocations{id}", HandlePatchHistoricalLocations},
			{models.HTTPOperationPut, "/v1.0/{c:.*}/HistoricalLocations{id}", HandlePatchHistoricalLocations},
		},
	}
}
