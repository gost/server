package odata

import "net/http"

// QueryResultFormat is used to return Observations in a data array format, a components section
// is returned in the response to describe the order of returned values.
type QueryResultFormat struct {
	QueryBase
	//ToDo: make a type out of format, not done because there is currently only the dataArray option
	format string
}

// Parse tries to parse the given data format, if the data format is not supported
// it will return an error
func (q *QueryResultFormat) Parse(value string) error {
	q.RawQuery = value
	if value != "dataArray" {
		return CreateQueryError(QueryResultFormatInvalid, http.StatusBadRequest, value)
	}

	q.format = value
	return nil
}

// IsValid always returns true, errors are already filtered out by parse
func (q *QueryResultFormat) IsValid() (bool, error) {
	return true, nil
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
