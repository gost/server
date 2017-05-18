package odata

import (
	"fmt"
	"net/http"
	"strings"
)

//ToDo: Implement sub properties such as Person/Name

// OrderType defines an Order By type
type OrderType string

func (o OrderType) ToString() string {
	return fmt.Sprintf("%s", o)
}

// OrderType is an "enumeration" of the OrderTypes
const (
	OrderTypeASC  OrderType = "asc"
	OrderTypeDESC OrderType = "desc"
)

// QueryOrderBy is used to define the order of the results set be ascending (asc) or descending (desc) order.
// orderby Is used to specify which properties are used to order the collection of entities identified by the resource path.
type QueryOrderBy struct {
	QueryBase
	Queries []OrderBy
}

// OrderBy describes the property on which property to order and the OrderType (asc or desc)
type OrderBy struct {
	Property  string
	OrderType OrderType
}

// Parse tries to parse the OrderBy query into a list of OrderBy queries
// if the OrderType is not "asc" or "desc" or no property is given then
// parse will return an error
func (q *QueryOrderBy) Parse(value string) error {
	q.RawQuery = value
	queries := strings.Split(value, ",")
	if len(queries) < 1 || queries[0] == "" {
		return CreateQueryError(QueryOrderByInvalid, http.StatusBadRequest, value)
	}

	for _, query := range queries {
		obParts := strings.Split(query, " ")
		if len(obParts) > 2 || (len(obParts) == 2 && (strings.ToLower(obParts[1]) != "asc" && strings.ToLower(obParts[1]) != "desc")) {
			return CreateQueryError(QueryOrderByInvalid, http.StatusBadRequest, value)
		}

		ob := OrderBy{
			Property: obParts[0],
		}

		if len(obParts) == 1 || strings.ToLower(obParts[1]) == "asc" {
			ob.OrderType = OrderTypeASC
		} else {
			ob.OrderType = OrderTypeDESC
		}

		q.Queries = append(q.Queries, ob)
	}

	return nil
}

// IsValid checks if the given property value in the request is valid for the
// used endpoint, returns an error if not supported
// values = available properties of an entity
func (q *QueryOrderBy) IsValid(values []string) (bool, error) {
	for _, query := range q.Queries {
		found := false
		for _, s := range values {
			if query.Property == s {
				found = true
			}
		}
		if !found {
			return false, CreateQueryError(QueryOrderByInvalid, http.StatusBadRequest, query.Property)
		}
	}

	return true, nil
}

// GetQueryOptionType returns the QueryOptionType for QueryOrderBy
func (q *QueryOrderBy) GetQueryOptionType() QueryOptionType {
	return QueryOptionOrderBy
}

// IsNil checks if *QueryOrderBy is nil
func (q *QueryOrderBy) IsNil() bool {
	if q == nil {
		return true
	}

	return false
}
