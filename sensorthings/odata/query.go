package odata

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
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
		default:
			errorList = append(errorList, errors.New(CreateQueryError(QueryUnknown, key)))
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

// QueryTop is a filter for limiting the number of returns to be returned.
// Specifies a non-negative integer that limits the number of entities
// returned within a collection. The service must return the number of
// available entities up to, but not exceeding, the specified value.
type QueryTop struct {
	limit int
}

// Parse $top values in QueryTop, returns error if the supplied value is
// invalid (non-integer or < 0)
func (q *QueryTop) Parse(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil || i < 0 {
		return errors.New(CreateQueryError(QueryTopInvalid, value))
	}

	q.limit = i
	return nil
}

// GetQueryOptionType returns the QueryOptionType for QueryTop
func (q *QueryTop) GetQueryOptionType() QueryOptionType {
	return QueryOptionTop
}

// IsNil checks if *QueryTop is nil
func (q *QueryTop) IsNil() bool {
	if q == nil {
		return true
	}

	return false
}

// QuerySkip is used for retrieving records from a specified index of the entire record set.
// If set the request will return table entries after the provided index value.
type QuerySkip struct {
	index int
}

// Parse $skip values in QuerySkip, returns error if the supplied value is
// invalid (non-integer or < 0)
func (q *QuerySkip) Parse(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil || i < 0 {
		return errors.New(CreateQueryError(QuerySkipInvalid, value))
	}

	q.index = i
	return nil
}

// GetQueryOptionType returns the QueryOptionType for QuerySkip
func (q *QuerySkip) GetQueryOptionType() QueryOptionType {
	return QueryOptionSkip
}

// IsNil checks if *QuerySkip is nil
func (q *QuerySkip) IsNil() bool {
	if q == nil {
		return true
	}

	return false
}

// QuerySelect is used to return only the entity property values desired, this is used
// help to reduce the amount of information in a response from the server.
// If set, the result will include the specified property of the SensorThing entity object.
type QuerySelect struct {
	params []string
}

// Parse $select values in QuerySelect, at this stage we don't know
// if the select params are valid, this depends on Select values available
// for the used endpoint
func (q *QuerySelect) Parse(value string) error {
	q.params = strings.Split(value, ",")
	return nil
}

// IsValid checks if the given $select values are supported for the endpoint
func (q *QuerySelect) IsValid(values []string) (bool, error) {
	//ToDo: check if select values are valid for endpoint
	return true, nil
}

// GetQueryOptionType returns the QueryOptionType for QuerySelect
func (q *QuerySelect) GetQueryOptionType() QueryOptionType {
	return QueryOptionSelect
}

// IsNil checks if *QuerySelect is nil
func (q *QuerySelect) IsNil() bool {
	if q == nil {
		return true
	}

	return false
}

// QueryExpand is used to return a linked entity member’s full details.
// Expand retrieves the specified named property and represents it inline to the base entity.
type QueryExpand struct {
	params []string
}

// Parse splits the given values by the , delimiter and stores the params, if the delimiter is not
// a comma the IsValid will filter it out later on
func (q *QueryExpand) Parse(value string) error {
	// Enhancement: Check if the params contain Things, DataStreams, Sensors, etc error if a non valid param is found return error
	q.params = strings.Split(value, ",")
	return nil
}

// IsValid checks if the endpoint supports the expand params given by the user
func (q *QueryExpand) IsValid(values []string) (bool, error) {
	//ToDo: Check if the endpoint supports the expand params given by the user
	return true, nil
}

// GetQueryOptionType returns the QueryOptionType for QueryExpand
func (q *QueryExpand) GetQueryOptionType() QueryOptionType {
	return QueryOptionExpand
}

// IsNil checks if *QueryExpand is nil
func (q *QueryExpand) IsNil() bool {
	if q == nil {
		return true
	}

	return false
}

// QueryOrderBy is used to define the order of the results set be ascending (asc) or descending (desc) order.
// orderby Is used to specify which properties are used to order the collection of entities identified by the resource path.
type QueryOrderBy struct {
	suffix   string
	property string
}

// Parse tries to parse the OrderBy query into a suffix and property value
// if the suffix is not "asc" or "desc" or no property is given then
// parse will return an error
func (q *QueryOrderBy) Parse(value string) error {
	ob := strings.Split(value, " ")
	if len(ob) != 2 || (ob[1] != "asc" && ob[1] != "desc") || len(ob[0]) < 1 {
		return errors.New(CreateQueryError(QueryOrderByInvalid, value))
	}

	q.property = ob[0]
	q.suffix = ob[1]
	return nil
}

// IsValid checks if the given property value in the request is valid for the
// used endpoint, returns an error if not supported
func (q *QueryOrderBy) IsValid(values []string) (bool, error) {
	//ToDo: implement check property to order by (is property supported on endpoint)
	return true, nil
}

// GetQueryOptionType returns the QueryOptionType for QueryOrderBy
func (q *QueryOrderBy) GetQueryOptionType() QueryOptionType {
	return QueryOptionOrderBy
}

// IsNil checks if *QueryOrderBy is nil
func (q *QueryOrderBy) IsNil() bool {
	if q == nil {
		return true
	}

	return false
}

// QueryCount is used to get a total count for the number of each entity in the system.
// Count is used to retrieve the total number of items in a collection matching the requested entity.
type QueryCount struct {
	count bool
}

// Parse tries to parse the given count query to a bool, if parse fails to convert
// it to a bool it will return an error
func (q *QueryCount) Parse(value string) error {
	b, e := strconv.ParseBool(value)
	if e != nil {
		return errors.New(CreateQueryError(QueryCountInvalid, value))
	}

	q.count = b
	return nil
}

// GetQueryOptionType returns the QueryOptionType for QueryCount
func (q *QueryCount) GetQueryOptionType() QueryOptionType {
	return QueryOptionCount
}

// IsNil checks if *QueryCount is nil
func (q *QueryCount) IsNil() bool {
	if q == nil {
		return true
	}

	return false
}

// QueryFilter is used to perform conditional operations on the parameter values
// Count is used to retrieve the total number of items in a collection matching the requested entity.
type QueryFilter struct {
	count bool
}

// Parse tries to parse the given filter
func (q *QueryFilter) Parse(value string) error {
	//ToDo: implement, documentation not clear enough on how to implement this
	return nil
}

// GetQueryOptionType returns the QueryOptionType for QueryFilter
func (q *QueryFilter) GetQueryOptionType() QueryOptionType {
	return QueryOptionFilter
}

// IsNil checks if *QueryFilter is nil
func (q *QueryFilter) IsNil() bool {
	if q == nil {
		return true
	}

	return false
}

// QueryResultFormat is used to return Observations in a data array format, a components section
// is returned in the response to describe the order of returned values.
type QueryResultFormat struct {
	//ToDo: make a type out of format, not done because there is currently only the dataArray option
	format string
}

// Parse tries to parse the given data format, if the data format is not supported
// it will return an error
func (q *QueryResultFormat) Parse(value string) error {
	if value != "dataArray" {
		return errors.New(CreateQueryError(QueryResultFormatInvalid, value))
	}

	q.format = value
	return nil
}

// GetQueryOptionType returns the QueryOptionType for QueryResultFormat
func (q *QueryResultFormat) GetQueryOptionType() QueryOptionType {
	return QueryOptionResultFormat
}

// IsNil checks if *QueryResultFormat is nil
func (q *QueryResultFormat) IsNil() bool {
	if q == nil {
		return true
	}

	return false
}
