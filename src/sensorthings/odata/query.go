package odata

import (
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
	QueryOptionValue
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
	QueryOptionValue:        "$value",
}

// String returns the string representation of the current QueryOptionType
func (q QueryOptionType) String() string {
	return QueryOptionValues[q]
}

// QueryBase is the base for all oData queries
type QueryBase struct {
	RawQuery string
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
	QueryOptionValue  bool
}

// CreateQueryOptions parses the requested request parameters into usable Query options
// and filters out invalid requests
func CreateQueryOptions(queryParams map[string]string) (*QueryOptions, []error) {
	qo := &QueryOptions{}

	var errorList []error
	err := &errorList

	for key, value := range queryParams {
		switch key {
		case QueryOptionExpand.String():
			qo.QueryExpand = &QueryExpand{}
			ParseQueryOption(value, qo.QueryExpand, err)
		case QueryOptionSelect.String():
			qo.QuerySelect = &QuerySelect{}
			ParseQueryOption(value, qo.QuerySelect, err)
		case QueryOptionOrderBy.String():
			qo.QueryOrderBy = &QueryOrderBy{}
			ParseQueryOption(value, qo.QueryOrderBy, err)
		case QueryOptionTop.String():
			qo.QueryTop = &QueryTop{}
			ParseQueryOption(value, qo.QueryTop, &errorList)
		case QueryOptionSkip.String():
			qo.QuerySkip = &QuerySkip{}
			ParseQueryOption(value, qo.QuerySkip, err)
		case QueryOptionCount.String():
			qo.QueryCount = &QueryCount{}
			ParseQueryOption(value, qo.QueryCount, err)
		case QueryOptionFilter.String():
			qo.QueryFilter = &QueryFilter{}
			ParseQueryOption(value, qo.QueryFilter, err)
		case QueryOptionResultFormat.String():
			qo.QueryResultFormat = &QueryResultFormat{}
			ParseQueryOption(value, qo.QueryResultFormat, err)
		case QueryOptionRef.String():
			qo.QueryOptionRef = true
		case QueryOptionValue.String():
			qo.QueryOptionValue = true
		default:
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
// to the supplied error list
func ParseQueryOption(value string, q QueryOption, errorList *[]error) {
	e := q.Parse(value)
	if e != nil {
		*errorList = append(*errorList, e)
	}
}
