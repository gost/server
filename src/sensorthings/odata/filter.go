package odata

import (
	"net/http"
)

// QueryFilter is used to perform conditional operations on the parameter values
// Count is used to retrieve the total number of items in a collection matching the requested entity.
type QueryFilter struct {
	QueryBase
	Predicate *Predicate
}

// Parse tries to parse the given filter
func (q *QueryFilter) Parse(value string) error {
	var err error
	q.RawQuery = value
	q.Predicate, err = ParseODATAFilter(value)
	if err != nil {
		return CreateQueryError(QueryfilterFormatInvalid, http.StatusBadRequest, value)
	}

	return nil
}

// IsValid always returns true, errors are already filtered out by parse
func (q *QueryFilter) IsValid() (bool, error) {
	return true, nil
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
