package odata

import (
	"log"
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

	printTest(q.Predicate)

	return nil
}

func printTest(p *Predicate) {
	log.Printf("operator: %v\n", p.Operator.ToString())

	switch v := p.Subject.(type) {
	case *Predicate:
		printTest(v)
	case string:
		log.Printf("subject: %v\n", v)
	case float64:
		log.Printf("subject: %v\n", v)
	}

	switch v := p.Value.(type) {
	case *Predicate:
		printTest(v)
	case string:
		log.Printf("value: %v\n", v)
	case float64:
		log.Printf("value: %v\n", v)
	}
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
