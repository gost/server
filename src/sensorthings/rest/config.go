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
			odata.QueryOptionExpand, odata.QueryOptionSelect,
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
			{models.HTTPOperationGet, "/v1.0/Things/{params}", HandleGetThings},
			{models.HTTPOperationGet, "/v1.0/Things{id}", HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/Things{id}/{params}", HandleGetThing},
			{models.HTTPOperationGet, "/v1.0/Locations{id}/Things", HandleGetThingsByLocation},
			{models.HTTPOperationGet, "/v1.0/Locations{id}/Things/{params}", HandleGetThingsByLocation},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}/Thing", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}/Thing/{params}", HandleGetThingByHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Thing", HandleGetThingByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Thing/{params}", HandleGetThingByDatastream},
			{models.HTTPOperationPost, "/v1.0/Things", HandlePostThing},
			{models.HTTPOperationDelete, "/v1.0/Things{id}", HandleDeleteThing},
			{models.HTTPOperationPatch, "/v1.0/Things{id}", HandlePatchThing},
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
			odata.QueryOptionExpand, odata.QueryOptionSelect,
		},
		SupportedExpandParams: []string{
			"Thing",
			"Sensor",
			"Observedproperty",
			"Observations",
		},
		SupportedSelectParams: []string{
			"description",
			"unitofmeasurement",
			"observationtype",
			"observedarea",
			"phenomenontime",
			"resulttime",
			"Thing",
			"Sensor",
			"ObservedProperty",
			"Observations",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/Datastreams", HandleGetDatastreams},
			{models.HTTPOperationGet, "/v1.0/Datastreams/{params}", HandleGetDatastreams},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/{params}", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/Observations{id}/Datastream", HandleGetDatastreamByObservation},
			{models.HTTPOperationGet, "/v1.0/Observations{id}/Datastream/{params}", HandleGetDatastreamByObservation},
			{models.HTTPOperationGet, "/v1.0/Things{id}/Datastreams", HandleGetDatastreamsByThing},
			{models.HTTPOperationGet, "/v1.0/Things{id}/Datastreams/{params}", HandleGetDatastreamsByThing},
			{models.HTTPOperationGet, "/v1.0/Sensors{id}/Datastreams", HandleGetDatastreamsBySensor},
			{models.HTTPOperationGet, "/v1.0/Sensors{id}/Datastreams/{params}", HandleGetDatastreamsBySensor},
			{models.HTTPOperationGet, "/v1.0/ObservedProperties{id}/Datastreams", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationGet, "/v1.0/ObservedProperties{id}/Datastreams/{params}", HandleGetDatastreamsByObservedProperty},
			{models.HTTPOperationPost, "/v1.0/Datastreams", HandlePostDatastream},
			{models.HTTPOperationPost, "/v1.0/Things{id}/Datastreams", HandlePostDatastreamByThing},
			{models.HTTPOperationDelete, "/v1.0/Datastreams{id}", HandleDeleteDatastream},
			{models.HTTPOperationPatch, "/v1.0/Datastreams{id}", HandlePatchDatastream},
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
			odata.QueryOptionExpand, odata.QueryOptionSelect,
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
			{models.HTTPOperationGet, "/v1.0/ObservedProperties/{params}", HandleGetObservedProperties},
			{models.HTTPOperationGet, "/v1.0/ObservedProperties{id}", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/ObservedProperties{id}/{params}", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/ObservedProperties", HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/ObservedProperties/{params}", HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationPost, "/v1.0/ObservedProperties", HandlePostObservedProperty},
			{models.HTTPOperationDelete, "/v1.0/ObservedProperties{id}", HandleDeleteObservedProperty},
			{models.HTTPOperationPatch, "/v1.0/ObservedProperties{id}", HandlePatchObservedProperty},
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
			"encodingtype",
			"location",
			"Things",
			"HistoricalLocations",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/Locations", HandleGetLocations},
			{models.HTTPOperationGet, "/v1.0/Locations/{params}", HandleGetLocations},
			{models.HTTPOperationGet, "/v1.0/Locations{id}", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/Locations{id}/{params}", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}/Locations", HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}/Locations/{params}", HandleGetLocationsByHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/Things{id}/Locations", HandleGetLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/Things{id}/Locations/{params}", HandleGetLocationsByThing},
			{models.HTTPOperationPost, "/v1.0/Locations", HandlePostLocation},
			{models.HTTPOperationPost, "/v1.0/Things{id}/Locations", HandlePostLocationByThing},
			{models.HTTPOperationDelete, "/v1.0/Locations{id}", HandleDeleteLocation},
			{models.HTTPOperationPatch, "/v1.0/Locations{id}", HandlePatchLocation},
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
			odata.QueryOptionExpand, odata.QueryOptionSelect,
		},
		SupportedExpandParams: []string{
			"Datastream",
		},
		SupportedSelectParams: []string{
			"description",
			"encodingtype",
			"metadata",
			"Datastreams",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/Sensors", HandleGetSensors},
			{models.HTTPOperationGet, "/v1.0/Sensors/{params}", HandleGetSensors},
			{models.HTTPOperationGet, "/v1.0/Sensors{id}", HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/Sensors{id}/{params}", HandleGetSensor},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Sensor", HandleGetSensorByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Sensor/{params}", HandleGetSensorByDatastream},
			{models.HTTPOperationPost, "/v1.0/Sensors", HandlePostSensors},
			{models.HTTPOperationDelete, "/v1.0/Sensors{id}", HandleDeleteSensor},
			{models.HTTPOperationPatch, "/v1.0/Sensors{id}", HandlePatchSensor},
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
			odata.QueryOptionExpand, odata.QueryOptionSelect,
		},
		SupportedExpandParams: []string{
			"Datastream",
			"FeatureOfInterest",
		},
		SupportedSelectParams: []string{
			"description",
			"encodingtype",
			"feature",
			"Observations",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/Observations", HandleGetObservations},
			{models.HTTPOperationGet, "/v1.0/Observations/{params}", HandleGetObservations},
			{models.HTTPOperationGet, "/v1.0/Observations{id}", HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/Observations{id}/{params}", HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Observations", HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Observations/{params}", HandleGetObservationsByDatastream},
			{models.HTTPOperationGet, "/v1.0/FeatureOfInterest{id}/Observations", HandleGetObservationsByFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/FeatureOfInterest{id}/Observations/{params}", HandleGetObservationsByFeatureOfInterest},
			{models.HTTPOperationPost, "/v1.0/Observations", HandlePostObservation},
			{models.HTTPOperationPost, "/v1.0/Datastreams{id}/Observations", HandlePostObservationByDatastream},
			{models.HTTPOperationDelete, "/v1.0/Observations{id}", HandleDeleteObservation},
			{models.HTTPOperationPatch, "/v1.0/Observations{id}", HandlePatchObservation},
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
			odata.QueryOptionExpand, odata.QueryOptionSelect,
		},
		SupportedExpandParams: []string{
			"Observation",
		},
		SupportedSelectParams: []string{
			"description",
			"encodingtype",
			"feature",
			"Observations",
		},
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0/FeaturesOfInterest", HandleGetFeatureOfInterests},
			{models.HTTPOperationGet, "/v1.0/FeaturesOfInterest/{params}", HandleGetFeatureOfInterests},
			{models.HTTPOperationGet, "/v1.0/FeaturesOfInterest{id}", HandleGetFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/FeaturesOfInterest{id}/{params}", HandleGetFeatureOfInterest},
			{models.HTTPOperationGet, "/v1.0/Observations{id}/FeatureOfInterest", HandleGetFeatureOfInterestByObservation},
			{models.HTTPOperationGet, "/v1.0/Observations{id}/FeatureOfInterest/{params}", HandleGetFeatureOfInterestByObservation},
			{models.HTTPOperationPost, "/v1.0/FeaturesOfInterest", HandlePostFeatureOfInterest},
			{models.HTTPOperationDelete, "/v1.0/FeaturesOfInterest{id}", HandleDeleteFeatureOfInterest},
			{models.HTTPOperationPatch, "/v1.0/FeaturesOfInterest{id}", HandlePatchFeatureOfInterest},
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
			odata.QueryOptionExpand, odata.QueryOptionSelect,
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
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations/{params}", HandleGetHistoricalLocations},
			{models.HTTPOperationGet, "/v1.0/Things{id}/HistoricalLocations", HandleGetHistoricalLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/Things{id}/HistoricalLocations/{params}", HandleGetHistoricalLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/Locations{id}/HistoricalLocations", HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/Locations{id}/HistoricalLocations/{params}", HandleGetHistoricalLocationsByLocation},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}", HandleGetHistoricalLocation},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}/{params}", HandleGetHistoricalLocation},
			{models.HTTPOperationDelete, "/v1.0/HistoricalLocations{id}", HandleDeleteHistoricalLocations},
			{models.HTTPOperationPatch, "/v1.0/HistoricalLocations{id}", HandlePatchHistoricalLocations},
		},
	}
}
