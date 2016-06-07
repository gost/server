package odata

import (
	"net/url"
	"strings"

	"net/http"
)

// QueryOption contains user requested query information for retrieving objects
// from the SensorThings API
type QueryOption interface {
	Parse(string) error
	GetQueryOptionType() QueryOptionType
	IsNil() bool
}

// QueryOptionType holds a supported query option
type QueryOptionType int

// List of supported query options
const (
	_ QueryOptionType = iota
	QueryOptionExpand
	QueryOptionSelect
	QueryOptionOrderBy
	QueryOptionTop
	QueryOptionSkip
	QueryOptionCount
	QueryOptionFilter
	QueryOptionResultFormat
	QueryOptionRef
)

// QueryOptionValues is a list of names mapped to their QueryOptionType
var QueryOptionValues = []string{
	QueryOptionExpand:       "$expand",
	QueryOptionSelect:       "$select",
	QueryOptionOrderBy:      "$orderby",
	QueryOptionTop:          "$top",
	QueryOptionSkip:         "$skip",
	QueryOptionCount:        "$count",
	QueryOptionFilter:       "$filter",
	QueryOptionResultFormat: "$resultFormat",
	QueryOptionRef:          "$ref",
}

// String returns the string representation of the current QueryOptionType
func (q QueryOptionType) String() string {
	return QueryOptionValues[q]
}

// QueryOptions holds the parsed query information requested by the user
type QueryOptions struct {
	QueryTop          *QueryTop
	QuerySkip         *QuerySkip
	QuerySelect       *QuerySelect
	QueryExpand       *QueryExpand
	QueryOrderBy      *QueryOrderBy
	QueryCount        *QueryCount
	QueryFilter       *QueryFilter
	QueryResultFormat *QueryResultFormat
	QueryOptionRef    bool
}

// CreateQueryOptions parses the requested request parameters into usable Query options
// and filters out invalid requests
func CreateQueryOptions(queryParams url.Values) (*QueryOptions, []error) {
	qo := &QueryOptions{}

	var errorList []error
	err := &errorList

	for key, values := range queryParams {
		value := strings.Join(values, " ")

		switch key {
		case QueryOptionExpand.String():
			qo.QueryExpand = &QueryExpand{}
			ParseQueryOption(value, qo.QueryExpand, err)
			break
		case QueryOptionSelect.String():
			qo.QuerySelect = &QuerySelect{}
			ParseQueryOption(value, qo.QuerySelect, err)
			break
		case QueryOptionOrderBy.String():
			qo.QueryOrderBy = &QueryOrderBy{}
			ParseQueryOption(value, qo.QueryOrderBy, err)
			break
		case QueryOptionTop.String():
			qo.QueryTop = &QueryTop{}
			ParseQueryOption(value, qo.QueryTop, &errorList)
			break
		case QueryOptionSkip.String():
			qo.QuerySkip = &QuerySkip{}
			ParseQueryOption(value, qo.QuerySkip, err)
			break
		case QueryOptionCount.String():
			qo.QueryCount = &QueryCount{}
			ParseQueryOption(value, qo.QueryCount, err)
			break
		case QueryOptionFilter.String():
			qo.QueryFilter = &QueryFilter{}
			ParseQueryOption(value, qo.QueryFilter, err)
			break
		case QueryOptionResultFormat.String():
			qo.QueryResultFormat = &QueryResultFormat{}
			ParseQueryOption(value, qo.QueryResultFormat, err)
			break
		case QueryOptionRef.String():
			qo.QueryOptionRef = true
			break
		default:
			// Req 21 If a service does not support a system query option, it SHALL fail any request that contains the
			// unsupported option and SHOULD return 501 Not Implemented.
			errorList = append(errorList, CreateQueryError(QueryUnknown, http.StatusNotImplemented, key))
		}
	}

	if len(errorList) > 0 {
		return nil, errorList
	}

	return qo, nil
}

// ParseQueryOption tries to parse the user supplied values into the desired QueryOption
// if an error occurred in the parsing process a new error will be added
// to the supplied errorlist
func ParseQueryOption(value string, q QueryOption, errorList *[]error) {
	e := q.Parse(value)
	if e != nil {
		*errorList = append(*errorList, e)
	}
}
