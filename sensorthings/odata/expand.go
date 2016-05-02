package odata

import "strings"

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
