package odata

import (
	"fmt"
	"github.com/geodan/gost/src/sensorthings/entities"
	"net/http"
	"strings"
)

// work in progress
type ExpandOperation struct {
	Entity          *entities.EntityType
	QueryOptions    *QueryOptions
	ExpandOperation *ExpandOperation
}

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

	if et, err := entities.EntityTypeFromString(fp); err != nil {
		return []error{fmt.Errorf("Unable to create expand query, unknown entity %s", fp)}
	} else {
		e.Entity = et
	}

	if len(queryString) > 0 {
		splitQuery := strings.Split(queryString, ";")
		values := map[string]string{}

		for _, q := range splitQuery {
			kvp := strings.Split(q, "=")
			if len(kvp) != 2 {
				return []error{fmt.Errorf("Invalid query (%s) inside $expand %s", queryString, e.Entity.ToString())}
			}
			values[kvp[0]] = kvp[1]
		}

		if qo, err := CreateQueryOptions(values); err != nil {
			return err
		} else {
			e.QueryOptions = qo
		}
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

	// Work in progress, remove Params if finished
	l1 := strings.Split(value, ",") // split layer 1, for example $expand=Observations/Things,Sensor,ObservedProperty
	for _, sl1 := range l1 {
		eo := ExpandOperation{}
		if err := eo.Create(sl1); err != nil {
			return err[0]
		} else {
			q.Operations = append(q.Operations, eo)
		}
	}

	// debug print
	//for _, test := range q.Operations {
	//	displayExpandOperation(test)
	//}

	return nil
}

// debug
func displayExpandOperation(e ExpandOperation) {
	fmt.Println("====== top level expand: " + *e.Entity + " ======")
	printQueries(e.QueryOptions)
	if e.ExpandOperation != nil {
		displayExpandOperationLower(e.ExpandOperation, 1)
	}
}

func displayExpandOperationLower(e *ExpandOperation, level int) {
	fmt.Printf("--- sub expand %v: %s ---\n", level, *e.Entity)
	printQueries(e.QueryOptions)
}

func printQueries(qo *QueryOptions) {
	if qo != nil {
		fmt.Println(" Queries:")
		test2 := *qo.QuerySelect
		fmt.Println("    HIEP HOI " + test2.RawQuery)
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
