package odata

import (
	"strings"

	"errors"
	"fmt"
)

// QuerySelect is used to return only the entity property values desired, this is used
// help to reduce the amount of information in a response from the server.
// If set, the result will include the specified property of the SensorThing entity object.
type QuerySelect struct {
	Params []string
}

// Parse $select values in QuerySelect, at this stage we don't know
// if the select params are valid, this depends on Select values available
// for the used endpoint
func (q *QuerySelect) Parse(value string) error {
	q.Params = strings.Split(value, ",")
	return nil
}

// IsValid checks if the given $select values are supported for the endpoint
func (q *QuerySelect) IsValid(values []string) (bool, error) {
	for _, rp := range q.Params {
		found := false
		for _, hp := range values {
			if rp == hp {
				found = true
			}
		}

		if !found {
			return false, errors.New(fmt.Sprintf("Parameter %s not supported", rp))
		}

	}

	return true, nil
}

// GetQueryOptionType returns the QueryOptionType for QuerySelect
func (q *QuerySelect) GetQueryOptionType() QueryOptionType {
	return QueryOptionSelect
}

// IsNil checks if *QuerySelect is nil
func (q *QuerySelect) IsNil() bool {
	if q == nil || len(q.Params) == 0 {
		return true
	}

	return false
}
