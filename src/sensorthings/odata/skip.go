package odata

import (
	"net/http"
	"strconv"
)

// QuerySkip is used for retrieving records from a specified index of the entire record set.
// If set the request will return table entries after the provided index value.
type QuerySkip struct {
	QueryBase
	Index int
}

// Parse $skip values in QuerySkip, returns error if the supplied value is
// invalid (non-integer or < 0)
func (q *QuerySkip) Parse(value string) error {
	q.RawQuery = value
	i, err := strconv.Atoi(value)
	if err != nil || i < 0 {
		return CreateQueryError(QuerySkipInvalid, http.StatusBadRequest, value)
	}

	q.Index = i
	return nil
}

// IsValid always returns true, errors are already filtered out by parse
func (q *QuerySkip) IsValid() (bool, error) {
	return true, nil
}

// GetQueryOptionType returns the QueryOptionType for QuerySkip
func (q *QuerySkip) GetQueryOptionType() QueryOptionType {
	return QueryOptionSkip
}

// IsNil checks if *QuerySkip is nil
func (q *QuerySkip) IsNil() bool {
	return (q == nil)
}
