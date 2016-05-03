package rest

import (
	"fmt"

	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
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

func createVersion(externalURL string) *Endpoint {
	return &Endpoint{
		Name: "Version",
		URL:  fmt.Sprintf("%s/%s", externalURL, "Version"),
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/Version", HandleVersion},
		},
	}
}

func createRoot(externalURL string) *Endpoint {
	return &Endpoint{
		Name: "Root",
		URL:  fmt.Sprintf("%s/%s", externalURL, "v1.0"),
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0", HandleAPIRoot},
		},
	}
}

func createThings(externalURL string) *Endpoint {
	return &Endpoint{
		Name: "Things",
		URL:  fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "Things")),
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
			{models.HTTPOperationGet, "/v1.0/Things{id}", HandleGetThing},
			{models.HTTPOperationPost, "/v1.0/Things", HandlePostThing},
			{models.HTTPOperationDelete, "/v1.0/Things{id}", HandleDeleteThing},
			{models.HTTPOperationPatch, "/v1.0/Things{id}", HandlePatchThing},
		},
	}
}

func createDatastreams(externalURL string) *Endpoint {
	return &Endpoint{
		Name: "Datastreams",
		URL:  fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "Datastreams")),
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
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}", HandleGetDatastream},
			{models.HTTPOperationGet, "/v1.0/Things{id}/Datastreams", HandleGetDatastreamsByThing},
			{models.HTTPOperationGet, "/v1.0/Sensors{id}/Datastreams", HandleGetDatastreamsBySensor},
			{models.HTTPOperationPost, "/v1.0/Datastreams", HandlePostDatastream},
			{models.HTTPOperationPost, "/v1.0/Things{id}/Datastreams", HandlePostDatastreamByThing},
			{models.HTTPOperationDelete, "/v1.0/Datastreams{id}", HandleDeleteDatastream},
			{models.HTTPOperationPatch, "/v1.0/Datastreams{id}", HandlePatchDatastream},
		},
	}
}

func createObservedProperties(externalURL string) *Endpoint {
	return &Endpoint{
		Name: "ObservedProperties",
		URL:  fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "ObservedProperties")),
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
			{models.HTTPOperationGet, "/v1.0/ObservedProperty{id}", HandleGetObservedProperty},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/ObservedProperty", HandleGetObservedPropertyByDatastream},
			{models.HTTPOperationPost, "/v1.0/ObservedProperty", HandlePostObservedProperty},
			{models.HTTPOperationDelete, "/v1.0/ObservedProperty{id}", HandleDeleteObservedProperty},
			{models.HTTPOperationPatch, "/v1.0/ObservedProperty{id}", HandlePatchObservedProperty},
		},
	}
}

func createLocations(externalURL string) *Endpoint {
	return &Endpoint{
		Name: "Locations",
		URL:  fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "Locations")),
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
			{models.HTTPOperationGet, "/v1.0/Locations{id}", HandleGetLocation},
			{models.HTTPOperationGet, "/v1.0/Things{id}/Locations", HandleGetLocationsByThing},
			{models.HTTPOperationPost, "/v1.0/Locations", HandlePostLocation},
			{models.HTTPOperationPost, "/v1.0/Things{id}/Locations", HandlePostLocationByThing},
			{models.HTTPOperationDelete, "/v1.0/Locations{id}", HandleDeleteLocation},
			{models.HTTPOperationPatch, "/v1.0/Locations{id}", HandlePatchLocation},
		},
	}
}

func createSensors(externalURL string) *Endpoint {
	return &Endpoint{
		Name: "Sensors",
		URL:  fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "Sensors")),
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
			{models.HTTPOperationGet, "/v1.0/Sensors{id}", HandleGetSensor},
			{models.HTTPOperationPost, "/v1.0/Sensors", HandlePostSensors},
			{models.HTTPOperationDelete, "/v1.0/Sensors{id}", HandleDeleteSensor},
			{models.HTTPOperationPatch, "/v1.0/Sensors{id}", HandlePatchSensor},
		},
	}
}

func createObservations(externalURL string) *Endpoint {
	return &Endpoint{
		Name: "Observations",
		URL:  fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "Observations")),
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
			{models.HTTPOperationGet, "/v1.0/Observations{id}", HandleGetObservation},
			{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Observations", HandleGetObservationsByDatastream},
			{models.HTTPOperationPost, "/v1.0/Observations", HandlePostObservation},
			{models.HTTPOperationPost, "/v1.0/Datastreams{id}/Observations", HandlePostObservationByDatastream},
			{models.HTTPOperationDelete, "/v1.0/Observations{id}", HandleDeleteObservation},
			{models.HTTPOperationPatch, "/v1.0/Observations{id}", HandlePatchObservation},
		},
	}
}

func createFeaturesOfInterest(externalURL string) *Endpoint {
	return &Endpoint{
		Name: "FeaturesOfInterest",
		URL:  fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "FeaturesOfInterest")),
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
			{models.HTTPOperationGet, "/v1.0/FeaturesOfInterest{id}", HandleGetFeatureOfInterest},
			{models.HTTPOperationPost, "/v1.0/FeaturesOfInterest", HandlePostFeatureOfInterest},
			{models.HTTPOperationDelete, "/v1.0/FeaturesOfInterest{id}", HandleDeleteFeatureOfInterest},
			{models.HTTPOperationPatch, "/v1.0/FeaturesOfInterest{id}", HandlePatchFeatureOfInterest},
		},
	}
}

func createHistoricalLocations(externalURL string) *Endpoint {
	return &Endpoint{
		Name: "HistoricalLocations",
		URL:  fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, fmt.Sprintf("%v", "HistoricalLocations")),
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
			{models.HTTPOperationGet, "/v1.0/Things{id}/HistoricalLocations", HandleGetHistoricalLocationsByThing},
			{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}", HandleGetHistoricalLocation},
			{models.HTTPOperationDelete, "/v1.0/HistoricalLocations{id}", HandleDeleteHistoricalLocations},
			{models.HTTPOperationPatch, "/v1.0/HistoricalLocations{id}", HandlePatchHistoricalLocations},
		},
	}
}
