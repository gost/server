package sensorthings

import (
	"errors"
	"fmt"
	"net/http"
)

// HTTPOperation describes the HTTP operation such as GET POST DELETE.
type HTTPOperation string

// HTTPOperation "enumeration".
const (
	HTTPOperationGet    HTTPOperation = "GET"
	HTTPOperationPost   HTTPOperation = "POST"
	HTTPOperationPatch  HTTPOperation = "PATCH"
	HTTPOperationDelete HTTPOperation = "DELETE"
)

// Endpoint contains all information for creating and handling a main SensorThings endpoint.
// A SensorThings endpoint contains multiple EndpointOperations
// Endpoint can be marshalled to JSON for returning endpoint information requested
// by the user: http://www.sensorup.com/docs/#resource-path
type Endpoint struct {
	Name                  string              `json:"name"` // Name of the endpoint
	URL                   string              `json:"url"`  // External URL to the endpoint
	Operations            []EndpointOperation `json:"-"`
	SupportedQueryOptions []QueryOptionType   `json:"-"`
	SupportedExpandParams []string            `json:"-"`
	SupportedSelectParams []string            `json:"-"`
}

// EndpointOperation contains the needed information to create an endpoint in the HTTP.Router
type EndpointOperation struct {
	OperationType HTTPOperation
	Path          string //relative path to the endpoint for example: /v1.0/myendpoint/
	Handler       HTTPHandler
}

// HTTPHandler func defines the format of the handler to process the incoming request
type HTTPHandler func(w http.ResponseWriter, r *http.Request, e *Endpoint, a *SensorThingsApi)

// AreQueryOptionsSupported checks if the endpoint supports the requested query and if
// the values are valid for the given endpoint
func (e *Endpoint) AreQueryOptionsSupported(queryOptions *QueryOptions) (bool, []error) {
	if queryOptions == nil {
		return true, nil
	}

	var errorList []error
	if queryOptions.QueryTop != nil {
		if !e.SupportsQueryOptionType(queryOptions.QueryTop.GetQueryOptionType()) {
			errorList = append(errorList, errors.New(CreateQueryError(QueryTopInvalid, e.Name)))
		}
	}

	if queryOptions.QuerySkip != nil {
		if !e.SupportsQueryOptionType(queryOptions.QuerySkip.GetQueryOptionType()) {
			errorList = append(errorList, errors.New(CreateQueryError(QuerySkipNotAvailable, queryOptions.QuerySkip.GetQueryOptionType().String(), e.Name)))
		}
	}

	if queryOptions.QuerySelect != nil {
		if !e.SupportsQueryOptionType(queryOptions.QuerySelect.GetQueryOptionType()) {
			//ToDo: Create error message
		}
	}

	if queryOptions.QueryExpand != nil {
		if !e.SupportsQueryOptionType(queryOptions.QueryExpand.GetQueryOptionType()) {
			//ToDo: Create error message
		}
	}

	if queryOptions.QueryOrderBy != nil {
		if !e.SupportsQueryOptionType(queryOptions.QueryOrderBy.GetQueryOptionType()) {
			//ToDo: Create error message
		}
	}

	if queryOptions.QueryCount != nil {
		if !e.SupportsQueryOptionType(queryOptions.QueryCount.GetQueryOptionType()) {
			//ToDo: Create error message
		}
	}

	if queryOptions.QueryFilter != nil {
		if !e.SupportsQueryOptionType(queryOptions.QueryFilter.GetQueryOptionType()) {
			//ToDo: Create error message
		}
	}

	if queryOptions.QueryResultFormat != nil {
		if !e.SupportsQueryOptionType(queryOptions.QueryResultFormat.GetQueryOptionType()) {
			//ToDo: Create error message
		}
	}

	if errorList != nil {
		return false, errorList
	}

	return true, nil
}

// SupportsQueryOptionType checks if a given QueryOptionType is configured for the endpoint
func (e *Endpoint) SupportsQueryOptionType(queryOptionType QueryOptionType) bool {
	for _, qo := range e.SupportedQueryOptions {
		if qo == queryOptionType {
			return true
		}
	}

	return false
}

// CreateEndPoints creates the pre-defined endpoint config, the config contains all endpoint info
// describing the SupportedQueryOptions (if needed) and EndpointOperation for each endpoint
// parameter externalURL is the URL where the GOST service can be reached, main endpoint urls
// are generated based upon this URL
func CreateEndPoints(externalURL string) []Endpoint {
	endpoints := []Endpoint{
		{
			Name: "Version",
			Operations: []EndpointOperation{
				{HTTPOperationGet, "/Version", HandleVersion},
			},
		},
		{
			Name: "Root",
			Operations: []EndpointOperation{
				{HTTPOperationGet, "/v1.0", HandleAPIRoot},
			},
		},
		{
			Name: "Things",
			SupportedQueryOptions: []QueryOptionType{
				QueryOptionTop, QueryOptionSkip, QueryOptionOrderBy, QueryOptionCount, QueryOptionResultFormat,
				QueryOptionExpand, QueryOptionSelect,
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
			Operations: []EndpointOperation{
				{HTTPOperationGet, "/v1.0/Things", HandleGetThings},
				{HTTPOperationGet, "/v1.0/Things{id}", HandleGetThingById},
				{HTTPOperationPost, "/v1.0/Things", HandlePostThing},
				{HTTPOperationDelete, "/v1.0/Things{id}", HandleDeleteThing},
				{HTTPOperationPatch, "/v1.0/Things{id}", HandlePatchThing},
			},
		},
		{
			Name: "Datastreams",
			SupportedQueryOptions: []QueryOptionType{
				QueryOptionTop, QueryOptionSkip, QueryOptionOrderBy, QueryOptionCount, QueryOptionResultFormat,
				QueryOptionExpand, QueryOptionSelect,
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
			Operations: []EndpointOperation{
				{HTTPOperationGet, "/v1.0/Datastreams", HandleGetDatastreams},
				{HTTPOperationGet, "/v1.0/Datastreams{id}", HandleGetDatastreamById},
				{HTTPOperationGet, "/v1.0/Things{id}/Datastreams", HandleGetDatastreamsByThing},
				{HTTPOperationPost, "/v1.0/Datastreams", HandlePostDatastream},
				{HTTPOperationPost, "/v1.0/Things{id}/Datastreams", HandlePostAndLinkDatastream},
				{HTTPOperationDelete, "/v1.0/Datastreams{id}", HandleDeleteDatastream},
				{HTTPOperationPatch, "/v1.0/Things{id}/Datastreams", HandlePatchDatastream},
			},
		},
		{
			Name: "ObservedProperties",
			SupportedQueryOptions: []QueryOptionType{
				QueryOptionTop, QueryOptionSkip, QueryOptionOrderBy, QueryOptionCount, QueryOptionResultFormat,
				QueryOptionExpand, QueryOptionSelect,
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
			Operations: []EndpointOperation{
				{HTTPOperationGet, "/v1.0/ObservedProperties", HandleGetObservedProperties},
				{HTTPOperationGet, "/v1.0/ObservedProperty{id}", HandleGetObservedPropertyById},
				{HTTPOperationGet, "/v1.0/Datastreams{id}/ObservedProperty", HandleGetObservedPropertyFromDatastream},
				{HTTPOperationPost, "/v1.0/ObservedProperty", HandlePostObservedProperty},
				{HTTPOperationDelete, "/v1.0/ObservedProperty{id}", HandleDeleteObservedProperty},
				{HTTPOperationPatch, "/v1.0/ObservedProperty{id}", HandlePatchObservedProperty},
			},
		},
		{
			Name: "Locations",
			SupportedQueryOptions: []QueryOptionType{
				QueryOptionTop, QueryOptionSkip, QueryOptionOrderBy, QueryOptionCount, QueryOptionResultFormat,
				QueryOptionExpand, QueryOptionSelect, QueryOptionFilter,
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
			Operations: []EndpointOperation{
				{HTTPOperationGet, "/v1.0/Locations", HandleGetLocations},
				{HTTPOperationGet, "/v1.0/Locations{id}", HandleGetLocationById},
				{HTTPOperationPost, "/v1.0/Locations", HandlePostLocation},
				{HTTPOperationPost, "/v1.0/Things{id}/Locations", HandlePostAndLinkLocation},
				{HTTPOperationDelete, "/v1.0/Locations{id}", HandleDeleteLocation},
				{HTTPOperationPatch, "/v1.0/Locations{id}", HandlePatchLocation},
			},
		},
		{
			Name: "Sensors",
			SupportedQueryOptions: []QueryOptionType{
				QueryOptionTop, QueryOptionSkip, QueryOptionOrderBy, QueryOptionCount, QueryOptionResultFormat,
				QueryOptionExpand, QueryOptionSelect,
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
			Operations: []EndpointOperation{
				{HTTPOperationGet, "/v1.0/Sensors", HandleGetSensors},
				{HTTPOperationGet, "/v1.0/Sensors{id}", HandleGetSensorById},
				{HTTPOperationPost, "/v1.0/Sensors", HandlePostSensors},
				{HTTPOperationDelete, "/v1.0/Sensors{id}", HandleDeleteSensor},
				{HTTPOperationPatch, "/v1.0/Sensors{id}", HandlePatchSensor},
			},
		},
		{
			Name: "Observations",
			SupportedQueryOptions: []QueryOptionType{
				QueryOptionTop, QueryOptionSkip, QueryOptionOrderBy, QueryOptionCount, QueryOptionResultFormat,
				QueryOptionExpand, QueryOptionSelect,
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
			Operations: []EndpointOperation{
				{HTTPOperationGet, "/v1.0/Observations", nil},
				{HTTPOperationGet, "/v1.0/Observations{id}", nil},
				{HTTPOperationGet, "/v1.0/Datastreams{id}/Observations", nil},
				{HTTPOperationPost, "/v1.0/Observations", nil},
				{HTTPOperationPost, "/v1.0/Datastreams{id}/Observations", nil},
				{HTTPOperationDelete, "/v1.0/Observations{id}", nil},
				{HTTPOperationPatch, "/v1.0/Observations{id}", nil},
			},
		},
		{
			Name: "FeaturesOfInterest",
			SupportedQueryOptions: []QueryOptionType{
				QueryOptionTop, QueryOptionSkip, QueryOptionOrderBy, QueryOptionCount, QueryOptionResultFormat,
				QueryOptionExpand, QueryOptionSelect,
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
			Operations: []EndpointOperation{
				{HTTPOperationGet, "/v1.0/FeaturesOfInterest", nil},
				{HTTPOperationGet, "/v1.0/FeaturesOfInterest{id}", nil},
				{HTTPOperationPost, "/v1.0/FeaturesOfInterest", nil},
				{HTTPOperationDelete, "/v1.0/FeaturesOfInterest{id}", nil},
				{HTTPOperationPatch, "/v1.0/FeaturesOfInterest{id}", nil},
			},
		},
		{
			Name: "HistoricalLocations",
			SupportedQueryOptions: []QueryOptionType{
				QueryOptionTop, QueryOptionSkip, QueryOptionOrderBy, QueryOptionCount, QueryOptionResultFormat,
				QueryOptionExpand, QueryOptionSelect,
			},
			SupportedExpandParams: []string{
				"locations",
				"thing",
			},
			SupportedSelectParams: []string{
				"time",
			},
			Operations: []EndpointOperation{
				{HTTPOperationGet, "/v1.0/HistoricalLocations", nil},
				{HTTPOperationGet, "/v1.0/HistoricalLocations{id}", nil},
				{HTTPOperationDelete, "/v1.0/HistoricalLocations{id}", nil},
				{HTTPOperationPatch, "/v1.0/HistoricalLocations{id}", nil},
			},
		},
	}

	// Generate url for endpoints
	for k := range endpoints {
		rp := endpoints[k]
		rp.URL = fmt.Sprintf("%s/%s/%s", externalURL, API_PREFIX, fmt.Sprintf("%v", rp.Name))
		endpoints[k] = rp
	}

	return endpoints
}
