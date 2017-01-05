package odata

import (
	"fmt"
	"github.com/geodan/gost/src/sensorthings/entities"
	"net/http"
	"strings"
)

// ExpandOperation holds information on a received $expand query
type ExpandOperation struct {
	Entity          entities.Entity
	QueryOptions    *QueryOptions
	ExpandOperation *ExpandOperation
}

// Create tries to construct the ExpandOperation from given query string
func (e *ExpandOperation) Create(eo string) []error {
	slashIndex := strings.Index(eo, "/")
	var fp, trail string

	if slashIndex != -1 {
		trail = eo[slashIndex+1:]
		fp = eo[:slashIndex]
	} else {
		fp = eo
	}

	var queryString string
	queryIndexStart := strings.Index(fp, "(")
	if queryIndexStart != -1 {
		queryString = fp[queryIndexStart+1 : len(fp)-1]
		fp = fp[:queryIndexStart]
	}

	et, err := entities.EntityFromString(fp)
	if err != nil {
		return []error{fmt.Errorf("Unable to create expand query, unknown entity %s", fp)}
	}

	e.Entity = et

	if len(queryString) > 0 {
		splitQuery := strings.Split(queryString, ";")
		values := map[string]string{}

		for _, q := range splitQuery {
			kvp := strings.Split(q, "=")
			if len(kvp) != 2 {
				return []error{fmt.Errorf("Invalid query (%s) inside $expand %s", queryString, e.Entity.GetEntityType().ToString())}
			}
			values[kvp[0]] = kvp[1]
		}

		qo, err := CreateQueryOptions(values)
		if err != nil {
			return err
		}
		e.QueryOptions = qo
	}

	if len(trail) > 0 {
		neo := &ExpandOperation{}
		neo.Create(trail)
		e.ExpandOperation = neo
	}

	return nil
}

// QueryExpand is used to return a linked entity member’s full details.
// Expand retrieves the specified named property and represents it inline to the base entity.
type QueryExpand struct {
	QueryBase
	Params     []string
	Operations []ExpandOperation
}

// Parse splits the given values by the , delimiter and stores the params, if the delimiter is not
// a comma the IsValid will filter it out later on
func (q *QueryExpand) Parse(value string) error {
	q.RawQuery = value
	q.Params = strings.Split(value, ",")

	l1 := strings.Split(value, ",") // split layer 1, for example $expand=Observations/Things,Sensor,ObservedProperty
	for _, sl1 := range l1 {
		eo := ExpandOperation{}
		if err := eo.Create(sl1); err != nil {
			return err[0]
		}

		q.Operations = append(q.Operations, eo)
	}

	return nil
}

// debug
func displayExpandOperation(e ExpandOperation) {
	fmt.Println("====== top level expand: " + e.Entity.GetEntityType().ToString() + " ======")
	printQueries(e.QueryOptions)
	if e.ExpandOperation != nil {
		displayExpandOperationLower(e.ExpandOperation, 1)
	}
}

func displayExpandOperationLower(e *ExpandOperation, level int) {
	fmt.Printf("--- sub expand %v: %s ---\n", level, e.Entity.GetEntityType().ToString())
	printQueries(e.QueryOptions)
}

func printQueries(qo *QueryOptions) {
	if qo != nil {
		fmt.Println(" Queries:")
		test2 := *qo.QuerySelect
		fmt.Println("    " + test2.RawQuery)
	} else {
		fmt.Println(" No queries defined")
	}
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
