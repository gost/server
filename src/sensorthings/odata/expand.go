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
	QueryOptions    QueryOptions
	ExpandOperation *ExpandOperation
}

func (e *ExpandOperation) Create(eo string) error {
	expandOperation := *e
	l := strings.Split(eo, "/") // split for example $expand=Things/Locations/HistoricalLocations
	for i, s := range l {
		if i == 0 { // first entry = w.Entity
			es := strings.Split(s, "(")
			if et, err := entities.EntityTypeFromString(es[0]); err != nil {
				return fmt.Errorf("Unable to create expand query, unknown entity %s", es[0])
			} else {
				expandOperation.Entity = et
				if len(es) > 1 { // ToDo: parse QueryOptions here found in () chars
					innerQuery := fmt.Sprintf("(%s", es[1])
					fmt.Println("inner query for " + et.ToString() + ": " + innerQuery)
				}
			}
		} else { // this is a multiple level expand

		}
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
	return nil

	// Work in progress, remove Params if finished
	l1 := strings.Split(value, ",") // split layer 1, for example $expand=Observations/Things,Sensor,ObservedProperty
	for _, sl1 := range l1 {
		eo := ExpandOperation{}
		if err := eo.Create(sl1); err != nil {
			return err
		} else {
			q.Operations = append(q.Operations, eo)
			name := *eo.Entity

			fmt.Printf("ExpandOperation created: %v\n", name)
		}
	}

	fmt.Println("Done parsing expand query")
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
