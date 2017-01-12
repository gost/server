package rest

import (
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
	"net/http"
)

// Endpoint contains all information for creating and handling a main SensorThings endpoint.
// A SensorThings endpoint contains multiple EndpointOperations
// Endpoint can be marshalled to JSON for returning endpoint information requested
// by the user: http://www.sensorup.com/docs/#resource-path
type Endpoint struct {
	Name                  string                     `json:"name"` // Name of the endpoint
	URL                   string                     `json:"url"`  // External URL to the endpoint
	EntityType            entities.EntityType        `json:"-"`
	OutputInfo            bool                       `json:"-"` //Output when BasePathInfo is requested by the user
	Operations            []models.EndpointOperation `json:"-"`
	SupportedQueryOptions []odata.QueryOptionType    `json:"-"`
	SupportedExpandParams []string                   `json:"-"`
	SupportedSelectParams []string                   `json:"-"`
}

// GetName returns the endpoint name
func (e *Endpoint) GetName() string {
	return e.Name
}

// ShowOutputInfo returns true if the endpoint should output his info when BasePathInfo is requested
func (e *Endpoint) ShowOutputInfo() bool {
	return e.OutputInfo
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
// ToDo: refactor AreQueryOptionsSupported remove dependencies
func (e *Endpoint) AreQueryOptionsSupported(queryOptions *odata.QueryOptions) (bool, []error) {
	if queryOptions == nil {
		return true, nil
	}

	var errorList []error
	qo := *queryOptions
	checkQueryOptionSupported(e, qo.QueryTop, &errorList, odata.CreateQueryError(odata.QueryNotAvailable, http.StatusNotImplemented, qo.QueryTop.GetQueryOptionType().String(), e.Name))
	checkQueryOptionSupported(e, qo.QuerySkip, &errorList, odata.CreateQueryError(odata.QueryNotAvailable, http.StatusNotImplemented, qo.QuerySkip.GetQueryOptionType().String(), e.Name))
	checkQueryOptionSupported(e, qo.QuerySelect, &errorList, odata.CreateQueryError(odata.QueryNotAvailable, http.StatusNotImplemented, qo.QuerySelect.GetQueryOptionType().String(), e.Name))
	checkQueryOptionSupported(e, qo.QueryExpand, &errorList, odata.CreateQueryError(odata.QueryNotAvailable, http.StatusNotImplemented, qo.QueryExpand.GetQueryOptionType().String(), e.Name))
	checkQueryOptionSupported(e, qo.QueryOrderBy, &errorList, odata.CreateQueryError(odata.QueryNotAvailable, http.StatusNotImplemented, qo.QueryOrderBy.GetQueryOptionType().String(), e.Name))
	checkQueryOptionSupported(e, qo.QueryCount, &errorList, odata.CreateQueryError(odata.QueryNotAvailable, http.StatusNotImplemented, qo.QueryCount.GetQueryOptionType().String(), e.Name))
	checkQueryOptionSupported(e, qo.QueryFilter, &errorList, odata.CreateQueryError(odata.QueryNotAvailable, http.StatusNotImplemented, qo.QueryFilter.GetQueryOptionType().String(), e.Name))
	checkQueryOptionSupported(e, qo.QueryResultFormat, &errorList, odata.CreateQueryError(odata.QueryNotAvailable, http.StatusNotImplemented, qo.QueryResultFormat.GetQueryOptionType().String(), e.Name))

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

	// check if query is valid for endpoint
	switch v := q.(type) {
	case *odata.QueryExpand:
		if _, err = v.IsValid(e.Name); err != nil {
			*errorList = append(errors, err)
		}
		for _, operation := range v.Operations {
			ep := Endpoints[operation.Entity.GetEntityType()]
			supported, qoe := ep.AreQueryOptionsSupported(operation.QueryOptions)
			if !supported {
				*errorList = append(errors, qoe...)
			}
		}
	case *odata.QuerySelect:
		if _, err = v.IsValid(e.SupportedSelectParams); err != nil {
			*errorList = append(errors, err)
		}
	case *odata.QueryOrderBy:
		//ToDo
	case *odata.QueryTop:
		if _, err = v.IsValid(); err != nil {
			*errorList = append(errors, err)
		}
	case *odata.QuerySkip:
		if _, err = v.IsValid(); err != nil {
			*errorList = append(errors, err)
		}
	case *odata.QueryCount:
		if _, err = v.IsValid(); err != nil {
			*errorList = append(errors, err)
		}
	case *odata.QueryFilter:
		if _, err = v.IsValid(); err != nil {
			*errorList = append(errors, err)
		}
	case *odata.QueryResultFormat:
		if _, err = v.IsValid(); err != nil {
			*errorList = append(errors, err)
		}
	default:
		//set error, unknown
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
