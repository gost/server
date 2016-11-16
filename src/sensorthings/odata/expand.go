package odata

import (
	"net/http"
	"strings"
)

// QueryExpand is used to return a linked entity member’s full details.
// Expand retrieves the specified named property and represents it inline to the base entity.
type QueryExpand struct {
	QueryBase
	Params []string
}

// Parse splits the given values by the , delimiter and stores the params, if the delimiter is not
// a comma the IsValid will filter it out later on
func (q *QueryExpand) Parse(value string) error {
	q.RawQuery = value
	q.Params = strings.Split(value, ",")
	return nil
}

// IsValid checks if the endpoint supports the expand params given by the user
func (q *QueryExpand) IsValid(values []string, endpointName string) (bool, error) {
	for _, value := range q.Params {
		found := false
		for _, param := range values {
			if param == value {
				found = true
				break
			}
		}

		if !found {
			return false, CreateQueryError(QueryExpandAvailable, http.StatusBadRequest, value, endpointName)
		}
	}

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
