package odata

import (
	"net/http"
	"strconv"
)

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
		return CreateQueryError(QuerySkipInvalid, http.StatusBadRequest, value)
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
