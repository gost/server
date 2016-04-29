package rest

import (
	"errors"

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
	qo := *queryOptions
	checkQueryOptionSupported(e, qo.QueryTop, &errorList, errors.New(odata.CreateQueryError(odata.QueryTopNotAvailable, qo.QueryTop.GetQueryOptionType().String(), e.Name)))
	checkQueryOptionSupported(e, qo.QuerySkip, &errorList, errors.New(odata.CreateQueryError(odata.QuerySkipNotAvailable, qo.QuerySkip.GetQueryOptionType().String(), e.Name)))

	//ToDo: Create error message for queries below
	checkQueryOptionSupported(e, qo.QuerySelect, &errorList, errors.New(odata.CreateQueryError(odata.QuerySkipNotAvailable, qo.QuerySkip.GetQueryOptionType().String(), e.Name)))
	checkQueryOptionSupported(e, qo.QueryExpand, &errorList, errors.New(odata.CreateQueryError(odata.QuerySkipNotAvailable, qo.QueryExpand.GetQueryOptionType().String(), e.Name)))
	checkQueryOptionSupported(e, qo.QueryOrderBy, &errorList, errors.New(odata.CreateQueryError(odata.QuerySkipNotAvailable, qo.QueryOrderBy.GetQueryOptionType().String(), e.Name)))
	checkQueryOptionSupported(e, qo.QueryCount, &errorList, errors.New(odata.CreateQueryError(odata.QuerySkipNotAvailable, qo.QueryCount.GetQueryOptionType().String(), e.Name)))
	checkQueryOptionSupported(e, qo.QueryFilter, &errorList, errors.New(odata.CreateQueryError(odata.QuerySkipNotAvailable, qo.QueryFilter.GetQueryOptionType().String(), e.Name)))
	checkQueryOptionSupported(e, qo.QueryResultFormat, &errorList, errors.New(odata.CreateQueryError(odata.QuerySkipNotAvailable, qo.QueryResultFormat.GetQueryOptionType().String(), e.Name)))

	if errorList != nil {
		return false, errorList
	}

	return true, nil
}

// checkQueryOptionSupported checks if an QueryOption is supported on an endpoint, if not
// an error will be added to the error list
func checkQueryOptionSupported(e *Endpoint, q odata.QueryOption, errorList *[]error, err error) {
	if q.IsNil() {
		return
	}

	errors := *errorList
	if !e.SupportsQueryOptionType(q.GetQueryOptionType()) {
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
