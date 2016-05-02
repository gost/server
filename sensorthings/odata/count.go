package odata

import (
	"net/http"
	"strconv"
)

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
		return CreateQueryError(QueryCountInvalid, http.StatusBadRequest, value)
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
