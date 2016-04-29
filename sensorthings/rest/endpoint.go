package rest

import (
	"errors"
	"fmt"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
)

// Endpoint contains all information for creating and handling a main SensorThings endpoint.
// A SensorThings endpoint contains multiple EndpointOperations
// Endpoint can be marshalled to JSON for returning endpoint information requested
// by the user: http://www.sensorup.com/docs/#resource-path
type Endpoint struct {
	Name                  string                     `json:"name"` // Name of the endpoint
	URL                   string                     `json:"url"`  // External URL to the endpoint
	Operations            []models.EndpointOperation `json:"-"`
	SupportedQueryOptions []odata.QueryOptionType    `json:"-"`
	SupportedExpandParams []string                   `json:"-"`
	SupportedSelectParams []string                   `json:"-"`
}

// GetName returns the endpoint name
func (e *Endpoint) GetName() string {
	return e.Name
}

// GetURL returns the external url
func (e *Endpoint) GetURL() string {
	return e.URL
}

// GetOperations returns all operations for this endpoint such as GET, POST
func (e *Endpoint) GetOperations() []models.EndpointOperation {
	return e.Operations
}

// GetSupportedQueryOptions returns al possible odata query options on this endpoint
func (e *Endpoint) GetSupportedQueryOptions() []odata.QueryOptionType {
	return e.SupportedQueryOptions
}

// GetSupportedExpandParams returns which entities can be expanded
func (e *Endpoint) GetSupportedExpandParams() []string {
	return e.SupportedExpandParams
}

// GetSupportedSelectParams returns the supported select parameters for this endpoint
func (e *Endpoint) GetSupportedSelectParams() []string {
	return e.SupportedSelectParams
}

// AreQueryOptionsSupported checks if the endpoint supports the requested query and if
// the values are valid for the given endpoint
func (e *Endpoint) AreQueryOptionsSupported(queryOptions *odata.QueryOptions) (bool, []error) {
	if queryOptions == nil {
		return true, nil
	}

	var errorList []error

	check(e, queryOptions.QueryTop, &errorList, errors.New(odata.CreateQueryError(odata.QueryTopInvalid, e.Name)), func() bool {
		return e.SupportsQueryOptionType(queryOptions.QueryTop.GetQueryOptionType())
	})

	check(e, queryOptions.QuerySkip, &errorList, errors.New(odata.CreateQueryError(odata.QuerySkipNotAvailable, queryOptions.QuerySkip.GetQueryOptionType().String(), e.Name)), func() bool {
		return e.SupportsQueryOptionType(queryOptions.QuerySkip.GetQueryOptionType())
	})

	//ToDo: Create error message for queries below
	check(e, queryOptions.QuerySelect, &errorList, errors.New(odata.CreateQueryError(odata.QuerySkipNotAvailable, queryOptions.QuerySkip.GetQueryOptionType().String(), e.Name)), func() bool {
		return e.SupportsQueryOptionType(queryOptions.QuerySelect.GetQueryOptionType())
	})

	check(e, queryOptions.QueryExpand, &errorList, errors.New(odata.CreateQueryError(odata.QuerySkipNotAvailable, queryOptions.QuerySkip.GetQueryOptionType().String(), e.Name)), func() bool {
		return e.SupportsQueryOptionType(queryOptions.QueryExpand.GetQueryOptionType())
	})

	check(e, queryOptions.QueryOrderBy, &errorList, errors.New(odata.CreateQueryError(odata.QuerySkipNotAvailable, queryOptions.QuerySkip.GetQueryOptionType().String(), e.Name)), func() bool {
		return e.SupportsQueryOptionType(queryOptions.QueryOrderBy.GetQueryOptionType())
	})

	check(e, queryOptions.QueryCount, &errorList, errors.New(odata.CreateQueryError(odata.QuerySkipNotAvailable, queryOptions.QuerySkip.GetQueryOptionType().String(), e.Name)), func() bool {
		return e.SupportsQueryOptionType(queryOptions.QueryCount.GetQueryOptionType())
	})

	check(e, queryOptions.QueryFilter, &errorList, errors.New(odata.CreateQueryError(odata.QuerySkipNotAvailable, queryOptions.QuerySkip.GetQueryOptionType().String(), e.Name)), func() bool {
		return e.SupportsQueryOptionType(queryOptions.QueryFilter.GetQueryOptionType())
	})

	check(e, queryOptions.QueryResultFormat, &errorList, errors.New(odata.CreateQueryError(odata.QuerySkipNotAvailable, queryOptions.QuerySkip.GetQueryOptionType().String(), e.Name)), func() bool {
		return e.SupportsQueryOptionType(queryOptions.QueryResultFormat.GetQueryOptionType())
	})

	if errorList != nil {
		return false, errorList
	}

	return true, nil
}

func check(e *Endpoint, i interface{}, errorList *[]error, err error, supported func() bool) {
	if i == nil {
		return
	}

	errors := *errorList
	if !supported() {
		*errorList = append(errors, err)
	}
}

// SupportsQueryOptionType checks if a given QueryOptionType is configured for the endpoint
func (e *Endpoint) SupportsQueryOptionType(queryOptionType odata.QueryOptionType) bool {
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
func CreateEndPoints(externalURL string) []models.Endpoint {
	endpoints := []models.Endpoint{
		&Endpoint{
			Name: "Version",
			URL:  fmt.Sprintf("%s/%s", externalURL, "Version"),
			Operations: []models.EndpointOperation{
				{models.HTTPOperationGet, "/Version", HandleVersion},
			},
		},
		&Endpoint{
			Name: "Root",
			URL:  fmt.Sprintf("%s/%s", externalURL, "v1.0"),
			Operations: []models.EndpointOperation{
				{models.HTTPOperationGet, "/v1.0", HandleAPIRoot},
			},
		},
		&Endpoint{
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
				{models.HTTPOperationGet, "/v1.0/Things{id}", HandleGetThingByID},
				{models.HTTPOperationPost, "/v1.0/Things", HandlePostThing},
				{models.HTTPOperationDelete, "/v1.0/Things{id}", HandleDeleteThing},
				{models.HTTPOperationPatch, "/v1.0/Things{id}", HandlePatchThing},
			},
		},
		&Endpoint{
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
				{models.HTTPOperationGet, "/v1.0/Datastreams{id}", HandleGetDatastreamByID},
				{models.HTTPOperationGet, "/v1.0/Things{id}/Datastreams", HandleGetDatastreamsByThing},
				{models.HTTPOperationPost, "/v1.0/Datastreams", HandlePostDatastream},
				{models.HTTPOperationPost, "/v1.0/Things{id}/Datastreams", HandlePostAndLinkDatastream},
				{models.HTTPOperationDelete, "/v1.0/Datastreams{id}", HandleDeleteDatastream},
				{models.HTTPOperationPatch, "/v1.0/Things{id}/Datastreams", HandlePatchDatastream},
			},
		},
		&Endpoint{
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
				{models.HTTPOperationGet, "/v1.0/ObservedProperty{id}", HandleGetObservedPropertyByID},
				{models.HTTPOperationGet, "/v1.0/Datastreams{id}/ObservedProperty", HandleGetObservedPropertyFromDatastream},
				{models.HTTPOperationPost, "/v1.0/ObservedProperty", HandlePostObservedProperty},
				{models.HTTPOperationDelete, "/v1.0/ObservedProperty{id}", HandleDeleteObservedProperty},
				{models.HTTPOperationPatch, "/v1.0/ObservedProperty{id}", HandlePatchObservedProperty},
			},
		},
		&Endpoint{
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
				{models.HTTPOperationGet, "/v1.0/Locations{id}", HandleGetLocationByID},
				{models.HTTPOperationPost, "/v1.0/Locations", HandlePostLocation},
				{models.HTTPOperationPost, "/v1.0/Things{id}/Locations", HandlePostAndLinkLocation},
				{models.HTTPOperationDelete, "/v1.0/Locations{id}", HandleDeleteLocation},
				{models.HTTPOperationPatch, "/v1.0/Locations{id}", HandlePatchLocation},
			},
		},
		&Endpoint{
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
				{models.HTTPOperationGet, "/v1.0/Sensors{id}", HandleGetSensorByID},
				{models.HTTPOperationPost, "/v1.0/Sensors", HandlePostSensors},
				{models.HTTPOperationDelete, "/v1.0/Sensors{id}", HandleDeleteSensor},
				{models.HTTPOperationPatch, "/v1.0/Sensors{id}", HandlePatchSensor},
			},
		},
		&Endpoint{
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
				{models.HTTPOperationGet, "/v1.0/Observations", nil},
				{models.HTTPOperationGet, "/v1.0/Observations{id}", nil},
				{models.HTTPOperationGet, "/v1.0/Datastreams{id}/Observations", nil},
				{models.HTTPOperationPost, "/v1.0/Observations", nil},
				{models.HTTPOperationPost, "/v1.0/Datastreams{id}/Observations", nil},
				{models.HTTPOperationDelete, "/v1.0/Observations{id}", nil},
				{models.HTTPOperationPatch, "/v1.0/Observations{id}", nil},
			},
		},
		&Endpoint{
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
				{models.HTTPOperationGet, "/v1.0/FeaturesOfInterest", nil},
				{models.HTTPOperationGet, "/v1.0/FeaturesOfInterest{id}", nil},
				{models.HTTPOperationPost, "/v1.0/FeaturesOfInterest", nil},
				{models.HTTPOperationDelete, "/v1.0/FeaturesOfInterest{id}", nil},
				{models.HTTPOperationPatch, "/v1.0/FeaturesOfInterest{id}", nil},
			},
		},
		&Endpoint{
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
				{models.HTTPOperationGet, "/v1.0/HistoricalLocations", nil},
				{models.HTTPOperationGet, "/v1.0/HistoricalLocations{id}", nil},
				{models.HTTPOperationDelete, "/v1.0/HistoricalLocations{id}", nil},
				{models.HTTPOperationPatch, "/v1.0/HistoricalLocations{id}", nil},
			},
		},
	}

	return endpoints
}
