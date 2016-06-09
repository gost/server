package odata

// QueryFilter is used to perform conditional operations on the parameter values
// Count is used to retrieve the total number of items in a collection matching the requested entity.
type QueryFilter struct {
	QueryBase
	count bool
}

// Parse tries to parse the given filter
func (q *QueryFilter) Parse(value string) error {
	q.RawQuery = value
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
